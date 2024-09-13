package utils

import (
	"reflect"
	"time"

	"github.com/food-order-server/config"
	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func CreateToken(user *models.User) (string, error) {
	config := config.LoadConfig()
	duration, err := time.ParseDuration(config.Token.TokenExpiredTime)
	if err != nil {
		return "", err
	}

	expiredTime := time.Now().Add(duration)
	claims := models.Claims{
		User: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyPassword(password string, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fiber.ErrBadRequest
	} else {
		return nil
	}
}

func HashPassword(password string) (string, error) {
	config := config.LoadConfig()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), config.EncryptSaltRounds)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func GenerateUuid() string {
	return uuid.NewString()
}

func CreateBSON(model interface{}) bson.M {
	result := bson.M{}

	v := reflect.ValueOf(model).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("bson")

		if tag == "" || reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue
		}

		result[tag] = field.Interface()
	}

	return bson.M{"$set": result}
}

func GetUpdateOption() *options.FindOneAndUpdateOptions {
	return options.FindOneAndUpdate().SetReturnDocument(options.After)
}
