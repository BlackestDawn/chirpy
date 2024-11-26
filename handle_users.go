package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BlackestDawn/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}

	data := new(params)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "Couldn't decode user data", err)
		return
	}

	cTime := time.Now()
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: cTime,
		UpdatedAt: cTime,
		Email:     data.Email,
	}

	retUser, err := c.dbQueries.CreateUser(context.Background(), newUser)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "could not create user", err)
		return
	}

	respondJSON(w, http.StatusCreated, retUser)
}
