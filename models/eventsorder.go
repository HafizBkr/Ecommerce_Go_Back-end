package models

import "time"

// TicketOrder représente une commande de ticket
type TicketOrder struct {
    ID             string    `json:"id" db:"id"`
    NumeroCommande string    `json:"numero_commande" db:"numero_commande"`
    UserID         string    `json:"user_id" db:"user_id"`
    EventID        string    `json:"event_id" db:"event_id"`
    EventTitle     string    `json:"event_title" db:"event_title"`
    StartDate      time.Time `json:"start_date" db:"start_date"`
    StartTime      string    `json:"start_time" db:"start_time"`
    Quantity       int       `json:"quantity" db:"quantity"`
    PrixTotal      float64   `json:"prix_total" db:"prix_total"`
    Status         string    `json:"status" db:"status"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
