package observers

import (
	"github.com/kidoman/embd"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type WaterFlowObserver struct {
	waterFlow     prometheus.Gauge
	interruptChan chan bool
	flowCounter   int
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

	go p.aggregateInterrupts()

	flow.Watch(embd.EdgeRising, func(flow embd.DigitalPin) {
		log.Info("Received pin interrupt")
		p.interruptChan <- true
	})
	if err != nil {
		log.Fatalf("Could not watch port: %v", err)
	}
	return err
}

func (p *WaterFlowObserver) Stop() error {
	return nil
}

func (p *WaterFlowObserver) aggregateInterrupts() {
	for {
		res := <-p.interruptChan
		if res {
			p.flowCounter++
		}
		log.Infof("flow reader counter: %v", p.flowCounter)
	}
}

func (p *WaterFlowObserver) Observe() error {
	p.waterFlow.Set(float64(p.flowCounter))
	return nil
}

func (p *WaterFlowObserver) Collector() prometheus.Collector {
	return p.waterFlow
}
