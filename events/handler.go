// events/handler.go
package events

import (
    "encoding/json"
    "net/http"
    "ecommerce-api/models"  // Ajustez le chemin selon votre projet
    "github.com/go-chi/chi/v5"
)

type EventHandler struct {
    repo *EventRepository
}

func NewEventHandler(repo *EventRepository) *EventHandler {
    return &EventHandler{repo: repo}
}

func (h *EventHandler) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
    var event models.Event
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }
    
    if err := h.repo.CreateEvent(event); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Événement créé avec succès",
        "status": "success",
    })
}

func (h *EventHandler) HandleGetAllEvents(w http.ResponseWriter, r *http.Request) {
    events, err := h.repo.GetAllEvents()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) HandleGetEventByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    event, err := h.repo.GetEventByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) HandleGetEventsByCategoryID(w http.ResponseWriter, r *http.Request) {
    categoryID := chi.URLParam(r, "id")
    if categoryID == "" {
        http.Error(w, "L'ID de la catégorie est requis", http.StatusBadRequest)
        return
    }
    events, err := h.repo.GetEventsByCategoryID(categoryID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(events)
}


func (h *EventHandler) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var event models.Event
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }
    
    event.ID = id
    if err := h.repo.UpdateEvent(event); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Événement mis à jour avec succès",
        "status": "success",
    })
}

func (h *EventHandler) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if err := h.repo.DeleteEvent(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Événement supprimé avec succès",
        "status": "success",
    })
}