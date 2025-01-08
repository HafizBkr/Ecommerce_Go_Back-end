package repository

import (
	"database/sql"
	"ecommerce-api/models"
	"fmt"
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
func (r *UserRepository) GetUserByGoogleID(googleID string) (*models.User, error) {
    query := `
        SELECT id, email, first_name, last_name, is_admin, points, 
               COALESCE(last_login, NOW()), status, created_at, updated_at,
               COALESCE(address, ''), COALESCE(phone_number, ''), 
               COALESCE(residence_city, ''), COALESCE(residence_country, ''),
               google_id
        FROM users 
        WHERE google_id = $1`

    user := &models.User{}
    err := r.db.QueryRow(query, googleID).Scan(
        &user.ID,
        &user.Email,
        &user.FirstName,
        &user.LastName,
        &user.IsAdmin,
        &user.Points,
        &user.LastLogin,
        &user.Status,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.Address,
        &user.PhoneNumber,
        &user.ResidenceCity,
        &user.ResidenceCountry,
        &user.GoogleID,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, err
    }
    return user, nil
}
// CreateUser insère un nouvel utilisateur dans la base de données
// CreateUser mise à jour pour inclure GoogleID
func (repo *UserRepository) CreateUser(user models.User) error {
    query := `
        INSERT INTO users (
            email, password_hash, first_name, last_name, is_admin, points, 
            status, created_at, updated_at, address, phone_number, 
            residence_city, residence_country, google_id
        )
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        )
    `
    _, err := repo.db.Exec(
        query,
        user.Email,
        user.PasswordHash,
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
        user.GoogleID,
    )
    return err
}

// GetUserByEmail récupère un utilisateur par son email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
    query := `
        SELECT id, email, first_name, last_name, is_admin, points, 
               COALESCE(last_login, NOW()), status, created_at, updated_at,
               COALESCE(address, ''), COALESCE(phone_number, ''), 
               COALESCE(residence_city, ''), COALESCE(residence_country, ''),
               COALESCE(google_id, '')
        FROM users 
        WHERE email = $1`

    user := &models.User{}
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Email,
        &user.FirstName,
        &user.LastName,
        &user.IsAdmin,
        &user.Points,
        &user.LastLogin,
        &user.Status,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.Address,
        &user.PhoneNumber,
        &user.ResidenceCity,
        &user.ResidenceCountry,
        &user.GoogleID,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, err
    }
    return user, nil
}


// UpdateUser met à jour les informations utilisateur
func (r *UserRepository) UpdateUser(user models.User) error {
    query := `
        UPDATE users 
        SET 
            first_name = $1,
            last_name = $2,
            address = $3,
            phone_number = $4,
            residence_city = $5,
            residence_country = $6,
            google_id = $7,
            updated_at = NOW()
        WHERE email = $8
        RETURNING id`

    err := r.db.QueryRow(
        query,
        user.FirstName,
        user.LastName,
        user.Address,
        user.PhoneNumber,
        user.ResidenceCity,
        user.ResidenceCountry,
        user.GoogleID,
        user.Email,
    ).Scan(&user.ID)

    if err != nil {
        return fmt.Errorf("error updating user: %v", err)
    }
    return nil
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

// GetUserByID récupère un utilisateur par son ID
// GetUserByID récupère un utilisateur par son ID
func (r *UserRepository) GetUserByID(userID string) (*models.User, error) {
	query := `
        SELECT id, email, first_name, last_name, is_admin, points, 
               COALESCE(last_login, NOW()), status, created_at, updated_at,
               COALESCE(address, ''), COALESCE(phone_number, ''), 
               COALESCE(residence_city, ''), COALESCE(residence_country, ''),
               google_id
        FROM users 
        WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.IsAdmin,
		&user.Points,
		&user.LastLogin,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Address,
		&user.PhoneNumber,
		&user.ResidenceCity,
		&user.ResidenceCountry,
		&user.GoogleID,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

