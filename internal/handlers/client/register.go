package client

import (
	"log/slog"
	"net/http"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/storage/models"
	"github.com/eXoterr/FLProject/internal/utils"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

func Register(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := utils.GetTokenFromBody(r)
		if err != nil {
			render.JSON(w, r, utils.InternalError(w))
			return
		}
		userId, err := utils.GetUserIDFromJWT(token, securityConf.JWT.Key)
		if err != nil {
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		client := &models.Client{
			UserID: userId,
		}

		store.First(client, "user_id = ?", userId)

		if client.ID != 0 {
			render.JSON(w, r, utils.Error(w, "this user is already registered as client", 403))

			return
		}

		store.Create(client)

		render.JSON(w, r, utils.Success(w, "client registered", 201))
	}
}
