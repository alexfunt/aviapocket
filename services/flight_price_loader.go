package services

import (
	"aviapocket/api"
	"aviapocket/models"
	"database/sql"
	"fmt"
	"log"
)

type FlightPriceLoader struct {
	ApiClient *api.Api
	DB        *sql.DB
}

// NewFlightPriceLoader - создает экземпляр FlightPriceLoader.
func NewFlightPriceLoader(apiClient *api.Api, db *sql.DB) *FlightPriceLoader {
	return &FlightPriceLoader{
		ApiClient: apiClient,
		DB:        db,
	}
}

// LoadFlightPrices - загружает рейсы из API и сохраняет их в базу данных.
func (loader *FlightPriceLoader) LoadFlightPrices(origin, destination, departDate, returnDate, currency string) error {
	// Получаем рейсы из API.
	flights, err := loader.ApiClient.FetchAndCompareFlights(origin, destination, departDate, returnDate, currency)
	if err != nil {
		return fmt.Errorf("failed to fetch flights: %w", err)
	}

	// Сохраняем каждый рейс в базу данных.
	for _, flight := range flights {
		flightPrice := models.FlightPrice{
			Origin:      flight.Origin,
			Destination: flight.Destination,
			DepartDate:  flight.DepartDate,
			ReturnDate:  flight.ReturnDate,
			Price:       flight.Price,
			Gate:        flight.Gate,
			FoundAt:     flight.FoundAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if err := loader.saveToDB(flightPrice); err != nil {
			log.Printf("Failed to save flight price: Origin=%s, Destination=%s, Error=%v",
				flightPrice.Origin, flightPrice.Destination, err)
		}
	}
	return nil
}

// saveToDB - сохраняет информацию о рейсе в базу данных.
func (loader *FlightPriceLoader) saveToDB(flight models.FlightPrice) error {
	_, err := loader.DB.Exec(
		`INSERT INTO flight_prices 
		 (origin, destination, depart_date, return_date, price, gate, found_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT DO NOTHING`, // Предотвращает дублирование записей.
		flight.Origin, flight.Destination, flight.DepartDate, flight.ReturnDate,
		flight.Price, flight.Gate, flight.FoundAt, // Это уже строка в нужном формате.
	)
	return err
}
