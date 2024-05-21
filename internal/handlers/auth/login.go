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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TokenClaims struct {
	Role string `json:"role"`
	ID   int    `json:"id"`
	jwt.RegisteredClaims
}

func Login(store *gorm.DB, logger *slog.Logger, securityConf config.Security) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &payloads.User{}

		err := utils.ValidateRequest(r, w, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("invalid request: %s", err))
			return
		}

		user := &models.User{}

		store.First(user, "email = ?", payload.Email)

		if user.ID == 0 {
			logger.Info("login with incorrect email")

			render.JSON(w, r, utils.Error(w, "incorrect email or password", 403))

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
		if err != nil {
			logger.Info("login with incorrect password")

			render.JSON(w, r, utils.Error(w, "incorrect email or password", 403))

			return
		}

		pair, err := utils.CreateNewTokenPair(securityConf, user.ID)
		if err != nil {
			logger.Error("unable to create token pair")

			render.JSON(w, r, utils.InternalError(w))

			return
		}

		// accessClaims := pair.AccessToken.Claims.(*utils.TokenClaims)
		signedAccessToken, err := pair.AccessToken.SignedString([]byte(securityConf.JWT.Key))
		if err != nil {
			logger.Error("unable to sign access token")

			render.JSON(w, r, utils.InternalError(w))

			return
		}

		// refreshClaims := pair.RefreshToken.Claims.(*utils.TokenClaims)
		signedRefreshToken, err := pair.RefreshToken.SignedString([]byte(securityConf.JWT.Key))
		if err != nil {
			logger.Error("unable to sign refresh token")

			render.JSON(w, r, utils.InternalError(w))

			return
		}

		// accessCookie := &http.Cookie{
		// 	Name:     "access",
		// 	Value:    signedAccessToken,
		// 	Expires:  accessClaims.ExpiresAt.Time,
		// 	SameSite: http.SameSiteNoneMode,
		// 	Secure:   true,
		// 	Path:     "",
		// 	HttpOnly: true,
		// }

		// refreshCookie := &http.Cookie{
		// 	Name:     "refresh",
		// 	Value:    signedRefreshToken,
		// 	Expires:  refreshClaims.ExpiresAt.Time,
		// 	SameSite: http.SameSiteNoneMode,
		// 	Secure:   true,
		// 	Path:     "",
		// 	HttpOnly: true,
		// }

		// http.SetCookie(w, accessCookie)
		// http.SetCookie(w, refreshCookie)

		refresh := &models.Token{
			Value: signedRefreshToken,
		}

		store.Create(refresh)

		payloadPair := &payloads.JWTTokenPair{
			AccessToken:  signedAccessToken,
			RefreshToken: signedRefreshToken,
		}

		render.JSON(w, r, utils.Success(w, payloadPair, 200))
	}
}
