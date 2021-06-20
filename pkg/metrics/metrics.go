package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Vector struct {
	DurationCount *prometheus.CounterVec
	DurationTotal *prometheus.CounterVec
}

// Metrics is global metrics vector
var Metrics = Vector{
	DurationCount: promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "duration_count",
		Help: "counter of duration seconds",
	}, []string{"type"}),
	DurationTotal: promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "duration_total",
		Help: "counter of total duration seconds",
	}, []string{"type"}),
}
