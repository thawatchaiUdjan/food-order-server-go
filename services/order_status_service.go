package services

import (
	"context"

	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderStatusService struct {
	collection *mongo.Collection
}

func CreateOrderStatusService(db *mongo.Database) *OrderStatusService {
	return &OrderStatusService{
		collection: db.Collection("order_statuses"),
	}
}

// @Summary Retrieve all order statuses
// @Description Fetch a list of all order statuses.
// @Tags Order Status
// @Success 200 {array} models.OrderStatus
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /order-status [get]
func (s *OrderStatusService) FindAll(c *fiber.Ctx) error {
	var results []models.OrderStatus
	cursor, err := s.collection.Find(context.TODO(), bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}

	if err := cursor.All(context.TODO(), &results); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(results)
}

func (s *OrderStatusService) FindDefaultStatus() (primitive.ObjectID, error) {
	orderStatus := new(models.OrderStatus)
	if err := s.collection.FindOne(context.TODO(), bson.M{"status_value": 0}).Decode(&orderStatus); err != nil {
		return primitive.NilObjectID, err
	}

	return orderStatus.ID, nil
}
