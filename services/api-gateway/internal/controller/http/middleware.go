package http

import (
	"net/http"
	"strings"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		publicPaths := []string{
			"/v1/auth/register",
			"/v1/auth/login",
			"/v1/auth/refresh",
			"/v1/auth/forgot-password",
			"/v1/auth/reset-password",
		}

		for _, publicPath := range publicPaths {
			if path == publicPath || strings.HasPrefix(path, publicPath+"/") {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: missing or invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		id, err := s.authClient.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-Id", id)

		next.ServeHTTP(w, r)
	})
}
