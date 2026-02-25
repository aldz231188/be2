package authz

import (
	"sync"
	"time"

	"go.uber.org/fx"
)

type RevokedSessions struct {
	mu    sync.RWMutex
	items map[string]time.Time
}

func NewRevokedSessions() *RevokedSessions {
	return &RevokedSessions{items: make(map[string]time.Time)}
}

func (r *RevokedSessions) Revoke(sessionID string, exp time.Time) {
	if sessionID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[sessionID] = exp
}

func (r *RevokedSessions) IsRevoked(sessionID string) bool {
	if sessionID == "" {
		return false
	}
	now := time.Now()

	r.mu.RLock()
	exp, ok := r.items[sessionID]
	r.mu.RUnlock()
	if !ok {
		return false
	}
	if now.Before(exp) {
		return true
	}

	r.mu.Lock()
	delete(r.items, sessionID)
	r.mu.Unlock()
	return false
}

var RevocationModule = fx.Provide(NewRevokedSessions)
