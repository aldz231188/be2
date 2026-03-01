package app

import (
	"be2/services/authsvc/internal/config"
	"be2/services/authsvc/internal/domain"
	"be2/services/authsvc/internal/jwtkeys"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	users        map[uuid.UUID]domain.User
	usersByLogin map[string]uuid.UUID
	createErr    error
	getErr       error
	incremented  uuid.UUID
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[uuid.UUID]domain.User), usersByLogin: make(map[string]uuid.UUID)}
}

func (r *fakeUserRepo) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	if r.getErr != nil {
		return domain.User{}, r.getErr
	}
	id, ok := r.usersByLogin[login]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return r.users[id], nil
}

func (r *fakeUserRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if r.getErr != nil {
		return domain.User{}, r.getErr
	}
	user, ok := r.users[id]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *fakeUserRepo) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	if r.createErr != nil {
		return domain.User{}, r.createErr
	}
	user.ID = uuid.New()
	r.users[user.ID] = user
	r.usersByLogin[user.Login] = user.ID
	return user, nil
}

func (r *fakeUserRepo) IncrementTokenVersion(ctx context.Context, id uuid.UUID) error {
	user := r.users[id]
	user.TokenVersion++
	r.users[id] = user
	r.incremented = id
	return nil
}

type fakeSessionRepo struct {
	sessions map[string]domain.Session
}

func newFakeSessionRepo() *fakeSessionRepo {
	return &fakeSessionRepo{sessions: make(map[string]domain.Session)}
}

func (r *fakeSessionRepo) CreateSession(ctx context.Context, session domain.Session) error {
	r.sessions[session.JTIHash] = session
	return nil
}

func (r *fakeSessionRepo) GetSessionByHash(ctx context.Context, hash string) (domain.Session, error) {
	s, ok := r.sessions[hash]
	if !ok {
		return domain.Session{}, domain.ErrSessionNotFound
	}
	return s, nil
}

func (r *fakeSessionRepo) RevokeSession(ctx context.Context, hash string) error {
	session, ok := r.sessions[hash]
	if !ok {
		return domain.ErrSessionNotFound
	}
	now := time.Now()
	session.RevokedAt = &now
	r.sessions[hash] = session
	return nil
}

func (r *fakeSessionRepo) RevokeSessionsByUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	for k, s := range r.sessions {
		if s.UserID == userID {
			s.RevokedAt = &now
			r.sessions[k] = s
		}
	}
	return nil
}

func newTestService(t *testing.T, users domain.UserRepo, sessions domain.SessionRepo) (AuthService, *jwtkeys.RSAKey) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	rsaKey, err := jwtkeys.NewRSAKey(&config.Secrets{JWTPrivateKey: string(pemKey)})
	if err != nil {
		t.Fatalf("build rsa key: %v", err)
	}

	cfg := config.Config{
		JWTIssuer:   "authsvc",
		JWTAudience: "be2",
	}
	return NewAuthService(users, sessions, rsaKey, cfg), rsaKey
}

func TestRegisterAndAuthenticate(t *testing.T) {
	users := newFakeUserRepo()
	sessions := newFakeSessionRepo()
	service, rsaKey := newTestService(t, users, sessions)

	pair, err := service.Register(context.Background(), "user", "password")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("expected tokens after registration")
	}

	parsed, err := jwt.ParseWithClaims(pair.AccessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return rsaKey.Public, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}), jwt.WithIssuer("authsvc"), jwt.WithAudience("be2"))
	if err != nil || !parsed.Valid {
		t.Fatalf("failed to parse access token: %v", err)
	}

	authPair, err := service.Authenticate(context.Background(), "user", "password")
	if err != nil {
		t.Fatalf("authenticate failed: %v", err)
	}
	if authPair.AccessToken == pair.AccessToken {
		t.Fatal("expected new tokens on authenticate")
	}
}

func TestAuthenticateInvalidPassword(t *testing.T) {
	users := newFakeUserRepo()
	sessions := newFakeSessionRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := domain.User{ID: uuid.New(), Login: "user", PasswordHash: string(hash)}
	users.users[user.ID] = user
	users.usersByLogin[user.Login] = user.ID

	service, _ := newTestService(t, users, sessions)
	if _, err := service.Authenticate(context.Background(), "user", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestRefreshAndSessionRevocation(t *testing.T) {
	users := newFakeUserRepo()
	sessions := newFakeSessionRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := domain.User{ID: uuid.New(), Login: "user", PasswordHash: string(hash)}
	users.users[user.ID] = user
	users.usersByLogin[user.Login] = user.ID

	service, _ := newTestService(t, users, sessions)
	pair, err := service.Authenticate(context.Background(), "user", "correct")
	if err != nil {
		t.Fatalf("authenticate failed: %v", err)
	}

	oldSessions := len(sessions.sessions)
	newPair, err := service.Refresh(context.Background(), pair.RefreshToken)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if len(sessions.sessions) <= oldSessions {
		t.Fatal("expected new session to be created")
	}
	if newPair.RefreshToken == pair.RefreshToken {
		t.Fatal("expected refreshed token to differ")
	}
}

func TestLogoutCurrent(t *testing.T) {
	users := newFakeUserRepo()
	sessions := newFakeSessionRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := domain.User{ID: uuid.New(), Login: "user", PasswordHash: string(hash)}
	users.users[user.ID] = user
	users.usersByLogin[user.Login] = user.ID

	service, _ := newTestService(t, users, sessions)
	pair, err := service.Authenticate(context.Background(), "user", "correct")
	if err != nil {
		t.Fatalf("authenticate failed: %v", err)
	}

	if err := service.Logout(context.Background(), pair.RefreshToken); err != nil {
		t.Fatalf("logout current failed: %v", err)
	}
}

func TestLogoutAll(t *testing.T) {
	users := newFakeUserRepo()
	sessions := newFakeSessionRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := domain.User{ID: uuid.New(), Login: "user", PasswordHash: string(hash)}
	users.users[user.ID] = user
	users.usersByLogin[user.Login] = user.ID

	service, _ := newTestService(t, users, sessions)
	pair, err := service.Authenticate(context.Background(), "user", "correct")
	if err != nil {
		t.Fatalf("authenticate failed: %v", err)
	}

	if err := service.LogoutAll(context.Background(), pair.RefreshToken); err != nil {
		t.Fatalf("logout all failed: %v", err)
	}
	if users.users[user.ID].TokenVersion != 1 {
		t.Fatal("expected token version to increment")
	}
}

func TestValidateAccessToken(t *testing.T) {
	users := newFakeUserRepo()
	sessions := newFakeSessionRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := domain.User{ID: uuid.New(), Login: "user", PasswordHash: string(hash)}
	users.users[user.ID] = user
	users.usersByLogin[user.Login] = user.ID

	service, _ := newTestService(t, users, sessions)
	pair, err := service.Authenticate(context.Background(), "user", "correct")
	if err != nil {
		t.Fatalf("authenticate failed: %v", err)
	}

	claims, err := service.ValidateAccessToken(context.Background(), pair.AccessToken)
	if err != nil {
		t.Fatalf("validate access failed: %v", err)
	}
	if claims.TokenType != tokenTypeAccess {
		t.Fatalf("expected access token type, got %s", claims.TokenType)
	}
}
