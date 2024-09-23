package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/models"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func FoodRoute(app *fiber.App, db *mongo.Database) {
	foodService := services.CreateFoodService(db)
	route := app.Group("/foods", middlewares.AuthToken)

	route.Get("/", foodService.FindAll)

	route.Post("/", middlewares.UploadFoodFile, func(c *fiber.Ctx) error {
		foodBody := new(models.FoodReq)
		id := c.Locals("id").(string)
		file := c.Locals("file").(string)

		if err := c.BodyParser(&foodBody); err != nil {
			return fiber.ErrBadRequest
		}

		if err := middlewares.Validate(foodBody); err != nil {
			return err
		}

		food, err := foodService.Create(foodBody, id, file)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(models.FoodDataRes{Food: *food, Message: "Food added successfully"})
	})

	route.Put("/:id", middlewares.UploadFoodFile, func(c *fiber.Ctx) error {
		foodBody := new(models.FoodReq)
		user := c.Locals("user").(models.UserReq)
		id := c.Params("id")
		file := c.Locals("file").(string)

		if err := c.BodyParser(&foodBody); err != nil {
			return fiber.ErrBadRequest
		}

		food, err := foodService.Update(user, id, foodBody, file)
		if err == fiber.ErrNotAcceptable {
			return fiber.NewError(fiber.StatusNotAcceptable, "Food is currently ordered, cannot update")
		} else if err == fiber.ErrUnauthorized {
			return fiber.NewError(fiber.StatusNotAcceptable, "No permission for this food. Please try another food")
		} else if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(models.FoodDataRes{Food: *food, Message: "Food item successfully updated"})
	})

	route.Delete("/:id", func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.UserReq)
		id := c.Params("id")

		if err := foodService.Remove(user, id); err != nil {
			if err == fiber.ErrNotAcceptable {
				return fiber.NewError(fiber.StatusNotAcceptable, "Food is currently ordered, cannot delete")
			} else if err == fiber.ErrUnauthorized {
				return fiber.NewError(fiber.StatusNotAcceptable, "No permission for this food. Please try another food")
			} else {
				return fiber.ErrInternalServerError
			}
		}
		return c.JSON(models.MessageRes{Message: "Food item successfully deleted"})
	})
}
