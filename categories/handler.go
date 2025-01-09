package categories

import (
	"ecommerce-api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// CategoryHandler gère les requêtes HTTP relatives aux catégories.
type CategoryHandler struct {
	repo *CategoryRepository
}

// NewCategoryHandler crée un nouveau handler pour les catégories.
func NewCategoryHandler(repo *CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// HandleCreateCategory gère la création d'une catégorie.
func (h *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	// Décodage du corps de la requête en un objet Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Appel au dépôt pour créer la catégorie
	if err := h.repo.CreateCategory(category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retour de la réponse
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Catégorie créée avec succès",
		"status":  "success",
	})
}

// HandleUpdateCategory gère la mise à jour d'une catégorie.
// HandleUpdateCategory gère la mise à jour d'une catégorie.
func (h *CategoryHandler) HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {
    // Récupérer l'ID à partir du paramètre de la requête
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "ID manquant", http.StatusBadRequest)
        return
    }

    // Conversion de l'ID de la catégorie en entier
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID de catégorie invalide", http.StatusBadRequest)
        return
    }

    // Décodage du corps de la requête en un objet Category
    var category models.Category
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }

    category.ID = id // Ajouter l'ID à la catégorie

    // Appel au dépôt pour mettre à jour la catégorie
    if err := h.repo.UpdateCategory(category); err != nil {
        http.Error(w, fmt.Sprintf("Échec de la mise à jour de la catégorie : %v", err), http.StatusInternalServerError)
        return
    }

    // Retour de la réponse
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catégorie mise à jour avec succès",
        "status":  "success",
    })
}

// HandleDeleteCategory gère la suppression d'une catégorie.
func (h *CategoryHandler) HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {
    // Récupérer l'ID à partir du paramètre de la requête
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "ID manquant", http.StatusBadRequest)
        return
    }

    // Conversion de l'ID de la catégorie en entier
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID de catégorie invalide", http.StatusBadRequest)
        return
    }

    // Appel au dépôt pour supprimer la catégorie
    if err := h.repo.DeleteCategory(id); err != nil {
        http.Error(w, fmt.Sprintf("Échec de la suppression de la catégorie : %v", err), http.StatusInternalServerError)
        return
    }

    // Retour de la réponse
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catégorie supprimée avec succès",
        "status":  "success",
    })
}


// HandleGetCategory gère la récupération d'une catégorie par son ID.
func (h *CategoryHandler) HandleGetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	// Conversion de l'ID de la catégorie en entier
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de catégorie invalide", http.StatusBadRequest)
		return
	}

	// Appel au dépôt pour récupérer la catégorie
	category, err := h.repo.GetCategoryByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Échec de la récupération de la catégorie : %v", err), http.StatusNotFound)
		return
	}

	// Retour de la catégorie en format JSON
	json.NewEncoder(w).Encode(category)
}

// HandleGetAllCategories gère la récupération de toutes les catégories.
func (h *CategoryHandler) HandleGetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Appel au dépôt pour récupérer toutes les catégories
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		http.Error(w, fmt.Sprintf("Échec de la récupération des catégories : %v", err), http.StatusInternalServerError)
		return
	}

	// Retour des catégories en format JSON
	json.NewEncoder(w).Encode(categories)
}
