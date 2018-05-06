package observers

import (
	"sync"

	"github.com/kidoman/embd/convertors/mcp3008"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type PlateModeObserver struct {
	plateMode prometheus.Gauge
	mpc       *mcp3008.MCP3008
	mut       *sync.Mutex
}

func NewPlateModeObserver(mpc *mcp3008.MCP3008, mut *sync.Mutex) *PlateModeObserver {
	return &PlateModeObserver{
		plateMode: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "plate_mode",
			Help: "Current plate mode high/low.",
		}),
		mpc: mpc,
		mut: mut,
	}
}

func (p *PlateModeObserver) Run() error {
	return nil
}

func (p *PlateModeObserver) Stop() error {
	return nil
}

func (p *PlateModeObserver) Observe() error {
	p.mut.Lock()
	defer p.mut.Unlock()

	log.Infof("Collecting: %v", "plate mode")

	amps, err := ReadACS712AC(p.mpc, 0)
	if err != nil {
		return err
	}

	log.Infof("plate mode amps: %v", amps)

	if amps > 0.1 {
		p.plateMode.Set(1)
	} else {
		p.plateMode.Set(0)
	}
	return nil
}

func (p *PlateModeObserver) Collector() prometheus.Collector {
	return p.plateMode
}
