package models

import "time"

// ListeSouhaits représente une entrée dans la liste des souhaits
type ListeSouhaits struct {
    ID        string    `db:"id" json:"id"`               // UUID pour l'id de la liste de souhaits
    GoogleID  string    `db:"google_id" json:"google_id"` // Identifiant de l'utilisateur
    ProduitID string    `db:"produit_id" json:"produit_id"` // Identifiant du produit
    CreatedAt time.Time `db:"created_at" json:"created_at"` // Date de création
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // Date de mise à jour
}

// ProduitSouhait représente les détails d'un produit dans la liste de souhaits
type ProduitSouhait struct {
    ID       string  `db:"id"`
    Nom      string  `db:"nom"`
    Prix     float64 `db:"prix"`
    Marque   string  `db:"marque"`
    Photos   []string `db:"photos"` // Si tu veux aussi garder toutes les photos
    Photo    string  // Champ supplémentaire pour la première photo
}