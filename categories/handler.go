package categories

import (
    "encoding/json"
    "fmt"
    "net/http"
    "ecommerce-api/models"
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
)

type CategoryHandler struct {
    repo *CategoryRepository
}

func NewCategoryHandler(repo *CategoryRepository) *CategoryHandler {
    return &CategoryHandler{repo: repo}
}

// HandleCreateCategory - Création d'une catégorie
func (h *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
    var category models.Category
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }
    
    if err := h.repo.CreateCategory(category); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catégorie créée avec succès",
        "status":  "success",
    })
}

// HandleGetAllCategories - Liste toutes les catégories
func (h *CategoryHandler) HandleGetAllCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := h.repo.GetAllCategories()
    if err != nil {
        http.Error(w, fmt.Sprintf("Échec de la récupération des catégories : %v", err), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

// HandleGetCategoryByID - Récupère une catégorie par ID
func (h *CategoryHandler) HandleGetCategoryByID(w http.ResponseWriter, r *http.Request) {
    // Extraction de l'ID
    id := chi.URLParam(r, "id")
    
    // Debug logging
    fmt.Printf("URL complète: %s\n", r.URL.Path)
    fmt.Printf("ID extrait: %s\n", id)
    
    if id == "" {
        http.Error(w, "ID non trouvé dans l'URL", http.StatusBadRequest)
        return
    }
    
    // Validation UUID
    _, err := uuid.Parse(id)
    if err != nil {
        http.Error(w, "Format d'ID invalide", http.StatusBadRequest)
        return
    }
    
    category, err := h.repo.GetCategoryByID(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("Catégorie non trouvée : %v", err), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(category)
}

// HandleUpdateCategory - Met à jour une catégorie
func (h *CategoryHandler) HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {
    // Extraction de l'ID
    id := chi.URLParam(r, "id")
    
    // Debug logging
    fmt.Printf("URL complète: %s\n", r.URL.Path)
    fmt.Printf("ID extrait: %s\n", id)
    
    if id == "" {
        http.Error(w, "ID non trouvé dans l'URL", http.StatusBadRequest)
        return
    }
    
    // Validation UUID
    _, err := uuid.Parse(id)
    if err != nil {
        http.Error(w, "Format d'ID invalide", http.StatusBadRequest)
        return
    }
    
    // Vérification de l'existence de la catégorie
    existingCategory, err := h.repo.GetCategoryByID(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("Catégorie non trouvée: %v", err), http.StatusNotFound)
        return
    }
    
    if existingCategory == nil {
        http.Error(w, "Catégorie non trouvée", http.StatusNotFound)
        return
    }
    
    // Décodage du corps de la requête
    var updatedCategory models.Category
    if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }
    
    // Mise à jour de l'ID et conservation des champs non modifiés
    updatedCategory.ID = id
    
    // Si certains champs sont vides dans la requête, on garde les valeurs existantes
    if updatedCategory.Nom == "" {
        updatedCategory.Nom = existingCategory.Nom
    }
    if updatedCategory.NombreProduits == 0 {
        updatedCategory.NombreProduits = existingCategory.NombreProduits
    }
    if updatedCategory.Statut == "" {
        updatedCategory.Statut = existingCategory.Statut
    }
    
    // Exécution de la mise à jour
    if err := h.repo.UpdateCategory(updatedCategory); err != nil {
        http.Error(w, fmt.Sprintf("Échec de la mise à jour : %v", err), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catégorie mise à jour avec succès",
        "status":  "success",
    })
}

// HandleDeleteCategory - Supprime une catégorie
func (h *CategoryHandler) HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {
    // Extraction de l'ID
    id := chi.URLParam(r, "id")
    
    // Debug logging
    fmt.Printf("URL complète: %s\n", r.URL.Path)
    fmt.Printf("ID extrait: %s\n", id)
    
    if id == "" {
        http.Error(w, "ID non trouvé dans l'URL", http.StatusBadRequest)
        return
    }
    
    // Validation UUID
    _, err := uuid.Parse(id)
    if err != nil {
        http.Error(w, "Format d'ID invalide", http.StatusBadRequest)
        return
    }
    
    if err := h.repo.DeleteCategory(id); err != nil {
        http.Error(w, fmt.Sprintf("Échec de la suppression : %v", err), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catégorie supprimée avec succès",
        "status":  "success",
    })
}