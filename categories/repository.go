package categories

import (
	"ecommerce-api/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(category models.Category) error {
    query := `INSERT INTO categories (nom, nombre_produits, statut, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
    
    err := r.db.QueryRow(query, category.Nom, category.NombreProduits, category.Statut, time.Now(), time.Now()).Scan(&category.ID)
    if err != nil {
        return err
    }

    return nil
}


func (r *CategoryRepository) UpdateCategory(category models.Category) error {
	query := `
        UPDATE categories 
        SET 
            nom = $1,
            nombre_produits = $2,
            statut = $3,
            updated_at = NOW()
        WHERE id = $4`

	result, err := r.db.Exec(
		query,
		category.Nom, // Utilisez pq.Array ici aussi
		category.NombreProduits,
		category.Statut,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating category: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *CategoryRepository) GetCategoryByID(id int) (*models.Category, error) {
	query := `
        SELECT id, nom, nombre_produits, 
               statut, created_at, updated_at
        FROM categories 
        WHERE id = $1`

	category := &models.Category{}
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Nom, // Correct use of pq.Array to scan into []string
		&category.NombreProduits,
		&category.Statut,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}

	return category, nil
}

func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	query := `
        SELECT id, nom, nombre_produits, 
               statut, created_at, updated_at
        FROM categories`

	categories := []models.Category{}
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("error getting categories: %v", err)
	}

	return categories, nil
}

func (r *CategoryRepository) DeleteCategory(id int) error {
    // Vérifiez d'abord si la catégorie existe avant de la supprimer
    var exists bool
    queryCheck := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`
    err := r.db.QueryRow(queryCheck, id).Scan(&exists)
    if err != nil {
        return fmt.Errorf("error checking if category exists: %v", err)
    }
    if !exists {
        return fmt.Errorf("category with id %d not found", id)
    }

    // Suppression de la catégorie
    query := `DELETE FROM categories WHERE id = $1`
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("error deleting category: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %v", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("no category deleted, it might have already been deleted")
    }

    return nil
}
