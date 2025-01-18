// order/repository.go
package order

import (
    "ecommerce-api/models"
    "fmt"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "time"
)

type Repository struct {
    db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) CreerCommande(userID string, produits []*models.CommandeProduit) (*models.Commande, error) {
    numeroCommande := fmt.Sprintf("CMD-%s-%s",
        time.Now().Format("20060102"),
        uuid.New().String()[:5])

    var montantTotal float64

    tx, err := r.db.Beginx()
    if err != nil {
        return nil, fmt.Errorf("erreur lors du début de la transaction: %v", err)
    }
    defer tx.Rollback()

    for _, produit := range produits {
        var prixUnite float64
        var stockDisponible int

        query := `SELECT prix, stock FROM produits WHERE id = $1`
        err := tx.QueryRow(query, produit.ProduitID).Scan(&prixUnite, &stockDisponible)
        if err != nil {
            return nil, fmt.Errorf("erreur lors de la récupération du produit avec l'ID %s : %v", produit.ProduitID, err)
        }

        if produit.Quantite > stockDisponible {
            return nil, fmt.Errorf("le produit %s n'a pas assez de stock disponible", produit.ProduitID)
        }

        produit.PrixUnite = prixUnite
        montantTotal += prixUnite * float64(produit.Quantite)
    }

    commande := &models.Commande{
        ID:             uuid.New().String(),
        NumeroCommande: numeroCommande,
        UserID:         userID,
        MontantTotal:   montantTotal,
        Status:         "en_attente",
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }

    _, err = tx.NamedExec(`
        INSERT INTO commandes (id, numero_commande, user_id, montant_total, status, created_at, updated_at)
        VALUES (:id, :numero_commande, :user_id, :montant_total, :status, :created_at, :updated_at)`,
        commande)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de l'insertion de la commande: %v", err)
    }

    for _, produit := range produits {
        produit.CommandeID = commande.ID
        _, err = tx.NamedExec(`
            INSERT INTO commande_produits (commande_id, produit_id, quantite, prix_unite)
            VALUES (:commande_id, :produit_id, :quantite, :prix_unite)`,
            produit)
        if err != nil {
            return nil, fmt.Errorf("erreur lors de l'insertion des produits: %v", err)
        }

        _, err = tx.Exec(`
            UPDATE produits 
            SET stock = stock - $1,
                updated_at = NOW()
            WHERE id = $2`,
            produit.Quantite, produit.ProduitID)
        if err != nil {
            return nil, fmt.Errorf("erreur lors de la mise à jour du stock: %v", err)
        }
    }

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
func (r *Repository) ListerCommandesParUtilisateur(userID string) ([]*models.Commande, error) {
    const query = `
        SELECT c.id, c.user_id, c.montant_total, c.status, c.created_at
        FROM commandes c
        WHERE c.user_id = $1
        ORDER BY c.created_at DESC
    `
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des commandes: %w", err)
    }
    defer rows.Close()

    var commandes []*models.Commande
    for rows.Next() {
        commande := &models.Commande{}
        if err := rows.Scan(
            &commande.ID,
            &commande.UserID,
            &commande.MontantTotal,
            &commande.Status,
            &commande.CreatedAt,
        ); err != nil {
            return nil, fmt.Errorf("erreur lors de la lecture des données: %w", err)
        }
        commandes = append(commandes, commande)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("erreur lors de l'itération des lignes: %w", err)
    }

    return commandes, nil
}


func (r *Repository) CreerCommandeTicket(userEmail, eventID string, quantity int) (*models.TicketOrder, error) {
    // Générer un numéro unique pour la commande
    numeroCommande := fmt.Sprintf("TICKET-%s-%s",
        time.Now().Format("20060102"),
        uuid.New().String()[:5])

    // Récupérer les détails de l'événement
    var eventDetails models.Event
    err := r.db.Get(&eventDetails, `
        SELECT id, title, start_date, start_time, price
        FROM events
        WHERE id = $1`, eventID)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des détails de l'événement : %v", err)
    }

    // Ignorer toute vérification de l'heure de début et récupérer la valeur telle quelle
    // L'heure de début pourra être invalide, mais elle sera stockée dans la base de données
    // sans modification.

    // Calculer le prix total
    montantTotal := float64(quantity) * eventDetails.Price

    // Créer la structure pour la commande
    ticketOrder := &models.TicketOrder{
        ID:             uuid.New().String(),
        NumeroCommande: numeroCommande,
        UserID:         userEmail,  // Utilisation de l'email utilisateur
        EventID:        eventID,
        EventTitle:     eventDetails.Title,
        StartDate:      eventDetails.StartDate,
        StartTime:      eventDetails.StartTime, // L'heure de début, même invalide, est insérée directement
        Quantity:       quantity,
        PrixTotal:      montantTotal,
        Status:         "en_attente",  // Statut initial de la commande
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }

    // Insérer la commande dans la table `ticket_orders`
    _, err = r.db.NamedExec(`
        INSERT INTO ticket_orders (id, numero_commande, user_id, event_id, event_title, start_date, start_time, quantity, prix_total, status, created_at, updated_at)
        VALUES (:id, :numero_commande, :user_id, :event_id, :event_title, :start_date, :start_time, :quantity, :prix_total, :status, :created_at, :updated_at)`,
        ticketOrder)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de l'insertion de la commande : %v", err)
    }

    return ticketOrder, nil
}


