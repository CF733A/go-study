package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type JWTManager struct {
	SecretKey    string
	TokenStorage *inMemoryTokenStorage
}

type JWTClaim struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		SecretKey:    secretKey,
		TokenStorage: NewInMemoryTokenStorage(),
	}
}

func (j *JWTManager) Generate(userID int64, username string) (string, error) {

	claim := &JWTClaim{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", fmt.Errorf("token signing error: %w", err)
	}

	expiresAt := time.Now().Add(time.Hour * 24)
	err = j.TokenStorage.StoreToken(tokenString, userID, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}

	return tokenString, nil
}

func (j *JWTManager) Validate(tokenString string) (int64, error) {
	userID, err := j.TokenStorage.GetUserIDByToken(tokenString)
	if err != nil {
		return 0, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		if claims.UserID != userID {
			return 0, ErrInvalidToken
		}
		return claims.UserID, nil
	}

	return 0, ErrInvalidToken
}

func (j *JWTManager) RevokeToken(token string) error {
	return j.TokenStorage.DeleteToken(token)
}

func (j *JWTManager) RevokeALLToken(userID int64) error {
	return j.TokenStorage.RevokeAllUserTokens(userID)
}
