package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func New(brokers []string, topic string) *Producer {
	return &Producer{
		w: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Producer) Send(ctx context.Context, _ string, key, value []byte) error {
	return p.w.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

func (p *Producer) Close() error { return p.w.Close() }
