package graph

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"

	"github.com/glower/kaze/graph/model"
)

// CreatePowerPlant is the resolver for the createPowerPlant field.
// It creates a new power plant with the given input data.
func (r *mutationResolver) CreatePowerPlant(ctx context.Context, input model.NewPowerPlantInput) (*model.PowerPlant, error) {
	slog.Debug("Creating a new power plant", "payload", input)

	// Validate the input data
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		slog.Error("Input validation failed", "error", err, "payload", input)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Call the service layer to create the power plant
	powerPlant, err := r.PowerPlantService.CreatePowerPlant(ctx, &model.PowerPlant{
		Name:      input.Name,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	})
	if err != nil {
		slog.Error("Failed to create power plant", "error", err, "payload", input)
		return nil, fmt.Errorf("failed to create power plant: %w", err)
	}

	return powerPlant, nil
}

// UpdatePowerPlant is the resolver for the updatePowerPlant field.
// It updates an existing power plant identified by the given ID with new input data.
func (r *mutationResolver) UpdatePowerPlant(ctx context.Context, id string, input model.UpdatePowerPlantInput) (*model.PowerPlant, error) {
	slog.Debug("Updating power plant", "id", id, "payload", input)

	// Validate the input data
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		slog.Error("Input validation failed", "error", err, "id", id, "payload", input)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Prepare the update payload
	updatePayload := &model.PowerPlant{ID: id}
	if input.Name != nil {
		updatePayload.Name = *input.Name
	}
	if input.Latitude != nil {
		updatePayload.Latitude = *input.Latitude
	}
	if input.Longitude != nil {
		updatePayload.Longitude = *input.Longitude
	}

	// Call the service layer to update the power plant
	updatedPowerPlant, err := r.PowerPlantService.UpdatePowerPlant(ctx, updatePayload)
	if err != nil {
		slog.Error("Failed to update power plant", "error", err, "id", id, "payload", updatePayload)
		return nil, fmt.Errorf("failed to update power plant: %w", err)
	}

	return updatedPowerPlant, nil
}
