// admin/handler.go
package admin

import (
    "encoding/json"
    "net/http"
    "time"
    "ecommerce-api/models"
    "github.com/golang-jwt/jwt"
)

type AdminHandler struct {
    repo *AdminRepository
}

func NewAdminHandler(repo *AdminRepository) *AdminHandler {
    return &AdminHandler{repo: repo}
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

// HandleAdminRegister gère l'enregistrement d'un nouvel administrateur
// HandleAdminRegister gère l'enregistrement d'un nouvel administrateur
func (h *AdminHandler) HandleAdminRegister(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Validation basique
    if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }

    // Création de l'admin
    admin := models.User{
        Email:       req.Email,
        PasswordHash: req.Password, // Le hash sera généré dans le repository
        FirstName:   req.FirstName,
        LastName:    req.LastName,
        IsAdmin:     true,
        Status:      "active",
    }

    err := h.repo.CreateAdmin(admin)
    if err != nil {
        http.Error(w, "Error creating admin", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Admin created successfully",
        "admin": map[string]interface{}{
            "email":      admin.Email,
            "first_name": admin.FirstName,
            "last_name":  admin.LastName,
            "is_admin":   true,
        },
    })
}
// HandleAdminLogin gère la connexion d'un administrateur
func (h *AdminHandler) HandleAdminLogin(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    admin, err := h.repo.LoginAdmin(req.Email, req.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Générer un JWT pour l'admin
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  admin.ID,
        "email":    admin.Email,
        "is_admin": true,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte("ADMIN_SECRET_KEY"))
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "token": tokenString,
        "admin": map[string]interface{}{
            "id":         admin.ID,
            "email":      admin.Email,
            "first_name": admin.FirstName,
            "last_name":  admin.LastName,
            "is_admin":   true,
        },
    })
}