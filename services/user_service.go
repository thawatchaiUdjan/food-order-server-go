package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/food-order-server/config"
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/models"
	"github.com/food-order-server/utils"
	"github.com/gofiber/fiber/v2"
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

// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags User
// @Param user body models.UserLoginReq true "User login information"
// @Success 200 {object} models.UserDataRes
// @Failure 500 {object} models.MessageRes
// @Router /user/login [post]
func (s *UserService) Login(c *fiber.Ctx) error {
	userBody := new(models.UserLoginReq)

	if err := c.BodyParser(userBody); err != nil {
		return fiber.ErrBadRequest
	}

	if err := middlewares.Validate(userBody); err != nil {
		return err
	}

	user, err := s.findUser(userBody.Username)
	if err == fiber.ErrBadRequest {
		return fiber.NewError(fiber.StatusUnauthorized, "Username or password invalid")
	} else if err != nil {
		return err
	}

	if err = utils.VerifyPassword(userBody.Password, user.Password); err == fiber.ErrBadRequest {
		return fiber.NewError(fiber.StatusUnauthorized, "Username or password invalid")
	} else if err != nil {
		return err
	}

	token, err := utils.CreateToken(user)
	if err != nil {
		return err
	}

	return c.JSON(&models.UserDataRes{User: *user, Token: token})
}

func (s *UserService) Register(userBody *models.UserRegisterReq) (*models.UserDataRes, error) {
	return s.findOrCreateUser(userBody.Username, userBody.Password, userBody.Name, "", false)
}

func (s *UserService) GoogleLogin(code string) (*models.UserDataRes, error) {
	google := config.LoadGoogle()
	authToken, err := google.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	payload, err := idtoken.Validate(context.TODO(), authToken.Extra("id_token").(string), google.ClientID)
	if err != nil {
		return nil, err
	}

	username := payload.Subject
	name := payload.Claims["name"].(string)
	image := payload.Claims["picture"].(string)

	return s.findOrCreateUser(username, "", name, image, true)
}

func (s *UserService) FacebookLogin(accessToken string) (*models.UserDataRes, error) {
	url := "https://graph.facebook.com/me?fields=id,name,email,first_name,last_name,picture&access_token=" + accessToken
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, http.ErrNotSupported
	}

	user := new(models.UserFacebook)
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, err
	}

	return s.findOrCreateUser(user.ID, "", user.Name, user.ProfilePicture.URL, true)
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

func (s *UserService) findOrCreateUser(username string, password string, name string, image string, isSocial bool) (*models.UserDataRes, error) {
	user, err := s.findUser(username)
	if user != nil {
		if isSocial {
			token, err := utils.CreateToken(user)
			if err != nil {
				return nil, err
			}
			return &models.UserDataRes{User: *user, Token: token}, nil
		} else {
			return nil, fiber.ErrConflict
		}
	} else if err != nil && err != fiber.ErrBadRequest {
		return nil, err
	}

	result, err := s.createUser(username, password, name, image)
	if err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(result)
	if err != nil {
		return nil, err
	}

	return &models.UserDataRes{User: *result, Token: token}, nil
}
