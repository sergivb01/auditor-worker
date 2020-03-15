package worker

import (
	"time"

	"go.uber.org/zap"
)

func Start() time.Time {
	return time.Now()
}

func (s *Worker) Track(name string, startTime time.Time) {
	s.log.Debug("Task executed",
		zap.String("name", name),
		zap.Duration("duration", time.Since(startTime)))
}
