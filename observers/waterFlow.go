package observers

import (
	"github.com/kidoman/embd"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type WaterFlowObserver struct {
	waterFlow   prometheus.Gauge
	flowCounter int
}

func NewWaterFlowObserver() *WaterFlowObserver {
	return &WaterFlowObserver{
		waterFlow: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "water_flow",
			Help: "The water flow reported by the flow meter.",
		}),
	}
}

func (p *WaterFlowObserver) Run() error {
	var interruptChan = make(chan bool)

	flow, err := embd.NewDigitalPin(26)
	if err != nil {
		log.Errorf("Could not open flow gpio: %v", err)
		return err
	}
	defer flow.Close()

	if err := flow.SetDirection(embd.In); err != nil {
		log.Errorf("Could not set pin direction: %v", err)
		return err
	}
	flow.ActiveLow(true)

	flow.Watch(embd.EdgeRising, func(flow embd.DigitalPin) {
		log.Info("Received pin interrupt")
		interruptChan <- true
	})
	if err != nil {
		log.Fatalf("Could not watch port: %v", err)
	}

	for {
		res := <-interruptChan
		if res {
			p.flowCounter++
		}
		log.Infof("flow reader counter: %v", p.flowCounter)
	}

	return err
}

func (p *WaterFlowObserver) Stop() error {
	return nil
}

func (p *WaterFlowObserver) Observe() error {
	log.Infof("Collecting: %v", "water flow")
	p.waterFlow.Set(float64(p.flowCounter))
	log.Infof("water flow: %v", p.flowCounter)
	return nil
}

func (p *WaterFlowObserver) Collector() prometheus.Collector {
	return p.waterFlow
}
