package observers

import (
	"sync"

	"github.com/kidoman/embd/convertors/mcp3008"
	"github.com/prometheus/client_golang/prometheus"
)

type PlateTempObserver struct {
	plateTemp prometheus.Gauge
	mpc       *mcp3008.MCP3008
	mut       *sync.Mutex
}

func NewPlateTempObserver(mpc *mcp3008.MCP3008, mut *sync.Mutex) *PlateTempObserver {
	return &PlateTempObserver{
		plateTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "plate_temp",
			Help: "",
		}),
		mpc: mpc,
		mut: mut,
	}
}

func (p *PlateTempObserver) Observe() error {
	p.mut.Lock()
	defer p.mut.Unlock()

	return nil
}

func (p *PlateTempObserver) Collector() prometheus.Collector {
	return p.plateTemp
}
