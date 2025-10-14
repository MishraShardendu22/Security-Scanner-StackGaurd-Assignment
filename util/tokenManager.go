package util

import (
	"strings"
	"sync"
)

// TokenManager manages multiple API tokens with round-robin rotation
type TokenManager struct {
	tokens       []string
	currentIndex int
	mu           sync.Mutex
}

var (
	globalTokenManager *TokenManager
	once               sync.Once
)

// InitTokenManager initializes the global token manager with comma-separated tokens
func InitTokenManager(tokenString string) {
	once.Do(func() {
		tokens := strings.Split(tokenString, ",")
		// Clean up tokens (trim whitespace)
		cleanedTokens := make([]string, 0, len(tokens))
		for _, token := range tokens {
			trimmed := strings.TrimSpace(token)
			if trimmed != "" {
				cleanedTokens = append(cleanedTokens, trimmed)
			}
		}

		globalTokenManager = &TokenManager{
			tokens:       cleanedTokens,
			currentIndex: 0,
		}
	})
}

// GetTokenManager returns the global token manager instance
func GetTokenManager() *TokenManager {
	if globalTokenManager == nil {
		panic("TokenManager not initialized. Call InitTokenManager first.")
	}
	return globalTokenManager
}

// GetCurrentToken returns the current token without rotating
func (tm *TokenManager) GetCurrentToken() string {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if len(tm.tokens) == 0 {
		return ""
	}
	return tm.tokens[tm.currentIndex]
}

// RotateToken moves to the next token using (n+1) % len logic
func (tm *TokenManager) RotateToken() string {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if len(tm.tokens) == 0 {
		return ""
	}

	// Move to next token: (currentIndex + 1) % length
	tm.currentIndex = (tm.currentIndex + 1) % len(tm.tokens)
	return tm.tokens[tm.currentIndex]
}

// GetTokenCount returns the total number of tokens available
func (tm *TokenManager) GetTokenCount() int {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return len(tm.tokens)
}

// GetCurrentIndex returns the current token index
func (tm *TokenManager) GetCurrentIndex() int {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.currentIndex
}
