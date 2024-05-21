package routing

import (
	"log/slog"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/handlers/client"
	"github.com/eXoterr/FLProject/internal/routing/middlewares"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func Client(parentRouter chi.Router, store *gorm.DB, config *config.Config, logger *slog.Logger) {
	router := chi.NewRouter()

	router.Use(middlewares.SetupAccessControl(config.Security, logger))

	router.Post("/register", client.Register(store, logger, config.Security))

	parentRouter.Mount("/client", router)
}
