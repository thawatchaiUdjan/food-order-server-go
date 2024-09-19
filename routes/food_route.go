package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/models"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodRoute(app *fiber.App, db *mongo.Database) {
	foodService := services.CreateFoodService(db)
	route := app.Group("/foods", middlewares.AuthToken)

	route.Get("/", func(c fiber.Ctx) error {
		foods, err := foodService.FindAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(foods)
	})

	route.Post("/", func(c fiber.Ctx) error {
		foodBody := new(models.FoodReq)
		id := c.Locals("id").(string)
		file := c.Locals("file").(string)

		if err := c.Bind().Body(foodBody); err != nil {
			return fiber.ErrBadRequest
		}

		food, err := foodService.Create(foodBody, id, file)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(models.FoodDataRes{Food: *food, Message: "Food added successfully"})
	}, middlewares.UploadFoodFile)

	route.Put("/:id", func(c fiber.Ctx) error {
		foodBody := new(models.FoodReq)
		id := c.Params("id")
		file := c.Locals("file").(string)

		if err := c.Bind().Body(foodBody); err != nil {
			return fiber.ErrBadRequest
		}

		food, err := foodService.Update(id, foodBody, file)
		if err == fiber.ErrNotAcceptable {
			return err
		} else if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(models.FoodDataRes{Food: *food, Message: "Food item successfully updated"})
	}, middlewares.UploadFoodFile)

	route.Delete("/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		if err := foodService.Remove(id); err != nil {
			if err == fiber.ErrNotAcceptable {
				return err
			} else {
				return fiber.ErrInternalServerError
			}
		}
		return c.JSON(models.MessageRes{Message: "Food item successfully deleted"})
	})
}
