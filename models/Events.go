// models/event.go
package models

import (
    "time"
)

type Event struct {
    ID             string    `db:"id" json:"id"`
    Title          string    `db:"title" json:"title"`
    Description    string    `db:"description" json:"description"`
    StartDate      time.Time `db:"start_date" json:"start_date"`
    EndDate        time.Time `db:"end_date" json:"end_date"`
    StartTime      string    `db:"start_time" json:"start_time"`
    Price          float64   `db:"price" json:"price"`
    EventTypeID    string    `db:"event_type_id" json:"event_type_id"`
    AvailableSeats int       `db:"available_seats" json:"available_seats"`
    ImageURL       string    `db:"image_url" json:"image_url"`
    Latitude       float64   `db:"latitude" json:"latitude"`
    Longitude      float64   `db:"longitude" json:"longitude"`
    CreatedAt      time.Time `db:"created_at" json:"created_at"`
    UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}