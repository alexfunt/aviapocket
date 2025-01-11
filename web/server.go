package web

import (
	"aviapocket/api"
	"aviapocket/db"
	"aviapocket/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App      *fiber.App
	Api      *api.Api
	DB       *db.DB
	Services *services.FlightPriceLoader
	Handler  *FlightHandler
}

func NewServer(apiClient *api.Api, dbClient *db.DB) *Server {
	flightPriceLoader := services.NewFlightPriceLoader(apiClient, dbClient.Connection)
	flightHandler := NewFlightHandler(flightPriceLoader)

	return &Server{
		App:      fiber.New(),
		Api:      apiClient,
		DB:       dbClient,
		Services: flightPriceLoader,
		Handler:  flightHandler,
	}
}

func (s *Server) Start(port string) {
	s.setupRoutes()

	log.Printf("Starting server on port %s", port)
	if err := s.App.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) setupRoutes() {
	s.App.Get("/api/flights", s.Handler.GetFlights)
}
