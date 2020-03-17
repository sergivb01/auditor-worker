package worker

import (
	"time"

	"go.uber.org/zap"
)

func Start() time.Time {
	return time.Now()
}

func (w *Worker) Track(name string, startTime time.Time) {
	w.log.Debug("Task executed",
		zap.String("name", name),
		zap.Duration("duration", time.Since(startTime)))
}
