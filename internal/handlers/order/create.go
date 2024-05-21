package order

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/handlers/payloads"
	"github.com/eXoterr/FLProject/internal/storage/models"
	"github.com/eXoterr/FLProject/internal/utils"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

func Create(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		payload := &payloads.Order{}

		err := utils.ValidateRequest(r, w, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("invalid request: %s", err))
			return
		}

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

		client := &models.Client{}
		store.First(client, "user_id = ?", userId)

		if client.ID == 0 {
			render.JSON(w, r, utils.Error(w, "user is not registered as client", 400))
			return
		}

		category := &models.Category{}
		result := store.First(category, "id = ?", payload.CategoryID)

		if result.Error != nil || category.ID == 0 {
			render.JSON(w, r, utils.Error(w, "specified category does not exists", 400))
			return
		}

		orderTags := make([]models.Tag, len(payload.Tags))
		for idx, tag := range payload.Tags {
			orderTags[idx] = models.Tag{
				ID: tag,
			}
		}

		if payload.Deadline.Before(time.Now()) {
			render.JSON(w, r, utils.Error(w, "deadline cannot be in past", 400))
			return
		}

		order := &models.Order{
			ClientID:    client.ID,
			Title:       payload.Title,
			Description: payload.Description,
			CategoryID:  payload.CategoryID,
			State:       0,
			BudgetMin:   payload.BudgetMin,
			BudgetMax:   payload.BudgetMax,
			Deadline:    payload.Deadline,
			Tags:        orderTags,
		}

		result = store.Create(order)
		if result.Error != nil {
			logger.Error(result.Error.Error())
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		render.JSON(w, r, utils.Success(w, "order created", 201))
	}
}
