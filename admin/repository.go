// admin/repository.go
package admin

import (
    "database/sql"
    "time"
    "ecommerce-api/models"
    "golang.org/x/crypto/bcrypt"
    "github.com/jmoiron/sqlx"
)

type AdminRepository struct {
    db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
    return &AdminRepository{db: db}
}

func (r *AdminRepository) CreateAdmin(admin models.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.PasswordHash), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    query := `
        INSERT INTO users (
            email, password_hash, first_name, last_name, is_admin,
            status, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, true, 'active', $5, $5)
    `
    _, err = r.db.Exec(
        query,
        admin.Email,
        string(hashedPassword),
        admin.FirstName,
        admin.LastName,
        time.Now(),
    )
    return err
}

func (r *AdminRepository) LoginAdmin(email, password string) (*models.User, error) {
    query := `
        SELECT id, email, password_hash, first_name, last_name, is_admin
        FROM users 
        WHERE email = $1 AND is_admin = true
    `
    
    var admin models.User
    err := r.db.QueryRow(query, email).Scan(
        &admin.ID,
        &admin.Email,
        &admin.PasswordHash,
        &admin.FirstName,
        &admin.LastName,
        &admin.IsAdmin,
    )
    
    if err == sql.ErrNoRows {
        return nil, err
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
    if err != nil {
        return nil, err
    }
    
    return &admin, nil
}