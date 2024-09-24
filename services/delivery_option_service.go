package services

import (
	"context"

	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeliveryOptionService struct {
	collection *mongo.Collection
}

func CreateDeliveryOptionService(db *mongo.Database) *DeliveryOptionService {
	return &DeliveryOptionService{
		collection: db.Collection("delivery_options"),
	}
}

// @Summary Retrieve all delivery options
// @Description Fetch a list of all delivery options sorted by delivery cost.
// @Tags Delivery Option
// @Success 200 {array} models.DeliveryOption
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /delivery [get]
func (s *DeliveryOptionService) FindAll(c *fiber.Ctx) error {
	sort := options.Find().SetSort(bson.D{{Key: "delivery_cost", Value: 1}})
	cursor, err := s.collection.Find(context.TODO(), bson.M{}, sort)
	if err != nil {
		return err
	}

	var deliveryOptions []models.DeliveryOption
	if err := cursor.All(context.TODO(), &deliveryOptions); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(deliveryOptions)
}
