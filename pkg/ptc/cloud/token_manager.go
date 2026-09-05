package cloud

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

type TokenManager struct {
	mu     sync.RWMutex
	tokens map[string]*oauth2.Token // Key: "tenantID/username"
}

func NewTokenManager() *TokenManager {
	return &TokenManager{
		tokens: make(map[string]*oauth2.Token),
	}
}

// GetToken returns a valid token for the given client, handling refresh and login transparently
func (tm *TokenManager) GetToken(ctx context.Context, c *Client) (*oauth2.Token, error) {
	if c == nil || c.creds == nil {
		return nil, fmt.Errorf("cannot get token for nil client or credentials")
	}

	key := fmt.Sprintf("%s/%s", c.creds.TenantID, c.creds.Username)
	buffer := 30 * time.Second

	// 1. Read Lock: Fast check for valid token in cache
	tm.mu.RLock()
	tok, exists := tm.tokens[key]
	if exists && isTokenValid(tok, buffer) {
		tm.mu.RUnlock()
		return tok, nil
	}
	tm.mu.RUnlock()

	// 2. Write Lock: Exclusive access for auth network calls
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Double-check after acquiring write lock
	tok, exists = tm.tokens[key]
	if exists && isTokenValid(tok, buffer) {
		return tok, nil
	}

	// Step A: Attempt refresh if token exists and has a refresh string
	if exists && tok.RefreshToken != "" {
		refreshedTok, err := c.refreshToken(ctx, tok.RefreshToken)
		if err == nil && isTokenValid(refreshedTok, 0) {
			tm.tokens[key] = refreshedTok
			return refreshedTok, nil
		}
	}

	// Step B: Fall back to fresh login method
	newTok, err := c.login(ctx)
	if err != nil {
		delete(tm.tokens, key) // Clear invalid state
		return nil, fmt.Errorf("auth failed for tenant %s user %s: %w", c.creds.TenantID, c.creds.Username, err)
	}

	tm.tokens[key] = newTok
	return newTok, nil
}

func isTokenValid(tok *oauth2.Token, buffer time.Duration) bool {
	if tok == nil || !tok.Valid() {
		return false
	}
	if !tok.Expiry.IsZero() && time.Now().Add(buffer).After(tok.Expiry) {
		return false
	}
	return true
}
