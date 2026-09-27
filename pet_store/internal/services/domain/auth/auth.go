package auth

import (
	"context"
	"errors"

	services "pet-store/internal/services"
)

type AuthService struct {
	userService services.UserService
	jwtManager  services.JWTManager
}

func NewAuthService(userService services.UserService, jwtManager services.JWTManager) *AuthService {
	return &AuthService{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

func (a *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.userService.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive() {
		return "", errors.New("user account is deactivated")
	}

	if err = a.jwtManager.RevokeALLToken(user.ID); err!=nil{
		return "", errors.New("revoke old token error")
	}

	token, err := a.jwtManager.Generate(user.ID, user.Username)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (a *AuthService) Logout(ctx context.Context, token string) error {
	return a.jwtManager.RevokeToken(token)
}

func (a *AuthService) ValidateToken(ctx context.Context, token string) (int64, error) {
	userID, err := a.jwtManager.Validate(token)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
