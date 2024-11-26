package main

import (
	"fmt"
	"strings"
)

func validateChirp(msg string) (string, error) {
	if len(msg) > maxChirpLength {
		return "", fmt.Errorf("chirp is too long")
	}

	bannedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedBody := censorMessage(msg, bannedWords)
	return cleanedBody, nil
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
