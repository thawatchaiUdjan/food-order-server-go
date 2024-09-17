package middlewares

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/food-order-server/config"
	"github.com/food-order-server/utils"
	"github.com/gofiber/fiber/v3"
)

func UploadFoodFile(c fiber.Ctx) error {
	cld := config.LoadCloudinary()
	config := config.LoadConfig()

	file, err := c.FormFile("food_image_url")
	if err != nil {
		return c.Next()
	}

	fileStream, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return fiber.ErrInternalServerError
	}
	defer fileStream.Close()

	foodId := utils.GenerateUuid()
	publicId := utils.GenerateHash(foodId)
	uploadResponse, err := cld.Upload.Upload(c.Context(), fileStream, uploader.UploadParams{
		Format:   config.UploadFile.Format,
		Folder:   config.UploadFile.FoodFolder,
		PublicID: publicId,
	})
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return fiber.ErrInternalServerError
	}

	c.Locals("file", uploadResponse.SecureURL)
	c.Locals("id", foodId)

	return c.Next()
}
