package routes

import (
	"github.com/food-order-server/middlewares"
	"github.com/food-order-server/models"
	"github.com/food-order-server/services"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoute(app *fiber.App, db *mongo.Database) {
	userService := services.CreateUserService(db)
	route := app.Group("/user")

	route.Get("/verify-token", func(c fiber.Ctx) error {
		return c.SendString("verify complete")
	}, middlewares.AuthToken)

	route.Get("/", func(c fiber.Ctx) error {
		req := c.Locals("user").(models.UserReq)

		user, err := userService.FindOne(req.User.UserID, req.Token)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(user)
	}, middlewares.AuthToken)

	route.Post("/login", func(c fiber.Ctx) error {
		userBody := new(models.UserLoginReq)

		if err := c.Bind().Body(userBody); err != nil {
			return fiber.ErrBadRequest
		}

		user, err := userService.Login(userBody)
		if err == fiber.ErrBadRequest {
			return fiber.NewError(fiber.StatusUnauthorized, "Username or password invalid")
		} else if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(user)
	})

	route.Post("/register", func(c fiber.Ctx) error {
		userBody := new(models.UserRegisterReq)

		if err := c.Bind().Body(userBody); err != nil {
			return fiber.ErrBadRequest
		}

		user, err := userService.Register(userBody)
		if err == fiber.ErrConflict {
			return fiber.NewError(fiber.StatusConflict, "Username is already in use")
		} else if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(user)
	})

	route.Post("/google-login", func(c fiber.Ctx) error {
		googleLoginBody := new(models.UserGoogleLoginReq)

		if err := c.Bind().Body(googleLoginBody); err != nil {
			return fiber.ErrBadRequest
		}

		user, err := userService.GoogleLogin(googleLoginBody.Code)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return c.JSON(user)
	})

	route.Put("/", func(c fiber.Ctx) error {
		userBody := new(models.User)
		req := c.Locals("user").(models.UserReq)
		file := c.Locals("file").(string)

		if err := c.Bind().Body(userBody); err != nil {
			return fiber.ErrBadRequest
		}

		result, err := userService.Update(req.User.UserID, userBody, file)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(result)
	}, middlewares.AuthToken, middlewares.UploadProfileFile)

	route.Delete("/", func(c fiber.Ctx) error {
		req := c.Locals("user").(models.UserReq)

		result, err := userService.Remove(req.User.UserID)
		if err == fiber.ErrNotAcceptable {
			return fiber.NewError(fiber.StatusNotAcceptable, "account have an order, cant be delete")
		} else if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(result)
	}, middlewares.AuthToken)
}
