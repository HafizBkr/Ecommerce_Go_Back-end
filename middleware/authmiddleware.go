package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"ecommerce-api/googleauth" // Assurez-vous d'utiliser le bon package
)

// AuthMiddleware vérifie que le token JWT est présent et valide
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraire le token du header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Extraire le token de la forme "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Token is missing", http.StatusUnauthorized)
			return
		}

		// Utiliser la fonction de validation pour valider le token JWT
		claims, err := googleauth.ValidateJWTToken(tokenString, "HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE")
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid authentication token: %v", err), http.StatusUnauthorized)
			return
		}

		// Ajouter les informations de l'utilisateur dans le contexte de la requête
		ctx := context.WithValue(r.Context(), "user_claims", claims)

		// Passer le contexte modifié à la prochaine fonction du handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
