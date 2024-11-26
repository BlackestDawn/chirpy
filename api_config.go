package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/BlackestDawn/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func NewApiConfig() (*apiConfig, error) {
	cfg := new(apiConfig)
	cfg.fileserverHits.Store(0)
	cfg.platform = os.Getenv("PLATFORM")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to DB: %w", err)
	}

	cfg.dbQueries = database.New(db)

	return cfg, nil
}
func (c *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitStr := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", c.fileserverHits.Load())
	w.Write([]byte(hitStr))
}

func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(c.platform) != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := c.dbQueries.ResetUsers(context.Background())
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "could not reset users", err)
		return
	}

	c.fileserverHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
