package events_category

import (
    "encoding/json"
    "net/http"
    "ecommerce-api/models"  // Assurez-vous d'ajuster le chemin d'import selon votre structure
    "github.com/go-chi/chi/v5"
)

type EventCategoryHandler struct {
    repo *EventCategoryRepository
}

func NewEventCategoryHandler(repo *EventCategoryRepository) *EventCategoryHandler {
    return &EventCategoryHandler{repo: repo}
}

func (h *EventCategoryHandler) HandleCreateEventCategory(w http.ResponseWriter, r *http.Request) {
    var category models.EventCategory
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }
    
    if err := h.repo.CreateEventCategory(category); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catégorie d'événement créée avec succès",
        "status": "success",
    })
}

func (h *EventCategoryHandler) HandleGetAllEventCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := h.repo.GetAllEventCategories()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

func (h *EventCategoryHandler) HandleGetEventCategoryByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    category, err := h.repo.GetEventCategoryByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(category)
}

func (h *EventCategoryHandler) HandleUpdateEventCategory(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var category models.EventCategory
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }
    
    category.ID = id
    if err := h.repo.UpdateEventCategory(category); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catégorie d'événement mise à jour avec succès",
        "status": "success",
    })
}

func (h *EventCategoryHandler) HandleDeleteEventCategory(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if err := h.repo.DeleteEventCategory(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catégorie d'événement supprimée avec succès",
        "status": "success",
    })
}