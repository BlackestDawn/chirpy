package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/BlackestDawn/chirpy/internal/auth"
)

// Handler: POST /api/refresh
func (c *apiConfig) handlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	type retVal struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "error fetching refresh token", err)
		return
	}

	userId, err := c.dbQueries.GetUserFromRefreshToken(context.Background(), token)
	if err != nil {
		log.Println("no user returned for refresh")
		respondJSONError(w, http.StatusUnauthorized, "invalid or non-existant refresh token", err)
		return
	}

	accountToken, err := auth.MakeJWT(userId, c.tokenSecret, time.Hour)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "error refreshing access token", err)
		return
	}

	respondJSON(w, http.StatusOK, retVal{Token: accountToken})
}

// Handler: POST /api/revoke
func (c *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "error fetching refresh token", err)
		return
	}

	_, err = c.dbQueries.RevokeRefreshToken(context.Background(), token)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "error revoking refresh token", err)
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
