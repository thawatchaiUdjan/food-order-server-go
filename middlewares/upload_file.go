package middlewares

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/food-order-server/config"
	"github.com/food-order-server/models"
	"github.com/food-order-server/utils"
	"github.com/gofiber/fiber/v2"
)

func uploadFile(c *fiber.Ctx, name string, id string, folder string) error {
	cld := config.LoadCloudinary()
	config := config.LoadConfig()

	file, err := c.FormFile(name)
	if err != nil {
		c.Locals("file", "")
		c.Locals("id", "")
		return c.Next()
	}

	fileStream, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return fiber.ErrInternalServerError
	}
	defer fileStream.Close()

	publicId := utils.GenerateHash(id)
	uploadResponse, err := cld.Upload.Upload(c.Context(), fileStream, uploader.UploadParams{
		Format:   config.UploadFile.Format,
		Folder:   folder,
		PublicID: publicId,
	})
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return fiber.ErrInternalServerError
	}

	c.Locals("file", uploadResponse.SecureURL)
	c.Locals("id", id)

	return c.Next()
}

func UploadFoodFile(c *fiber.Ctx) error {
	config := config.LoadConfig()
	name := "food_image_url"
	folder := config.UploadFile.FoodFolder
	id := c.Params("id")
	if id == "" {
		id = utils.GenerateUuid()
	}

	if err := uploadFile(c, name, id, folder); err != nil {
		return err
	}
	return nil
}

func UploadProfileFile(c *fiber.Ctx) error {
	config := config.LoadConfig()
	name := "profile_image_url"
	folder := config.UploadFile.ProfileFolder
	id := c.Locals("user").(models.UserReq).User.UserID

	if err := uploadFile(c, name, id, folder); err != nil {
		return err
	}
	return nil
}
