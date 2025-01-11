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

func NewFlightPriceLoader(apiClient *api.Api, db *sql.DB) *FlightPriceLoader {
	return &FlightPriceLoader{
		ApiClient: apiClient,
		DB:        db,
	}
}

func (loader *FlightPriceLoader) LoadFlightPrices(flight models.Flight) error {

	flights, err := loader.ApiClient.FetchAndCompareFlights(flight.Origin, flight.Destination, flight.DepartDate, flight.ReturnDate, "rub")
	if err != nil {
		return fmt.Errorf("failed to fetch flights: %w", err)
	}

	tx, err := loader.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	for _, flight := range flights {
		flightPrice := models.Flight{
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

func (loader *FlightPriceLoader) saveToDB(flight models.Flight) error {
	log.Printf("Saving flight: %+v", flight)

	if flight.ReturnDate == "" {
		flight.ReturnDate = ""
	}
	if flight.Gate == "" {
		flight.Gate = ""
	}

	_, err := loader.DB.Exec(
		`INSERT INTO flight_prices 
         (origin, destination, depart_date, return_date, price, gate, found_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7)
         ON CONFLICT DO NOTHING`,
		flight.Origin, flight.Destination, flight.DepartDate, flight.ReturnDate,
		flight.Price, flight.Gate, flight.FoundAt,
	)
	if err != nil {
		log.Printf("Error saving flight: %v", err)
	}
	return err
}
