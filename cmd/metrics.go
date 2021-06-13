package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var durationCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "duration_count",
	Help: "counter of duration seconds",
},
	[]string{"type"})
