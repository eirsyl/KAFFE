package observers

import (
	"sync"

	"github.com/kidoman/embd/convertors/mcp3008"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type WaterContainerObserver struct {
	waterContainer prometheus.Gauge
	mpc            *mcp3008.MCP3008
	mut            *sync.Mutex
}

func NewWaterContainerObserver(mpc *mcp3008.MCP3008, mut *sync.Mutex) *WaterContainerObserver {
	return &WaterContainerObserver{
		waterContainer: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "water_container",
			Help: "Water container switch on/off",
		}),
		mpc: mpc,
		mut: mut,
	}
}

func (p *WaterContainerObserver) Run() error {
	return nil
}

func (p *WaterContainerObserver) Stop() error {
	return nil
}

func (p *WaterContainerObserver) Observe() error {
	p.mut.Lock()
	defer p.mut.Unlock()

	log.Infof("Collecting: %v", "water_container")

	amps, err := ReadACS712AC(p.mpc, 2)
	if err != nil {
		return err
	}

	log.Infof("water container amps: %v", amps)

	if amps > 0.1 {
		p.waterContainer.Set(1)
	} else {
		p.waterContainer.Set(0)
	}
	return nil
}

func (p *WaterContainerObserver) Collector() prometheus.Collector {
	return p.waterContainer
}
