package services

import (
	"context"

	"github.com/food-order-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodOptionService struct {
	collection *mongo.Collection
}

func CreateFoodOptionService(db *mongo.Database) *FoodOptionService {
	return &FoodOptionService{
		collection: db.Collection("food_options"),
	}
}

func (s *FoodOptionService) FindAll() ([]models.FoodOption, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "option_choices"},
			{Key: "localField", Value: "option_choices"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "option_choices"},
		}}},
	}

	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var foodOptions []models.FoodOption
	if err := cursor.All(context.TODO(), &foodOptions); err != nil {
		return nil, err
	}

	return foodOptions, nil
}
