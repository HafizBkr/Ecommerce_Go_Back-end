package googleauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

// Middleware pour vérifier le token JWT
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'en-tête Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Extraire le token de l'en-tête (format "Bearer <token>")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Token is missing", http.StatusUnauthorized)
			return
		}

		// Analyser le token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Vérifier que la méthode de signature est correcte
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Ajouter les informations utilisateur au contexte de la requête
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)

		// Passer à la route suivante avec les informations d'utilisateur dans le contexte
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
