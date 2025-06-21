package usecase

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"kulturaGo/events-service/internal/domain"
	"time"
)

type Exporter struct {
	Repo       domain.EventRepo
	Producer   KafkaProducer
	OutTopic   string
	Batch      int
	Interval   time.Duration
	Logger     logrus.FieldLogger
	CursorRepo CursorRepo
}

type KafkaProducer interface {
	Send(ctx context.Context, topic string, key []byte, value []byte) error
}

type CursorRepo interface {
	Get(ctx context.Context, consumer string) (time.Time, error)
	Update(ctx context.Context, consumer string, ts time.Time) error
}

func (e *Exporter) Run(ctx context.Context) {
	t := time.NewTicker(e.Interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			e.exportBatch(ctx)
		}
	}
}

func (e *Exporter) exportBatch(ctx context.Context) {
	lastTs, _ := e.CursorRepo.Get(ctx, "events-exporter")

	events, err := e.Repo.IterChangedSince(ctx, lastTs, e.Batch)
	if err != nil {
		e.Logger.Error("iter", zap.Error(err))
		return
	}
	if len(events) == 0 {
		return
	}

	for _, ev := range events {
		b, _ := json.Marshal(ev)
		if err = e.Producer.Send(ctx, e.OutTopic, []byte(ev.Slug), b); err != nil {
			e.Logger.Error("produce", zap.Error(err))
			return
		}
		lastTs = ev.CreatedAt
	}
	if err = e.CursorRepo.Update(ctx, "events-exporter", lastTs); err != nil {
		e.Logger.Error("cursor", zap.Error(err))
	}
}
