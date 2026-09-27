package auth

import (
	"errors"
	"geo-service/internal/entities"
	"geo-service/internal/service"
	"geo-service/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepo storage.UserRepository
	tokenGen service.TokenGenerator
}

func NewAuthService(userRepo storage.UserRepository, tokenGen service.TokenGenerator) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo: userRepo,
		tokenGen: tokenGen,
	}
}

func (a *AuthServiceImpl) Register(username, password string) error {
	existingUser, _ := a.userRepo.FindByUsername(username)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}


	user := &entities.User{
		Username:     username,
		PasswordHash: string(hashedPassBytes),
	}

	return a.userRepo.Create(user)
}

func (a *AuthServiceImpl) Login(username, password string) (string, error) {

	user, err := a.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}


	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}


	return a.tokenGen.Generate(user)
}

func (a *AuthServiceImpl) ValidateToken(token string) (*entities.User, error) {
	return a.tokenGen.Validate(token)
}
