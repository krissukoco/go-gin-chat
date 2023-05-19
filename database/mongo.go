package database

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo() (*mongo.Database, error) {
	mongoUri, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		return nil, errors.New("MONGO_URI is not set")
	}
	dbName, exists := os.LookupEnv("MONGO_DBNAME")
	if !exists {
		return nil, errors.New("MONGO_DBNAME is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)

	return db, nil
}
