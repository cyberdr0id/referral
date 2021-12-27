package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cyberdr0id/referral/internal/context"
)

const (
	authHeaderMessage     = "authorization header required"
	emptyTokenMessage     = "JWT token cannot be empty"
	invalidSecurityScheme = "invalid security scheme"
	invalidAuthHeaderKey  = "invalid authorization header value"

	authHeaderKey = "Authorization"
	bearerScheme  = "Bearer"
)

// AuthorizationMiddleware checks if user is authorized.
func (s *Server) AuthorizationMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get(authHeaderKey)
		if headerValue == "" {
			sendResponse(rw, ErrorResponse{Message: authHeaderMessage}, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(headerValue, " ")
		if len(parts) != 2 {
			sendResponse(rw, ErrorResponse{Message: invalidAuthHeaderKey}, http.StatusUnauthorized)
			return
		}

		securityScheme := parts[0]
		token := parts[1]
		if securityScheme != bearerScheme {
			sendResponse(rw, ErrorResponse{Message: invalidSecurityScheme}, http.StatusUnauthorized)
			return
		}
		if token == "" {
			sendResponse(rw, ErrorResponse{Message: emptyTokenMessage}, http.StatusUnauthorized)
			return
		}

		claims, err := s.Auth.ParseToken(token)
		if err != nil {
			sendResponse(rw, ErrorResponse{Message: fmt.Errorf("cannot parse JWT token: %w", err).Error()}, http.StatusUnauthorized)
			return
		}

		currentUserID = claims.Subject

		ctx := context.Set(r.Context(), currentUserID)

		nextHandler.ServeHTTP(rw, r.WithContext(ctx))
	})
}
