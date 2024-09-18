package services

import (
	"context"

	"github.com/food-order-server/models"
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

func (s *DeliveryOptionService) FindAll() ([]models.DeliveryOption, error) {
	sort := options.Find().SetSort(bson.D{{Key: "delivery_cost", Value: 1}})
	cursor, err := s.collection.Find(context.TODO(), bson.M{}, sort)
	if err != nil {
		return nil, err
	}

	var deliveryOptions []models.DeliveryOption
	if err := cursor.All(context.TODO(), &deliveryOptions); err != nil {
		return nil, err
	}

	return deliveryOptions, nil
}
