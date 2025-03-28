<<<<<<< HEAD
# Back


go mod tidy pour installer les dependances 
utiliser air si air est installer pour executer le projet 
acceder au cmd/api en suite go run main.go.
=======
# Projet Backend

## Description
Ce projet est un backend développé en **Go** (Golang) pour une application de gestion d'événements et de commandes. Il comprend des fonctionnalités d'authentification, de gestion des utilisateurs, des événements, des commandes, des produits et plus encore.

## Structure du projet

```
.
├── admin/               # Gestion des administrateurs
├── categories/          # Gestion des catégories
├── cmd/api/            # Point d'entrée principal de l'API
│   ├── main.go         # Démarrage du serveur
│   ├── static/        
│   │   └── index.html  # Fichier statique
│   └── tmp/            # Fichiers temporaires
├── config/             # Configuration de la base de données
├── email/              # Gestion des emails
├── events/             # Gestion des événements
├── events_categories/  # Association événements - catégories
├── googleauth/         # Authentification Google et JWT
├── middleware/         # Middleware d'authentification
├── migrations/         # Scripts de migration de la base de données
├── models/             # Modèles de la base de données
├── order/              # Gestion des commandes
├── panier/             # Gestion du panier
├── pkg/                # Packages utilitaires
├── products/           # Gestion des produits
├── repository/         # Repository pour les intégrations externes
├── scripts/            # Scripts divers
├── search/             # Fonctionnalités de recherche
├── tmp/                # Logs et fichiers temporaires
└── user/               # Gestion des utilisateurs
```

## Installation
### Prérequis
- **Go** (>=1.18)
- **PostgreSQL** pour la base de données

### Étapes d'installation
1. Cloner le dépôt :
   ```sh
   git clone https://github.com/GestionDeProjetESGIS/Backend.git
   cd Backend
   ```
2. Installer les dépendances :
   ```sh
   go mod tidy
   ```
3. Configurer la base de données dans `config/database.go`
4. Exécuter les migrations :
   ```sh
   go run migrations/schema.sql
   ```
5. Démarrer le serveur :
   ```sh
   go run cmd/api/main.go
   ```

## API
Ce projet expose une API RESTful. La documentation des endpoints est disponible via **Postman** dans le dossier `postman/`.

### Routes principales
| Méthode | Endpoint             | Description                      |
|---------|----------------------|----------------------------------|
| POST    | `/auth/login`        | Connexion utilisateur           |
| POST    | `/auth/register`     | Inscription utilisateur         |
| GET     | `/events`            | Liste des événements            |
| GET     | `/products`          | Liste des produits              |
| POST    | `/order`             | Passer une commande             |

## Contribuer
1. Forker le projet
2. Créer une branche feature
3. Faire un commit et un push
4. Créer une pull request

## License
Ce projet est sous licence **MIT**.

---
Auteur : **Boukari Hafiz**

>>>>>>> e17b96c (commit for the ecomerce project)
