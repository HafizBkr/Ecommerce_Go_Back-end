package models

import (
    "time"
)

// Ticket représente un ticket acheté pour un événement.
type Ticket struct {
    ID          string    `json:"id" db:"id"`
    EventID     string    `json:"event_id" db:"event_id"` // L'ID de l'événement auquel ce ticket est lié
    UserID      string    `json:"user_id" db:"user_id"`   // L'ID de l'utilisateur ayant acheté le ticket
    Price       float64   `json:"price" db:"price"`       // Le prix du ticket
    Status      string    `json:"status" db:"status"`     // Le statut du ticket (par exemple "validé" ou "annulé")
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
