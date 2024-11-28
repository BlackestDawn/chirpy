package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/BlackestDawn/chirpy/internal/auth"
	"github.com/BlackestDawn/chirpy/internal/database"
	"github.com/google/uuid"
)

// Handler: POST /api/chirps
func (c *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "error getting bearer token", err)
		return
	}
	userID, err := auth.ValidateJWT(authToken, c.tokenSecret)
	if err != nil {
		respondJSONError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	data := new(params)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(data)
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "could not decode message data", err)
		return
	}

	data.Body, err = validateChirp(data.Body)
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "chirp validation failure", err)
		return
	}

	newPost := database.CreatePostParams{
		Body:   data.Body,
		UserID: userID,
	}
	retVal, err := c.dbQueries.CreatePost(context.Background(), newPost)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "failed storing chirp", err)
		return
	}

	respondJSON(w, http.StatusCreated, retVal)
}

// Handler: GET /api/chirps
func (c *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	retVal, err := c.dbQueries.ListPosts(context.Background())
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "error fetching posts", err)
		return
	}

	respondJSON(w, http.StatusOK, retVal)
}

// Handler: GET /api/chirps/{chirpID}
func (c *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "error parsing UUID", err)
		return
	}

	post, err := c.dbQueries.GetPostByID(context.Background(), postID)
	if err != nil {
		respondJSONError(w, http.StatusNotFound, "error fetching post", err)
		return
	}

	respondJSON(w, http.StatusOK, post)
}

// Handler: DELETE /api/chirps/{chirpID}
func (c *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "error parsing UUID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJSONError(w, http.StatusUnauthorized, "error verifying bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, c.tokenSecret)
	if err != nil {
		respondJSONError(w, http.StatusUnauthorized, "error verifying access token", err)
		return
	}

	post, err := c.dbQueries.GetPostByID(context.Background(), postID)
	if err != nil {
		respondJSONError(w, http.StatusNotFound, "error fetching post data", err)
		return
	}
	if post.UserID != userID {
		respondJSONError(w, http.StatusForbidden, "ID mismatch", nil)
		return
	}

	err = c.dbQueries.DeletePostByID(context.Background(), postID)
	respondSimple(w, http.StatusNoContent, "", "plain")
}
