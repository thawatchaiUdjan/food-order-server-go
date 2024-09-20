package utils

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"reflect"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/food-order-server/config"
	"github.com/food-order-server/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GenerateHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashSum := hash.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(hashSum)
}

func CreateBSON(model interface{}) bson.M {
	result := bson.M{}

	v := reflect.ValueOf(model).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("bson")
		tag = strings.Split(tag, ",")[0]

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

func ConvertToObjectIDs(array []string) []primitive.ObjectID {
	var objectIDs []primitive.ObjectID
	for _, idStr := range array {
		objectID, _ := primitive.ObjectIDFromHex(idStr)
		objectIDs = append(objectIDs, objectID)
	}
	return objectIDs
}

func DeleteFile(imageURL string, imageFolder string) error {
	if imageURL == "" {
		return nil
	}

	cld := config.LoadCloudinary()
	parts := strings.Split(imageURL, "/")
	publicId := strings.Split(parts[len(parts)-1], ".")[0]

	if imageFolder != "" {
		publicId = imageFolder + "/" + publicId
	}

	if _, err := cld.Upload.Destroy(context.TODO(), uploader.DestroyParams{
		PublicID: publicId,
	}); err != nil {
		return err
	}

	return nil
}
