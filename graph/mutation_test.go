package graph

import (
	"context"
	"testing"

	"github.com/glower/kaze/graph/model"
	"github.com/glower/kaze/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePowerPlant(t *testing.T) {
	ctx := context.Background()

	t.Run("fail due to invalid input", func(t *testing.T) {
		resolver, mockService := setupTests(t)

		// Invalid input, as required fields are missing
		input := model.NewPowerPlantInput{}

		_, err := resolver.CreatePowerPlant(ctx, input)
		assert.Error(t, err)
		mockService.AssertNotCalled(t, "CreatePowerPlant", mock.Anything)
	})

	t.Run("fail due to service error", func(t *testing.T) {
		resolver, mockService := setupTests(t)

		input := model.NewPowerPlantInput{Name: "Solar Plant", Latitude: 1.234, Longitude: 5.678}
		mockService.On("CreatePowerPlant", ctx, mock.Anything).Return(nil, assert.AnError).Once()

		_, err := resolver.CreatePowerPlant(ctx, input)
		assert.NotEmpty(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("succeed creating power plant", func(t *testing.T) {
		resolver, mockService := setupTests(t)

		input := model.NewPowerPlantInput{Name: "Solar Plant", Latitude: 1.234, Longitude: 5.678}
		output := &model.PowerPlant{ID: "1", Name: "Solar Plant", Latitude: 1.234, Longitude: 5.678}
		mockService.On("CreatePowerPlant", ctx, mock.Anything).Return(output, nil).Once()

		result, err := resolver.CreatePowerPlant(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, output, result)
		mockService.AssertExpectations(t)
	})
}

func TestUpdatePowerPlant(t *testing.T) {
	ctx := context.Background()

	t.Run("fail due to invalid input", func(t *testing.T) {
		resolver, mockService := setupTests(t)
		// Invalid input: wrong format
		input := model.UpdatePowerPlantInput{
			Latitude: floatPointer(2000),
		}

		_, err := resolver.UpdatePowerPlant(ctx, "1", input)
		assert.Error(t, err)
		mockService.AssertNotCalled(t, "UpdatePowerPlant", mock.Anything)
	})

	t.Run("fail due to service error", func(t *testing.T) {
		resolver, mockService := setupTests(t)
		input := model.UpdatePowerPlantInput{
			Name:      stringPointer("Updated Plant"),
			Latitude:  floatPointer(2.345),
			Longitude: floatPointer(6.789),
		}

		mockService.On("UpdatePowerPlant", ctx, mock.Anything).Return(nil, assert.AnError).Once()

		_, err := resolver.UpdatePowerPlant(ctx, "1", input)
		assert.Error(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("succeed updating power plant", func(t *testing.T) {
		resolver, mockService := setupTests(t)
		input := model.UpdatePowerPlantInput{
			Name:      stringPointer("Updated Plant"),
			Latitude:  floatPointer(2.345),
			Longitude: floatPointer(6.789),
		}
		updatedPlant := &model.PowerPlant{ID: "1", Name: "Updated Plant", Latitude: 2.345, Longitude: 6.789}
		mockService.On("UpdatePowerPlant", ctx, mock.Anything).Return(updatedPlant, nil).Once()

		result, err := resolver.UpdatePowerPlant(ctx, "1", input)
		assert.NoError(t, err)
		assert.Equal(t, updatedPlant, result)
		mockService.AssertExpectations(t)
	})
}

func stringPointer(s string) *string {
	return &s
}

func floatPointer(f float64) *float64 {
	return &f
}

func setupTests(t *testing.T) (MutationResolver, *mocks.PowerPlantService) {
	mockService := mocks.NewPowerPlantService(t)
	resolver := &Resolver{PowerPlantService: mockService}

	return resolver.Mutation(), mockService
}
