// order/repository.go
package order

import (
	"database/sql"
	"ecommerce-api/models"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
    db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{db: db}
}
func (r *Repository) CreerCommande(userID string, produits []*models.CommandeProduit) (*models.Commande, error) {
    // Validation des entrées
    if len(produits) == 0 {
        return nil, fmt.Errorf("la commande doit contenir au moins un produit")
    }
    if userID == "" {
        return nil, fmt.Errorf("userID ne peut pas être vide")
    }

    // Générer le numéro de commande
    numeroCommande := fmt.Sprintf("CMD-%s-%s",
        time.Now().Format("20060102"),
        uuid.New().String()[:8])

    var montantTotal float64

    // Démarrer la transaction
    tx, err := r.db.Beginx()
    if err != nil {
        return nil, fmt.Errorf("erreur lors du début de la transaction: %v", err)
    }
    defer tx.Rollback()

    // Préparer la requête pour optimiser les performances
    stmt, err := tx.Preparex(`SELECT nom, prix, stock FROM produits WHERE id = $1 FOR UPDATE`)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la préparation de la requête: %v", err)
    }
    defer stmt.Close()

    var produitsDetails []models.ProduitDetail

    // Vérifier tous les produits avant de faire des modifications
    for _, produit := range produits {
        if produit.Quantite <= 0 {
            return nil, fmt.Errorf("la quantité doit être supérieure à 0 pour le produit %s", produit.ProduitID)
        }

        var nom string
        var prixUnite float64
        var stockDisponible int

        err := stmt.QueryRow(produit.ProduitID).Scan(&nom, &prixUnite, &stockDisponible)
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("le produit %s n'existe pas dans la base de données", produit.ProduitID)
        } else if err != nil {
            return nil, fmt.Errorf("erreur lors de la lecture du produit %s: %v", produit.ProduitID, err)
        }

        if stockDisponible < produit.Quantite {
            return nil, fmt.Errorf("stock insuffisant pour le produit %s (demandé: %d, disponible: %d)",
                nom, produit.Quantite, stockDisponible)
        }

        produit.PrixUnite = prixUnite
        montantTotal += prixUnite * float64(produit.Quantite)

        produitsDetails = append(produitsDetails, models.ProduitDetail{
            Nom:       nom,
            PrixUnite: prixUnite,
            Quantite:  produit.Quantite,
        })
    }

    // Créer la commande
    commande := &models.Commande{
        ID:             uuid.New().String(),
        NumeroCommande: numeroCommande,
        UserID:         userID,
        MontantTotal:   montantTotal,
        Status:         "en_attente",
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
        Produits:       produitsDetails,
    }

    // Insérer la commande
    _, err = tx.NamedExec(`
        INSERT INTO commandes (id, numero_commande, user_id, montant_total, status, created_at, updated_at)
        VALUES (:id, :numero_commande, :user_id, :montant_total, :status, :created_at, :updated_at)`,
        commande)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de l'insertion de la commande: %v", err)
    }

    // Insérer les produits de la commande et mettre à jour les stocks
    for _, produit := range produits {
        // Insérer dans commande_produits
        produit.CommandeID = commande.ID
        _, err = tx.NamedExec(`
            INSERT INTO commande_produits (commande_id, produit_id, quantite, prix_unite)
            VALUES (:commande_id, :produit_id, :quantite, :prix_unite)`,
            produit)
        if err != nil {
            return nil, fmt.Errorf("erreur lors de l'insertion du produit dans la commande: %v", err)
        }

        // Mettre à jour le stock
        result, err := tx.Exec(`
            UPDATE produits 
            SET stock = stock - $1,
                updated_at = NOW()
            WHERE id = $2 AND stock >= $1`,
            produit.Quantite, produit.ProduitID)
        if err != nil {
            return nil, fmt.Errorf("erreur lors de la mise à jour du stock: %v", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            return nil, fmt.Errorf("erreur lors de la vérification de la mise à jour: %v", err)
        }
        if rowsAffected == 0 {
            return nil, fmt.Errorf("impossible de mettre à jour le stock du produit %s", produit.ProduitID)
        }
    }

    // Valider la transaction
    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("erreur lors de la validation de la transaction: %v", err)
    }

    return commande, nil
}
func (r *Repository) GetCommandesByUser(userID string) ([]*models.Commande, error) {
    var commandes []*models.Commande

    query := `
        SELECT id, numero_commande, user_id, montant_total, status, created_at, updated_at
        FROM commandes
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

    if err := r.db.Select(&commandes, query, userID); err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des commandes pour l'utilisateur %s: %v", userID, err)
    }

    return commandes, nil
}
// ListerCommandesParUtilisateur récupère toutes les commandes passées par un utilisateur spécifique.
// Dans repository.go
func (r *Repository) ListerCommandesParUtilisateur(userID string) ([]*models.Commande, error) {
    const query = `
        SELECT 
            c.id AS commande_id, 
            c.numero_commande,
            c.user_id, 
            c.montant_total, 
            c.status, 
            c.created_at,
            c.updated_at,
            p.nom AS produit_nom,
            p.marque AS produit_marque,
            p.modele AS produit_modele,
            cp.prix_unite AS produit_prix,
            cp.quantite AS produit_quantite,
            p.etat AS produit_etat,
            p.localisation AS produit_localisation,
            p.photos AS produit_photos
        FROM commandes c
        LEFT JOIN commande_produits cp ON cp.commande_id = c.id
        LEFT JOIN produits p ON p.id = cp.produit_id
        WHERE c.user_id = $1
        ORDER BY c.created_at DESC
    `

    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des commandes: %w", err)
    }
    defer rows.Close()

    commandesMap := make(map[string]*models.Commande)
    for rows.Next() {
        var (
            commandeID   string
            numeroCommande string
            montantTotal float64
            status       string
            createdAt    time.Time
            updatedAt    time.Time
            produitNom   sql.NullString
            produitMarque sql.NullString
            produitModele sql.NullString
            produitPrix  sql.NullFloat64
            produitQuantite sql.NullInt64
            produitEtat  sql.NullString
            produitLocalisation sql.NullString
            produitPhotos sql.NullString
        )

        err := rows.Scan(
            &commandeID, &numeroCommande, &userID, &montantTotal, &status, &createdAt, &updatedAt,
            &produitNom, &produitMarque, &produitModele, &produitPrix, &produitQuantite, 
            &produitEtat, &produitLocalisation, &produitPhotos,
        )
        if err != nil {
            return nil, fmt.Errorf("erreur lors du scan des commandes: %w", err)
        }

        if _, exists := commandesMap[commandeID]; !exists {
            commandesMap[commandeID] = &models.Commande{
                ID:             commandeID,
                NumeroCommande: numeroCommande,
                UserID:         userID,
                MontantTotal:   montantTotal,
                Status:         status,
                CreatedAt:      createdAt,
                UpdatedAt:      updatedAt,
                Produits:       []models.ProduitDetail{},
            }
        }

        if produitPhotos.Valid {
            var raw string = produitPhotos.String
            if !strings.HasPrefix(raw, "[") { // Vérifie si ce n'est pas un tableau JSON
                raw = fmt.Sprintf(`["%s"]`, strings.Join(strings.Split(raw, ", "), `","`))
            }
            var photos pq.StringArray
            err := json.Unmarshal([]byte(raw), &photos)
            if err != nil {
                return nil, fmt.Errorf("erreur lors de la conversion des photos pour le produit %s: %w", produitNom.String, err)
            }
        
        

            commandesMap[commandeID].Produits = append(commandesMap[commandeID].Produits, models.ProduitDetail{
                Nom:         produitNom.String,
                Model:       produitModele.String,
                PrixUnite:   produitPrix.Float64,
                Quantite:    int(produitQuantite.Int64),
                Etat:        produitEtat.String,
                Localisation: produitLocalisation.String,
                Photos:      photos,
            })
        }
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("erreur lors de l'itération des commandes: %w", err)
    }

    commandes := make([]*models.Commande, 0, len(commandesMap))
    for _, commande := range commandesMap {
        commandes = append(commandes, commande)
    }

    return commandes, nil
}


func (r *Repository) ListerToutesCommandes() ([]*models.CommandeDetail, error) {
    var commandes []*models.CommandeDetail

    query := `
        SELECT 
            c.id, 
            c.numero_commande, 
            u.google_id, 
            u.first_name, 
            u.last_name, 
            u.email, 
            c.montant_total, 
            c.status, 
            c.created_at, 
            c.updated_at
        FROM commandes c
        JOIN users u ON c.user_id = u.google_id  -- Jointure sur google_id
        ORDER BY c.created_at DESC`

    err := r.db.Select(&commandes, query)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des commandes : %w", err)
    }

    return commandes, nil
}

