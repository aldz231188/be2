package ports

import "context"

type AuthService interface {
	// Create(ctx context.Context, userid, name, surename string) (string, error)
	Register(ctx context.Context, login, password string) (*TokenPair, error)
	// Login(ctx context.Context) (*LoginResponse, error)
	// Refresh(ctx context.Context) (*RefreshResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	// LogoutAll(ctx context.Context) (*empty.Empty, error)
}
