package app

import (
	"be2/internal/config"
	"be2/internal/domain"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	TokenType    string `json:"type"`
	TokenVersion int32  `json:"tv"`
	SessionID    string `json:"sid"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (TokenPair, error)
	LogoutCurrent(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
	ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error)
}

type authService struct {
	users      domain.UserRepo
	sessions   domain.SessionRepo
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

var _ AuthService = (*authService)(nil)

func NewAuthService(users domain.UserRepo, sessions domain.SessionRepo, secrets *config.Secrets) AuthService {
	return &authService{
		users:      users,
		sessions:   sessions,
		secret:     []byte(secrets.JWTSecret),
		accessTTL:  15 * time.Minute,
		refreshTTL: 30 * 24 * time.Hour,
	}
}

func (s *authService) Authenticate(ctx context.Context, username, password string) (TokenPair, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return TokenPair{}, ErrInvalidCredentials
		}
		return TokenPair{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	return s.issueSessionTokens(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	claims, err := s.parseToken(refreshToken, tokenTypeRefresh)
	if err != nil {
		return TokenPair{}, err
	}

	user, err := s.ensureSessionValid(ctx, claims)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.sessions.RevokeSession(ctx, hashJTI(claims.SessionID)); err != nil {
		return TokenPair{}, err
	}

	return s.issueSessionTokens(ctx, user)
}

func (s *authService) LogoutCurrent(ctx context.Context, refreshToken string) error {
	claims, err := s.parseToken(refreshToken, tokenTypeRefresh)
	if err != nil {
		return err
	}

	if _, err := s.ensureSessionValid(ctx, claims); err != nil {
		return err
	}

	return s.sessions.RevokeSession(ctx, hashJTI(claims.SessionID))
}

func (s *authService) LogoutAll(ctx context.Context, refreshToken string) error {
	claims, err := s.parseToken(refreshToken, tokenTypeRefresh)
	if err != nil {
		return err
	}

	user, err := s.ensureSessionValid(ctx, claims)
	if err != nil {
		return err
	}

	if err := s.sessions.RevokeSessionsByUser(ctx, user.ID); err != nil {
		return err
	}
	return s.users.IncrementTokenVersion(ctx, user.ID)
}

func (s *authService) ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error) {
	claims, err := s.parseToken(token, tokenTypeAccess)
	if err != nil {
		return nil, err
	}

	if _, err := s.ensureSessionValid(ctx, claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *authService) issueSessionTokens(ctx context.Context, user domain.User) (TokenPair, error) {
	sessionID := uuid.NewString()
	refreshExpires := time.Now().Add(s.refreshTTL)
	if err := s.sessions.CreateSession(ctx, domain.Session{
		JTIHash:   hashJTI(sessionID),
		UserID:    user.ID,
		ExpiresAt: refreshExpires,
	}); err != nil {
		return TokenPair{}, err
	}

	now := time.Now()
	access, err := s.signClaims(TokenClaims{
		TokenType:    tokenTypeAccess,
		TokenVersion: user.TokenVersion,
		SessionID:    sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        sessionID,
		},
	})
	if err != nil {
		return TokenPair{}, err
	}

	refresh, err := s.signClaims(TokenClaims{
		TokenType:    tokenTypeRefresh,
		TokenVersion: user.TokenVersion,
		SessionID:    sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(refreshExpires),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        sessionID,
		},
	})
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *authService) ensureSessionValid(ctx context.Context, claims *TokenClaims) (domain.User, error) {
	session, err := s.sessions.GetSessionByHash(ctx, hashJTI(claims.SessionID))
	if err != nil {
		return domain.User{}, ErrInvalidToken
	}

	if session.UserID.String() != claims.Subject {
		return domain.User{}, ErrInvalidToken
	}

	if session.RevokedAt != nil || time.Now().After(session.ExpiresAt) {
		return domain.User{}, ErrInvalidToken
	}

	user, err := s.users.GetByID(ctx, session.UserID)
	if err != nil {
		return domain.User{}, ErrInvalidToken
	}

	if user.TokenVersion != claims.TokenVersion {
		return domain.User{}, ErrInvalidToken
	}

	return user, nil
}

func (s *authService) parseToken(token, expectedType string) (*TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(*TokenClaims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *authService) signClaims(claims TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func hashJTI(jti string) string {
	sum := sha256.Sum256([]byte(jti))
	return hex.EncodeToString(sum[:])
}
