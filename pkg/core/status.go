package core

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Status struct {
	Counter      prometheus.Counter `json:"-"`
	TotalCounter prometheus.Counter `json:"-"`
	Duration     time.Duration      `json:"duration,omitempty" yaml:"duration"`
}

func NewStatus(c, tc prometheus.Counter, duration time.Duration) *Status {
	return &Status{
		Counter:      c,
		TotalCounter: tc,
		Duration:     duration,
	}
}
