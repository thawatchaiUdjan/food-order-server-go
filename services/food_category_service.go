package services

import (
	"context"

	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v2"
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

// @Summary Retrieve all food categories
// @Description Fetch a list of all food categories.
// @Tags Food Category
// @Success 200 {array} models.FoodCategory
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /category [get]
func (s *FoodCategoryService) FindAll(c *fiber.Ctx) error {
	cursor, err := s.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	var foodCategories []models.FoodCategory
	if err := cursor.All(context.TODO(), &foodCategories); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(foodCategories)
}
