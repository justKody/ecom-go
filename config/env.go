package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser        string
	DBPassword    string
	DBAddress     string
	DBName        string
	JWTExpiration int
	JWTSecret     string
}

var Envs = initConfig()

func initConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		PublicHost:    getEnv("PUBLIC_HOST", "http://localhost"),
		Port:          getEnv("PORT", "8080"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "Root@123"),
		DBAddress:     fmt.Sprintf("%s:%s", getEnv("DB_ADDRESS", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:        getEnv("DB_NAME", "ecom"),
		JWTExpiration: getEnvInt("JWT_EXPIRATION", 30),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback	
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Println("Error converting JWT expiration to int:", err)
			return 0
		}
		return intValue
	}
	return fallback
}
