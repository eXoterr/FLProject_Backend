package category

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

		payload := &payloads.Category{}

		err := utils.ValidateRequest(r, w, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("invalid request: %s", err))
			return
		}

		existingCategory := &models.Category{}
		store.First(existingCategory, "name = ?", payload.Name)

		if existingCategory.ID != 0 {
			render.JSON(w, r, utils.Error(w, "category with this name is already exists", 400))
			return
		}

		category := &models.Category{
			Name: payload.Name,
		}
		result := store.Create(category)

		if result.Error != nil || category.ID == 0 {
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		render.JSON(w, r, utils.Success(w, "category created", 201))
	}
}
