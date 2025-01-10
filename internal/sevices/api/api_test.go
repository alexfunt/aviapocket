package api

import "testing"

func TestGetFlightsByOriginAndDestination(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %s", err)
	}

	origin := "LED"
	destination := "KLF"
	currency := "rub"
	departDate := "2025-01-18"

	t.Logf("Fetching flights with origin=%s", origin)
	originResult, err := client.GetPriceTrendsWithParams(map[string]string{
		"origin":      origin,
		"currency":    currency,
		"depart_date": departDate,
	})
	if err != nil {
		t.Fatalf("Failed to fetch flights by origin: %s", err)
	}

	t.Logf("Filtering flights with destination=%s", destination)
	matchingFlights := filterFlightsByDestination(originResult, destination)

	// Проверка совпадений
	if len(matchingFlights) == 0 {
		t.Fatalf("No matching flights found with destination=%s", destination)
	}

	t.Logf("Matching flights: %+v", matchingFlights)
}

func filterFlightsByDestination(flights []map[string]interface{}, destination string) []map[string]interface{} {
	filtered := []map[string]interface{}{}
	for _, flight := range flights {
		if flight["destination"] == destination {
			filtered = append(filtered, flight)
		}
	}
	return filtered
}
