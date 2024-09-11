package main

import (
	"log"

	"github.com/food-order-server/config"
	"github.com/food-order-server/db"
	"github.com/food-order-server/routes"
	"github.com/gofiber/fiber/v3"
)

func main() {
	config := config.LoadConfig()
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Fail to connect database: ", err)
	}

	app := fiber.New()

	routes.FoodRoute(app, db)

	log.Fatal(app.Listen(":" + config.Port))
}
