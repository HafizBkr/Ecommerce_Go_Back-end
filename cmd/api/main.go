package main

import (
	"log"
	"net/http"
	"os"

	"ecommerce-api/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}

	DATABASE_URL := os.Getenv("DNS_LINK")
	if DATABASE_URL == "" {
		log.Fatal("La chaîne de connexion à la base de données est manquante.")
	}
	config.InitDatabase(DATABASE_URL)

	r := chi.NewRouter()
	// Utiliser les middleware pour la gestion des requêtes
	r.Use(middleware.Logger)     // Journalisation des requêtes
	r.Use(middleware.Recoverer)  // Récupération en cas d'erreurs
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Durée maximale de mise en cache
	}))

	// Définir une route de test
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// Démarrer le serveur HTTP
	log.Println("Démarrage du serveur sur le port 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur : %v", err)
	}
}
