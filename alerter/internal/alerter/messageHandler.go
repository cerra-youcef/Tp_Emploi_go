package alerter

import (
	"alerter/internal/helpers"
	"alerter/internal/mailer"
	"alerter/internal/models"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

// Builds and sends emails
func ProcessMessage(msg jetstream.Msg, cfg helpers.Config) {
	var alert models.Alert
	if err := json.Unmarshal(msg.Data(), &alert); err != nil {
		slog.Error("Failed to unmarshal alert", "error", err)
		if err := msg.Nak(); err != nil {
			slog.Error("Failed to NAK message", "error", err)
		}
		return
	}

	data := mailer.TemplateData{
		EventName: alert.Event,
		Date:      alert.Start,
		Location:  alert.Location,
		Changes:   alert.Changes,
	}
	//build mail content
	templateMap := map[string]string{
		"event.created": "templates/created.html",
		"event.deleted": "templates/deleted.html",
		"event.updated": "templates/updated.html",
	}

	// Générer le contenu de l'email
	templatePath, exists := templateMap[alert.Type]
	if !exists {
		slog.Error("unknown alert type: " + alert.Type)
		return
	}

	subject, content, err := mailer.GetEmailContent(templatePath, data, alert.Type)
	if err != nil {
		slog.Error("error while generating email template: " + err.Error())
		return
	}

	for _, resource := range alert.Resources {
		fetchedAlerts, err := helpers.FetchAlertsByResource(cfg.ConfigURL, resource)
		if err != nil {
			slog.Error("failed to get alerts: ", "error", err)
			return
		}

		for _, fetchedAlert := range fetchedAlerts {
			to := fetchedAlert
			err = mailer.SendEmail(to, subject, content, cfg.MailToken, cfg.ApiURL)
			if err != nil {
				slog.Error("error while sending email " + err.Error())
				return
			}

			if err := msg.Ack(); err != nil {
				slog.Error("Failed to ACK message", "error", err)
			}
		}
	}
}
