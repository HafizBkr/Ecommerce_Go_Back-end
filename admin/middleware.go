// admin/middleware.go
package admin

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt"
)

func AdminAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("ADMIN_SECRET_KEY"), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !claims["is_admin"].(bool) {
            http.Error(w, "Unauthorized: admin access required", http.StatusForbidden)
            return
        }

        ctx := context.WithValue(r.Context(), "admin_claims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}