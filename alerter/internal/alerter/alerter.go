package alerter

import (
	"alerter/internal/mailer"
	"alerter/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName     = "ALERTS"
	consumerName   = "alertsConsumer"
	defaultSubject = "ALERTS.>"
)

var (
	js jetstream.JetStream
	nc *nats.Conn
)

type Config struct {
	ConfigURL    string
	TimetableURL string
	MailToken    string
	ApiURL       string
}

// InitNATS initializes NATS connection and creates stream if not exists
func InitNATS() error {
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		return errors.New("failed to connect to NATS: " + err.Error())
	}

	js, err = jetstream.New(nc)
	if err != nil {
		return errors.New("failed to create JetStream context: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if stream exists first
	_, err = js.Stream(ctx, streamName)
	if err != nil {
		_, err = js.CreateStream(ctx, jetstream.StreamConfig{
			Name:     streamName,
			Subjects: []string{defaultSubject},
		})
		if err != nil {
			return errors.New("failed to create stream: " + err.Error())
		}
		slog.Info("NATS stream created", "stream", streamName)
	}

	slog.Info("NATS initialized successfully")
	return nil
}

// GetConsumer gets or creates a durable consumer
func GetConsumer(ctx context.Context) (jetstream.Consumer, error) {
	stream, err := js.Stream(ctx, streamName)
	if err != nil {
		return nil, errors.New("failed to get stream: " + err.Error())
	}

	consumer, err := stream.Consumer(ctx, consumerName)
	if err == nil {
		slog.Info("Using existing consumer", "consumer", consumerName)
		return consumer, nil
	}

	slog.Info("Creating new consumer", "consumer", consumerName)
	return stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       consumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})
}

// ConsumeMessages starts consuming messages with context support
func ConsumeMessages(ctx context.Context, cfg Config) error {
	consumer, err := GetConsumer(ctx)
	if err != nil {
		return errors.New("failed to get consumer: " + err.Error())
	}

	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		processMessage(msg, cfg)
	})
	if err != nil {
		return errors.New("failed to start consuming: " + err.Error())
	}
	defer cc.Stop()

	slog.Info("Started consuming messages", "subject", defaultSubject)

	// Wait for context cancellation
	<-ctx.Done()
	slog.Info("Stopping message consumption due to context cancellation")
	return nil
}

func processMessage(msg jetstream.Msg, cfg Config) {
	var alert models.Alert
	if err := json.Unmarshal(msg.Data(), &alert); err != nil {
		slog.Error("Failed to unmarshal alert", "error", err)
		if err := msg.Nak(); err != nil {
			slog.Error("Failed to NAK message", "error", err)
		}
		return
	}

	//build mail
	//get to mail :

	//mail data
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
		fetchedAlerts, err := FetchAlertsByResource(cfg.ConfigURL, resource)
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

func FetchAlertsByResource(apiURL string, resource int) ([]string, error) {
	// Convert resource ID to string for the URL
	resp, err := http.Get(fmt.Sprintf("%s/alerts?ucaID=%d", apiURL, resource))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a slice of structs containing just the email
	var alerts []struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &alerts); err != nil {
		return nil, err
	}

	// Extract just the emails into a string slice
	emails := make([]string, len(alerts))
	for i, alert := range alerts {
		emails[i] = alert.Email
	}

	return emails, nil
}
