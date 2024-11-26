package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BlackestDawn/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	data := new(params)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "could not decode message data", err)
	}

	data.Body, err = validateChirp(data.Body)
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "chirp validation failure", err)
	}

	cTime := time.Now()
	newPost := database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: cTime,
		UpdatedAt: cTime,
		Body:      data.Body,
		UserID:    data.UserID,
	}
	retVal, err := c.dbQueries.CreatePost(context.Background(), newPost)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "failed storing chirp", err)
	}

	respondJSON(w, http.StatusCreated, retVal)
}

func (c *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	retVal, err := c.dbQueries.ListPosts(context.Background())
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "error fetching posts", err)
	}

	respondJSON(w, http.StatusOK, retVal)
}

func (c *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "error parsing UUID", err)
	}

	post, err := c.dbQueries.GetPostByID(context.Background(), postID)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "error fetching post", err)
	}

	respondJSON(w, http.StatusOK, post)
}
