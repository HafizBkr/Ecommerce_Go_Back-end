package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"ecommerce-api/internal/domain/models"
	"ecommerce-api/internal/domain/repositories" // Assurez-vous que cet import est utilisé correctement
)

type AuthHandler struct {
	Repo  repository.UserRepository
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{Repo: *repo}
}

func (h *AuthHandler) HandleOAuthRedirect(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	url := config.AuthCodeURL("random", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code manquant dans la requête", http.StatusBadRequest)
		return
	}

	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Erreur lors de l'échange du token : "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des informations utilisateur", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Erreur lors du décodage des informations utilisateur", http.StatusInternalServerError)
		return
	}

	// Découper le nom complet en prénom et nom
	names := splitName(userInfo.Name)
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}

	// Vérifier si l'utilisateur existe déjà
	user, err := h.Repo.GetUserByEmail(userInfo.Email)
	if err != nil || user == nil {
		// Si l'utilisateur n'existe pas, le créer
		newUser := &models.User{
			Email:     userInfo.Email,
			FirstName: firstName,
			LastName:  lastName,
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := h.Repo.CreateUser(*newUser); err != nil {
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}
		user = newUser
	}

	// Mettre à jour le dernier login
	user.LastLogin = time.Now()
	if err := h.Repo.UpdateUser(*user); err != nil {
		http.Error(w, "Erreur lors de la mise à jour des informations utilisateur", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/complete-profile?email=%s", user.Email), http.StatusSeeOther)
}

func splitName(fullName string) []string {
	return strings.Fields(fullName)
}

func (h *AuthHandler) HandleCompleteProfile(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email manquant", http.StatusBadRequest)
		return
	}

	// Formulaire pour compléter le profil
	tmpl := template.Must(template.ParseFiles("templates/complete-profile.html"))
	if err := tmpl.Execute(w, struct{ Email string }{Email: email}); err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) SaveUserProfile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur lors du traitement des données", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	phoneNumber := r.FormValue("phone_number")
	city := r.FormValue("residence_city")
	country := r.FormValue("residence_country")

	if email == "" || phoneNumber == "" || city == "" || country == "" {
		http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
		return
	}

	err := h.Repo.SaveUserProfile(email, phoneNumber, city, country)
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde du profil", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
