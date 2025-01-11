package events

import (
    "database/sql"
    "fmt"
    "time"
    "ecommerce-api/models"  // Ajustez le chemin selon votre projet
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
)

type EventRepository struct {
    db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
    return &EventRepository{db: db}
}

func (r *EventRepository) CreateEvent(event models.Event) error {
    event.ID = uuid.New().String()
    query := `
        INSERT INTO events (
            id, title, description, start_date, end_date, 
            start_time, price, event_type_id, available_seats,
            image_url, latitude, longitude, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        )`
    
    _, err := r.db.Exec(
        query,
        event.ID, event.Title, event.Description, event.StartDate,
        event.EndDate, event.StartTime, event.Price, event.EventTypeID,
        event.AvailableSeats, event.ImageURL, event.Latitude,
        event.Longitude, time.Now(), time.Now(),
    )
    
    if err != nil {
        return fmt.Errorf("erreur lors de la création de l'événement: %v", err)
    }
    return nil
}

func (r *EventRepository) GetAllEvents() ([]models.Event, error) {
    var events []models.Event
    query := `SELECT * FROM events ORDER BY start_date`
    
    err := r.db.Select(&events, query)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des événements: %v", err)
    }
    return events, nil
}

func (r *EventRepository) GetEventByID(id string) (*models.Event, error) {
    var event models.Event
    query := `SELECT * FROM events WHERE id = $1`
    
    err := r.db.Get(&event, query, id)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("événement non trouvé")
    }
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération de l'événement: %v", err)
    }
    return &event, nil
}

func (r *EventRepository) GetEventsByCategoryID(categoryID string) ([]models.Event, error) {
    if _, err := uuid.Parse(categoryID); err != nil {
        return nil, fmt.Errorf("ID de catégorie invalide : %v", err)
    }
    var events []models.Event
    query := `SELECT * FROM events WHERE event_type_id = $1 ORDER BY start_date`
    
    err := r.db.Select(&events, query, categoryID)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des événements: %v", err)
    }
    return events, nil
}


func (r *EventRepository) UpdateEvent(event models.Event) error {
    query := `
        UPDATE events SET 
            title = $1,
            description = $2,
            start_date = $3,
            end_date = $4,
            start_time = $5,
            price = $6,
            event_type_id = $7,
            available_seats = $8,
            image_url = $9,
            latitude = $10,
            longitude = $11,
            updated_at = $12
        WHERE id = $13`
    
    result, err := r.db.Exec(
        query,
        event.Title, event.Description, event.StartDate,
        event.EndDate, event.StartTime, event.Price,
        event.EventTypeID, event.AvailableSeats, event.ImageURL,
        event.Latitude, event.Longitude, time.Now(), event.ID,
    )
    
    if err != nil {
        return fmt.Errorf("erreur lors de la mise à jour: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erreur lors de la vérification de la mise à jour: %v", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("aucun événement trouvé avec l'ID %s", event.ID)
    }
    return nil
}

func (r *EventRepository) DeleteEvent(id string) error {
    query := `DELETE FROM events WHERE id = $1`
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("erreur lors de la suppression: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("erreur lors de la vérification de la suppression: %v", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("aucun événement trouvé avec l'ID %s", id)
    }
    return nil
}