package routing

import (
	"log/slog"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupHandlers(router *chi.Mux, store *gorm.DB, logger *slog.Logger, config *config.Config) {
	authRoutes(router, store, logger, config)
	Client(router, store, config, logger)
	Order(router, store, config, logger)
}
