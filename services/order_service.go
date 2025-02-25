package services

import (
	"context"
	"math"
	"time"

	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/models"
	"github.com/food-order-server/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// @Summary Retrieve a specific food order
// @Description Fetch a food order by user ID.
// @Tags Order
// @Success 200 {object} models.FoodOrderRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /orders [get]
func (s *OrderService) FindOne(c *fiber.Ctx) error {
	req := c.Locals("user").(models.UserReq)
	id := req.User.UserID
	foodOrder, err := s.findFoodOrderByUserId(id)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(foodOrder)
}

// @Summary Retrieve all orders
// @Description Fetch a list of all orders.
// @Tags Order
// @Success 200 {array} models.OrderAll
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /orders/all-order [get]
func (s *OrderService) FindAll(c *fiber.Ctx) error {
	orders, err := s.findOrders()
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return fiber.ErrInternalServerError
	}

	for i, order := range orders {
		foods, err := s.findFoodByOrderId(order.OrderID)
		if err != nil {
			return err
		}
		orders[i].Foods = foods
	}

	return c.JSON(orders)
}

// @Summary Create a new food order
// @Description Place a food order and update the user's balance.
// @Tags Order
// @Param order body models.OrderReq true "Order request body"
// @Success 200 {object} models.OrderDataRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /orders [post]
func (s *OrderService) Create(c *fiber.Ctx) error {
	req := c.Locals("user").(models.UserReq)

	orderReq := new(models.OrderReq)
	if err := c.BodyParser(&orderReq); err != nil {
		return fiber.ErrBadRequest
	}

	if err := middlewares.Validate(orderReq); err != nil {
		return err
	}

	user := req.User
	order := orderReq.Order
	foods := orderReq.Foods

	if user.Balance >= order.TotalPrice {
		status, err := s.orderStatusService.FindDefaultStatus()
		if err != nil {
			return err
		}

		order.OrderID = utils.GenerateUuid()
		order.UserID = user.UserID
		order.OrderStatus = status
		order.UpdatedAt = time.Now()
		order.CreatedAt = time.Now()

		balance := math.Round((user.Balance-order.TotalPrice)*100) / 100
		userData, err := s.userService.UpdateUser(user.UserID, &models.User{Balance: balance})
		if err != nil {
			return err
		}

		if _, err := s.collection.InsertOne(context.TODO(), order); err != nil {
			return err
		}

		for _, food := range foods {
			if err := s.createOrderFood(order.OrderID, food.Food.FoodID, &food); err != nil {
				return err
			}
		}

		foodOrder, err := s.findFoodOrderByUserId(order.UserID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(&models.OrderDataRes{FoodOrder: *foodOrder, User: *userData, Message: "Order item successfully added"})
	} else {
		return fiber.NewError(fiber.StatusNotAcceptable, "Balance is not enough to order")
	}
}

// @Summary Update order status
// @Description Change the status of an existing order by order ID.
// @Tags Order
// @Param id path string true "Order ID"
// @Param status path string true "Order status id"
// @Success 200 {object} models.OrderUpdateRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /orders/{id}/{status} [put]
func (s *OrderService) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	status := c.Params("status")
	statusId, _ := primitive.ObjectIDFromHex(status)
	order := &models.OrderCreate{
		OrderStatus: statusId,
		UpdatedAt:   time.Now(),
	}

	update := utils.CreateBSON(order)
	option := utils.GetUpdateOption()

	if err := s.collection.FindOneAndUpdate(context.TODO(), bson.M{"order_id": id}, update, option).Decode(&order); err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "Order to update not found")
		} else {
			return fiber.ErrInternalServerError
		}
	}

	return c.JSON(&models.OrderUpdateRes{Order: *order, Message: "Order status updated successfully"})
}

// @Summary Remove an order
// @Description Delete an order by order ID.
// @Tags Order
// @Param id path string true "Order ID"
// @Success 200 {object} models.MessageRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /orders/{id} [delete]
func (s *OrderService) Remove(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := s.collection.DeleteOne(context.TODO(), bson.M{"order_id": id}); err != nil {
		return err
	}

	if err := s.removeOrderFood(id); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(models.MessageRes{Message: "Order successfully canceled"})
}

func (s *OrderService) FindOrderFood(id string) error {
	orderFood := new(models.OrderFood)
	if err := s.foodCollection.FindOne(context.TODO(), bson.M{"food_id": id}).Decode(&orderFood); err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.ErrNotFound
		} else {
			return err
		}
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
	orderFood := &models.OrderFoodCreate{
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

func (s *OrderService) removeOrderFood(id string) error {
	if _, err := s.foodCollection.DeleteMany(context.TODO(), bson.M{"order_id": id}); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) findOrders() ([]models.OrderAll, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "user_id"},
			{Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "order_statuses"},
			{Key: "localField", Value: "order_status"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "status"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "delivery_options"},
			{Key: "localField", Value: "order_delivery_option"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "delivery_option"},
		}}},
		bson.D{{Key: "$unwind", Value: "$user"}},
		bson.D{{Key: "$unwind", Value: "$status"}},
		bson.D{{Key: "$unwind", Value: "$delivery_option"}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: bson.D{{Key: "$mergeObjects", Value: bson.A{
			"$$ROOT", bson.D{{Key: "$mergeObjects", Value: bson.A{"$user", "$status"}}},
		}}}}}}},
	}

	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []models.OrderAll
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
