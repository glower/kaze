package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/glower/kaze/pkg/config"
	"github.com/glower/kaze/pkg/database"
	"github.com/glower/kaze/pkg/handler"
	"github.com/glower/kaze/pkg/repository"
	"github.com/glower/kaze/pkg/service"
)

func initLog() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

func main() {
	initLog()

	conf := config.NewConfig()
	db, err := database.NewConnection(conf)
	if err != nil {
		slog.Error("can't create new database connection", "error", err)
		return
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		slog.Error("can't init migration", "error", err)
		return
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to get current working directory:", "error", err)
		return
	}

	// Construct the path to the migration files
	migrationPath := filepath.Join(cwd, "migrations") // Adjust the path as needed

	dbMigrationPath := fmt.Sprintf("file://%s", migrationPath)
	slog.Debug("migration path is", "path", dbMigrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		dbMigrationPath,
		"kaze", driver)
	if err != nil {
		slog.Error("can't initiate the database migration", "error", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("can't run migration", "error", err)
		return
	}

	slog.Info("Migration executed successfully")

	powerPlantRepo := repository.NewPowerPlantRepository(db)
	openMeteoRepo := repository.NewOpenMeteoRepository(conf.OpenMeteoAPIKey)
	powerPlantService := service.NewPowerPlantService(powerPlantRepo, openMeteoRepo)

	server := handler.NewServer(powerPlantService)
	mux := server.SetupRoutes()

	// Start the server
	slog.Info("Starting GraphQL server on http://localhost:8080/")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
