package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"user-service/internal/events"
	kfk "user-service/pkg/kafka"
	"user-service/internal/models"

	"github.com/segmentio/kafka-go"
)

type KafkaHandler struct {
	repo     UserProfileRepository
	producer *kfk.Producer
}

func NewKafkaHandler(repo UserProfileRepository, producer *kfk.Producer) *KafkaHandler {
	return &KafkaHandler{repo: repo, producer: producer}
}

func (h *KafkaHandler) Handle(ctx context.Context, msg kafka.Message) error {

	var event events.UserRegisteredEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.toDLQ(ctx, msg)
		return nil
	}

	h.delay(msg.Topic)

	if err := h.handleUserRegistered(ctx, event); err != nil {
		return h.retry(ctx, event)
	}

	return nil
}

func (h *KafkaHandler) handleUserRegistered(
	ctx context.Context,
	event events.UserRegisteredEvent,
) error {
	_, err := h.repo.CreateUser(ctx, &models.UserProfile{
		UserID: event.Payload.UserID,
		Email:  event.Payload.Email,
	})
	return err
}

func (h *KafkaHandler) retry(ctx context.Context, event events.UserRegisteredEvent) error {

	event.Meta.RetryCount++

	var topic string

	switch event.Meta.RetryCount {
	case 1:
		topic = h.producer.Topic + ".retry.1"
	case 2:
		topic = h.producer.Topic + ".retry.2"
	case 3:
		topic = h.producer.Topic + ".retry.3"
	default:
		return h.producer.Publish(
			ctx,
			h.producer.Topic + ".dlq",
			event.Payload.UserID,
			event,
		)
	}

	log.Printf("retry %d for user %s",
		event.Meta.RetryCount,
		event.Payload.UserID,
	)

	return h.producer.Publish(
		ctx,
		topic,
		event.Payload.UserID,
		event,
	)
}

func (h *KafkaHandler) toDLQ(ctx context.Context, msg kafka.Message) {
	_ = h.producer.Publish(
		ctx,
		h.producer.Topic + ".dlq",
		string(msg.Key),
		msg.Value,
	)
}

func (h *KafkaHandler) delay(topic string) {
	switch topic {
	case h.producer.Topic + ".retry.1":
		time.Sleep(10 * time.Second)
	case h.producer.Topic + ".retry.1":
		time.Sleep(1 * time.Minute)
	case h.producer.Topic + ".retry.1":
		time.Sleep(10 * time.Minute)
	}
}
