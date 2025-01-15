package panier

import (
    "ecommerce-api/googleauth" // Importer votre package googleauth
    "fmt"
    "net/http"
    "strings"
    "encoding/json"
)

type PanierHandler struct {
    repo *Repository
}

func NewPanierHandler(repo *Repository) *PanierHandler {
    return &PanierHandler{repo: repo}
}

// HandleAjouterProduit ajoute un produit au panier après validation du token JWT.
// HandleAjouterProduit adds a product to the cart after validating the JWT token.
func (h *PanierHandler) HandleAjouterProduit(w http.ResponseWriter, r *http.Request) {
    // Get the Bearer token from the Authorization header
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
        return
    }
    token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the JWT token and retrieve claims
			claims, err := googleauth.ValidateJWTToken(token, "HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE")
		if err != nil {
			http.Error(w, fmt.Sprintf("Token JWT invalide : %v", err), http.StatusUnauthorized)
			return
		}

		fmt.Printf("Claims : %+v\n", claims)  // Affiche les informations du token

		// Récupérer l'ID Google à partir du champ "user_id" du token
		googleID, ok := claims["user_id"].(string)
		if !ok || googleID == "" {
			http.Error(w, "L'ID Google est manquant dans le token", http.StatusBadRequest)
			return
		}

	

    // Decode the product from the request body
    var req struct {
        ProduitID string `json:"produit_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProduitID == "" {
        http.Error(w, "Invalid request format or missing ProduitID", http.StatusBadRequest)
        return
    }

    // Add the product to the user's cart (using repository method)
    if err := h.repo.AjouterProduitAuPanier(googleID, req.ProduitID); err != nil {
        http.Error(w, "Error adding product to cart", http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Product added to cart",
    })

	fmt.Printf("Claims: %+v\n", claims)

}

// HandleAfficherPanier récupère le contenu du panier d'un utilisateur.
func (h *PanierHandler) HandleAfficherPanier(w http.ResponseWriter, r *http.Request) {
    // Récupérer et valider le token JWT
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
        return
    }
    token := strings.TrimPrefix(authHeader, "Bearer ")

    claims, err := googleauth.ValidateJWTToken(token, "HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE") // Utilisez la clé secrète définie dans googleauth
    if err != nil {
        http.Error(w, fmt.Sprintf("Invalid JWT token: %v", err), http.StatusUnauthorized)
        return
    }

    googleID, ok := claims["user_id"].(string)
    if !ok || googleID == "" {
        http.Error(w, "Google ID is missing in the token", http.StatusBadRequest)
        return
    }

    // Récupérer le contenu du panier
    panier, err := h.repo.ObtenirPanierParUserID(googleID)
    if err != nil {
        http.Error(w, "Erreur lors de la récupération du panier", http.StatusInternalServerError)
        return
    }

    // Répondre avec le contenu du panier
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "success",
        "data":   panier,
    })
}

func (h *PanierHandler) HandleEnleverDuPanier(w http.ResponseWriter, r *http.Request) {
    // Récupérer et valider le token JWT
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
        return
    }
    token := strings.TrimPrefix(authHeader, "Bearer ")

    claims, err := googleauth.ValidateJWTToken(token, "HDBCSOAVNOAHBVIJVNYWUONCPOIEUIBVE") // Utilisez la clé secrète définie dans googleauth
    if err != nil {
        http.Error(w, fmt.Sprintf("Invalid JWT token: %v", err), http.StatusUnauthorized)
        return
    }

    googleID, ok := claims["user_id"].(string)
    if !ok || googleID == "" {
        http.Error(w, "Google ID is missing in the token", http.StatusBadRequest)
        return
    }

    // Decode the request body to get the product ID to remove
    var req struct {
        ProduitID string `json:"produit_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProduitID == "" {
        http.Error(w, "Invalid request format or missing ProduitID", http.StatusBadRequest)
        return
    }

    // Call the repository method to remove the product from the cart
    if err := h.repo.EnleverDuPanier(googleID, req.ProduitID); err != nil {
        http.Error(w, "Error removing product from cart", http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Product removed from cart",
    })
}