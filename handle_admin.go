package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Handler: GET /api/healthz
func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Handler: GET /admin/metrics
func (c *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitStr := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", c.fileserverHits.Load())
	w.Write([]byte(hitStr))
}

// Handler: POST /admin/reset
func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(c.platform) != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := c.dbQueries.ResetUsers(context.Background())
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "could not reset database", err)
		return
	}

	c.fileserverHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
