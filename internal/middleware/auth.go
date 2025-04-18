package middleware

import (
	"Sector/internal/auth"
	"context"
	"net/http"
	"strings"
)

// AuthUserCtxKey is the context key for the authenticated user
type AuthUserCtxKey struct{}

// UserInfo represents authenticated user information
type UserInfo struct {
	UserID   string
	Username string
}

// JWTAuth is a middleware that validates JWT tokens
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for specific endpoints
		if isExemptPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract the token from the Authorization header
		tokenString, err := auth.ExtractTokenFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user information to the request context
		userInfo := UserInfo{
			UserID:   claims.UserID,
			Username: claims.Username,
		}

		ctx := context.WithValue(r.Context(), AuthUserCtxKey{}, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extracts the UserInfo from the request context
func GetUserFromContext(ctx context.Context) (UserInfo, bool) {
	userInfo, ok := ctx.Value(AuthUserCtxKey{}).(UserInfo)
	return userInfo, ok
}

// isExemptPath determines if a path should be exempt from authentication
func isExemptPath(path string) bool {
	exemptPaths := []string{
		"/v1/api/",
		"/v1/api/v1/api/",
		"/v1/api/v1/api/health",
		"/v1/api/v1/api/challenge",
		"/v1/api/v1/api/login",
		"/v1/swagger-ui/",
		"/v1/swagger.json",
	}

	for _, exemptPath := range exemptPaths {
		if path == exemptPath || strings.HasPrefix(path, exemptPath) {
			return true
		}
	}

	return false
}
