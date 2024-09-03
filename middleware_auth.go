package main

import (
	"fmt"
	"net/http"

	"github.com/elishambadi/go-rss-agg/internal/auth"
	"github.com/elishambadi/go-rss-agg/internal/database"
)

// Overloads the http.Handler func to add database.User
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// This function gets the user before calling the actual function
// Steps
// 1. You write your function as authedHandler
// 2. You call it inside middlewareAuth which converts it to expected http.HnadlerFunc
// 3. Enjoy
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("error getting APIKey from header: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("error getting user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
