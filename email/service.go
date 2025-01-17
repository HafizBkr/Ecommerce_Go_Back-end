// email/service.go
package email

import (
    "bytes"
    "ecommerce-api/models"
    "fmt"
    "html/template"
    "net/smtp"
)

type Config struct {
    Host     string
    Port     string
    Username string
    Password string
    FromName string
    FromEmail string
}

type Service struct {
    config Config
}

func NewEmailService(config Config) *Service {
    return &Service{
        config: config,
    }
}

// EnvoyerEmailConfirmationCommande envoie un email de confirmation après une commande
func (s *Service) EnvoyerEmailConfirmationCommande(commande *models.Commande, email string) error {
    // Préparer le template HTML
    htmlBody, err := s.renderTemplate(emailConfirmationTemplate, map[string]interface{}{
        "NumeroCommande": commande.NumeroCommande,
        "MontantTotal":   fmt.Sprintf("%.2f €", commande.MontantTotal),
        "Date":           commande.CreatedAt.Format("02/01/2006 15:04"),
        "Status":         commande.Status,
    })
    if err != nil {
        return fmt.Errorf("erreur lors du rendu du template: %v", err)
    }

    // Configurer les en-têtes de l'email
    headers := make(map[string]string)
    headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
    headers["To"] = email
    headers["Subject"] = fmt.Sprintf("Confirmation de votre commande %s", commande.NumeroCommande)
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/html; charset=UTF-8"

    // Construire le message
    message := ""
    for k, v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + htmlBody

    // Configurer l'authentification SMTP
    auth := smtp.PlainAuth(
        "",
        s.config.Username,
        s.config.Password,
        s.config.Host,
    )

    // Envoyer l'email
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

// renderTemplate rend un template HTML avec les données fournies
func (s *Service) renderTemplate(templateText string, data interface{}) (string, error) {
    t, err := template.New("email").Parse(templateText)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    if err := t.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}

// Template HTML pour l'email de confirmation
const emailConfirmationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Confirmation de commande</title>
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
        <h1>Confirmation de commande</h1>
    </div>
    
    <div class="content">
        <p>Merci pour votre commande !</p>
        
        <div class="details">
            <h3>Détails de votre commande :</h3>
            <p><strong>Numéro de commande :</strong> {{.NumeroCommande}}</p>
            <p><strong>Montant total :</strong> {{.MontantTotal}}</p>
            <p><strong>Date :</strong> {{.Date}}</p>
            <p><strong>Statut :</strong> {{.Status}}</p>
        </div>
        
        <p>Nous vous informerons par email lorsque votre commande sera expédiée.</p>
    </div>
    
    <div class="footer">
        <p>Cet email a été envoyé automatiquement, merci de ne pas y répondre.</p>
        <p>© 2024 Votre E-commerce. Tous droits réservés.</p>
    </div>
</body>
</html>
`