package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type GEConfig struct {
	APP_PORT string

	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	JWT_Secret string

	Cloudinary CloudinaryConfig
	Midtrans   MidtransConfig
	OpenAi     OpenAi
	Google     Google
}

type CloudinaryConfig struct {
	CloudName        string
	ApiKeyStorage    string
	ApiSecretStorage string
}

type MidtransConfig struct {
	ClientKey string
	ServerKey string
}

type Google struct {
	ClientID  string
	ClientKey string
}

type OpenAi struct {
	ApiKey string
}

func InitConfig() *GEConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic("Failed loading .env file")
	}

	var res = new(GEConfig)

	res.APP_PORT = os.Getenv("APP_PORT")

	res.DB_HOST = os.Getenv("DB_HOST")
	res.DB_PORT, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	res.DB_USER = os.Getenv("DB_USER")
	res.DB_PASSWORD = os.Getenv("DB_PASS")
	res.DB_NAME = os.Getenv("DB_NAME")

	res.JWT_Secret = os.Getenv("JWT_SECRET")

	res.Cloudinary.CloudName = os.Getenv("CLOUD_NAME")
	res.Cloudinary.ApiKeyStorage = os.Getenv("STORAGE_API_KEY")
	res.Cloudinary.ApiSecretStorage = os.Getenv("STORAGE_API_SECRET")

	res.Midtrans.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	res.Midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")

	res.OpenAi.ApiKey = os.Getenv("OPENAI_API_KEY")

	res.Google.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	res.Google.ClientKey = os.Getenv("GOOGLE_CLIENT_SECRET")

	return res
}
