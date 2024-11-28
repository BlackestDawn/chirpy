package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/BlackestDawn/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	tokenSecret    string
}

func NewApiConfig() (*apiConfig, error) {
	cfg := new(apiConfig)
	cfg.fileserverHits.Store(0)
	cfg.platform = os.Getenv("PLATFORM")
	cfg.tokenSecret = os.Getenv("SECRET")

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
