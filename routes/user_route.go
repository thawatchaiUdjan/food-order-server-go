package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoute(app *fiber.App, db *mongo.Database) {
	userService := services.CreateUserService(db)
	route := app.Group("/user")

	route.Get("/verify-token", middlewares.AuthToken, userService.VerifyToken)
	route.Get("/", middlewares.AuthToken, userService.FindOne)
	route.Post("/login", userService.Login)
	route.Post("/register", userService.Register)
	route.Post("/google-login", userService.GoogleLogin)
	route.Post("/facebook-login", userService.FacebookLogin)
	route.Put("/", middlewares.AuthToken, middlewares.UploadProfileFile, userService.Update)
	route.Delete("/", middlewares.AuthToken, userService.Remove)
}
