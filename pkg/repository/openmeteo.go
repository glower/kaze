package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherForecastResponse struct {
	Hourly struct {
		Time             []string  `json:"time"`
		Precipitation    []float64 `json:"precipitation"`
		WindSpeed10m     []float64 `json:"wind_speed_10m"`
		Temperature2m    []float64 `json:"temperature_2m"`
		WindDirection10m []float64 `json:"wind_direction_10m"`
	} `json:"hourly"`
}

type ElevationResponse struct {
	Elevation []float64 `json:"elevation"`
}

type OpenMeteoRepository interface {
	GetElevation(ctx context.Context, latitude, longitude float64) (float64, error)
	GetWeatherForecast(ctx context.Context, latitude, longitude float64) (*WeatherForecastResponse, error)
}

type openMeteoRepo struct {
	apiKey string
}

func NewOpenMeteoRepository(apiKey string) OpenMeteoRepository {
	return &openMeteoRepo{
		apiKey: apiKey,
	}
}

func (r *openMeteoRepo) GetElevation(ctx context.Context, latitude, longitude float64) (float64, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/elevation?latitude=%f&longitude=%f&apikey=%s", latitude, longitude, r.apiKey)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making request to Open-Meteo API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-OK response from Open-Meteo API: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	var elevationResponse ElevationResponse
	if err := json.Unmarshal(body, &elevationResponse); err != nil {
		return 0, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	if len(elevationResponse.Elevation) == 0 {
		return 0, fmt.Errorf("no elevation data received")
	}

	return elevationResponse.Elevation[0], nil
}

func (r *openMeteoRepo) GetWeatherForecast(ctx context.Context, latitude, longitude float64) (*WeatherForecastResponse, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation,wind_speed_10m,wind_direction_10m", latitude, longitude)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to Open-Meteo API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response from Open-Meteo API: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var forecastResponse WeatherForecastResponse
	if err := json.Unmarshal(body, &forecastResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	return &forecastResponse, nil
}
