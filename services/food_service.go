package services

import (
	"context"
	"time"

	"github.com/food-order-server/config"
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/models"
	"github.com/food-order-server/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodService struct {
	collection   *mongo.Collection
	orderService *OrderService
}

func CreateFoodService(db *mongo.Database) *FoodService {
	return &FoodService{
		collection:   db.Collection("foods"),
		orderService: CreateOrderService(db),
	}
}

// @Summary Get all foods
// @Description Retrieve a list of all foods from the database
// @Tags Food
// @Success 200 {array} models.Food
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /foods [get]
func (s *FoodService) FindAll(c *fiber.Ctx) error {
	var results []models.Food
	cursor, err := s.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	var foods []bson.M
	if err = cursor.All(context.TODO(), &foods); err != nil {
		return err
	}

	for _, v := range foods {
		id := v["food_id"].(string)
		food, err := s.findFood(id)
		if err != nil {
			return err
		}
		results = append(results, *food)
	}
	return c.JSON(results)
}

// @Summary Create a new food item
// @Description Create a new food item with the provided details.
// @Tags Food
// @Param food body models.FoodReq true "Food request body"
// @Success 200 {object} models.FoodDataRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /foods [post]
func (s *FoodService) Create(c *fiber.Ctx) error {
	foodBody := new(models.FoodReq)
	id := c.Locals("id").(string)
	file := c.Locals("file").(string)

	if err := c.BodyParser(&foodBody); err != nil {
		return fiber.ErrBadRequest
	}

	if err := middlewares.Validate(foodBody); err != nil {
		return err
	}

	if id == "" {
		id = utils.GenerateUuid()
	}

	food := &models.FoodCreate{
		FoodID:            id,
		FoodName:          foodBody.FoodName,
		FoodPrice:         foodBody.FoodPrice,
		FoodPriceDiscount: foodBody.FoodPriceDiscount,
		FoodDescription:   foodBody.FoodDescription,
		FoodImageURL:      file,
		CategoryID:        foodBody.CategoryID,
		FoodOptions:       utils.ConvertToObjectIDs(foodBody.FoodOptions),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if _, err := s.collection.InsertOne(context.TODO(), food); err != nil {
		return err
	}

	result, err := s.findFood(id)
	if err != nil {
		return err
	}

	return c.JSON(models.FoodDataRes{Food: *result, Message: "Food added successfully"})
}

// @Summary Update a food item
// @Description Modify food item details by ID.
// @Tags Food
// @Param food body models.FoodReq true "Food request body"
// @Param id path string true "Food ID"
// @Success 200 {object} models.FoodDataRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /foods/{id} [put]
func (s *FoodService) Update(c *fiber.Ctx) error {
	foodBody := new(models.FoodReq)
	id := c.Params("id")
	user := c.Locals("user").(models.UserReq)
	file := c.Locals("file").(string)

	if err := s.orderService.FindOrderFood(id); err == nil {
		return fiber.NewError(fiber.StatusNotAcceptable, "Food is currently ordered, cannot update")
	} else if err != fiber.ErrNotFound {
		return err
	}

	if err := s.checkFoodPermission(user.User.Role, id); err != nil {
		return fiber.NewError(fiber.StatusNotAcceptable, "No permission for this food. Please try another food")
	}

	if err := c.BodyParser(&foodBody); err != nil {
		return fiber.ErrBadRequest
	}

	food := &models.FoodCreate{
		FoodName:          foodBody.FoodName,
		FoodPrice:         foodBody.FoodPrice,
		FoodPriceDiscount: foodBody.FoodPriceDiscount,
		FoodDescription:   foodBody.FoodDescription,
		CategoryID:        foodBody.CategoryID,
		FoodOptions:       utils.ConvertToObjectIDs(foodBody.FoodOptions),
		UpdatedAt:         time.Now(),
	}

	if file != "" {
		food.FoodImageURL = file
	}

	update := utils.CreateBSON(food)

	if _, err := s.collection.UpdateOne(context.TODO(), bson.M{"food_id": id}, update); err != nil {
		return err
	}

	result, err := s.findFood(id)
	if err != nil {
		return err
	}

	return c.JSON(models.FoodDataRes{Food: *result, Message: "Food item successfully updated"})
}

// @Summary Remove a food item
// @Description Delete a food item by ID.
// @Tags Food
// @Param id path string true "Food ID"
// @Success 200 {object} models.MessageRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /foods/{id} [delete]
func (s *FoodService) Remove(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserReq)
	id := c.Params("id")

	err := s.orderService.FindOrderFood(id)
	if err == nil {
		return fiber.NewError(fiber.StatusNotAcceptable, "Food is currently ordered, cannot delete")
	} else if err != fiber.ErrNotFound {
		return err
	}

	if err := s.checkFoodPermission(user.User.Role, id); err != nil {
		return fiber.NewError(fiber.StatusNotAcceptable, "No permission for this food. Please try another food")
	}

	config := config.LoadConfig()
	food := new(models.FoodCreate)
	if err := s.collection.FindOneAndDelete(context.TODO(), bson.M{"food_id": id}).Decode(&food); err != nil {
		return err
	}

	if err := utils.DeleteFile(food.FoodImageURL, config.UploadFile.FoodFolder); err != nil {
		return err
	}

	return c.JSON(models.MessageRes{Message: "Food item successfully deleted"})
}

func (s *FoodService) findFood(id string) (*models.Food, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "food_id", Value: id}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "food_options"},
			{Key: "localField", Value: "food_options"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "food_options"},
			{Key: "pipeline", Value: bson.A{bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "option_choices"},
					{Key: "localField", Value: "option_choices"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "option_choices"},
				}},
			}}},
		}}},
	}

	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []models.Food
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return &results[0], nil
}

func (s *FoodService) checkFoodPermission(role string, id string) error {
	if food, err := s.findFood(id); err != nil {
		return err
	} else if food.Permission == "owner" && role != "owner" {
		return fiber.ErrUnauthorized
	}
	return nil
}
