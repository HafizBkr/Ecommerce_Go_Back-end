package middlewares

import (
	"context"
	"ecommerce-api/googleauth"
	"net/http"
	"strings"
	"fmt"
)

// GoogleAuthMiddleware vérifie que le token Google est présent et valide
func GoogleAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraire le token du header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Extraire le token de la forme "Bearer <token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Token is missing", http.StatusUnauthorized)
			return
		}

		// Valider le token Google
		claims, err := googleauth.ValidateGoogleToken(r.Context(), token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid Google token: %v", err), http.StatusUnauthorized)
			return
		}

		// Ajouter les informations de l'utilisateur dans le contexte de la requête
		ctx := context.WithValue(r.Context(), "google_claims", claims)

		// Passer le contexte modifié à la prochaine fonction du handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
