package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/KazakMaxim/auth-jwt-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthMongo struct {
	db *mongo.Database
}

func NewAuthMongo(db *mongo.Database) *AuthMongo {
	return &AuthMongo{db: db}
}

func (m *AuthMongo) SaveUser(user models.User) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := UsersCollection(m.db).InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (m *AuthMongo) GetUser(username, password string) (string, error) {
	var user bson.M
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	getErr := UsersCollection(m.db).FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if getErr != nil {
		return "", getErr
	}

	hashPassword, ok := user["password"]
	hashPasswordStr := fmt.Sprintf("%v", hashPassword)
	if !ok {
		return "", getErr
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(hashPasswordStr), []byte(password))
	if compareErr != nil {
		return "", compareErr
	}

	userGuid := fmt.Sprintf("%v", user["user_guid"])

	return userGuid, nil
}

func (m *AuthMongo) FindUserByUsername(username string) bool {
	var user bson.M
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	getErr := UsersCollection(m.db).FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if getErr != nil {
		return false
	}

	return true
}

func (m *AuthMongo) FindUserByGuid(userGuid string) bool {
	var user bson.M
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	getErr := UsersCollection(m.db).FindOne(ctx, bson.M{"user_guid": userGuid}).Decode(&user)
	if getErr != nil {
		return false
	}

	return true
}
