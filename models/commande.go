package models

import (
    "time"
)

type ProduitDetail struct {
    Nom       string  `json:"nom" db:"nom"`
    PrixUnite float64 `json:"prix_unite" db:"prix_unite"`
    Quantite  int     `json:"quantite" db:"quantite"`
    Model     string  `json:"model" db:"model"`         // Ajout du modèle
    Etat         string  `json:"etat" db:"etat"`
    Localisation string  `json:"localisation" db:"localisation"`
    Photos       []string `json:"photos" db:"photos"`
}

type Commande struct {
    ID             string          `json:"id" db:"id"`
    NumeroCommande string          `json:"numero_commande" db:"numero_commande"`
    UserID         string          `json:"user_id" db:"user_id"`
    MontantTotal   float64         `json:"montant_total" db:"montant_total"`
    Status         string          `json:"status" db:"status"`
    CreatedAt      time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
    Produits       []ProduitDetail `json:"produits" db:"-"` // Le "-" indique d'ignorer ce champ pour la BD
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

type CommandeDetail struct {
    ID             string    `db:"id" json:"id"`
    NumeroCommande string    `db:"numero_commande" json:"numero_commande"`
    GoogleID       string    `db:"google_id" json:"google_id"`
    FirstName      string    `db:"first_name" json:"first_name"`
    LastName       string    `db:"last_name" json:"last_name"`
    Email          string    `db:"email" json:"email"`
    MontantTotal   float64   `db:"montant_total" json:"montant_total"`
    Status         string    `db:"status" json:"status"`
    CreatedAt      time.Time `db:"created_at" json:"created_at"`
    UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

