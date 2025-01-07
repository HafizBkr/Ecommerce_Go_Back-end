package models

import (
	"time"
)

type User struct {
	ID               int       `json:"id" db:"id"`
	Email            string    `json:"email" db:"email"`
	PasswordHash     string    `json:"-" db:"password_hash"`  // Ne pas exposer le mot de passe
	FirstName        string    `json:"first_name" db:"first_name"`
	LastName         string    `json:"last_name" db:"last_name"`
	IsAdmin          bool      `json:"is_admin" db:"is_admin"`
	Points           int       `json:"points" db:"points"`
	LastLogin        time.Time `json:"last_login" db:"last_login"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	Address          string    `json:"address" db:"address"`  // Adresse complète
	PhoneNumber      string    `json:"phone_number" db:"phone_number"`  // Numéro de téléphone
	ResidenceCity    string    `json:"residence_city" db:"residence_city"`  // Ville de résidence
	ResidenceCountry string    `json:"residence_country" db:"residence_country"`  // Pays de résidence
}
