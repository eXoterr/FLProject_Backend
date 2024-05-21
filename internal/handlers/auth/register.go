package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/handlers/payloads"
	"github.com/eXoterr/FLProject/internal/storage/models"
	"github.com/eXoterr/FLProject/internal/utils"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &payloads.User{}

		err := utils.ValidateRequest(r, w, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("invalid request: %s", err))
			return
		}

		user := &models.User{}
		store.First(user, "email = ?", payload.Email)

		if user.ID != 0 {
			render.JSON(w, r, utils.Error(w, "this email is already registered", 403))

			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), securityConf.BCryptCost)
		if err != nil {
			logger.Error("unable to hash password", err)

			render.JSON(w, r, utils.InternalError(w))

			return
		}

		user.Email = payload.Email
		user.Password = string(hashedPassword)

		store.Create(user)

		render.JSON(w, r, utils.Success(w, "user registered", 201))
	}
}
