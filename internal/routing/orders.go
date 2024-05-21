package routing

import (
	"log/slog"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/handlers/order"
	"github.com/eXoterr/FLProject/internal/routing/middlewares"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func Order(parentRouter chi.Router, store *gorm.DB, config *config.Config, logger *slog.Logger) {
	router := chi.NewRouter()

	router.Use(middlewares.SetupAccessControl(config.Security, logger))

	router.Post("/create", order.Create(store, logger, config.Security))
	router.Get("/list", order.Search(store, logger, config.Security))

	// Additional modules
	Categories(router, store, config, logger)
	Tags(router, store, config, logger)

	parentRouter.Mount("/order", router)
}
