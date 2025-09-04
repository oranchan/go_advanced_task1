package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// Load variables from .env if present; fallback to OS env if missing
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}
}

func main() {
	//QueryBlockInfoAt()
	//transfer()
	//DeployCounter()
	//InteractWithCounter()
}
