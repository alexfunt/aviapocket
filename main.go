package main

import (
	"aviapocket/api"
	"aviapocket/db"
	"aviapocket/services"
	"aviapocket/web"
	_ "fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {

	apiKey := os.Getenv("AVIASALES_API_KEY")
	if apiKey == "" {
		log.Fatalf("API key is missing")
	}

	client := api.NewClient(apiKey)

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	loader := services.NewFlightPriceLoader(client, database)

	flightHandler := web.NewFlightHandler(loader)

	app := fiber.New()
	web.SetupRouter(app, flightHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
