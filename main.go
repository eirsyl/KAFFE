package main

import (
	"flag"
	"os"
	"time"

	"os/signal"

	"syscall"

	"sync"

	"kaffe/observers"

	"net/http"

	"net"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/convertors/mcp3008"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	var slackToken = flag.String("slacktoken", "", "slack bot token")
	var slackChannel = flag.String("slackchannel", "#general", "slack channel")
	flag.Parse()

	if *pushgateway == "" {
		log.Fatalf("The pushgateway flag cannot be empty")
	}

	if *slackToken != "" && *slackChannel != "" {
		go func() {
			var lastIP net.IP
			for {
				time.Sleep(60 * time.Second)
				ip, err := GetOutboundIP()
				if err != nil {
					log.Warn("Could not find outbound ip: %v", err)
					continue
				}
				if !lastIP.Equal(ip) {
					err = PostToSlack(*slackToken, *slackChannel, ip)
					if err != nil {
						log.Warn("Could not post message to slack: %v", err)
					}
					lastIP = ip
				}
			}
		}()
	}

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	if err := embd.InitSPI(); err != nil {
		panic(err)
	}
	defer embd.CloseSPI()

	const (
		channel = 0
		speed   = 1000000
		bpw     = 8
		delay   = 0
	)

	spiBus := embd.NewSPIBus(embd.SPIMode0, channel, speed, bpw, delay)
	defer spiBus.Close()
	adc := mcp3008.New(mcp3008.SingleMode, spiBus)

	var mutex = &sync.Mutex{}

	metrics := []MetricObserver{
		observers.NewPlateModeObserver(adc, mutex),
		observers.NewPowerObserver(adc, mutex),
		observers.NewWaterContainerObserver(adc, mutex),
		observers.NewPlateTempObserver(adc, mutex),
		observers.NewWaterFlowObserver(adc, mutex),
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
				log.Warn("Could not push metrics: %v", err)
			}
		}
	}()

	go func() {
		addr := ":8081"
		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		log.Infof("Prometheus handler is listening on %v", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			failure <- err
		}
	}()

	go func() {
		// Testing code for flow sensor
		flow, err := embd.NewDigitalPin(26)
		if err != nil {
			log.Errorf("Could not open flow gpio: %v", err)
			failure <- err
		}
		defer flow.Close()

		if err := flow.SetDirection(embd.In); err != nil {
			log.Errorf("Could not set pin direction: %v", err)
			failure <- err
		}
		flow.ActiveLow(true)

		reading := make(chan bool)
		go func(c chan bool) {
			err := flow.Watch(embd.EdgeRising, func(flow embd.DigitalPin) {
				log.Info("Received pin interrupt")
				reading <- true
			})
			if err != nil {
				log.Errorf("Could not watch port: %v", err)
				failure <- err
			}
		}(reading)

		var flowCounter int
		for {
			res := <-reading
			if res {
				flowCounter++
			}
			log.Infof("flow reader counter: %v", flowCounter)
		}
	}()

	select {
	case sig := <-terminate:
		log.Errorf("Received signal: %v", sig)
	case err := <-failure:
		log.Errorf("Internal error: %v", err)
	}
}
