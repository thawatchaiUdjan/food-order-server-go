package services

import (
	"context"
	"time"

	"github.com/food-order-server/config"
	"github.com/food-order-server/models"
	"github.com/food-order-server/utils"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/idtoken"
)

type UserService struct {
	collection      *mongo.Collection
	orderCollection *mongo.Collection
}

func CreateUserService(db *mongo.Database) *UserService {
	return &UserService{
		collection:      db.Collection("users"),
		orderCollection: db.Collection("orders"),
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

	result, err := s.createUser(userBody.Username, userBody.Password, userBody.Name, "")
	if err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(result)
	if err != nil {
		return nil, err
	}

	return &models.UserDataRes{User: *result, Token: token}, nil
}

func (s *UserService) GoogleLogin(code string) (*models.UserDataRes, error) {
	googleAuth := config.LoadGoogleAuth()
	authToken, err := googleAuth.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	payload, err := idtoken.Validate(context.TODO(), authToken.Extra("id_token").(string), googleAuth.ClientID)
	if err != nil {
		return nil, err
	}

	username := payload.Subject
	user, err := s.findUser(username)
	if user != nil {
		token, err := utils.CreateToken(user)
		if err != nil {
			return nil, err
		}

		return &models.UserDataRes{User: *user, Token: token}, nil
	} else if err != nil && err != fiber.ErrBadRequest {
		return nil, err
	}

	name := payload.Claims["name"].(string)
	image := payload.Claims["picture"].(string)
	user, err = s.createUser(username, "", name, image)
	if err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.UserDataRes{User: *user, Token: token}, nil
}

func (s *UserService) Update(id string, userBody *models.User, file string) (*models.UserUpdateRes, error) {
	if file != "" {
		userBody.ProfileImageURL = file
	}

	result, err := s.UpdateUser(id, userBody)
	if err != nil {
		return nil, err
	}
	return &models.UserUpdateRes{User: *result, Message: "User data successfully updated"}, nil
}

func (s *UserService) Remove(id string) (*models.MessageRes, error) {
	order := new(models.OrderCreate)
	if err := s.orderCollection.FindOne(context.TODO(), bson.M{"user_id": id}).Decode(&order); err == nil {
		return nil, fiber.ErrNotAcceptable
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	config := config.LoadConfig()
	user := new(models.User)
	if err := s.collection.FindOneAndDelete(context.TODO(), bson.M{"user_id": id}).Decode(&user); err != nil {
		return nil, err
	}

	if err := utils.DeleteFile(user.ProfileImageURL, config.UploadFile.ProfileFolder); err != nil {
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

func (s *UserService) createUser(username string, password string, name string, image string) (*models.User, error) {
	if password != "" {
		if hashedPassword, err := utils.HashPassword(password); err != nil {
			return nil, err
		} else {
			password = hashedPassword
		}
	}

	user := &models.User{
		UserID:          utils.GenerateUuid(),
		Username:        username,
		Password:        password,
		Name:            name,
		Role:            "user",
		Balance:         0,
		ProfileImageURL: image,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
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
