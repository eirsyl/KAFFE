package observers

import (
	"sync"

	"github.com/kidoman/embd/convertors/mcp3008"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type PowerObserver struct {
	power prometheus.Gauge
	mpc   *mcp3008.MCP3008
	mut   *sync.Mutex
}

func NewPowerObserver(mpc *mcp3008.MCP3008, mut *sync.Mutex) *PowerObserver {
	return &PowerObserver{
		power: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "power",
			Help: "Power switch on/off",
		}),
		mpc: mpc,
		mut: mut,
	}
}

func (p *PowerObserver) Run() error {
	return nil
}

func (p *PowerObserver) Stop() error {
	return nil
}

func (p *PowerObserver) Observe() error {
	p.mut.Lock()
	defer p.mut.Unlock()

	log.Infof("Collecting: %v", "power")

	amps, err := ReadACS712AC(p.mpc, 1)
	if err != nil {
		return err
	}

	log.Infof("power amps: %v", amps)

	if amps > 0.1 {
		p.power.Set(1)
	} else {
		p.power.Set(0)
	}
	return nil
}

func (p *PowerObserver) Collector() prometheus.Collector {
	return p.power
}
