package categories

import (
	"database/sql"
	"ecommerce-api/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(category models.Category) error {
	// Utilisation du UUID pour générer un ID unique
	query := `INSERT INTO categories (nom, nombre_produits, statut, created_at, updated_at) 
		      VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, category.Nom, category.NombreProduits, category.Statut, time.Now(), time.Now()).Scan(&category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	query := `SELECT id, nom, nombre_produits, statut, created_at, updated_at FROM categories`
	categories := []models.Category{}
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("error getting categories: %v", err)
	}
	return categories, nil
}

func (repo *CategoryRepository) GetCategoryByID(id string) (*models.Category, error) {
	// Exemple de requête pour récupérer la catégorie par son ID
	var category models.Category
	err := repo.db.QueryRow("SELECT id, nom, nombre_produits, statut, created_at, updated_at FROM categories WHERE id = $1", id).Scan(
		&category.ID,
		&category.Nom,
		&category.NombreProduits,
		&category.Statut,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Aucune catégorie trouvée avec l'ID %s", id)
		}
		return nil, fmt.Errorf("Erreur lors de la récupération de la catégorie : %v", err)
	}
	return &category, nil
}


func (r *CategoryRepository) UpdateCategory(category models.Category) error {
	// Validation du format UUID
	_, err := uuid.Parse(category.ID)
	if err != nil {
		return fmt.Errorf("ID invalide : %v", err)
	}

	query := `
        UPDATE categories 
        SET nom = $1, nombre_produits = $2, 
            statut = $3, updated_at = $4
        WHERE id = $5`
    
	_, err = r.db.Exec(query, category.Nom, category.NombreProduits, category.Statut, time.Now(), category.ID)
	if err != nil {
		return fmt.Errorf("échec de la mise à jour de la catégorie : %v", err)
	}

	return nil
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	// Validation du format UUID
	_, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("ID invalide : %v", err)
	}

	query := `DELETE FROM categories WHERE id = $1`

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("échec de la suppression de la catégorie : %v", err)
	}

	return nil
}
