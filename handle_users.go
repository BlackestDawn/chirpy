package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BlackestDawn/chirpy/internal/auth"
	"github.com/BlackestDawn/chirpy/internal/database"
)

// Handler: POST /api/users
func (c *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type retVal struct {
		cleanUser
	}

	data := new(params)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "Couldn't decode user data", err)
		return
	}

	hashedPass, err := auth.HashPassword(data.Password)
	if err != nil {
		respondJSONError(w, http.StatusBadRequest, "error hashing password", err)
	}
	newUser := database.CreateUserParams{
		Email:    data.Email,
		Password: hashedPass,
	}

	retUser, err := c.dbQueries.CreateUser(context.Background(), newUser)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "could not create user", err)
		return
	}

	respondJSON(w, http.StatusCreated, retVal{
		cleanUser: cleanUser{
			ID:        retUser.ID,
			CreatedAt: retUser.CreatedAt,
			UpdatedAt: retUser.UpdatedAt,
			Email:     retUser.Email,
		},
	})
}

// Handler: POST /api/login
func (c *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type RetVal struct {
		cleanUser
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	data := new(params)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "Couldn't decode user data", err)
		return
	}

	user, err := c.dbQueries.GetUserByEmail(context.Background(), data.Email)
	if err != nil {
		respondJSONError(w, http.StatusUnauthorized, "incorrect password or login", nil)
		return
	}

	err = auth.CheckPasswordHash(data.Password, user.Password)
	if err != nil {
		respondJSONError(w, http.StatusUnauthorized, "incorrect password or login", nil)
		return
	}

	token, err := auth.MakeJWT(user.ID, c.tokenSecret, time.Hour)
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, "error creating access token", err)
		return
	}

	refreshToken, err := c.dbQueries.GetValidRefreshTokenForUser(context.Background(), user.ID)
	if err != nil {

		refreshToken, err = auth.MakeRefreshToken()
		if err != nil {
			respondJSONError(w, http.StatusInternalServerError, "error generating refresh token", err)
			return
		}

		_, err = c.dbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
			Token:  refreshToken,
			UserID: user.ID,
		})
		if err != nil {
			respondJSONError(w, http.StatusInternalServerError, "error storing refresh token", err)
			return
		}
	}

	respondJSON(w, http.StatusOK, RetVal{
		cleanUser: cleanUser{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
