package services

import (
	"context"
	"time"

	"github.com/food-order-server/models"
	"github.com/food-order-server/utils"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	collection *mongo.Collection
}

func CreateUserService(db *mongo.Database) *UserService {
	return &UserService{
		collection: db.Collection("users"),
	}
}

func (s *UserService) FindOne(id string, token string) (*models.UserDataRes, error) {
	user := new(models.User)
	if err := s.collection.FindOne(context.TODO(), bson.M{"user_id": id}).Decode(&user); err != nil {
		return nil, err
	}

	return &models.UserDataRes{User: *user, Token: token}, nil
}

func (s *UserService) Login(userBody *models.UserLoginReq) (*models.UserDataRes, error) {
	user, err := s.findUser(userBody.Username)
	if err != nil {
		return nil, err
	}

	if err = utils.VerifyPassword(userBody.Password, user.Password); err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.UserDataRes{User: *user, Token: token}, nil
}

func (s *UserService) Register(userBody *models.UserRegisterReq) (*models.UserDataRes, error) {
	user, err := s.findUser(userBody.Username)
	if user != nil {
		return nil, fiber.ErrConflict
	} else if err != nil && err != fiber.ErrBadRequest {
		return nil, err
	}

	result, err := s.createUser(userBody.Username, userBody.Password, userBody.Name)
	if err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(result)
	if err != nil {
		return nil, err
	}

	return &models.UserDataRes{User: *result, Token: token}, nil
}

func (s *UserService) Update(id string, userBody *models.User) (*models.UserDataRes, error) {
	result, err := s.UpdateUser(id, userBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) Remove(id string) (*models.MessageRes, error) {
	user := new(models.User)
	if err := s.collection.FindOneAndDelete(context.TODO(), bson.M{"user_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &models.MessageRes{Message: "Delete account successfully"}, nil
}

func (s *UserService) UpdateUser(id string, userBody *models.User) (*models.UserDataRes, error) {
	user := new(models.User)
	update := utils.CreateBSON(userBody)
	option := utils.GetUpdateOption()
	if err := s.collection.FindOneAndUpdate(context.TODO(), bson.M{"user_id": id}, update, option).Decode(&user); err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.UserDataRes{User: *user, Token: token}, nil
}

func (s *UserService) createUser(username string, password string, name string) (*models.User, error) {
	if password != "" {
		if hashedPassword, err := utils.HashPassword(password); err != nil {
			return nil, err
		} else {
			password = hashedPassword
		}
	}

	user := &models.User{
		UserID:    utils.GenerateUuid(),
		Username:  username,
		Password:  password,
		Name:      name,
		Role:      "user",
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := s.collection.InsertOne(context.TODO(), user); err != nil {
		return nil, err
	}

	return s.findUser(username)
}

func (s *UserService) findUser(username string) (*models.User, error) {
	user := new(models.User)
	if err := s.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user); err != nil {
		return nil, fiber.ErrBadRequest
	}
	return user, nil
}
