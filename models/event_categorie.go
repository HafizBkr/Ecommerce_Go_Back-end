package models

import "time"

type EventCategory struct {
    ID        string    `db:"id" json:"id"`
    Label     string    `db:"label" json:"label"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}