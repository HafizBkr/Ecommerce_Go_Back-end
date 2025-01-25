// models/product.go
package models

import "time"

type Product struct {
    ID          string    `db:"id" json:"id"`
    Nom         string    `db:"nom" json:"nom"`
    Prix        float64   `db:"prix" json:"prix"`
    Stock       int       `db:"stock" json:"stock"`
    Etat        string    `db:"etat" json:"etat"`
    Photos      []string  `db:"photos" json:"photos"`
    CategorieID string    `db:"categorie_id" json:"categorie_id"`
    CategorieNom  string    `json:"categorie_nom"`
    Localisation string   `db:"localisation" json:"localisation"`
    Description string    `db:"description" json:"description"`
    NombreVues  int       `db:"nombre_vues" json:"nombre_vues"`
    Disponible  bool      `db:"disponible" json:"disponible"`
    Marque      string    `db:"marque" json:"marque"`
    Modele      string    `db:"modele" json:"modele"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}