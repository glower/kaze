package handler

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/glower/kaze/graph"
	"github.com/glower/kaze/pkg/service"
)

// Server struct represents the GraphQL server
type Server struct {
	powerPlantService service.PowerPlantService
}

// NewServer creates a new GraphQL server
func NewServer(powerPlantService service.PowerPlantService) *Server {
	return &Server{
		powerPlantService: powerPlantService,
	}
}

func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Create an instance of Resolver and inject the service
	resolver := &graph.Resolver{
		PowerPlantService: s.powerPlantService,
	}

	// Setup GraphQL handler
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	mux.Handle("/graphql", srv)

	// Setup the GraphQL playground handler
	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	return mux
}
