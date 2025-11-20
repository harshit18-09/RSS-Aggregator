package main

import (
	"context"
	"net/http"

	"github.com/harshit18-09/RSS-Aggregator/internal/auth"
)

type contextKey string

const userContextKey contextKey = "user"

// MiddlewareAuth authenticates the request using an API key. On success it
// attaches the db.User to the request context under the userContextKey
// and calls the next handler. It returns an http.Handler
func (apiCfg *apiConfig) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, "Unauthorized: "+err.Error())
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, "Failed to get user")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
