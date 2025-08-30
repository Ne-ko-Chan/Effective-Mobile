package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	ConnString string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using fallback values: ", err)
	}
	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		ConnString: fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s port=%s",
			getEnv("DB_HOST", "postgres"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_NAME", "service"),
			getEnv("DB_SSLMODE", "disable"),
		  getEnv("DB_PASSWORD", "admin"),
		  getEnv("DB_PORT", "5432")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
