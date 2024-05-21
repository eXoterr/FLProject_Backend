package routing

import (
	"log/slog"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/handlers/tags"
	"github.com/eXoterr/FLProject/internal/routing/middlewares"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func Tags(parentRouter chi.Router, store *gorm.DB, config *config.Config, logger *slog.Logger) {
	router := chi.NewRouter()

	router.Use(middlewares.SetupAccessControl(config.Security, logger))

	router.Post("/create", tags.Create(store, logger, config.Security))
	router.Get("/list", tags.GetList(store, logger, config.Security))
	router.Get("/search", tags.Search(store, logger, config.Security))

	parentRouter.Mount("/tags", router)
}
