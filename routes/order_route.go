package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/models"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderRoute(app *fiber.App, db *mongo.Database) {
	orderService := services.CreateOrderService(db)
	route := app.Group("/orders", middlewares.AuthToken)

	route.Get("/", func(c fiber.Ctx) error {
		req := c.Locals("user").(models.UserReq)

		order, err := orderService.FindOne(req.User.UserID)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(order)
	})
}
