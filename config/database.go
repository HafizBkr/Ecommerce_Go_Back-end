package config

import (
	"log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB est une variable globale qui contiendra la connexion à la base de données.
var DB *sqlx.DB

// InitDatabase initialise la connexion à la base de données et configure le pool de connexions.
func InitDatabase(postgresDNS string) {
	// Connexion à la base de données PostgreSQL avec les informations de connexion fournies.
	var err error
	DB, err = sqlx.Connect("postgres", postgresDNS)
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données : %v", err)
	}

	// Test de la connexion pour s'assurer que tout fonctionne correctement.
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données : %v", err)
	}

	// Si tout est ok, on affiche un message de succès.
	log.Println("Connexion à la base de données PostgreSQL réussie.")
}

// Close ferme la connexion à la base de données.
func Close() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Erreur lors de la fermeture de la connexion à la base de données : %v", err)
	}
	log.Println("Connexion à la base de données fermée.")
}