package repository

import (
	"context"
	"time"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDb() (*mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("db.uri")))
	if err != nil {
		return nil, err
	}

	database := client.Database(viper.GetString("db.name"))

	return database, nil
}

func UsersCollection(database *mongo.Database) *mongo.Collection {
	UsersCollection := database.Collection("Users")
	return UsersCollection
}

func RefreshCollection(database *mongo.Database) *mongo.Collection {
	refreshCollection := database.Collection("Refresh_tokens")
	return refreshCollection
}
