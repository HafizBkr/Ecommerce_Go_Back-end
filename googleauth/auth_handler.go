package googleauth

import (
	"context"
	"ecommerce-api/models"
	"ecommerce-api/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

type GoogleAuthHandler struct {
	repo *repository.UserRepository
}

// NewGoogleAuthHandler creates a new GoogleAuthHandler
func NewGoogleAuthHandler(repo *repository.UserRepository) *GoogleAuthHandler {
	return &GoogleAuthHandler{repo: repo}
}

// HandleOAuthRedirect initiates the Google OAuth2 login flow
func HandleOAuthRedirect(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	url := config.AuthCodeURL("state", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleAuthCallback handles the callback from Google OAuth2
// HandleAuthCallback handles the callback from Google OAuth2
// Modifiez la signature pour en faire une méthode de GoogleAuthHandler
// HandleAuthCallback handles the callback from Google OAuth2
func (h *GoogleAuthHandler) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Authorization code is missing", http.StatusBadRequest)
        return
    }

    clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
    clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")

    config := oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        RedirectURL:  "http://localhost:8080/auth/callback", // Redirige vers le backend pour la callback
        Scopes:       []string{"openid", "email", "profile"},
        Endpoint:     google.Endpoint,
    }

    token, err := config.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
        return
    }

    idToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "ID token not found in token response", http.StatusInternalServerError)
        return
    }

    claims, err := ValidateGoogleToken(context.Background(), idToken)
    if err != nil {
        http.Error(w, "Invalid ID token", http.StatusUnauthorized)
        return
    }

    email := claims["email"].(string)
    googleID := claims["sub"].(string)

    // Vérifier si l'utilisateur existe déjà
    _, err = h.repo.GetUserByGoogleID(googleID)
    isNewUser := err != nil // Si err != nil, l'utilisateur est nouveau

    var redirectURL string
    if isNewUser {
        // Rediriger vers le formulaire de complétion de profil sur le frontend
        // Inclure le `id_token` pour une utilisation sur le frontend
        redirectURL = fmt.Sprintf("http://localhost:3000/complete-profile?id_token=%s", idToken)
    } else {
        // Générer un JWT pour un utilisateur existant
        jwtToken, err := GenerateJWT(googleID, email)
        if err != nil {
            http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
            return
        }

        // Rediriger l'utilisateur connecté vers son tableau de bord
        redirectURL = fmt.Sprintf("http://localhost:3000/profile?jwt_token=%s", jwtToken)
    }

    // Effectuer la redirection
    http.Redirect(w, r, redirectURL, http.StatusFound)
}


// ValidateGoogleToken validates the ID token from Google
func ValidateGoogleToken(ctx context.Context, token string) (map[string]interface{}, error) {
	audience := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	payload, err := idtoken.Validate(ctx, token, audience)
	if err != nil {
		return nil, fmt.Errorf("invalid Google token: %v", err)
	}

	return payload.Claims, nil
}


type UpdateProfileRequest struct {
	Address          string `json:"address"`
	PhoneNumber      string `json:"phone_number"`
	ResidenceCity    string `json:"residence_city"`
	ResidenceCountry string `json:"residence_country"`
}

// HandleCompleteProfile handles profile completion or updates
func (h *GoogleAuthHandler) HandleCompleteProfile(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
        return
    }
    token := strings.TrimPrefix(authHeader, "Bearer ")

    claims, err := ValidateGoogleToken(ctx, token)
    if err != nil {
        http.Error(w, "Invalid Google token", http.StatusUnauthorized)
        return
    }

    email := claims["email"].(string)
    firstName := claims["given_name"].(string)
    lastName := claims["family_name"].(string)
    googleID := claims["sub"].(string)
	profilePicture := claims["picture"].(string)  // Récupération du Google ID depuis le token

    // Debug log
    fmt.Printf("Processing update for email: %s, Google ID: %s\n", email, googleID)

    // Essayer de trouver l'utilisateur par Google ID d'abord
    user, err := h.repo.GetUserByGoogleID(googleID)
    if err != nil {
        // Si non trouvé par Google ID, essayer par email
        user, err = h.repo.GetUserByEmail(email)
        if err != nil {
            fmt.Printf("Creating new user for email: %s\n", email)
            user = &models.User{
                Email:     email,
                FirstName: firstName,
                LastName:  lastName,
                Status:    "active",
                GoogleID:  googleID,
				ProfilePicture: profilePicture,
            }
            if err := h.repo.CreateUser(*user); err != nil {
                fmt.Printf("Error creating user: %v\n", err)
                http.Error(w, "Failed to create user", http.StatusInternalServerError)
                return
            }
        } else {
            // Si trouvé par email mais pas de Google ID, mettre à jour avec le Google ID
            user.GoogleID = googleID
        }
    }

    var payload UpdateProfileRequest
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        fmt.Printf("Error decoding payload: %v\n", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Mettre à jour les informations de l'utilisateur
    user.Address = payload.Address
    user.PhoneNumber = payload.PhoneNumber
    user.ResidenceCity = payload.ResidenceCity
    user.ResidenceCountry = payload.ResidenceCountry

    // Mettre à jour le profil utilisateur dans la base de données
    if err := h.repo.UpdateUser(*user); err != nil {
        fmt.Printf("Error updating user: %v\n", err)
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }



    // Répondre avec un message de succès et le JWT
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Profile updated successfully",
        "status":  "success",
        "google_id": googleID,
        "updated_fields": fmt.Sprintf("address: %s, phone: %s, city: %s, country: %s",
            user.Address, user.PhoneNumber, user.ResidenceCity, user.ResidenceCountry),  // Ajouter le JWT dans la réponse
    })
}


// GetUserHandler retrieves a user by their Google ID or email
func (h *GoogleAuthHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le token Bearer du header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Valider le token JWT et récupérer les claims
	claims, err := ValidateJWTToken(token, "HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE") // Assure-toi d'utiliser la clé secrète appropriée
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JWT token: %v", err), http.StatusUnauthorized)
		return
	}

	// Extraire le userID et l'email depuis les claims
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "User ID is missing in the token", http.StatusBadRequest)
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "Email is missing in the token", http.StatusBadRequest)
		return
	}

	// Chercher l'utilisateur dans la base de données par userID
	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		// Si l'utilisateur n'est pas trouvé par ID, chercher par email
		user, err = h.repo.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	}

	// Répondre avec les données de l'utilisateur
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding user data: %v", err), http.StatusInternalServerError)
	}
}