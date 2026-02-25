package authz

import (
	"context"
	"fmt"
	"strconv"
	"time"

	keyfunc "github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"

	"be2/internal/config"
)

type Claims struct {
	UID   int64    `json:"uid"`
	Ver   int      `json:"ver"`
	Roles []string `json:"roles,omitempty"`
	Type  string   `json:"TokenType,omitempty"`
	jwt.RegisteredClaims
}

type Validator struct {
	kf  keyfunc.Keyfunc // интерфейс, не *указатель
	iss string
	aud string
}

func NewValidator(lc fx.Lifecycle, cfg config.Config) (*Validator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// keyfunc v3: создаём интерфейс Keyfunc, который сам обновляет JWKS в фоне.
	kf, err := keyfunc.NewDefaultCtx(ctx, []string{cfg.JWKSURL})
	if err != nil {
		cancel()
		return nil, err
	}
	lc.Append(fx.Hook{OnStop: func(context.Context) error { cancel(); return nil }})

	return &Validator{
		kf:  kf,
		iss: cfg.JWTIssuer,
		aud: cfg.JWTAudience,
	}, nil
}

func (v *Validator) ParseAccess(tokenStr string) (*Claims, error) {
	var c Claims

	opts := []jwt.ParserOption{
		jwt.WithIssuer(v.iss),
		jwt.WithLeeway(10 * time.Second),
		jwt.WithValidMethods([]string{"RS256", "RS512", "EdDSA"}),
	}
	if v.aud != "" {
		opts = append(opts, jwt.WithAudience(v.aud))
	}

	// ВАЖНО: используем пакетную функцию, не (*Parser).ParseWithClaims
	tok, err := jwt.ParseWithClaims(tokenStr, &c, v.kf.Keyfunc, opts...)
	if err != nil || !tok.Valid {
		return nil, err
	}

	if c.Type != "" && c.Type != "access" {
		return nil, fmt.Errorf("invalid token type: %s", c.Type)
	}

	if c.UID == 0 {
		uid, convErr := strconv.ParseInt(c.Subject, 10, 64)
		if convErr != nil || uid <= 0 {
			return nil, fmt.Errorf("invalid token subject")
		}
		c.UID = uid
	}

	return &c, nil
}

var Module = fx.Provide(NewValidator)
