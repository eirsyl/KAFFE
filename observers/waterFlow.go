package observers

import (
	"sync"

	"github.com/kidoman/embd/convertors/mcp3008"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type WaterFlowObserver struct {
	waterFlow prometheus.Gauge
	mpc       *mcp3008.MCP3008
	mut       *sync.Mutex
}

func NewWaterFlowObserver(mpc *mcp3008.MCP3008, mut *sync.Mutex) *WaterFlowObserver {
	return &WaterFlowObserver{
		waterFlow: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "water_flow",
			Help: "Flow test",
		}),
		mpc: mpc,
		mut: mut,
	}
}

func (p *WaterFlowObserver) Observe() error {
	p.mut.Lock()
	defer p.mut.Unlock()

	log.Infof("Collecting: %v", "water flow")

	readValue, err := p.mpc.AnalogValueAt(4)
	if err != nil {
		return err
	}

	log.Infof("water flow: %v", readValue)

	p.waterFlow.Set(float64(readValue))
	return nil
}

func (p *WaterFlowObserver) Collector() prometheus.Collector {
	return p.waterFlow
}
