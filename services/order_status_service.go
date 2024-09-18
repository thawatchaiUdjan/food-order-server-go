package services

import (
	"context"

	"github.com/food-order-server/models"
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

func (s *OrderStatusService) FindAll() ([]models.OrderStatus, error) {
	var results []models.OrderStatus
	cursor, err := s.collection.Find(context.TODO(), bson.M{})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *OrderStatusService) FindDefaultStatus() (primitive.ObjectID, error) {
	orderStatus := new(models.OrderStatus)
	if err := s.collection.FindOne(context.TODO(), bson.M{"status_value": 0}).Decode(&orderStatus); err != nil {
		return primitive.NilObjectID, err
	}

	return orderStatus.ID, nil
}
