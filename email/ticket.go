// email/service.go
package email

import (
	"ecommerce-api/models"
	"fmt"
	"net/smtp"
)



func (s *Service) EnvoyerEmailConfirmationTicket(ticket *models.TicketOrder, email string) error {
    // Préparer le template HTML
    htmlBody, err := s.renderTemplate(emailConfirmationTicketTemplate, map[string]interface{}{
        "NumeroCommande": ticket.NumeroCommande,
        "EventTitle":     ticket.EventTitle,
        "Quantity":       ticket.Quantity,
        "MontantTotal":   fmt.Sprintf("%.2f CFA", ticket.PrixTotal),
        "Date":           ticket.StartDate.Format("02/01/2006"),
        "Heure":          ticket.StartTime,
        "Status":         ticket.Status,
    })
    if err != nil {
        return fmt.Errorf("erreur lors du rendu du template: %v", err)
    }

    headers := make(map[string]string)
    headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
    headers["To"] = email
    headers["Subject"] = fmt.Sprintf("Confirmation de vos tickets - %s", ticket.NumeroCommande)
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/html; charset=UTF-8"

    message := ""
    for k, v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + htmlBody

    auth := smtp.PlainAuth(
        "",
        s.config.Username,
        s.config.Password,
        s.config.Host,
    )

    err = smtp.SendMail(
        fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
        auth,
        s.config.FromEmail,
        []string{email},
        []byte(message),
    )
    if err != nil {
        return fmt.Errorf("erreur lors de l'envoi de l'email: %v", err)
    }

    return nil
}

const emailConfirmationTicketTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Confirmation de tickets</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background-color: #4CAF50;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 5px;
        }
        .content {
            padding: 20px;
            background-color: #f9f9f9;
            border-radius: 5px;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            margin-top: 20px;
            padding: 20px;
            font-size: 12px;
            color: #666;
        }
        .details {
            margin: 20px 0;
            padding: 15px;
            background-color: white;
            border-radius: 5px;
            border: 1px solid #ddd;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Confirmation de vos tickets</h1>
    </div>
    
    <div class="content">
        <p>Merci pour votre réservation !</p>
        
        <div class="details">
            <h3>Détails de votre réservation :</h3>
            <p><strong>Numéro de commande :</strong> {{.NumeroCommande}}</p>
            <p><strong>Événement :</strong> {{.EventTitle}}</p>
            <p><strong>Nombre de tickets :</strong> {{.Quantity}}</p>
            <p><strong>Montant total :</strong> {{.MontantTotal}}</p>
            <p><strong>Date :</strong> {{.Date}}</p>
            <p><strong>Heure :</strong> {{.Heure}}</p>
            <p><strong>Statut :</strong> {{.Status}}</p>
        </div>
        
        <p>Conservez précieusement ce mail, il servira de justificatif lors de l'événement.</p>
    </div>
    
    <div class="footer">
        <p>Cet email a été envoyé automatiquement, merci de ne pas y répondre.</p>
        <p>© 2024 Événements Lomé. Tous droits réservés.</p>
    </div>
</body>
</html>
`