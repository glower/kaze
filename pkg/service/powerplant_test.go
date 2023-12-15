package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glower/kaze/graph/model"
	"github.com/glower/kaze/mocks"
	"github.com/glower/kaze/pkg/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePowerPlant(t *testing.T) {
	t.Run("failed due to database error", func(t *testing.T) {
		service, mockDB, _ := setupTests(t)
		validPlant := &model.PowerPlant{Name: "Valid Plant", Latitude: 10.0, Longitude: 20.0}

		mockDB.On("Create", mock.Anything, validPlant).Return(nil, assert.AnError)

		_, err := service.CreatePowerPlant(context.Background(), validPlant)
		assert.Error(t, err)

		mockDB.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		service, mockDB, _ := setupTests(t)
		validPlant := &model.PowerPlant{Name: "Valid Plant", Latitude: 10.0, Longitude: 20.0}

		mockDB.On("Create", mock.Anything, validPlant).Return(validPlant, nil)

		result, err := service.CreatePowerPlant(context.Background(), validPlant)
		assert.NoError(t, err, "Expected no error on successful creation")
		assert.Equal(t, validPlant, result, "Expected created plant to match input")

		mockDB.AssertExpectations(t)
	})
}

func TestGetPowerPlant(t *testing.T) {
	t.Run("failed to retrieve non-existent plant", func(t *testing.T) {
		service, mockDB, _ := setupTests(t)

		mockDB.On("GetByID", mock.Anything, "non-existent-id").Return(nil, fmt.Errorf("not found"))

		_, err := service.GetPowerPlant(context.Background(), "non-existent-id", false, false)
		assert.Error(t, err)

		mockDB.AssertExpectations(t)
	})

	t.Run("failed to fetch elevation data", func(t *testing.T) {
		service, mockDB, mockOpenMeteo := setupTests(t)

		validPlant := &model.PowerPlant{ID: "1", Name: "Valid Plant"}

		mockDB.On("GetByID", mock.Anything, "1").Return(validPlant, nil)
		mockOpenMeteo.On("GetElevation", mock.Anything, validPlant.Latitude, validPlant.Longitude).Return(0.0, fmt.Errorf("elevation error"))

		_, err := service.GetPowerPlant(context.Background(), "1", true, false)
		assert.Error(t, err)

		mockDB.AssertExpectations(t)
		mockOpenMeteo.AssertExpectations(t)
	})

	t.Run("success with elevation and weather forecasts", func(t *testing.T) {
		service, mockDB, mockOpenMeteo := setupTests(t)

		validPlant := &model.PowerPlant{ID: "1", Name: "Valid Plant"}

		mockDB.On("GetByID", mock.Anything, "1").Return(validPlant, nil)
		mockOpenMeteo.On("GetElevation", mock.Anything, validPlant.Latitude, validPlant.Longitude).Return(100.0, nil)
		mockOpenMeteo.On("GetWeatherForecast", mock.Anything, validPlant.Latitude, validPlant.Longitude).
			Return(&repository.WeatherForecastResponse{
				Hourly: repository.Hourly{
					Time:             []string{"12"},
					Precipitation:    []float64{100.1}, // to test HasPrecipitationToday
					WindSpeed10m:     []float64{10},
					Temperature2m:    []float64{10},
					WindDirection10m: []float64{10},
				},
			}, nil)

		result, err := service.GetPowerPlant(context.Background(), "1", true, true)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.HasPrecipitationToday)

		mockDB.AssertExpectations(t)
		mockOpenMeteo.AssertExpectations(t)
	})
}

func TestListPowerPlants(t *testing.T) {

	t.Run("failed due to database error", func(t *testing.T) {
		service, mockDB, _ := setupTests(t)

		mockDB.On("List", mock.Anything, 0, 10).Return(nil, 0, fmt.Errorf("database error"))
		
		_, err := service.ListPowerPlants(context.Background(), 1, 10, false, false)
		assert.Error(t, err)
		
		mockDB.AssertExpectations(t)
	})

	t.Run("success with elevation and weather forecasts", func(t *testing.T) {
		service, mockDB, mockOpenMeteo := setupTests(t)

		plants := []model.PowerPlant{{ID: "1"}, {ID: "2"}}
		
		mockDB.On("List", mock.Anything, 0, 10).Return(plants, len(plants), nil)
		mockOpenMeteo.On("GetElevation", mock.Anything, mock.AnythingOfType("float64"), mock.AnythingOfType("float64")).Return(100.0, nil).Twice()
		mockOpenMeteo.On("GetWeatherForecast", mock.Anything, mock.AnythingOfType("float64"), mock.AnythingOfType("float64")).
			Return(&repository.WeatherForecastResponse{}, nil).Twice()

		result, err := service.ListPowerPlants(context.Background(), 1, 10, true, true)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.TotalCount)
		assert.Len(t, result.PowerPlants, 2)

		mockDB.AssertExpectations(t)
		mockOpenMeteo.AssertExpectations(t)
	})
}

func setupTests(t *testing.T) (PowerPlantService, *mocks.PowerPlantRepository, *mocks.OpenMeteoRepository) {
	mockDB := mocks.NewPowerPlantRepository(t)
	mockOpenMeteo := mocks.NewOpenMeteoRepository(t)

	service := NewPowerPlantService(mockDB, mockOpenMeteo)

	return service, mockDB, mockOpenMeteo
}
