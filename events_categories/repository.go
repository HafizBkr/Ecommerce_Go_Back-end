// events_category/repository.go
package events_category

import (
    "database/sql"
    "fmt"
    "time"
    "ecommerce-api/models" // Assurez-vous d'ajuster le chemin d'import selon votre structure
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
)

type EventCategoryRepository struct {
    db *sqlx.DB
}

func NewEventCategoryRepository(db *sqlx.DB) *EventCategoryRepository {
    return &EventCategoryRepository{db: db}
}

func (r *EventCategoryRepository) CreateEventCategory(category models.EventCategory) error {
    category.ID = uuid.New().String()
    query := `INSERT INTO event_categories (id, label, created_at, updated_at) 
              VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, category.ID, category.Label, time.Now(), time.Now())
    if err != nil {
        return fmt.Errorf("erreur lors de la création de la catégorie: %v", err)
    }
    return nil
}

func (r *EventCategoryRepository) GetAllEventCategories() ([]models.EventCategory, error) {
    query := `SELECT id, label, created_at, updated_at FROM event_categories`
    categories := []models.EventCategory{}
    err := r.db.Select(&categories, query)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des catégories: %v", err)
    }
    return categories, nil
}

func (r *EventCategoryRepository) GetEventCategoryByID(id string) (*models.EventCategory, error) {
    var category models.EventCategory
    query := `SELECT id, label, created_at, updated_at FROM event_categories WHERE id = $1`
    err := r.db.Get(&category, query, id)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("catégorie non trouvée avec l'ID %s", id)
    }
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération de la catégorie: %v", err)
    }
    return &category, nil
}

func (r *EventCategoryRepository) UpdateEventCategory(category models.EventCategory) error {
    query := `UPDATE event_categories 
              SET label = $1, updated_at = $2
              WHERE id = $3`
    result, err := r.db.Exec(query, category.Label, time.Now(), category.ID)
    if err != nil {
        return fmt.Errorf("erreur lors de la mise à jour: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erreur lors de la vérification de la mise à jour: %v", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("aucune catégorie trouvée avec l'ID %s", category.ID)
    }
    return nil
}

func (r *EventCategoryRepository) DeleteEventCategory(id string) error {
    query := `DELETE FROM event_categories WHERE id = $1`
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("erreur lors de la suppression: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erreur lors de la vérification de la suppression: %v", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("aucune catégorie trouvée avec l'ID %s", id)
    }
    return nil
}
