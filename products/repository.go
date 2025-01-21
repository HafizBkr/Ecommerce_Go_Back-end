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
        SELECT
            p.id,
            p.nom,
            p.prix,
            p.stock,
            p.etat,
            p.photos,
            p.categorie_id,
            c.nom AS categorie_nom,  -- Récupérer le nom de la catégorie
            p.localisation,
            p.description,
            p.nombre_vues,
            p.disponible,
            p.marque,
            p.modele,
            p.created_at,
            p.updated_at
        FROM
            produits p
        JOIN
            categories c
        ON
            p.categorie_id = c.id
    `

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des produits : %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var product models.Product
        var photos []string
        var categorieNom string

        err := rows.Scan(
            &product.ID, &product.Nom, &product.Prix, &product.Stock,
            &product.Etat, pq.Array(&photos), &product.CategorieID,
            &categorieNom, // Récupérer le nom de la catégorie
            &product.Localisation, &product.Description, &product.NombreVues,
            &product.Disponible, &product.Marque, &product.Modele,
            &product.CreatedAt, &product.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("erreur lors du scan des produits : %v", err)
        }
        product.Photos = photos
        product.CategorieNom = categorieNom // Stocker le nom de la catégorie dans le produit
        products = append(products, product)
    }

    return products, nil
}



func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
    query := `
        SELECT
            p.id,
            p.nom,
            p.prix,
            p.stock,
            p.etat,
            p.photos,
            p.categorie_id,
            c.nom AS categorie_nom,  -- Récupérer le nom de la catégorie
            p.localisation,
            p.description,
            p.nombre_vues,
            p.disponible,
            p.marque,
            p.modele,
            p.created_at,
            p.updated_at
        FROM
            produits p
        JOIN
            categories c
        ON
            p.categorie_id = c.id
        WHERE
            p.id = $1
    `

    var product models.Product
    var photos []string
    var categorieNom string

    row := r.db.QueryRow(query, id)
    err := row.Scan(
        &product.ID, &product.Nom, &product.Prix, &product.Stock,
        &product.Etat, pq.Array(&photos), &product.CategorieID,
        &categorieNom, // Récupérer le nom de la catégorie
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
    product.CategorieNom = categorieNom // Stocker le nom de la catégorie dans le produit
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
        SELECT 
            p.id, 
            p.nom, 
            p.prix, 
            p.stock, 
            p.etat, 
            p.photos, 
            p.categorie_id, 
            c.nom AS categorie_nom,  -- Récupérer le nom de la catégorie
            p.localisation, 
            p.description, 
            p.nombre_vues, 
            p.disponible, 
            p.marque, 
            p.modele, 
            p.created_at, 
            p.updated_at
        FROM produits p
        JOIN categories c
        ON p.categorie_id = c.id
        WHERE p.categorie_id = $1`

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
            &product.CategorieNom, // Stocker le nom de la catégorie ici
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


// products/repository.go
func (r *ProductRepository) GetFilteredProducts(filters models.ProductFilters) ([]models.Product, error) {
    query := `
        SELECT id, nom, prix, stock, etat, photos, categorie_id,
               localisation, description, nombre_vues, disponible,
               marque, modele, created_at, updated_at
        FROM produits
        WHERE 1=1`
    
    var args []interface{}
    argCount := 1

    // Construction dynamique de la requête avec les filtres
    if filters.PrixMin != nil {
        query += fmt.Sprintf(" AND prix >= $%d", argCount)
        args = append(args, *filters.PrixMin)
        argCount++
    }

    if filters.PrixMax != nil {
        query += fmt.Sprintf(" AND prix <= $%d", argCount)
        args = append(args, *filters.PrixMax)
        argCount++
    }

    if len(filters.Marque) > 0 {
        query += fmt.Sprintf(" AND marque = ANY($%d)", argCount)
        args = append(args, pq.Array(filters.Marque))
        argCount++
    }

    if len(filters.Etat) > 0 {
        query += fmt.Sprintf(" AND etat = ANY($%d)", argCount)
        args = append(args, pq.Array(filters.Etat))
        argCount++
    }

    if len(filters.Localisation) > 0 {
        query += fmt.Sprintf(" AND localisation = ANY($%d)", argCount)
        args = append(args, pq.Array(filters.Localisation))
        argCount++
    }

    if filters.CategorieID != "" {
        query += fmt.Sprintf(" AND categorie_id = $%d", argCount)
        args = append(args, filters.CategorieID)
        argCount++
    }

    if filters.Disponible != nil {
        query += fmt.Sprintf(" AND disponible = $%d", argCount)
        args = append(args, *filters.Disponible)
        argCount++
    }

    if filters.SearchTerm != "" {
        searchTerm := "%" + filters.SearchTerm + "%"
        query += fmt.Sprintf(` AND (
            nom ILIKE $%d OR 
            description ILIKE $%d OR 
            marque ILIKE $%d OR 
            modele ILIKE $%d
        )`, argCount, argCount, argCount, argCount)
        args = append(args, searchTerm)
        argCount++
    }

    // Exécution de la requête
    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des produits filtrés : %v", err)
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
            return nil, fmt.Errorf("erreur lors du scan des produits : %v", err)
        }
        product.Photos = photos
        products = append(products, product)
    }

    return products, nil
}

func (r *ProductRepository) SearchProducts(searchTerm string) ([]models.Product, error) {
    query := `
        SELECT id, nom, prix, stock, etat, photos, categorie_id,
               localisation, description, nombre_vues, disponible,
               marque, modele, created_at, updated_at
        FROM produits
        WHERE 
            nom ILIKE $1 OR 
            description ILIKE $1 OR 
            marque ILIKE $1 OR 
            modele ILIKE $1
        ORDER BY 
            CASE 
                WHEN nom ILIKE $2 THEN 1
                WHEN marque ILIKE $2 THEN 2
                WHEN modele ILIKE $2 THEN 3
                WHEN description ILIKE $2 THEN 4
                ELSE 5
            END,
            nombre_vues DESC
        LIMIT 20`

    searchPattern := "%" + searchTerm + "%"
    exactPattern := searchTerm

    rows, err := r.db.Query(query, searchPattern, exactPattern)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la recherche des produits : %v", err)
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
            return nil, fmt.Errorf("erreur lors du scan des produits : %v", err)
        }
        product.Photos = photos
        products = append(products, product)
    }

    return products, nil
}
