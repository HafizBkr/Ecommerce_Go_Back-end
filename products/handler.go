// products/handler.go
package products

import (
    "encoding/json"
    "fmt"
    "net/http"
    "ecommerce-api/models"
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
)

type ProductHandler struct {
    repo *ProductRepository
}

func NewProductHandler(repo *ProductRepository) *ProductHandler {
    return &ProductHandler{repo: repo}
}

func (h *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
    var product models.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }
    
    if err := h.repo.CreateProduct(product); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Produit créé avec succès",
        "status":  "success",
    })
}

func (h *ProductHandler) HandleGetAllProducts(w http.ResponseWriter, r *http.Request) {
    products, err := h.repo.GetAllProducts()
    if err != nil {
        http.Error(w, fmt.Sprintf("Échec de la récupération des produits : %v", err), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) HandleGetProductByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    if id == "" {
        http.Error(w, "ID non trouvé dans l'URL", http.StatusBadRequest)
        return
    }
    
    if _, err := uuid.Parse(id); err != nil {
        http.Error(w, "Format d'ID invalide", http.StatusBadRequest)
        return
    }
    
    product, err := h.repo.GetProductByID(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("Produit non trouvé : %v", err), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    if id == "" {
        http.Error(w, "ID non trouvé dans l'URL", http.StatusBadRequest)
        return
    }
    
    if _, err := uuid.Parse(id); err != nil {
        http.Error(w, "Format d'ID invalide", http.StatusBadRequest)
        return
    }
    
    existingProduct, err := h.repo.GetProductByID(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("Produit non trouvé: %v", err), http.StatusNotFound)
        return
    }
    
    var updatedProduct models.Product
    if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }
    
    updatedProduct.ID = id
    
    // Conserver les valeurs existantes si non fournies
    if updatedProduct.Nom == "" {
        updatedProduct.Nom = existingProduct.Nom
    }
    if updatedProduct.Prix == 0 {
        updatedProduct.Prix = existingProduct.Prix
    }
    // ... et ainsi de suite pour les autres champs
    
    if err := h.repo.UpdateProduct(updatedProduct); err != nil {
        http.Error(w, fmt.Sprintf("Échec de la mise à jour : %v", err), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Produit mis à jour avec succès",
        "status":  "success",
    })
}

func (h *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    if id == "" {
        http.Error(w, "ID non trouvé dans l'URL", http.StatusBadRequest)
        return
    }
    
    if _, err := uuid.Parse(id); err != nil {
        http.Error(w, "Format d'ID invalide", http.StatusBadRequest)
        return
    }
    
    if err := h.repo.DeleteProduct(id); err != nil {
        http.Error(w, fmt.Sprintf("Échec de la suppression : %v", err), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Produit supprimé avec succès",
        "status":  "success",
    })
}

func (h *ProductHandler) HandleGetProductsByCategory(w http.ResponseWriter, r *http.Request) {
    categoryID := chi.URLParam(r, "categoryID")

    // Vérification du format de l'UUID
    if _, err := uuid.Parse(categoryID); err != nil {
        http.Error(w, "ID de catégorie invalide", http.StatusBadRequest)
        return
    }

    products, err := h.repo.GetProductsByCategory(categoryID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Erreur lors de la récupération des produits : %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}
