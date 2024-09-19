package services

import (
	"context"
	"time"

	"github.com/food-order-server/models"
	"github.com/food-order-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodService struct {
	collection *mongo.Collection
}

func CreateFoodService(db *mongo.Database) *FoodService {
	return &FoodService{
		collection: db.Collection("foods"),
	}
}

func (s *FoodService) FindAll() ([]models.Food, error) {
	var results []models.Food
	cursor, err := s.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var foods []bson.M
	if err = cursor.All(context.TODO(), &foods); err != nil {
		return nil, err
	}

	for _, v := range foods {
		id := v["food_id"].(string)
		food, err := s.findFood(id)
		if err != nil {
			return nil, err
		}
		results = append(results, *food)
	}
	return results, nil
}

func (s *FoodService) Create(foodBody *models.FoodReq, id string, file string) (*models.Food, error) {
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
		return nil, err
	}

	return s.findFood(id)
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
