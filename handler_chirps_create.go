package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	params := chirpParamaters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	_, err = cfg.db.GetUserByID(r.Context(), params.User_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "User not in database", err)
		return
	}

	err = ValidateChirp(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't validate chirp", err)
		return
	}

	chirpParams := chirpParamaters{
		Body:    params.Body,
		User_id: params.User_id, //uuid.NullUUID{UUID: params.User_id, Valid: true},
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
