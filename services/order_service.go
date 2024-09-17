package services

import (
	"context"

	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderService struct {
	collection *mongo.Collection
}

func CreateOrderService(db *mongo.Database) *OrderService {
	return &OrderService{
		collection: db.Collection("orders"),
	}
}

func (s *OrderService) FindAll() ([]models.Order, error) {
	// cursor, err := s.collection.Find(context.TODO(), bson.M{})
	// if err != nil {
	// 	return nil, err
	// }

	// var foodCategories []models.FoodCategory
	// if err := cursor.All(context.TODO(), &foodCategories); err != nil {
	// 	return nil, err
	// }

	// return foodCategories, nil
	return nil, fiber.ErrInternalServerError
}

func (s *OrderService) FindOne(id string) (*models.Order, error) {
	order := new(models.Order)
	if err := s.collection.FindOne(context.TODO(), bson.M{"user_id": id}).Decode(&order); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return order, nil
}
