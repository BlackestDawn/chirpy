package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/BlackestDawn/chirpy/internal/database"
)

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding JSON: %s\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(jsonData)
}

func respondJSONError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	type errVal struct {
		Error string `json:"error"`
	}

	respondJSON(w, code, errVal{Error: msg})
}
