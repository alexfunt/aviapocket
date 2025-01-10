package api

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

type Api struct {
	APIKey     string
	PartnerID  string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() (*Api, error) {
	err := godotenv.Load("/Users/xyziungvyncki/IntelliJ Projects/aviapocket/.env")
	if err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	apiKey := os.Getenv("AVIASALES_API_KEY")
	partnerID := os.Getenv("AVIASALES_PARTNER_ID")

	if apiKey == "" || partnerID == "" {
		return nil, fmt.Errorf("API Key or Partner ID is missing")
	}

	return &Api{
		APIKey:     apiKey,
		PartnerID:  partnerID,
		BaseURL:    "https://api.travelpayouts.com",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *Api) GetPriceTrends(origin, destination string, currency string) (map[string]interface{}, error) {
	var url string
	if origin != "" && destination != "" {
		url = fmt.Sprintf("%s/prices/latest?origin=%s&destination=%s&period_type=%s&sorting=%s&currency=%s&market=%s&limit=%d&page=%d&token=%s",
			c.BaseURL, origin, destination, "month", "price", currency, "ru", 30, 1, c.APIKey)
	} else if origin != "" {
		url = fmt.Sprintf("%s/prices/latest?origin=%s&period_type=%s&sorting=%s&currency=%s&market=%s&limit=%d&page=%d&token=%s",
			c.BaseURL, origin, "month", "price", currency, "ru", 30, 1, c.APIKey)
	} else if destination != "" {
		url = fmt.Sprintf("%s/prices/latest?destination=%s&period_type=%s&sorting=%s&currency=%s&market=%s&limit=%d&page=%d&token=%s",
			c.BaseURL, destination, "month", "price", currency, "ru", 30, 1, c.APIKey)
	} else {
		return nil, fmt.Errorf("Either origin or destination must be provided")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "flight-price-trends-app")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Api request failed with status code %s", resp.Status)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Api) GetPriceTrendsWithParams(params map[string]string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v2/prices/latest", c.BaseURL)
	query := url + "?"
	for key, value := range params {
		query += fmt.Sprintf("%s=%s&", key, value)
	}
	query += fmt.Sprintf("token=%s", c.APIKey)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "flight-price-trend-app")

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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}
