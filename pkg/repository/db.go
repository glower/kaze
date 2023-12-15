package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/glower/kaze/graph/model"
	"github.com/jmoiron/sqlx"
)

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

func (r *powerPlantRepo) Create(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error) {
	query := `INSERT INTO power_plants (name, latitude, longitude) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, plant.Name, plant.Latitude, plant.Longitude)

	var id int
	if err := row.Scan(&id); err != nil {
		log.Printf("Failed to insert power plant: %v", err)
		return nil, err
	}

	plant.ID = fmt.Sprintf("%d", id)
	return plant, nil
}

func (r *powerPlantRepo) GetByID(ctx context.Context, id string) (*model.PowerPlant, error) {
	var plant model.PowerPlant
	query := `SELECT id, name, latitude, longitude FROM power_plants WHERE id = $1`
	if err := r.db.GetContext(ctx, &plant, query, id); err != nil {
		log.Printf("Failed to get power plant by ID: %v", err)
		return nil, err
	}

	return &plant, nil
}

func (r *powerPlantRepo) Update(ctx context.Context, plant *model.PowerPlant) (*model.PowerPlant, error) {
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
		log.Printf("Failed to update power plant: %v", err)
		return nil, fmt.Errorf("error updating power plant: %w", err)
	}

	return plant, nil
}

func (r *powerPlantRepo) List(ctx context.Context, offset, limit int) ([]model.PowerPlant, int, error) {
	var powerPlants []model.PowerPlant
	var total int

	countQuery := "SELECT COUNT(*) FROM power_plants"
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		log.Printf("Error getting total number of power plants: %v", err)
		return nil, 0, fmt.Errorf("error getting total number of power plants: %w", err)
	}

	listQuery := `SELECT id, name, latitude, longitude FROM power_plants ORDER BY id LIMIT $1 OFFSET $2`
	if err := r.db.SelectContext(ctx, &powerPlants, listQuery, limit, offset); err != nil {
		log.Printf("Error querying power plants: %v", err)
		return nil, 0, fmt.Errorf("error querying power plants: %w", err)
	}

	return powerPlants, total, nil
}
