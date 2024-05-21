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
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Refresh(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &payloads.JWTRefresh{}

		err := utils.ValidateRequest(r, w, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("invalid request: %s", err))
			return
		}

		p := jwt.NewParser()

		refreshToken, err := p.Parse(payload.RefreshToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(securityConf.JWT.Key), nil
		})

		if err != nil {
			render.JSON(w, r, utils.Error(w, "invalid refresh token", 401))
			return
		}

		tokenClaims := refreshToken.Claims.(jwt.MapClaims)

		if tokenClaims["type"] != "refresh" {
			render.JSON(w, r, utils.Error(w, "provided token is not a refresh token", 403))
			return
		}

		token := &models.Token{}
		store.First(token, "value = ?", payload.RefreshToken)

		if len(token.Value) == 0 {
			render.JSON(w, r, utils.Error(w, "invalid refresh token", 401))
			return
		}

		store.Delete(token, "value = ?", token.Value)

		pair, err := utils.CreateNewTokenPair(securityConf, uint(tokenClaims["id"].(float64)))
		if err != nil {
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		signedAccessToken, err := utils.SignToken(securityConf.JWT.Key, pair.AccessToken)
		if err != nil {
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		signedRefreshToken, err := utils.SignToken(securityConf.JWT.Key, pair.RefreshToken)
		if err != nil {
			render.JSON(w, r, utils.InternalError(w))
			return
		}

		newRefreshToken := &models.Token{
			Value: signedRefreshToken,
		}

		store.Create(newRefreshToken)

		tokenPair := payloads.JWTTokenPair{
			AccessToken:  signedAccessToken,
			RefreshToken: signedRefreshToken,
		}

		render.JSON(w, r, utils.Success(w, tokenPair, 200))

	}
}
