package web

import (
	"aviapocket/models"
	"aviapocket/services"
	"github.com/gofiber/fiber/v2"
	"log"
)

type FlightHandler struct {
	FlightPriceLoader *services.FlightPriceLoader
}

func NewFlightHandler(loader *services.FlightPriceLoader) *FlightHandler {
	return &FlightHandler{FlightPriceLoader: loader}
}

func (h *FlightHandler) GetFlights(c *fiber.Ctx) error {
	origin := c.Query("origin")
	destination := c.Query("destination")
	departDate := c.Query("depart_date")
	returnDate := c.Query("return_date", "")
	currency := c.Query("currency", "USD")

	if origin == "" || destination == "" || departDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "origin, destination, and depart_date are required",
		})
	}

	flights, err := h.FlightPriceLoader.ApiClient.FetchAndCompareFlights(origin, destination, departDate, returnDate, currency)
	if err != nil {
		log.Printf("Error fetching flights: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch flights",
		})
	}

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

		err := h.FlightPriceLoader.LoadFlightPrices(flightPrice)
		if err != nil {
			log.Printf("Error saving flight: %v", err)
		}
	}

	return c.JSON(flights)
}
