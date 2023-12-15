package repository

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/glower/kaze/graph/model"
	"github.com/jmoiron/sqlx"
)

//go:generate go run github.com/vektra/mockery/v2@v2 --name=PowerPlantRepository --filename=power_plant_repository.go --output=../../mocks/
type PowerPlantRepository interface {
	Create(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error)
	GetByID(ctx context.Context, id string) (*model.PowerPlant, error)
	Update(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error)
	List(ctx context.Context, offset, limit int) ([]model.PowerPlant, int, error)
}

type powerPlantRepo struct {
	db *sqlx.DB
}

func NewPowerPlantRepository(db *sqlx.DB) PowerPlantRepository {
	return &powerPlantRepo{
		db: db,
	}
}

// Create inserts a new power plant into the database.
func (r *powerPlantRepo) Create(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error) {
	slog.Debug("Inserting new power plant", "name", plant.Name)

	query := `INSERT INTO power_plants (name, latitude, longitude) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, plant.Name, plant.Latitude, plant.Longitude)

	var id int
	if err := row.Scan(&id); err != nil {
		slog.Error("Failed to insert power plant", "error", err)
		return nil, err
	}

	plant.ID = fmt.Sprintf("%d", id)
	return plant, nil
}

// GetByID retrieves a power plant by its ID.
func (r *powerPlantRepo) GetByID(ctx context.Context, id string) (*model.PowerPlant, error) {
	slog.Debug("Retrieving power plant", "id", id)

	var plant model.PowerPlant
	query := `SELECT id, name, latitude, longitude FROM power_plants WHERE id = $1`
	if err := r.db.GetContext(ctx, &plant, query, id); err != nil {
		slog.Error("Failed to get power plant by ID", "error", err)
		return nil, err
	}

	return &plant, nil
}

// Update modifies an existing power plant.
func (r *powerPlantRepo) Update(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error) {
	slog.Debug("Updating power plant", "id", plant.ID)

	setParts := []string{}
	params := map[string]interface{}{}

	if plant.Name != "" {
		setParts = append(setParts, "name = :name")
		params["name"] = plant.Name
	}
	if plant.Latitude != 0 {
		setParts = append(setParts, "latitude = :latitude")
		params["latitude"] = plant.Latitude
	}
	if plant.Longitude != 0 {
		setParts = append(setParts, "longitude = :longitude")
		params["longitude"] = plant.Longitude
	}

	if len(setParts) == 0 {
		return plant, nil // No update needed
	}

	query := fmt.Sprintf("UPDATE power_plants SET %s WHERE id = :id", strings.Join(setParts, ", "))
	params["id"] = plant.ID

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		slog.Error("Failed to update power plant", "error", err)
		return nil, fmt.Errorf("error updating power plant: %w", err)
	}

	return plant, nil
}

// List fetches a list of power plants with pagination.
func (r *powerPlantRepo) List(ctx context.Context, offset, limit int) ([]model.PowerPlant, int, error) {
	slog.Debug("Listing power plants", "offset", offset, "limit", limit)

	var powerPlants []model.PowerPlant
	var total int

	countQuery := "SELECT COUNT(*) FROM power_plants"
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		slog.Error("Error getting total number of power plants", "error", err)
		return nil, 0, fmt.Errorf("error getting total number of power plants: %w", err)
	}

	slog.Debug("total number of all power plants", "total", countQuery)

	listQuery := `SELECT id, name, latitude, longitude FROM power_plants ORDER BY id LIMIT $1 OFFSET $2`
	if err := r.db.SelectContext(ctx, &powerPlants, listQuery, limit, offset); err != nil {
		slog.Error("Error querying power plants", "error", err)
		return nil, 0, fmt.Errorf("error querying power plants: %w", err)
	}

	return powerPlants, total, nil
}
