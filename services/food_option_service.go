package services

import (
	"context"

	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v2"
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

// @Summary Retrieve all food options
// @Description Fetch a list of all food options with their associated choices.
// @Tags Food Option
// @Success 200 {array} models.FoodOption
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /food-option [get]
func (s *FoodOptionService) FindAll(c *fiber.Ctx) error {
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
		return err
	}

	var foodOptions []models.FoodOption
	if err := cursor.All(context.TODO(), &foodOptions); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(foodOptions)
}
