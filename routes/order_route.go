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

		foodOrder, err := orderService.FindOne(req.User.UserID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotAcceptable, err.Error())
		}
		return c.JSON(foodOrder)
	})

	route.Post("/", func(c fiber.Ctx) error {
		req := c.Locals("user").(models.UserReq)

		orderReq := new(models.OrderReq)
		if err := c.Bind().Body(orderReq); err != nil {
			return fiber.ErrBadRequest
		}

		foodOrder, err := orderService.Create(&req.User, orderReq)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(foodOrder)
	})

	route.Delete("/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		if err := orderService.Remove(id); err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(models.MessageRes{Message: "Order successfully canceled"})
	})
}
