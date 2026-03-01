package app

import (
	"be2/services/authsvc/internal/config"
	"be2/services/authsvc/internal/domain"
	"be2/services/authsvc/internal/jwtkeys"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
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
	AccessToken      string
	AccessExpiresAt  int64
	RefreshToken     string
	RefreshExpiresAt int64
	SessionId        string
}

type TokenClaims struct {
	TokenType    string
	TokenVersion int32
	SessionID    string
	jwt.RegisteredClaims
}

type AuthService interface {
	Authenticate(ctx context.Context, login, password string) (TokenPair, error)
	Register(ctx context.Context, login, password string) (TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
	ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error)
}

type authService struct {
	users      domain.UserRepo
	sessions   domain.SessionRepo
	privateKey *jwtkeys.RSAKey
	issuer     string
	audience   string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

var _ AuthService = (*authService)(nil)

func NewAuthService(users domain.UserRepo, sessions domain.SessionRepo, keys *jwtkeys.RSAKey, cfg config.Config) AuthService {
	return &authService{
		users:      users,
		sessions:   sessions,
		privateKey: keys,
		issuer:     cfg.JWTIssuer,
		audience:   cfg.JWTAudience,
		accessTTL:  15 * time.Minute,
		refreshTTL: 30 * 24 * time.Hour,
	}
}

func (s *authService) Authenticate(ctx context.Context, login, password string) (TokenPair, error) {
	user, err := s.users.GetByLogin(ctx, login)
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

func (s *authService) Register(ctx context.Context, login, password string) (TokenPair, error) {
	if strings.TrimSpace(login) == "" || strings.TrimSpace(password) == "" {
		return TokenPair{}, ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, err
	}

	user, err := s.users.CreateUser(ctx, domain.User{Login: login, PasswordHash: string(hash)})
	if err != nil {
		return TokenPair{}, err
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

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
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
	accessClaims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        sessionID,
		Issuer:    s.issuer,
	}
	if s.audience != "" {
		accessClaims.Audience = []string{s.audience}
	}
	access, err := s.signClaims(TokenClaims{
		TokenType:        tokenTypeAccess,
		TokenVersion:     user.TokenVersion,
		SessionID:        sessionID,
		RegisteredClaims: accessClaims,
	})
	if err != nil {
		return TokenPair{}, err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(refreshExpires),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        sessionID,
		Issuer:    s.issuer,
	}
	if s.audience != "" {
		refreshClaims.Audience = []string{s.audience}
	}
	refresh, err := s.signClaims(TokenClaims{
		TokenType:        tokenTypeRefresh,
		TokenVersion:     user.TokenVersion,
		SessionID:        sessionID,
		RegisteredClaims: refreshClaims,
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
	if s.privateKey == nil || s.privateKey.Public == nil {
		return nil, ErrInvalidToken
	}
	parserOptions := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(s.issuer),
	}
	if s.audience != "" {
		parserOptions = append(parserOptions, jwt.WithAudience(s.audience))
	}

	parsed, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, ErrInvalidToken
		}
		if kid, ok := token.Header["kid"].(string); ok && kid != s.privateKey.KID {
			return nil, ErrInvalidToken
		}
		return s.privateKey.Public, nil
	}, parserOptions...)
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
	if s.privateKey == nil || s.privateKey.Private == nil {
		return "", fmt.Errorf("missing signing key")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.privateKey.KID
	signed, err := token.SignedString(s.privateKey.Private)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func hashJTI(jti string) string {
	sum := sha256.Sum256([]byte(jti))
	return hex.EncodeToString(sum[:])
}
