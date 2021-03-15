package repository

import (
	"github.com/KazakMaxim/auth-jwt-go/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	SaveUser(user models.User) error
	GetUser(username, password string) (string, error)
	FindUserByUsername(username string) bool
	FindUserByGuid(userGuid string) bool
}

type Tokens interface {
	SaveRefresh(refresh models.RefreshTable) error
	FindRefersh(refreshToken, userGuid string) error
	DeleteRefersh(id string) error
}

type Repository struct {
	Authorization
	Tokens
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(db),
		Tokens:        NewTokensMongo(db),
	}
}
