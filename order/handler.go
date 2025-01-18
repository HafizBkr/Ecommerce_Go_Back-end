package order

import (
    "encoding/json"
    "ecommerce-api/models"
    "ecommerce-api/googleauth"
    "fmt"
    "net/http"
    "strings"
)

type Handler struct {
    repo         *Repository
    emailService models.EmailService
}

func NewHandler(repo *Repository, emailService models.EmailService) *Handler {
    return &Handler{
        repo:         repo,
        emailService: emailService,
    }
}

// HandleCreerCommande gère la création d'une commande.
func (h *Handler) HandleCreerCommande(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
        return
    }

    token := strings.TrimPrefix(authHeader, "Bearer ")
    claims, err := googleauth.ValidateJWTToken(token, "HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE")
    if err != nil {
        http.Error(w, fmt.Sprintf("Token JWT invalide : %v", err), http.StatusUnauthorized)
        return
    }

    googleID, ok := claims["user_id"].(string)
    if !ok || googleID == "" {
        http.Error(w, "L'ID utilisateur est manquant dans le token", http.StatusBadRequest)
        return
    }

    email, ok := claims["email"].(string)
    if !ok || email == "" {
        http.Error(w, "L'email est manquant dans le token", http.StatusBadRequest)
        return
    }

    var req struct {
        Produits []*models.CommandeProduit `json:"produits"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Format de requête invalide", http.StatusBadRequest)
        return
    }

    if len(req.Produits) == 0 {
        http.Error(w, "La commande doit contenir au moins un produit", http.StatusBadRequest)
        return
    }

    commande, err := h.repo.CreerCommande(googleID, req.Produits)
    if err != nil {
        http.Error(w, fmt.Sprintf("Erreur lors de la création de la commande: %v", err), http.StatusInternalServerError)
        return
    }

    err = h.emailService.EnvoyerEmailConfirmationCommande(commande, email)
    if err != nil {
        fmt.Printf("Erreur lors de l'envoi de l'email: %v\n", err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Commande créée avec succès",
        "data":    commande,
    })
}

// HandleListerCommandes gère la récupération de toutes les commandes d'un utilisateur.
func (h *Handler) HandleListerCommandes(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
        return
    }

    token := strings.TrimPrefix(authHeader, "Bearer ")
    claims, err := googleauth.ValidateJWTToken(token, "HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE")
    if err != nil {
        http.Error(w, fmt.Sprintf("Token JWT invalide : %v", err), http.StatusUnauthorized)
        return
    }

    googleID, ok := claims["user_id"].(string)
    if !ok || googleID == "" {
        http.Error(w, "L'ID utilisateur est manquant dans le token", http.StatusBadRequest)
        return
    }

    commandes, err := h.repo.ListerCommandesParUtilisateur(googleID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Erreur lors de la récupération des commandes: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "data":    commandes,
    })
}

func (h *Handler) HandleCreerCommandeTicket(w http.ResponseWriter, r *http.Request) {
    // Vérification de la méthode HTTP
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    // Vérification du token d'authentification
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header manquant", http.StatusUnauthorized)
        return
    }

    token := strings.TrimPrefix(authHeader, "Bearer ")
    claims, err := googleauth.ValidateJWTToken(token, "HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE")
    if err != nil {
        http.Error(w, fmt.Sprintf("Token JWT invalide : %v", err), http.StatusUnauthorized)
        return
    }

    // Extraction de l'ID utilisateur et de l'email des claims
    userID, ok := claims["user_id"].(string)
    if !ok || userID == "" {
        http.Error(w, "ID utilisateur manquant dans le token", http.StatusBadRequest)
        return
    }

    email, ok := claims["email"].(string)
    if !ok || email == "" {
        http.Error(w, "Email manquant dans le token", http.StatusBadRequest)
        return
    }

    // Décodage du corps de la requête
    var req struct {
        EventID  string `json:"event_id"`
        Quantity int    `json:"quantity"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Format de requête invalide", http.StatusBadRequest)
        return
    }

    // Validation des données
    if req.EventID == "" {
        http.Error(w, "L'ID de l'événement est requis", http.StatusBadRequest)
        return
    }

    if req.Quantity <= 0 {
        http.Error(w, "La quantité doit être supérieure à zéro", http.StatusBadRequest)
        return
    }

    // Création de la commande
    commande, err := h.repo.CreerCommandeTicket(userID, req.EventID, req.Quantity)
    if err != nil {
        http.Error(w, fmt.Sprintf("Erreur lors de la création de la commande : %v", err), http.StatusInternalServerError)
        return
    }

    // Envoi de l'email de confirmation
    err = h.emailService.EnvoyerEmailConfirmationTicket(commande, email)
    if err != nil {
        fmt.Printf("Erreur lors de l'envoi de l'email : %v\n", err)
    }

    // Réponse au client
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Commande créée avec succès",
        "data":    commande,
    })
}
