package main

import (
	"chirpy/internal/auth"
	"net/http"
	"time"
)

type ResponseToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	token, err := cfg.db.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token in db", err)
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token has expired", err)
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user", err)
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.jwt_secret, auth.AccessTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect access token validation", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ResponseToken{
		Token: newAccessToken,
	})

}
