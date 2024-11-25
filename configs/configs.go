package configs

import (
	"os"
	"strconv"
)

type GEConfig struct {
	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	JWT_Secret string

	Midtrans MidtransConfig
	OpenAi   OpenAi
}

type MidtransConfig struct {
	ClientKey string
	ServerKey string
}

type OpenAi struct {
	ApiKey string
}

func InitConfig() *GEConfig {
	var res = new(GEConfig)

	res.DB_HOST = os.Getenv("DB_HOST")
	res.DB_PORT, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	res.DB_USER = os.Getenv("DB_USER")
	res.DB_PASSWORD = os.Getenv("DB_PASS")
	res.DB_NAME = os.Getenv("DB_NAME")

	res.JWT_Secret = os.Getenv("JWT_SECRET")

	res.Midtrans.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	res.Midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")

	res.OpenAi.ApiKey = os.Getenv("OPENAI_API_KEY")

	return res
}
