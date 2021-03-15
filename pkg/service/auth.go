package service

import (
	"fmt"

	"github.com/KazakMaxim/auth-jwt-go/models"
	"github.com/KazakMaxim/auth-jwt-go/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (s *AuthService) CreateUser(user models.User) error {
	FindUserByUsername := s.repos.FindUserByUsername(user.Username)
	if FindUserByUsername {
		return fmt.Errorf("Username is already in use")
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bcryptPassword)

	return s.repos.SaveUser(user)
}

func (s *AuthService) AuthUser(username, password string) (string, error) {
	userGuid, getErr := s.repos.GetUser(username, password)
	if getErr != nil {
		return "", getErr
	}

	return userGuid, nil
}
