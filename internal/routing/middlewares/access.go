package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/utils"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

func SetupAccessControl(conf config.Security, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			p := jwt.NewParser()

			tokenStr, err := utils.GetTokenFromBody(r)
			if err != nil {
				logger.Error(err.Error())
				render.JSON(w, r, utils.Error(w, "invalid auth token", 401))
				return
			}

			_, err = p.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(conf.JWT.Key), nil
			})

			if err != nil {
				logger.Error(err.Error())
				render.JSON(w, r, utils.Error(w, "invalid auth token", 401))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
