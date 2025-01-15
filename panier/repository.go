package panier

import (
	"ecommerce-api/models"
	"fmt"
	"log" // Le package log standard de Go

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
    db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) AjouterProduitAuPanier(googleID, produitID string) error {
    _, err := r.db.Exec(`
        INSERT INTO panier (user_id, produit_id, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        ON CONFLICT (user_id, produit_id) DO NOTHING`,
        googleID, produitID)
    if err != nil {
        return fmt.Errorf("impossible d'ajouter le produit au panier: %v", err)
    }
    return nil
}


func (r *Repository) ObtenirPanierParUserID(userID string) ([]*models.ProduitSouhait, error) {
    rows, err := r.db.Queryx(`
        SELECT pr.id, pr.nom, pr.prix, pr.marque, pr.photos
        FROM panier p
        JOIN produits pr ON p.produit_id = pr.id
        WHERE p.user_id = $1`, userID)

    if err != nil {
        log.Printf("Error executing query: %v", err)
        return nil, err
    }
    defer rows.Close()

    var produits []*models.ProduitSouhait
    for rows.Next() {
        var produit models.ProduitSouhait
        var photos pq.StringArray // Déclare un tableau de chaînes pour les photos

        // Scan des colonnes
        if err := rows.Scan(&produit.ID, &produit.Nom, &produit.Prix, &produit.Marque, &photos); err != nil {
            log.Printf("Error scanning row: %v", err)
            return nil, err
        }

        // Récupérer la première photo
        if len(photos) > 0 {
            produit.Photo = photos[0] // Assigner la première photo à la structure
        }

        produits = append(produits, &produit)
    }

    // Vérification si aucun produit n'a été trouvé
    if len(produits) == 0 {
        return nil, fmt.Errorf("aucun produit trouvé dans le panier")
    }

    return produits, nil
}


func (r *Repository) EnleverDuPanier(userID, produitID string) error {
    // Exécution de la requête de suppression
    result, err := r.db.Exec(`
        DELETE FROM panier
        WHERE user_id = $1 AND produit_id = $2`,
        userID, produitID)

    if err != nil {
        log.Printf("Error executing delete query: %v", err)
        return err
    }

    // Vérification si une ligne a été supprimée
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error getting rows affected: %v", err)
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("aucun produit trouvé dans le panier pour l'utilisateur %s", userID)
    }

    return nil
}


