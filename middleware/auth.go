package middleware

import (
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"rsvp/models"
	"rsvp/utils"
	"strings"
	"time"
)

var JwtAuthentication = func() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for authentication
			response := make(map[string]interface{})
			tokenHeader := r.Header.Get("Authorization")

			// If token is missing, then return error code 403 Unauthorized
			if tokenHeader == "" {

				response = utils.Message(false, http.StatusUnauthorized, "Missing auth token")
				utils.Respond(w, response)
				return
			}

			// Check if the token format is correct, ie. Bearer {token}
			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				response = utils.Message(false, http.StatusUnauthorized, "Invalid auth token format.")
				utils.Respond(w, response)
				return
			}

			tokenPart := splitted[1] // Grab the second part
			tk := &models.Token{}

			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("token_password")), nil
			})

			if err != nil {
				response = utils.Message(false, http.StatusUnauthorized, "Failed to check token: "+err.Error())
				utils.Respond(w, response)
				return
			}

			if !token.Valid {
				response = utils.Message(false, http.StatusUnauthorized, "Token is not valid.")
				utils.Respond(w, response)
				return
			}

			if time.Now().After(tk.Expiry) {
				response = utils.Message(false, http.StatusUnauthorized, "Token has expired. Please login again.")
				utils.Respond(w, response)
				return
			}

			// Set the user ID in the context
			ctx := context.WithValue(r.Context(), "user", tk.UserId)
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
		})
	}
}
