package services

import (
	"context"
	"math"
	"time"

	"github.com/food-order-server/models"
	"github.com/food-order-server/utils"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderService struct {
	collection         *mongo.Collection
	foodCollection     *mongo.Collection
	orderStatusService *OrderStatusService
	userService        *UserService
}

func CreateOrderService(db *mongo.Database) *OrderService {
	return &OrderService{
		collection:         db.Collection("orders"),
		foodCollection:     db.Collection("order_foods"),
		orderStatusService: CreateOrderStatusService(db),
		userService:        CreateUserService(db),
	}
}

func (s *OrderService) Create(user *models.User, orderReq *models.OrderReq) (*models.OrderDataRes, error) {
	order := orderReq.Order
	foods := orderReq.Foods

	if user.Balance >= order.TotalPrice {
		status, err := s.orderStatusService.FindDefaultStatus()
		if err != nil {
			return nil, err
		}

		order.OrderID = utils.GenerateUuid()
		order.UserID = user.UserID
		order.OrderStatus = status
		order.UpdatedAt = time.Now()
		order.CreatedAt = time.Now()

		balance := math.Round((user.Balance-order.TotalPrice)*100) / 100
		userData, err := s.userService.UpdateUser(user.UserID, &models.User{Balance: balance})
		if err != nil {
			return nil, err
		}

		if _, err := s.collection.InsertOne(context.TODO(), order); err != nil {
			return nil, err
		}

		for _, food := range foods {
			if err := s.createOrderFood(order.OrderID, food.Food.FoodID, &food); err != nil {
				return nil, err
			}
		}

		foodOrder, err := s.findFoodOrderByUserId(order.UserID)
		if err != nil {
			return nil, err
		}

		return &models.OrderDataRes{FoodOrder: *foodOrder, User: *userData, Message: "Order item successfully added"}, nil
	} else {
		return nil, fiber.ErrNotAcceptable
	}
}

func (s *OrderService) FindAll() ([]models.Order, error) {
	return nil, fiber.ErrInternalServerError
}

func (s *OrderService) FindOne(id string) (*models.FoodOrderRes, error) {
	foodOrder, err := s.findFoodOrderByUserId(id)
	if err != nil {
		return nil, err
	}

	return foodOrder, nil
}

func (s *OrderService) Remove(id string) error {
	if _, err := s.collection.DeleteOne(context.TODO(), bson.M{"order_id": id}); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) findOrderByUserId(id string) (*models.Order, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: id}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "order_statuses"},
			{Key: "localField", Value: "order_status"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "order_status"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "delivery_options"},
			{Key: "localField", Value: "order_delivery_option"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "order_delivery_option"},
		}}},
		bson.D{{Key: "$unwind", Value: "$order_status"}},
		bson.D{{Key: "$unwind", Value: "$order_delivery_option"}},
	}

	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []models.Order
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, fiber.NewError(fiber.StatusNotAcceptable, "Error on cursor: "+err.Error())
	}

	if len(results) <= 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &results[0], nil
}

func (s *OrderService) findFoodByOrderId(id string) ([]models.OrderFood, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "foods"},
			{Key: "localField", Value: "food_id"},
			{Key: "foreignField", Value: "food_id"},
			{Key: "as", Value: "food"},
		}}},
		bson.D{{Key: "$unwind", Value: "$food"}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: bson.D{{Key: "$mergeObjects", Value: bson.A{"$$ROOT", "$food"}}}}}}},
	}

	cursor, err := s.foodCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []models.OrderFood
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *OrderService) findFoodOrderByUserId(id string) (*models.FoodOrderRes, error) {
	order, err := s.findOrderByUserId(id)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	foods, err := s.findFoodByOrderId(order.OrderID)
	if err != nil {
		return nil, err
	}

	return &models.FoodOrderRes{Order: *order, Foods: foods}, nil
}

func (s *OrderService) createOrderFood(orderId string, foodId string, food *models.FoodCart) error {
	orderFood := &models.OrderFood{
		OrderID:          orderId,
		FoodID:           foodId,
		FoodAmount:       food.Amount,
		FoodTotalPrice:   food.Total,
		FoodOptionString: food.Option.OptionString,
		FoodOptionNote:   food.Option.OptionNote,
		UpdatedAt:        time.Now(),
		CreatedAt:        time.Now(),
	}
	if _, err := s.foodCollection.InsertOne(context.TODO(), orderFood); err != nil {
		return err
	}
	return nil
}
