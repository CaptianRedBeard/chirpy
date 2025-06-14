package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrive chirps", err)
		return
	}

	var responseSlice []chirpResponse

	for _, chirp := range chirps {

		response := chirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		responseSlice = append(responseSlice, response)

	}

	respondWithJSON(w, http.StatusOK, responseSlice)

}
