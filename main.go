package main

import (
	"log"

	"github.com/food-order-server/config"
	"github.com/food-order-server/db"
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	config := config.LoadConfig()
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Fail to connect database: ", err)
	}

	app := fiber.New(fiber.Config{
		StructValidator: middlewares.Validator(),
		ErrorHandler:    middlewares.ErrorHandler,
	})
	app.Use(cors.New())

	routes.FoodRoute(app, db)
	routes.UserRoute(app, db)
	routes.FoodCategoryRoute(app, db)

	log.Fatal(app.Listen(":" + config.Port))
}
