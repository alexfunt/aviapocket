package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Api struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Api {
	if apiKey == "" {
		panic("API key is required")
	}
	return &Api{
		APIKey:     apiKey,
		BaseURL:    "https://api.travelpayouts.com",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Flight - структура для представления рейса.
type Flight struct {
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	DepartDate  string    `json:"depart_date"`
	ReturnDate  string    `json:"return_date"`
	Price       int       `json:"value"`
	Gate        string    `json:"gate"`
	FoundAt     time.Time `json:"found_at"`
}

// GetFlights - получает данные о рейсах по параметрам.
func (c *Api) GetFlights(params map[string]string) ([]Flight, error) {
	if params["origin"] == "" && params["destination"] == "" {
		return nil, errors.New("either origin or destination must be provided")
	}

	url := fmt.Sprintf("%s/v2/prices/latest?", c.BaseURL)
	for key, value := range params {
		url += fmt.Sprintf("%s=%s&", key, value)
	}
	url += fmt.Sprintf("token=%s", c.APIKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "flight-price-comparator")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %s", resp.Status)
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var flights []Flight
	for _, item := range result.Data {
		foundAt, err := time.Parse("2006-01-02T15:04:05", item["found_at"].(string))
		if err != nil {
			return nil, fmt.Errorf("invalid found_at format: %w", err)
		}

		flight := Flight{
			Origin:      item["origin"].(string),
			Destination: item["destination"].(string),
			DepartDate:  item["depart_date"].(string),
			ReturnDate:  item["return_date"].(string),
			Price:       int(item["value"].(float64)),
			Gate:        item["gate"].(string),
			FoundAt:     foundAt,
		}
		flights = append(flights, flight)
	}

	return flights, nil
}

// CompareFlights - сравнивает рейсы на схожесть.
func CompareFlights(originFlights, destinationFlights []Flight) []Flight {
	flightMap := make(map[string]Flight)
	var similarFlights []Flight

	// Добавляем рейсы с `origin` в карту.
	for _, flight := range originFlights {
		key := fmt.Sprintf("%s-%s-%s", flight.Origin, flight.Destination, flight.DepartDate)
		flightMap[key] = flight
	}

	// Сравниваем рейсы с `destination` на схожесть.
	for _, flight := range destinationFlights {
		key := fmt.Sprintf("%s-%s-%s", flight.Origin, flight.Destination, flight.DepartDate)
		if _, exists := flightMap[key]; exists {
			similarFlights = append(similarFlights, flight)
		}
	}

	return similarFlights
}

// FetchAndCompareFlights - основная функция получения и сравнения рейсов.
func (c *Api) FetchAndCompareFlights(origin, destination, departDate, returnDate, currency string) ([]Flight, error) {
	originFlights, err := c.GetFlights(map[string]string{
		"origin":      origin,
		"depart_date": departDate,
		"currency":    currency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flights with origin: %w", err)
	}

	destinationFlights, err := c.GetFlights(map[string]string{
		"destination": destination,
		"return_date": returnDate,
		"currency":    currency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flights with destination: %w", err)
	}

	similarFlights := CompareFlights(originFlights, destinationFlights)
	return similarFlights, nil
}
