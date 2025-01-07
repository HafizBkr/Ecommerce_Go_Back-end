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
func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
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
		RedirectURL:  "http://localhost:8080/auth/callback",
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

	// Print the ID token
	fmt.Println("ID Token:", idToken)

	claims, err := ValidateGoogleToken(context.Background(), idToken)
	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	email := claims["email"].(string)
	name := claims["name"].(string)

	fmt.Fprintf(w, "Bienvenue %s (%s) !", name, email)
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

	user, err := h.repo.GetUserByEmail(email)
	if err != nil {
		user = &models.User{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			Status:    "active",
		}
		if err := h.repo.CreateUser(*user); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	}

	var payload UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println("Error decoding payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received payload: %+v\n", payload)

	// Update user information
	user.Address = payload.Address
	user.PhoneNumber = payload.PhoneNumber
	user.ResidenceCity = payload.ResidenceCity
	user.ResidenceCountry = payload.ResidenceCountry

	if err := h.repo.UpdateUser(*user); err != nil {
		fmt.Println("Error updating user:", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
