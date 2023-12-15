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

func TestCreateUpdatePowerPlant(t *testing.T) {
	query := `mutation ($input: NewPowerPlantInput!) { createPowerPlant(input: $input) { id name latitude longitude } }`
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"name":      "Berlin Wind Farm",
			"latitude":  52.636083,
			"longitude": 13.42977,
		},
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	assert.NoError(t, err)

	resp, err := http.Post("http://localhost:8080/graphql", "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response struct {
		Data struct {
			CreatePowerPlant struct {
				ID        string  `json:"id"`
				Name      string  `json:"name"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"createPowerPlant"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Data.CreatePowerPlant.ID)
	assert.Equal(t, "Berlin Wind Farm", response.Data.CreatePowerPlant.Name)
	assert.Equal(t, 52.636083, response.Data.CreatePowerPlant.Latitude)
	assert.Equal(t, 13.42977, response.Data.CreatePowerPlant.Longitude)

	// update the power plant name
	newName := "Berlin/Pankow Wind Farm Updated"

	query = `mutation UpdatePlant($id: ID!, $input: UpdatePowerPlantInput!) { 
        updatePowerPlant(id: $id, input: $input) { id name latitude longitude } 
    }`

	variables = map[string]interface{}{
		"id": response.Data.CreatePowerPlant.ID,
		"input": map[string]interface{}{
			"name": newName,
		},
	}

	requestBody, err = json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	assert.NoError(t, err)

	resp, err = http.Post("http://localhost:8080/graphql", "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	var updateResponse struct {
		Data struct {
			UpdatePowerPlant struct {
				ID        string  `json:"id"`
				Name      string  `json:"name"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"updatePowerPlant"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&updateResponse)
	assert.NoError(t, err)

	// Check if the name has been updated
	assert.Equal(t, newName, updateResponse.Data.UpdatePowerPlant.Name)
	assert.Equal(t, 52.636083, updateResponse.Data.UpdatePowerPlant.Latitude)
	assert.Equal(t, 13.42977, updateResponse.Data.UpdatePowerPlant.Longitude)
}
