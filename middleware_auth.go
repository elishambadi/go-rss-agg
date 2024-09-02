package main

import (
	"fmt"
	"net/http"

	"github.com/elishambadi/go-rss-agg/internal/auth"
	"github.com/elishambadi/go-rss-agg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("error creating user: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("error getting user from DB: %v", err))
			return
		}

		handler(w, r, user)
	}
}
