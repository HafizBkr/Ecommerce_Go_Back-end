package models

import (
    "time"
)

type User struct {
    ID               int       `json:"id" db:"id"`
    GoogleID         string    `json:"google_id" db:"google_id"`  // ID Google
    Email            string    `json:"email" db:"email"`
    PasswordHash     string    `json:"-" db:"password_hash"`
    FirstName        string    `json:"first_name" db:"first_name"`
    LastName         string    `json:"last_name" db:"last_name"`
    IsAdmin          bool      `json:"is_admin" db:"is_admin"`
    Points           int       `json:"points" db:"points"`
    LastLogin        time.Time `json:"last_login" db:"last_login"`
    Status           string    `json:"status" db:"status"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
    Address          string    `json:"address" db:"address"`
    PhoneNumber      string    `json:"phone_number" db:"phone_number"`
    ResidenceCity    string    `json:"residence_city" db:"residence_city"`
    ResidenceCountry string    `json:"residence_country" db:"residence_country"`
}
