package middlewares

import (
	"net/http"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/go-chi/cors"
)

func SetupCORS(cfg config.CORS) func(http.Handler) http.Handler {
	handler := cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Origins,
		AllowedMethods:   cfg.Methods,
		AllowedHeaders:   cfg.Headers,
		AllowCredentials: cfg.Credentials,
	})

	return handler
}
