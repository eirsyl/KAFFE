package main

import (
	"flag"
	"os"
	"time"

	"os/signal"

	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"
)

// Push metrics from a registry to a push gateway.
func pushMetrics(pushgateway string, registry *prometheus.Registry) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	log.Infof("Pushing metrics to: %v", pushgateway)
	return push.AddFromGatherer(
		"moccamaster", map[string]string{"instance": hostname},
		pushgateway,
		registry,
	)
}

func main() {
	var pushgateway = flag.String("pushgateway", "http://127.0.0.1:9091", "pushgateway url")
	flag.Parse()

	if *pushgateway == "" {
		log.Fatalf("The pushgateway flag cannot be empty")
	}

	metrics := []MetricObserver{
		NewPlateTempObserver(),
		NewPlateModeObserver(),
	}

	registry := prometheus.NewRegistry()
	for _, observer := range metrics {
		log.Infof("Adding observer: %v", observer)
		registry.MustRegister(observer.Collector())
	}

	var failure = make(chan error, 1)
	for _, observer := range metrics {
		go func(ob MetricObserver) {
			for {
				err := ob.Observe()
				if err != nil {
					failure <- err
				}
				time.Sleep(10 * time.Second)
			}
		}(observer)
	}

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			time.Sleep(30 * time.Second)
			err := pushMetrics(*pushgateway, registry)
			if err != nil {
				failure <- err
			}
		}
	}()

	select {
	case sig := <-terminate:
		log.Errorf("Received signal: %v", sig)
	case err := <-failure:
		log.Errorf("Internal error: %v", err)
	}
}
