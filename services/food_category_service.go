package services

import (
	"context"

	"github.com/food-order-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodCategoryService struct {
	collection *mongo.Collection
}

func CreateFoodCategoryService(db *mongo.Database) *FoodCategoryService {
	return &FoodCategoryService{
		collection: db.Collection("food_categories"),
	}
}

func (s *FoodCategoryService) FindAll() ([]models.FoodCategory, error) {
	cursor, err := s.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var foodCategories []models.FoodCategory
	if err := cursor.All(context.TODO(), &foodCategories); err != nil {
		return nil, err
	}

	return foodCategories, nil
}
