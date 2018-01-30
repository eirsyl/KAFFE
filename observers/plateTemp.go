package observers

import (
	"sync"

	"github.com/kidoman/embd/convertors/mcp3008"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
			Help: "Plate temperature in celsius",
		}),
		mpc: mpc,
		mut: mut,
	}
}

func (p *PlateTempObserver) Run() error {
	return nil
}

func (p *PlateTempObserver) Stop() error {
	return nil
}

func (p *PlateTempObserver) Observe() error {
	p.mut.Lock()
	defer p.mut.Unlock()

	log.Infof("Collecting: %v", "plate temp")

	readValue, err := p.mpc.AnalogValueAt(3)
	if err != nil {
		return err
	}

	log.Infof("plate temp: %v", readValue)

	p.plateTemp.Set(float64(readValue))
	return nil
}

func (p *PlateTempObserver) Collector() prometheus.Collector {
	return p.plateTemp
}
