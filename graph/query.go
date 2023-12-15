package graph

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/glower/kaze/graph/model"
)

// PowerPlant is the resolver for the powerPlant field.
// It retrieves a power plant by its ID, including elevation and weather forecasts if requested.
func (r *queryResolver) PowerPlant(ctx context.Context, id string) (*model.PowerPlant, error) {
	slog.Debug("Retrieving power plant", "id", id)

	withElevation := IsFieldRequested(ctx, "elevation")
	withWeatherForecasts := IsFieldRequested(ctx, "weatherForecasts")

	powerPlant, err := r.PowerPlantService.GetPowerPlant(ctx, id, withElevation, withWeatherForecasts)
	if err != nil {
		slog.Error("Failed to retrieve power plant", "error", err, "id", id)
		return nil, fmt.Errorf("error retrieving power plant by ID: %w", err)
	}

	return powerPlant, nil
}

// ListPowerPlants is the resolver for the listPowerPlants field.
// It retrieves a list of power plants, supporting pagination and optional inclusion of elevation and weather forecasts.
func (r *queryResolver) ListPowerPlants(ctx context.Context, page *int, pageSize *int) (*model.PowerPlantList, error) {
	slog.Debug("Retrieving a list of power plants", "page", *page, "pageSize", *pageSize)

	withElevation := IsFieldRequested(ctx, "powerPlants.elevation")
	withWeatherForecasts := IsFieldRequested(ctx, "powerPlants.weatherForecasts")

	return r.PowerPlantService.ListPowerPlants(ctx, toIntWithDefault(page, 1), toIntWithDefault(pageSize, 10), withElevation, withWeatherForecasts)
}
