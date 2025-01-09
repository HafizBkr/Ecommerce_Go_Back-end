package main

import (
	"ecommerce-api/admin"
	"ecommerce-api/categories"
	"ecommerce-api/config"
	"ecommerce-api/googleauth"
	middlewares "ecommerce-api/middleware"
	"ecommerce-api/repository"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Charger les variables d'environnement
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Vérification des variables d'environnement nécessaires
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment variables")
	}
	DATABASE_URL := os.Getenv("DNS_LINK")
	if DATABASE_URL == "" {
		log.Fatal("La chaîne de connexion à la base de données est manquante.")
	}

	// Initialisation de la base de données
	config.InitDatabase(DATABASE_URL)

	// Création du routeur
	r := chi.NewRouter()

	// Configuration du CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Serve les fichiers statiques
	staticDir := "./static"
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Page d'accueil
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/index.html")
	})

	// Initialisation du middleware GoogleAuth et du gestionnaire
	googleAuthMiddleware := middlewares.GoogleAuthMiddleware
	authMiddleware := middlewares.AuthMiddleware
	AdminMiddleware := admin.AdminAuthMiddleware
	// AdminMiddleware := admin.AdminAuthMiddleware
	userRepo := repository.NewUserRepository(config.DB)
	authHandler := googleauth.NewGoogleAuthHandler(userRepo)
	adminRepo := admin.NewAdminRepository(config.DB)
	adminHandler := admin.NewAdminHandler(adminRepo)
	categoryRepo := categories.NewCategoryRepository(config.DB)
	categoryHandler :=categories.NewCategoryHandler(categoryRepo)
	// Appliquer le middleware sur l'endpoint "/complete-profile"
	// Authentification Google - routes de callback
	r.Get("/oauth-test", googleauth.HandleOAuthRedirect)
	r.Get("/auth/callback", googleauth.HandleAuthCallback)
	//Route pour completer le profile du user classique
	r.Route("/complete-profile", func(r chi.Router) {
		r.Use(googleAuthMiddleware)                    // Appliquer le middleware d'authentification
		r.Post("/", authHandler.HandleCompleteProfile) // Associer la méthode POST à la fonction de gestion
	})
	//Route pour recuperer les info du user
	r.Route("/user/info", func(r chi.Router) {
		r.Use(authMiddleware)                  // Appliquer le middleware d'authentification
		r.Get("/", authHandler.GetUserHandler) // Lier le gestionnaire pour récupérer les infos de l'utilisateur
	})
	//Route pour gerer l'authnetification de l'admin
	
	r.Route("/admin", func(r chi.Router) {
		r.Post("/register", adminHandler.HandleAdminRegister)
		r.Post("/login", adminHandler.HandleAdminLogin)
	})

	//Route pour cree les categorie Uniquement pour l'admin
	r.Route("/categories", func(r chi.Router) {
		r.Use(AdminMiddleware) // Authentification pour les administrateurs
		r.Post("/", categoryHandler.HandleCreateCategory)
		r.Put("/{id}", categoryHandler.HandleUpdateCategory)
		r.Delete("/{id}", categoryHandler.HandleDeleteCategory)
	})

	r.Get("/categories/{id}", categoryHandler.HandleGetCategory)
	r.Get("/categories", categoryHandler.HandleGetAllCategories)

	// Démarrage du serveur
	server := http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", port),
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	log.Println("Server started on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
