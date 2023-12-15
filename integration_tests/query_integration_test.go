//go:build integration
// +build integration

package integrationtests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPowerPlant(t *testing.T) {
	plantID := "1"

	query := `query getPowerPlant($id: ID!) { 
        powerPlant(id: $id) { 
            id name latitude longitude hasPrecipitationToday elevation 
            weatherForecasts { time temperature precipitation windSpeed windDirection } 
        } 
    }`
	variables := map[string]interface{}{
		"id": plantID,
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	assert.NoError(t, err)

	resp, err := http.Post("http://localhost:8080/graphql", "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	var response struct {
		Data struct {
			PowerPlant struct {
				ID                    string  `json:"id"`
				Name                  string  `json:"name"`
				Latitude              float64 `json:"latitude"`
				Longitude             float64 `json:"longitude"`
				HasPrecipitationToday bool    `json:"hasPrecipitationToday"`
				Elevation             float64 `json:"elevation"`
				WeatherForecasts      []struct {
					Time          string  `json:"time"`
					Temperature   float64 `json:"temperature"`
					Precipitation float64 `json:"precipitation"`
					WindSpeed     float64 `json:"windSpeed"`
					WindDirection float64 `json:"windDirection"`
				} `json:"weatherForecasts"`
			} `json:"powerPlant"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, plantID, response.Data.PowerPlant.ID)
	assert.Equal(t, "Futaba Solar Plant", response.Data.PowerPlant.Name)
	assert.Equal(t, float64(34), response.Data.PowerPlant.Elevation)
	assert.NotEmpty(t, response.Data.PowerPlant.WeatherForecasts)
}

func TestListPowerPlants(t *testing.T) {
	query := `query ListPowerPlants($page: Int, $pageSize: Int) { 
        listPowerPlants(page: $page, pageSize: $pageSize) { 
            powerPlants { id name latitude longitude } 
            totalCount 
        } 
    }`
	variables := map[string]interface{}{
		"page":     1,
		"pageSize": 10,
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	assert.NoError(t, err)

	resp, err := http.Post("http://localhost:8080/graphql", "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	var response struct {
		Data struct {
			ListPowerPlants struct {
				PowerPlants []struct {
					ID        string  `json:"id"`
					Name      string  `json:"name"`
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"powerPlants"`
				TotalCount int `json:"totalCount"`
			} `json:"listPowerPlants"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.NotEqual(t, 0, response.Data.ListPowerPlants.TotalCount)
	assert.Equal(t, 10, len(response.Data.ListPowerPlants.PowerPlants))

	assert.Equal(t, "1", response.Data.ListPowerPlants.PowerPlants[0].ID)
	assert.Equal(t, "Futaba Solar Plant", response.Data.ListPowerPlants.PowerPlants[0].Name)
	assert.Equal(t, 37.4513, response.Data.ListPowerPlants.PowerPlants[0].Latitude)
	assert.Equal(t, 141.0334, response.Data.ListPowerPlants.PowerPlants[0].Longitude)
}
