package main

import (
	"log"

	"github.com/food-order-server/config"
	"github.com/food-order-server/db"
	_ "github.com/food-order-server/docs"
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @title Food Order App API
// @version 1.0.0
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config := config.LoadConfig()
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Fail to connect database: ", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
	})
	app.Use(cors.New())
	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.FoodRoute(app, db)
	routes.UserRoute(app, db)
	routes.FoodCategoryRoute(app, db)
	routes.OrderRoute(app, db)
	routes.OrderStatusRoute(app, db)
	routes.DeliveryOptionRoute(app, db)
	routes.FoodOptionRoute(app, db)

	log.Fatal(app.Listen(":" + config.Port))
}
