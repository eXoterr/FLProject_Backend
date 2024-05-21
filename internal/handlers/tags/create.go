package tags

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/handlers/payloads"
	"github.com/eXoterr/FLProject/internal/storage/models"
	"github.com/eXoterr/FLProject/internal/utils"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

func Create(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		payload := &payloads.Tag{}

		err := utils.ValidateRequest(r, w, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("invalid request: %s", err))
			return
		}

		existingTag := &models.Tag{}
		store.First(existingTag, "name = ?", payload.Name)

		if existingTag.ID != 0 {
			render.JSON(w, r, utils.Error(w, "tag with this name is already exists", 400))
			return
		}

		tag := &models.Tag{
			Name: payload.Name,
		}
		result := store.Create(tag)

		if result.Error != nil || tag.ID == 0 {
			logger.Error(result.Error.Error())
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		render.JSON(w, r, utils.Success(w, "tag created", 201))
	}
}
