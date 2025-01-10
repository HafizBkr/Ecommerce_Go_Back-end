// products/repository.go
package products

import (
    "database/sql"
    "ecommerce-api/models"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"
)

type ProductRepository struct {
    db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
    return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(product models.Product) error {
    query := `
        INSERT INTO produits (
            id, nom, prix, stock, etat, photos, categorie_id,
            localisation, description, nombre_vues, disponible,
            marque, modele, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
        ) RETURNING id`
    
    id := uuid.New().String()
    now := time.Now()
    
    _, err := r.db.Exec(
        query,
        id, product.Nom, product.Prix, product.Stock,
        product.Etat, pq.Array(product.Photos), product.CategorieID,
        product.Localisation, product.Description, 0, true,
        product.Marque, product.Modele, now, now,
    )
    
    if err != nil {
        return fmt.Errorf("erreur lors de la création du produit : %v", err)
    }
    return nil
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
    var products []models.Product
    query := `
        SELECT id, nom, prix, stock, etat, photos, categorie_id,
               localisation, description, nombre_vues, disponible,
               marque, modele, created_at, updated_at
        FROM produits`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des produits : %v", err)
    }
    defer rows.Close()

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
            return nil, fmt.Errorf("erreur lors du scan des produits : %v", err)
        }
        product.Photos = photos
        products = append(products, product)
    }

    return products, nil
}


func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
    query := `
        SELECT id, nom, prix, stock, etat, photos, categorie_id,
               localisation, description, nombre_vues, disponible,
               marque, modele, created_at, updated_at
        FROM produits WHERE id = $1`

    var product models.Product
    var photos []string

    row := r.db.QueryRow(query, id)
    err := row.Scan(
        &product.ID, &product.Nom, &product.Prix, &product.Stock,
        &product.Etat, pq.Array(&photos), &product.CategorieID,
        &product.Localisation, &product.Description, &product.NombreVues,
        &product.Disponible, &product.Marque, &product.Modele,
        &product.CreatedAt, &product.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("aucun produit trouvé avec l'ID %s", id)
        }
        return nil, fmt.Errorf("erreur lors de la récupération du produit : %v", err)
    }

    product.Photos = photos
    return &product, nil
}


func (r *ProductRepository) UpdateProduct(product models.Product) error {
    query := `
        UPDATE produits 
        SET nom = $1, prix = $2, stock = $3, etat = $4,
            photos = $5, categorie_id = $6, localisation = $7,
            description = $8, disponible = $9, marque = $10,
            modele = $11, updated_at = $12
        WHERE id = $13`
    
    _, err := r.db.Exec(
        query,
        product.Nom, product.Prix, product.Stock, product.Etat,
        pq.Array(product.Photos), product.CategorieID, product.Localisation,
        product.Description, product.Disponible, product.Marque,
        product.Modele, time.Now(), product.ID,
    )
    
    if err != nil {
        return fmt.Errorf("erreur lors de la mise à jour du produit : %v", err)
    }
    return nil
}

func (r *ProductRepository) DeleteProduct(id string) error {
    query := `DELETE FROM produits WHERE id = $1`
    
    _, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("erreur lors de la suppression du produit : %v", err)
    }
    return nil
}

func (r *ProductRepository) GetProductsByCategory(categoryID string) ([]models.Product, error) {
    var products []models.Product
    query := `
        SELECT id, nom, prix, stock, etat, photos, categorie_id,
               localisation, description, nombre_vues, disponible,
               marque, modele, created_at, updated_at
        FROM produits
        WHERE categorie_id = $1`

    rows, err := r.db.Query(query, categoryID)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des produits par catégorie : %v", err)
    }
    defer rows.Close()

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
            return nil, fmt.Errorf("erreur lors du scan des produits : %v", err)
        }
        product.Photos = photos
        products = append(products, product)
    }

    // Vérifier les erreurs de la boucle rows.Next()
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("erreur lors de l'itération sur les résultats des produits : %v", err)
    }

    return products, nil
}

