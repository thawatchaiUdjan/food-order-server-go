package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func LoadCloudinary() *cloudinary.Cloudinary {
	config := LoadConfig()
	cld, err := cloudinary.NewFromParams(config.Cloudinary.ApiName, config.Cloudinary.ApiKey, config.Cloudinary.ApiSecret)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	return cld
}
