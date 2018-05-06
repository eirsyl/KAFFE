package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricObserver interface {
	// Start background tasks like scraping GPIO pins.
	Run() error
	// Stop background tasks gracefully.
	Stop() error
	// This function is called each X seconds, usually used to sample values.
	Observe() error
	// Return an prometheus collector that the prometheus register can register and export.
	Collector() prometheus.Collector
}
