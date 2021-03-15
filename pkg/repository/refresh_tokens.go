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

type TokensMongo struct {
	db *mongo.Database
}

func NewTokensMongo(db *mongo.Database) *TokensMongo {
	return &TokensMongo{db: db}
}

func (m *TokensMongo) SaveRefresh(refresh models.RefreshTable) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := RefreshCollection(m.db).InsertOne(ctx, refresh)
	if err != nil {
		return err
	}

	return nil
}

func (m *TokensMongo) FindRefersh(refreshToken, userGuid string) error {
	var refresh bson.M
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	getErr := RefreshCollection(m.db).FindOne(ctx, bson.M{"user_guid": userGuid}).Decode(&refresh)
	if getErr != nil {
		return getErr
	}

	hashToken, ok := refresh["refresh_token"]
	hashTokenStr := fmt.Sprintf("%v", hashToken)
	if !ok {
		return fmt.Errorf("Token is not found")
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(hashTokenStr), []byte(refreshToken))
	if compareErr != nil {
		return compareErr
	}

	return nil
}

func (m *TokensMongo) DeleteRefersh(userGuid string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := RefreshCollection(m.db).DeleteMany(ctx, bson.M{"user_guid": userGuid})
	if err != nil {
		return err
	}

	return nil
}
