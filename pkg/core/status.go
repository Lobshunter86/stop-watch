package core

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Status struct {
	Counter  prometheus.Counter
	Duration time.Duration `json:"duration,omitempty" yaml:"duration"`
}

func NewStatus(c prometheus.Counter, duration time.Duration) Status {
	return Status{
		Counter:  c,
		Duration: duration,
	}
}
