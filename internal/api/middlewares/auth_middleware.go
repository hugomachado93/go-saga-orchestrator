package middlewares

import (
	"main/internal/services"
	"net/http"
)

type AuthMiddleware struct {
	cs *services.ClientService
}

func NewAuthmiddleware(cs *services.ClientService) *AuthMiddleware {
	return &AuthMiddleware{cs: cs}
}

func (am *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xapk := r.Header.Get("x-api-key")
		if !am.cs.IsClientAuthorized(xapk) {
			http.Error(w, "User with x-api-key not authorized", http.StatusUnauthorized)
			return
		}
		// Call the next middleware function or final handler
		next.ServeHTTP(w, r)
	})
}
