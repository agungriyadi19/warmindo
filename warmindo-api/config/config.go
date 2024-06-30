package config

import (
	"os"
)

var (
	POSTGRES_USER     = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB       = os.Getenv("POSTGRES_DB")
	CLIENT_URL        = os.Getenv("CLIENT_URL")
	SERVER_PORT       = os.Getenv("SERVER_PORT")
	JWT_KEY           = os.Getenv("JWT_KEY")
)
