package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricObserver interface {
	Observe() error
	Collector() prometheus.Collector
}
