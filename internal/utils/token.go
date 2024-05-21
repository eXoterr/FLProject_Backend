package utils

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type JWTPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}

func GetTokenFromBody(r *http.Request) (string, error) {
	headers := r.Header
	tokenHeader := headers[http.CanonicalHeaderKey("authorization")]
	if len(tokenHeader) == 0 {
		return "", errors.New("no token privided")
	}

	if len(tokenHeader) != 1 {
		return "", errors.New("auth header invalid length")
	}

	tokenHeader = strings.Split(tokenHeader[0], " ")

	if len(tokenHeader) != 2 {
		return "", errors.New("auth header data invalid length")
	}

	if strings.ToLower(tokenHeader[0]) != "bearer" {
		return "", errors.New("auth header invalid auth type")
	}

	return tokenHeader[1], nil

}

func GetUserIDFromJWT(token string, secret string) (uint, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}
	return jwtToken.Claims.(*TokenClaims).ID, nil
}

func SignToken(secret string, token *jwt.Token) (string, error) {
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func CreateNewTokenPair(conf config.Security, userId uint) (*JWTPair, error) {
	now := time.Now()
	accessClaims := &TokenClaims{
		ID:   userId,
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(conf.JWT.AccessTokenTTL) * time.Minute)), // offload to cfg
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		accessClaims,
	)

	refreshClaims := &TokenClaims{
		ID:   userId,
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(conf.JWT.RefreshTokenTTL) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		refreshClaims,
	)

	pair := &JWTPair{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
	return pair, nil
}
