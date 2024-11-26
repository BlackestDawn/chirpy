package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}
	type retVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	data := new(params)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err != nil {
		returnJSONError(w, http.StatusInternalServerError, "Couldn't decode data", err)
		return
	}

	if len(data.Body) > maxChirpLength {
		returnJSONError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	bannedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	returnJSON(w, http.StatusOK, retVal{CleanedBody: censorMessage(data.Body, bannedWords)})
}

func censorMessage(str string, censorList map[string]struct{}) string {
	words := strings.Split(str, " ")

	for i, word := range words {
		if _, ok := censorList[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
