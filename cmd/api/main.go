package main

import (
	"ecommerce-api/email"
	"ecommerce-api/order"
	"ecommerce-api/admin"
	"ecommerce-api/categories"
	"ecommerce-api/config"
	"ecommerce-api/events"
	events_category "ecommerce-api/events_categories"
	"ecommerce-api/googleauth"
	middlewares "ecommerce-api/middleware"
	"ecommerce-api/panier"
	"ecommerce-api/products"
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


	  

		emailConfig := email.Config{
			Host:      "smtp.gmail.com",  // ou votre serveur SMTP
			Port:      "587",
			Username:  "mongodb200@gmail.com",
			Password:  "vshq dmbi skwa fnlz", // Utilisez un mot de passe d'application pour Gmail
			FromName:  "Votre E-commerce",
			FromEmail: "mongodb200@gmail.com",
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
	userRepo := repository.NewUserRepository(config.DB)
	authHandler := googleauth.NewGoogleAuthHandler(userRepo)
	adminRepo := admin.NewAdminRepository(config.DB)
	adminHandler := admin.NewAdminHandler(adminRepo)
	categoryRepo := categories.NewCategoryRepository(config.DB)
	categoryHandler := categories.NewCategoryHandler(categoryRepo)
	productRepo := products.NewProductRepository(config.DB)
	productHandler := products.NewProductHandler(productRepo)
	eventCategoriesRepo :=events_category.NewEventCategoryRepository(config.DB)
	eventCategoriesHandler:=events_category.NewEventCategoryHandler(eventCategoriesRepo)
	eventRepo:= events.NewEventRepository(config.DB)
	eventHanlder := events.NewEventHandler(eventRepo)
	panierRepo := panier.NewRepository(config.DB)
    panierHandler := panier.NewPanierHandler(panierRepo)

	emailService:=email.NewEmailService(emailConfig)
	commandeRepo:= order.NewRepository(config.DB)
	CommandeHandler :=order.NewHandler(commandeRepo ,emailService)





	// Appliquer le middleware sur l'endpoint "/complete-profile"
	// Authentification Google - routes de callback
	r.Get("/oauth-test", googleauth.HandleOAuthRedirect)
	r.Get("/auth/callback", authHandler.HandleAuthCallback)
	//Route pour completer le profile du user classique
	r.Route("/complete-profile", func(r chi.Router) {
		r.Use(googleAuthMiddleware)                    // Appliquer le middleware d'authentification
		r.Post("/", authHandler.HandleCompleteProfile) // Associer la méthode POST à la fonction de gestion
	})
	//Route pour recuperer les info du user 
	r.Route("/user/info", func(r chi.Router) {
		r.Use(authMiddleware)                  // Appliquer le middleware d'authentification Pour le login du user
		r.Get("/", authHandler.GetUserHandler) // Lier le gestionnaire pour récupérer les infos de l'utilisateur
	})
	//Route pour gerer l'authnetification de l'admin

	r.Route("/admin", func(r chi.Router) {
		r.Post("/register", adminHandler.HandleAdminRegister)
		r.Post("/login", adminHandler.HandleAdminLogin)
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", categoryHandler.HandleGetAllCategories)    // Obtenir toutes les catégories
		r.Get("/{id}", categoryHandler.HandleGetCategoryByID) // Obtenir une catégorie par ID

		r.With(AdminMiddleware).Route("/", func(r chi.Router) {
			r.Post("/", categoryHandler.HandleCreateCategory)       // Créer une catégorie
			r.Put("/{id}", categoryHandler.HandleUpdateCategory)    // Mettre à jour une catégorie
			r.Delete("/{id}", categoryHandler.HandleDeleteCategory) // Supprimer une catégorie
		})
	})
	r.Route("/products", func(r chi.Router) {
		r.With(AdminMiddleware).Route("/", func(r chi.Router) {
			r.Post("/", productHandler.HandleCreateProduct)
			r.Delete("/{id}", productHandler.HandleDeleteProduct)
			r.Put("/{id}", productHandler.HandleUpdateProduct)
		})
		r.Get("/{id}", productHandler.HandleGetProductByID)
		r.Get("/", productHandler.HandleGetAllProducts)
		r.Get("/by-category/{categoryID}", productHandler.HandleGetProductsByCategory)
		r.Get("/filter", productHandler.HandleFilterProducts)
		r.Get("/search", productHandler.HandleSearchProducts)
	})
	r.Route("/event-categories", func(r chi.Router) {
		r.With(AdminMiddleware).Route("/", func(r chi.Router) {
			r.Post("/", eventCategoriesHandler.HandleCreateEventCategory)
			r.Put("/{id}", eventCategoriesHandler.HandleUpdateEventCategory)
			r.Delete("/{id}", eventCategoriesHandler.HandleDeleteEventCategory)
		})
		r.Get("/", eventCategoriesHandler.HandleGetAllEventCategories)
		r.Get("/{id}", eventCategoriesHandler.HandleGetEventCategoryByID)
	})
	r.Route("/events", func(r chi.Router) {
		r.With(AdminMiddleware).Route("/", func(r chi.Router) {
			r.Post("/", eventHanlder.HandleCreateEvent)
			r.Put("/{id}", eventHanlder.HandleUpdateEvent)
			r.Delete("/{id}", eventHanlder.HandleDeleteEvent)
		})
		r.Get("/", eventHanlder.HandleGetAllEvents)
		r.Get("/{id}", eventHanlder.HandleGetEventByID)
		r.Get("/category/{id}", eventHanlder.HandleGetEventsByCategoryID)
	})

	r.Route("/liste-souhaits", func(r chi.Router) {
		r.Use(authMiddleware) // Middleware d'authentification
		r.Get("/", panierHandler.HandleAfficherPanier)
		r.Post("/ajouter", panierHandler.HandleAjouterProduit)
		r.Delete("/enlever", panierHandler.HandleEnleverDuPanier)
	})

	r.Route("/commandes", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", CommandeHandler.HandleCreerCommande)
		r.Post("/ticket",CommandeHandler.HandleCreerCommandeTicket)
		r.Get("/", CommandeHandler.HandleListerCommandes)
	})
	
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
