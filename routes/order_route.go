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
			return fiber.ErrInternalServerError
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

	route.Put(":id/:status", func(c fiber.Ctx) error {
		id := c.Params("id")
		status := c.Params("status")

		result, err := orderService.UpdateStatus(id, status)
		if err == fiber.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Order to update not found")
		} else if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(result)
	})

	route.Delete("/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		if err := orderService.Remove(id); err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(models.MessageRes{Message: "Order successfully canceled"})
	})

	route.Get("/all-order", func(c fiber.Ctx) error {
		orders, err := orderService.FindAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(orders)
	})
}
