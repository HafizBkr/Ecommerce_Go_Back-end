package models

import (
    "time"
)

type Commande struct {
    ID            string    `json:"id" db:"id"`
    NumeroCommande string    `json:"numero_commande" db:"numero_commande"`
    UserID        string    `json:"user_id" db:"user_id"`
    MontantTotal  float64   `json:"montant_total" db:"montant_total"`
    Status        string    `json:"status" db:"status"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CommandeProduit struct {
    CommandeID string  `json:"commande_id" db:"commande_id"`
    ProduitID  string  `json:"produit_id" db:"produit_id"`
    Quantite   int     `json:"quantite" db:"quantite"`
    PrixUnite  float64 `json:"prix_unite" db:"prix_unite"`
}



type EmailService interface {
    EnvoyerEmailConfirmationCommande(commande *Commande, email string) error
}


