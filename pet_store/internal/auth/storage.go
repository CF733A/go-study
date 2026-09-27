package auth

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
)

type tokenData struct {
	UserID    int64
	ExpiresAt time.Time
}

type inMemoryTokenStorage struct {
	mu     sync.RWMutex
	Tokens map[string]tokenData
}

func NewInMemoryTokenStorage() *inMemoryTokenStorage {
	return &inMemoryTokenStorage{
		Tokens: make(map[string]tokenData),
	}
}

func (s *inMemoryTokenStorage) StoreToken(token string, userID int64, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Tokens[token] = tokenData{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return nil
}

func (s *inMemoryTokenStorage) GetUserIDByToken(token string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.Tokens[token]
	if !exists {
		return 0, ErrTokenNotFound
	}

	if time.Now().After(data.ExpiresAt) {
		return 0, ErrTokenExpired
	}

	return data.UserID, nil
}

func (s *inMemoryTokenStorage) DeleteToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Tokens, token)
	return nil
}

func (s *inMemoryTokenStorage) RevokeAllUserTokens(userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tokensToDelete := []string{}
	for token, data := range s.Tokens {
		if data.UserID == userID {
			tokensToDelete = append(tokensToDelete, token)
		}
	}

	for _, token := range tokensToDelete {
		delete(s.Tokens, token)
	}

	return nil
}
