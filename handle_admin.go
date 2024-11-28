package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Handler: GET /api/healthz
func handlerHealth(w http.ResponseWriter, r *http.Request) {
	respondSimple(w, http.StatusOK, http.StatusText(http.StatusOK), "plain")
}

// Handler: GET /admin/metrics
func (c *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	hitStr := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", c.fileserverHits.Load())
	respondSimple(w, http.StatusOK, hitStr, "html")
}

// Handler: POST /admin/reset
func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(c.platform) != "dev" {
		respondSimple(w, http.StatusForbidden, "Reset only allowed in DEV environment", "plain")
		return
	}

	err := c.dbQueries.ResetUsers(r.Context())
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "could not reset database", err)
		return
	}

	c.fileserverHits.Store(0)

	respondSimple(w, http.StatusOK, "Hit counter and database reset", "plain")
}
