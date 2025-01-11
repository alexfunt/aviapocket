package main

import (
	"aviapocket/api"
	"aviapocket/db"
	"aviapocket/services"
	"fmt"
	"log"
)

func main() {
	client := api.NewClient("1f9498bff1444abf819c027acbd3f4d9")
	similarFlights, err := client.FetchAndCompareFlights("LED", "CEK", "2025-01-12", "2025-01-26", "rub")
	if err != nil {
		log.Fatalf("Error fetching and comparing flights: %v", err)
	}

	for _, flight := range similarFlights {
		fmt.Printf("Flight: %+v\n", flight)
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	loader := services.NewFlightPriceLoader(client, database)
	if err := loader.LoadFlightPrices("LED", "CEK", "2025-01-12", "2025-01-26", "rub"); err != nil {
		log.Fatalf("Failed to load flight prices: %v", err)
	}

	log.Println("Flight prices loaded successfully")
}
