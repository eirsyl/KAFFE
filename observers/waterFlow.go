package observers

import (
	"sync"

	"github.com/kidoman/embd/convertors/mcp3008"
	"github.com/prometheus/client_golang/prometheus"
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
			Help: "",
		}),
		mpc: mpc,
		mut: mut,
	}
}

func (p *WaterFlowObserver) Observe() error {
	p.mut.Lock()
	defer p.mut.Unlock()

	return nil
}

func (p *WaterFlowObserver) Collector() prometheus.Collector {
	return p.waterFlow
}
