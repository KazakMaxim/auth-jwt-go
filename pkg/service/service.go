package service

import (
	"github.com/KazakMaxim/auth-jwt-go/models"
	"github.com/KazakMaxim/auth-jwt-go/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) error
	AuthUser(username, password string) (string, error)
}

type Tokens interface {
	GetTokensByGuid(userGuid string) ([]string, error)
	GenerateTokens(userGuid string) ([]string, error)
	NewTokens(tokens models.Tokens) ([]string, error)
}

type Service struct {
	Authorization
	Tokens
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
		Tokens:        NewTokensService(repos),
	}
}
