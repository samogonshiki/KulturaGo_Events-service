package scheduler

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"kulturaGo/events-service/internal/domain"
	"time"
)

type Deactivator struct {
	Repo     domain.EventRepo
	Interval time.Duration
	Logger   logrus.FieldLogger
}

func (d *Deactivator) Run(ctx context.Context) {
	t := time.NewTicker(d.Interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n, err := d.Repo.DeactivatePast(ctx)
			if err != nil {
				d.Logger.Error("deactivate", zap.Error(err))
				continue
			}
			if n > 0 {
				d.Logger.Info("events deactivated", zap.Int64("count", n))
			}
		}
	}
}
