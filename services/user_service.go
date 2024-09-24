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

// @Summary Verify user token
// @Description Check the validity of the user token.
// @Tags User
// @Success 200 {string} string
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /user/verify-token [get]
func (s *UserService) VerifyToken(c *fiber.Ctx) error {
	return c.SendString("verify complete")
}

// @Summary Retrieve a user by ID
// @Description Fetch user details based on user ID.
// @Tags User
// @Success 200 {object} models.UserDataRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /user [get]
func (s *UserService) FindOne(c *fiber.Ctx) error {
	req := c.Locals("user").(models.UserReq)
	user := new(models.User)
	if err := s.collection.FindOne(context.TODO(), bson.M{"user_id": req.User.UserID}).Decode(&user); err != nil {
		return err
	}

	return c.JSON(&models.UserDataRes{User: *user, Token: req.Token})
}

// @Summary User login
// @Description Authenticate a user and return user, token
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

// @Summary Register a new user
// @Description Create a new user account with the provided details.
// @Tags User
// @Param user body models.UserRegisterReq true "User register information"
// @Success 200 {object} models.UserDataRes
// @Failure 500 {object} models.MessageRes
// @Router /user/register [post]
func (s *UserService) Register(c *fiber.Ctx) error {
	userBody := new(models.UserRegisterReq)

	if err := c.BodyParser(&userBody); err != nil {
		return fiber.ErrBadRequest
	}

	if err := middlewares.Validate(userBody); err != nil {
		return err
	}

	user, err := s.findOrCreateUser(userBody.Username, userBody.Password, userBody.Name, "", false)
	if err == fiber.ErrConflict {
		return fiber.NewError(fiber.StatusConflict, "Username is already in use")
	} else if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(user)
}

func (s *UserService) GoogleLogin(c *fiber.Ctx) error {
	googleLoginBody := new(models.UserGoogleLoginReq)

	if err := c.BodyParser(&googleLoginBody); err != nil {
		return fiber.ErrBadRequest
	}

	if err := middlewares.Validate(googleLoginBody); err != nil {
		return err
	}

	google := config.LoadGoogle()
	code := googleLoginBody.Code
	authToken, err := google.Exchange(context.TODO(), code)
	if err != nil {
		return err
	}

	payload, err := idtoken.Validate(context.TODO(), authToken.Extra("id_token").(string), google.ClientID)
	if err != nil {
		return err
	}

	username := payload.Subject
	name := payload.Claims["name"].(string)
	image := payload.Claims["picture"].(string)
	user, err := s.findOrCreateUser(username, "", name, image, true)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(user)
}

func (s *UserService) FacebookLogin(c *fiber.Ctx) error {
	facebookLoginBody := new(models.UserFacebookLoginReq)

	if err := c.BodyParser(&facebookLoginBody); err != nil {
		return fiber.ErrBadRequest
	}

	if err := middlewares.Validate(facebookLoginBody); err != nil {
		return err
	}

	accessToken := facebookLoginBody.AccessToken
	url := "https://graph.facebook.com/me?fields=id,name,email,first_name,last_name,picture&access_token=" + accessToken
	res, err := http.Get(url)
	if err != nil {
		return err
	} else if res.StatusCode != http.StatusOK {
		return http.ErrNotSupported
	}

	payload := new(models.UserFacebook)
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return err
	}

	user, err := s.findOrCreateUser(payload.ID, "", payload.Name, payload.ProfilePicture.URL, true)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(user)
}

// @Summary Update user information
// @Description Modify existing user details by user ID.
// @Tags User
// @Param user body models.User true "User update information"
// @Success 200 {object} models.UserUpdateRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /user [put]
func (s *UserService) Update(c *fiber.Ctx) error {
	userBody := new(models.User)
	req := c.Locals("user").(models.UserReq)
	file := c.Locals("file").(string)
	id := req.User.UserID

	if err := c.BodyParser(&userBody); err != nil {
		return fiber.ErrBadRequest
	}

	if err := middlewares.Validate(userBody); err != nil {
		return err
	}

	if file != "" {
		userBody.ProfileImageURL = file
	}

	result, err := s.UpdateUser(id, userBody)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(&models.UserUpdateRes{User: *result, Message: "User data successfully updated"})
}

// @Summary Remove a user account
// @Description Delete a user account by user ID.
// @Tags User
// @Success 200 {object} models.MessageRes
// @Failure 500 {object} models.MessageRes
// @Security BearerAuth
// @Router /user [delete]
func (s *UserService) Remove(c *fiber.Ctx) error {
	req := c.Locals("user").(models.UserReq)
	id := req.User.UserID

	order := new(models.OrderCreate)
	if err := s.orderCollection.FindOne(context.TODO(), bson.M{"user_id": id}).Decode(&order); err == nil {
		return fiber.NewError(fiber.StatusNotAcceptable, "account have an order, cant be delete")
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	config := config.LoadConfig()
	user := new(models.User)
	if err := s.collection.FindOneAndDelete(context.TODO(), bson.M{"user_id": id}).Decode(&user); err != nil {
		return fiber.ErrInternalServerError
	}

	if err := utils.DeleteFile(user.ProfileImageURL, config.UploadFile.ProfileFolder); err != nil {
		return err
	}

	return c.JSON(&models.MessageRes{Message: "Delete account successfully"})
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
