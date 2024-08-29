package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, httpCode int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application-json")
	w.WriteHeader(httpCode) //code
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, httpCode int, msg string) {
	if httpCode > 499 {
		log.Println("Responding with 500 err: ", msg)
	}

	type errorResponse struct {
		Error string `json:"errors"` // Reflect it into a key of a JSON object with such structure
	}

	respondWithJSON(w, httpCode, errorResponse{
		Error: msg,
	})
}
