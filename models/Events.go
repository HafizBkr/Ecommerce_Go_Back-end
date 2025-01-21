package models

type Event struct {
    ID             string  `db:"id" json:"id"`
    Title          string  `db:"title" json:"title"`
    Description    string  `db:"description" json:"description"`
    StartDate      string  `db:"start_date" json:"start_date"`  // Type changé en string
    EndDate        string  `db:"end_date" json:"end_date"`      // Type changé en string
    StartTime      string  `db:"start_time" json:"start_time"`  // Type changé en string
    Price          float64 `db:"price" json:"price"`
    EventTypeID    string  `db:"event_type_id" json:"event_type_id"`
    AvailableSeats int     `db:"available_seats" json:"available_seats"`
    ImageURL       string  `db:"image_url" json:"image_url"`
    Latitude       float64 `db:"latitude" json:"latitude"`
    Longitude      float64 `db:"longitude" json:"longitude"`
    CreatedAt      string  `db:"created_at" json:"created_at"`
    UpdatedAt      string  `db:"updated_at" json:"updated_at"`
}
