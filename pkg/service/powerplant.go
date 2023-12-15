package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/glower/kaze/graph/model"
	"github.com/glower/kaze/pkg/repository"
)

// PowerPlantService defines the interface for operations on power plants.
//
//go:generate go run github.com/vektra/mockery/v2@v2 --name=PowerPlantService --filename=power_plant_service.go --output=../../mocks/
type PowerPlantService interface {
	CreatePowerPlant(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error)
	UpdatePowerPlant(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error)
	GetPowerPlant(ctx context.Context, id string, withElevation, withWeatherForecasts bool) (*model.PowerPlant, error)
	ListPowerPlants(ctx context.Context, page, pageSize int, withElevation, withWeatherForecasts bool) (*model.PowerPlantList, error)
}

// powerPlantService provides services related to power plants.
type powerPlantService struct {
	dbRepo        repository.PowerPlantRepository
	openMeteoRepo repository.OpenMeteoRepository
}

// NewPowerPlantService creates a new instance of PowerPlantService.
func NewPowerPlantService(dbRepo repository.PowerPlantRepository, openMeteoRepo repository.OpenMeteoRepository) PowerPlantService {
	return &powerPlantService{
		dbRepo:        dbRepo,
		openMeteoRepo: openMeteoRepo,
	}
}

// CreatePowerPlant handles the creation of a new power plant.
func (s *powerPlantService) CreatePowerPlant(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error) {
	slog.Debug("Creating a new power plant", "name", plant.Name)
	return s.dbRepo.Create(ctx, plant)
}

// UpdatePowerPlant handles updating an existing power plant.
func (s *powerPlantService) UpdatePowerPlant(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error) {
	slog.Debug("Updating power plant", "id", plant.ID)
	_, err := s.dbRepo.Update(ctx, plant)
	if err != nil {
		return nil, fmt.Errorf("could not update power plant: %w", err)
	}
	// Retrieve and return the updated plant
	return s.dbRepo.GetByID(ctx, plant.ID)
}

// GetPowerPlant retrieves a specific power plant by its ID.
func (s *powerPlantService) GetPowerPlant(ctx context.Context, id string, withElevation, withWeatherForecasts bool) (*model.PowerPlant, error) {
	slog.Debug("Retrieving power plant", "id", id, "withElevation", withElevation, "withWeatherForecasts", withWeatherForecasts)
	plant, err := s.dbRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Fetch additional data based on request
	if withElevation {
		if err := s.fetchElevation(ctx, plant); err != nil {
			return nil, err
		}
	}

	if withWeatherForecasts {
		if err := s.fetchWeatherForecasts(ctx, plant); err != nil {
			return nil, err
		}
	}

	return plant, nil
}

// ListPowerPlants retrieves a list of power plants with optional elevation and weather forecast data.
func (s *powerPlantService) ListPowerPlants(ctx context.Context, page, pageSize int, withElevation, withWeatherForecasts bool) (*model.PowerPlantList, error) {
	offset := (page - 1) * pageSize
	slog.Debug("Listing power plants", "page", page, "pageSize", pageSize)
	plants, total, err := s.dbRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	var plantPointers []*model.PowerPlant
	for _, plant := range plants {
		if withElevation {
			if err := s.fetchElevation(ctx, &plant); err != nil {
				return nil, err
			}
		}
		if withWeatherForecasts {
			if err := s.fetchWeatherForecasts(ctx, &plant); err != nil {
				return nil, err
			}
		}
		plantPointers = append(plantPointers, &plant)
	}

	return &model.PowerPlantList{
		PowerPlants: plantPointers,
		TotalCount:  total,
	}, nil
}

func (s *powerPlantService) fetchElevation(ctx context.Context, plant *model.PowerPlant) error {
	elevation, err := s.openMeteoRepo.GetElevation(ctx, plant.Latitude, plant.Longitude)
	if err != nil {
		return fmt.Errorf("can't get elevation data from the api: %w", err)
	}
	plant.Elevation = elevation
	return nil
}

// Helper functions for fetching additional data (elevation and weather forecasts) are defined below...

func (s *powerPlantService) fetchWeatherForecasts(ctx context.Context, plant *model.PowerPlant) error {
	forecast, err := s.openMeteoRepo.GetWeatherForecast(ctx, plant.Latitude, plant.Longitude)
	if err != nil {
		return fmt.Errorf("can't get weather forecast data from the api: %w", err)
	}

	mappedForecast, hasPrecipitationToday := mapHourlyWeatherDataToForecasts(forecast)
	plant.WeatherForecasts = mappedForecast
	plant.HasPrecipitationToday = hasPrecipitationToday

	return nil
}

func mapHourlyWeatherDataToForecasts(response *repository.WeatherForecastResponse) ([]*model.WeatherForecast, bool) {
	var forecasts []*model.WeatherForecast
	hasPrecipitationToday := false

	for i, timeStr := range response.Hourly.Time {
		if response.Hourly.Precipitation[i] > 0 {
			hasPrecipitationToday = true
		}

		forecast := &model.WeatherForecast{
			Time:          timeStr,
			Temperature:   response.Hourly.Temperature2m[i],
			WindSpeed:     response.Hourly.WindSpeed10m[i],
			Precipitation: response.Hourly.Precipitation[i],
			WindDirection: response.Hourly.WindDirection10m[i],
		}
		forecasts = append(forecasts, forecast)
	}

	return forecasts, hasPrecipitationToday
}
