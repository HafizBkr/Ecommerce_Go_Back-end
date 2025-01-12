// models/filters.go
package models

type ProductFilters struct {
    PrixMin         *float64 `json:"prix_min,omitempty"`
    PrixMax         *float64 `json:"prix_max,omitempty"`
    Marque          []string `json:"marque,omitempty"`
    Etat            []string `json:"etat,omitempty"`
    Localisation    []string `json:"localisation,omitempty"`
    CategorieID     string   `json:"categorie_id,omitempty"`
    Disponible      *bool    `json:"disponible,omitempty"`
    SearchTerm      string   `json:"search_term,omitempty"`
}