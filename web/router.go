package web

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, flightHandler *FlightHandler) {
	app.Get("/api/flights", flightHandler.GetFlights)
}
