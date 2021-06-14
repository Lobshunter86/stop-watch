package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricsVector struct {
	durationCount *prometheus.CounterVec
	durationTotal *prometheus.CounterVec
}

var metrics = metricsVector{
	durationCount: promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "duration_count",
		Help: "counter of duration seconds",
	}, []string{"type"}),
	durationTotal: promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "duration_total",
		Help: "counter of total duration seconds",
	}, []string{"type"}),
}
