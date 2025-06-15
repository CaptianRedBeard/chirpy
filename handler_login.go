package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserLogin struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := paramaters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// fmt.Printf("JWT Secret length: %d\n", len(cfg.jwt_secret))

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwt_secret, auth.AccessTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect access token validation", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating refresh token", err)
		return
	}

	refreshExpiration := time.Now().Add(auth.RefreshTokenExpiration)
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: refreshExpiration,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error calculating refresh token expiration", err)
		return
	}

	respondWithJSON(w, http.StatusOK, UserLogin{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})

}
