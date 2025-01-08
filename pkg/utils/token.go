package utils

import (
 // Remplacez par votre package exact
	"ecommerce-api/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("votre-cle-secrete") // À stocker dans une variable d'environnement

func GenerateJWT(user *models.User) (string, error) {
    claims := jwt.MapClaims{
        "id":                user.ID,
        "google_id":         user.GoogleID,
        "email":             user.Email,
        "first_name":        user.FirstName,
        "last_name":         user.LastName,
        "is_admin":          user.IsAdmin,
        "points":            user.Points,
        "status":            user.Status,
        "address":           user.Address,
        "phone_number":      user.PhoneNumber,
        "residence_city":    user.ResidenceCity,
        "residence_country": user.ResidenceCountry,
        "last_login":        user.LastLogin.Unix(),
        "created_at":        user.CreatedAt.Unix(),
        "updated_at":        user.UpdatedAt.Unix(),
        "exp":               time.Now().Add(24 * time.Hour).Unix(), // Expire dans 24h
        "iat":               time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
