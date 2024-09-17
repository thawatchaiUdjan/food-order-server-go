package config

import (
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host string
	Name string
}

type GoogleConfig struct {
	ClientID    string
	SecretID    string
	RedirectURL string
}

type CloudinaryConfig struct {
	ApiName   string
	ApiKey    string
	ApiSecret string
}

type FacebookConfig struct {
	ClientID string
	SecretID string
}

type TokenConfig struct {
	TokenExpiredError string
	TokenExpiredTime  string
}

type UploadFileConfig struct {
	Format        string
	FoodFolder    string
	ProfileFolder string
}

type Config struct {
	Port              string
	JWTSecret         string
	Database          DatabaseConfig
	Cloudinary        CloudinaryConfig
	Google            GoogleConfig
	Facebook          FacebookConfig
	Token             TokenConfig
	UploadFile        UploadFileConfig
	EncryptSaltRounds int
}

func LoadConfig() Config {
	godotenv.Load()

	return Config{
		Port:      os.Getenv("APP_PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Database: DatabaseConfig{
			Host: os.Getenv("DB_CONNECTION_STRING_MONGO"),
			Name: "food_order_db",
		},
		Cloudinary: CloudinaryConfig{
			ApiName:   os.Getenv("CLOUDINARY_NAME"),
			ApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
			ApiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
		},
		Google: GoogleConfig{
			ClientID:    os.Getenv("GOOGLE_CLIENT_ID"),
			SecretID:    os.Getenv("GOOGLE_SECRET_ID"),
			RedirectURL: "postmessage",
		},
		Facebook: FacebookConfig{
			ClientID: os.Getenv("FACEBOOK_CLIENT_ID"),
			SecretID: os.Getenv("FACEBOOK_SECRET_ID"),
		},
		Token: TokenConfig{
			TokenExpiredError: "TokenExpiredError",
			TokenExpiredTime:  "6h",
		},
		UploadFile: UploadFileConfig{
			Format:        "png",
			FoodFolder:    "foods",
			ProfileFolder: "profiles",
		},
		EncryptSaltRounds: 10,
	}
}
