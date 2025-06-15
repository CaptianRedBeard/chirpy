package main

import (
	"chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevokeAccessToken(w http.ResponseWriter, r *http.Request) {

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt update token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
