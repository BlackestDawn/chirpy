package main

import (
	"encoding/json"
	"net/http"

	"github.com/BlackestDawn/chirpy/internal/auth"
	"github.com/google/uuid"
)

// Handler: POST /api/polka/webhooks
func (c *apiConfig) handlerPolkaUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondSimple(w, http.StatusUnauthorized, "", "plain")
		return
	}
	if apiKey != c.apiConfig {
		respondSimple(w, http.StatusUnauthorized, "", "plain")
		return
	}

	data := new(params)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(data)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "Couldn't decode user data", err)
		return
	}

	if data.Event != "user.upgraded" {
		respondSimple(w, http.StatusNoContent, "", "plain")
		return
	}

	err = c.dbQueries.UpgradeUserToRed(r.Context(), data.Data.UserID)
	if err != nil {
		respondSimple(w, http.StatusNotFound, "", "plain")
		return
	}

	respondSimple(w, http.StatusNoContent, "", "plain")
}
