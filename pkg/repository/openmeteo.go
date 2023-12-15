package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// WeatherForecastResponse represents the response structure for weather forecasts.
type WeatherForecastResponse struct {
	Hourly Hourly `json:"hourly"`
}

type Hourly struct {
	Time             []string  `json:"time"`
	Precipitation    []float64 `json:"precipitation"`
	WindSpeed10m     []float64 `json:"wind_speed_10m"`
	Temperature2m    []float64 `json:"temperature_2m"`
	WindDirection10m []float64 `json:"wind_direction_10m"`
}

// ElevationResponse represents the response structure for elevation data.
type ElevationResponse struct {
	Elevation []float64 `json:"elevation"`
}

// OpenMeteoRepository defines the interface for interacting with the Open-Meteo API.
//
//go:generate go run github.com/vektra/mockery/v2@v2 --name=OpenMeteoRepository --filename=open_meteo_repository.go --output=../../mocks/
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

// GetElevation retrieves elevation data from the Open-Meteo API.
func (r *openMeteoRepo) GetElevation(ctx context.Context, latitude, longitude float64) (float64, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/elevation?latitude=%f&longitude=%f&apikey=%s", latitude, longitude, r.apiKey)
	slog.Debug("Fetching elevation data", "url", url)

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

// GetWeatherForecast retrieves weather forecast data from the Open-Meteo API.
func (r *openMeteoRepo) GetWeatherForecast(ctx context.Context, latitude, longitude float64) (*WeatherForecastResponse, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation,wind_speed_10m,wind_direction_10m", latitude, longitude)
	slog.Debug("Fetching weather forecast data", "url", url)

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
