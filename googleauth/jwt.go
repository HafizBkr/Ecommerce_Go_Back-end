package googleauth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Clé secrète pour signer les tokens JWT
var jwtKey = []byte("HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE")

// Structure qui représente les informations contenues dans le token JWT
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// Fonction pour générer un token JWT
func GenerateJWT(userID, email string) (string, error) {
	// Définir la durée d'expiration du token (par exemple, 24 heures)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "ton_application", // Nom de ton application
		},
	}

	// Créer un nouvel objet token avec les claims et la méthode de signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signer le token avec la clé secrète et générer la chaîne de caractères
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}




// ValidateJWTToken va valider le token JWT
func ValidateJWTToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Valider la méthode de signature (HS256 dans cet exemple)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
