package routing

import (
	"log/slog"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/handlers/auth"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func authRoutes(parentRouter chi.Router, store *gorm.DB, logger *slog.Logger, config *config.Config) {
	router := chi.NewRouter()
	router.Post("/register", auth.Register(store, logger, config.Security))
	router.Post("/login", auth.Login(store, logger, config.Security))
	router.Post("/refresh", auth.Refresh(store, logger, config.Security))
	parentRouter.Mount("/auth", router)
}
