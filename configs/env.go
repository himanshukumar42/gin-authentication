package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Loaded the .env file")
	return os.Getenv("MONGOURI")
}
