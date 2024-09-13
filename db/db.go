package db

import (
	"context"

	"github.com/food-order-server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Database, error) {
	config := config.LoadConfig()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Database.Host))
	if err != nil {
		return nil, err
	}
	return client.Database(config.Database.Name), nil
}
