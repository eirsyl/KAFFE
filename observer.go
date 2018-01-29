package main

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type MetricObserver interface {
	Observe() error
	Collector() prometheus.Collector
}

type PlateTempObserver struct {
	plateTemp prometheus.Gauge
}

func NewPlateTempObserver() *PlateTempObserver {
	return &PlateTempObserver{
		plateTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "plate_temp",
			Help: "Current temperature of the plate.",
		}),
	}
}

func (p *PlateTempObserver) Observe() error {
	log.Infof("Collecting: %v", "plate temp")
	p.plateTemp.Add(63)
	return nil
}

func (p *PlateTempObserver) Collector() prometheus.Collector {
	return p.plateTemp
}

type PlateModeObserver struct {
	plateMode prometheus.Gauge
}

func NewPlateModeObserver() *PlateModeObserver {
	return &PlateModeObserver{
		plateMode: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "plate_mode",
			Help: "Current plate mode high/low.",
		}),
	}
}

func (p *PlateModeObserver) Observe() error {
	log.Infof("Collecting: %v", "plate mode")

	high := true
	if high {
		p.plateMode.Set(1)
	} else {
		p.plateMode.Set(0)
	}
	return nil
}

func (p *PlateModeObserver) Collector() prometheus.Collector {
	return p.plateMode
}
