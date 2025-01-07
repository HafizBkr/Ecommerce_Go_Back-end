package repository

import (
	"ecommerce-api/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// UserRepository gère les interactions avec la base de données pour les utilisateurs
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository crée une nouvelle instance de UserRepository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser insère un nouvel utilisateur dans la base de données
func (repo *UserRepository) CreateUser(user models.User) error {
	query := `
		INSERT INTO users (
			email, password_hash, first_name, last_name, is_admin, points, status, created_at, updated_at, 
			address, phone_number, residence_city, residence_country
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`
	_, err := repo.db.Exec(
		query,
		user.Email,
		user.PasswordHash, // Ajouter un mot de passe par défaut ou gérer autrement
		user.FirstName,
		user.LastName,
		user.IsAdmin,
		user.Points,
		user.Status,
		time.Now(),
		time.Now(),
		user.Address,
		user.PhoneNumber,
		user.ResidenceCity,
		user.ResidenceCountry,
	)
	return err
}

// GetUserByEmail récupère un utilisateur par son email
func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := repo.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser met à jour les informations utilisateur
func (repo *UserRepository) UpdateUser(user models.User) error {
	query := `
		UPDATE users
		SET 
			first_name = $1, 
			last_name = $2, 
			is_admin = $3, 
			points = $4, 
			status = $5, 
			address = $6, 
			phone_number = $7, 
			residence_city = $8, 
			residence_country = $9, 
			updated_at = $10
		WHERE id = $11
	`
	_, err := repo.db.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.IsAdmin,
		user.Points,
		user.Status,
		user.Address,
		user.PhoneNumber,
		user.ResidenceCity,
		user.ResidenceCountry,
		time.Now(),
		user.ID,
	)
	return err
}

// SaveUserProfile met à jour les informations du profil utilisateur
func (repo *UserRepository) SaveUserProfile(email, address, phoneNumber, city, country string) error {
	query := `
		UPDATE users
		SET 
			address = $1,
			phone_number = $2, 
			residence_city = $3, 
			residence_country = $4, 
			updated_at = NOW()
		WHERE email = $5
	`
	_, err := repo.db.Exec(query, address, phoneNumber, city, country, email)
	return err
}
