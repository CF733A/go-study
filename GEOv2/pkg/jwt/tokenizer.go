package jwt

import (
    "errors"
    "fmt"
    "geo-service/internal/entities"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

type JWTTokenizer struct {
    secretKey string
}

type JWTClaim struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func NewJWTTokenizer(secretKey string) *JWTTokenizer {
    return &JWTTokenizer{
        secretKey: secretKey,
    }
}

func (j *JWTTokenizer) Generate(user *entities.User) (string, error) {
    claim := &JWTClaim{
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   user.Username,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    tokenString, err := token.SignedString([]byte(j.secretKey))
    if err != nil {
        return "", fmt.Errorf("token signing error: %w", err)
    }

    return tokenString, nil
}

func (j *JWTTokenizer) Validate(tokenString string) (*entities.User, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
        return []byte(j.secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
        return &entities.User{
            Username: claims.Username,
        }, nil
    }

    return nil, errors.New("invalid token")
}