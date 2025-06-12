package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLUser     string
	MySQLPassword string
	MySQLHost     string
	MySQLPort     int
	MySQLDatabase string
	RabbitMQUri   string
}

func LoadConfig() Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	mysqlPort, err := strconv.Atoi(getEnv("MYSQL_PORT", "3306"))
	if err != nil {
		log.Fatalf("Invalid MYSQL_PORT: %v", err)
	}

	return Config{
		MySQLUser:     getEnv("MYSQL_USER", "root"),
		MySQLPassword: getEnv("MYSQL_PASSWORD", "root"),
		MySQLHost:     getEnv("MYSQL_HOST", "mysql"),
		MySQLPort:     mysqlPort,
		MySQLDatabase: getEnv("MYSQL_DATABASE", "fasttrack"),
		RabbitMQUri:   getEnv("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/"),
	}
}

// helper function to read environment variables with default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
