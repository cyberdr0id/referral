package handler

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	emptyHeaderMessage = "authorization header cannot be empty"
	emptyTokenMessage  = "JWT token cannot be empty"

	authHeader   = "Authorization"
	bearerHeader = "Bearer "
)

// AuthorizationMiddleware checks if user is authorized.
func (s *Server) AuthorizationMiddleware(nextHandler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authHeader)
		if header == "" {
			sendResponse(rw, ErrorResponse{Message: emptyHeaderMessage}, http.StatusUnauthorized)
		}

		token := strings.Split(header, bearerHeader)[1]
		if token == "" {
			sendResponse(rw, ErrorResponse{Message: emptyTokenMessage}, http.StatusUnauthorized)
		}

		claims, err := s.Auth.ParseToken(token)
		if err != nil {
			sendResponse(rw, ErrorResponse{Message: fmt.Errorf("cannot parse JWT token: %w", err).Error()}, http.StatusInternalServerError)
		}

		currentUserID = claims.Subject

		nextHandler.ServeHTTP(rw, r)
	})
}
