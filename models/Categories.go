package models

import "time"

// Category représente une catégorie dans la base de données
type Category struct {
    ID             string    `db:"id" json:"id"`
    Nom            string    `db:"nom" json:"nom"`
    NombreProduits int       `db:"nombre_produits" json:"nombre_produits"`
    Statut         string    `db:"statut" json:"statut"`
    CreatedAt      time.Time `db:"created_at" json:"created_at"`
    UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

