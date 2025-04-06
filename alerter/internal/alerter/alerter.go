package alerter

import (
	"alerter/internal/helpers"
	"context"
	"errors"
	"log/slog"
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
func ConsumeMessages(ctx context.Context, cfg helpers.Config) error {
	consumer, err := GetConsumer(ctx)
	if err != nil {
		return errors.New("failed to get consumer: " + err.Error())
	}

	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		ProcessMessage(msg, cfg)
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
