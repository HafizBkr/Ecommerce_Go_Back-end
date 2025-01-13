
package search

import (
	"database/sql"
	"ecommerce-api/models"
	"fmt"

	"github.com/lib/pq"
)

type SearchEngine struct {
    db *sql.DB
}

func NewSearchEngine(db *sql.DB) *SearchEngine {
    return &SearchEngine{
        db: db,
    }
}

func (s *SearchEngine) Search(opts SearchOptions) (SearchResult, error) {
    // Valider et définir les valeurs par défaut
    if opts.Page <= 0 {
        opts.Page = 1
    }
    if opts.PageSize <= 0 {
        opts.PageSize = 20
    }
    if opts.PageSize > 100 {
        opts.PageSize = 100
    }

    offset := (opts.Page - 1) * opts.PageSize
    searchTerm := "%" + opts.Query + "%"

    // Compter le total des résultats
    var total int
    countQuery := `
        SELECT COUNT(*) 
        FROM produits
        WHERE 
            disponible = true AND
            (
                nom ILIKE $1 OR
                description ILIKE $1 OR
                marque ILIKE $1 OR
                modele ILIKE $1
            )
    `
    
    err := s.db.QueryRow(countQuery, searchTerm).Scan(&total)
    if err != nil {
        return SearchResult{}, fmt.Errorf("erreur lors du comptage des résultats : %v", err)
    }

    // Récupérer les produits
    searchQuery := `
        SELECT 
            id, nom, prix, stock, etat, photos, categorie_id,
            localisation, description, nombre_vues, disponible,
            marque, modele, created_at, updated_at
        FROM produits
        WHERE 
            disponible = true AND
            (
                nom ILIKE $1 OR
                description ILIKE $1 OR
                marque ILIKE $1 OR
                modele ILIKE $1
            )
        ORDER BY 
            CASE 
                WHEN nom ILIKE $1 THEN 1
                WHEN marque ILIKE $1 THEN 2
                WHEN modele ILIKE $1 THEN 3
                ELSE 4
            END,
            nombre_vues DESC
        LIMIT $2 OFFSET $3
    `

    rows, err := s.db.Query(searchQuery, searchTerm, opts.PageSize, offset)
    if err != nil {
        return SearchResult{}, fmt.Errorf("erreur lors de la recherche des produits : %v", err)
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var product models.Product
        var photos []string
        err := rows.Scan(
            &product.ID, &product.Nom, &product.Prix, &product.Stock,
            &product.Etat, pq.Array(&photos), &product.CategorieID,
            &product.Localisation, &product.Description, &product.NombreVues,
            &product.Disponible, &product.Marque, &product.Modele,
            &product.CreatedAt, &product.UpdatedAt,
        )
        if err != nil {
            return SearchResult{}, fmt.Errorf("erreur lors du scan des produits : %v", err)
        }
        product.Photos = photos
        products = append(products, product)
    }

    return SearchResult{
        Products: products,
        Total:    total,
        Page:     opts.Page,
        PageSize: opts.PageSize,
    }, nil
}