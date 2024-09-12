package services

import (
	"context"

	"github.com/food-order-server/config"
	"github.com/food-order-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodService struct {
	collection *mongo.Collection
}

func CreateFoodService(client *mongo.Client) *FoodService {
	config := config.LoadConfig()
	return &FoodService{
		collection: client.Database(config.Database.Name).Collection("foods"),
	}
}

func (s *FoodService) GetFoods() ([]models.Food, error) {
	cursor, err := s.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var foods []models.Food
	if err = cursor.All(context.TODO(), &foods); err != nil {
		return nil, err
	}
	return foods, nil
}
