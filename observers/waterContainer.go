package observers

import (
	"sync"

	"bytes"
	"net/http"

	"github.com/kidoman/embd/convertors/mcp3008"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type WaterContainerObserver struct {
	waterContainer prometheus.Gauge
	mpc            *mcp3008.MCP3008
	mut            *sync.Mutex
	hubot          string
	token          string
	containerEmpty bool
}

func NewWaterContainerObserver(mpc *mcp3008.MCP3008, mut *sync.Mutex, hubot string, token string) *WaterContainerObserver {
	return &WaterContainerObserver{
		waterContainer: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "water_container",
			Help: "Water container switch on/off",
		}),
		mpc:            mpc,
		mut:            mut,
		hubot:          hubot,
		token:          token,
		containerEmpty: true,
	}
}

func (p *WaterContainerObserver) Run() error {
	return nil
}

func (p *WaterContainerObserver) Stop() error {
	return nil
}

/*
Post brewing done to hubot when the water container changes from non-empty to empty.
*/
func (p *WaterContainerObserver) postHubot(empty bool) {
	if p.containerEmpty == false && empty == true {
		log.Info("posting brewing done to hubot")

		var body = []byte(`{}`)
		req, err := http.NewRequest("POST", p.hubot, bytes.NewBuffer(body))
		req.Header.Set("Authorization", p.token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Warnf("could not post to hubot: %v", err)
		} else {
			log.Info("hubot status: %d", resp.Status)
		}
		defer resp.Body.Close()
	}
	p.containerEmpty = empty
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
		p.postHubot(false)
	} else {
		p.waterContainer.Set(0)
		p.postHubot(true)
	}
	return nil
}

func (p *WaterContainerObserver) Collector() prometheus.Collector {
	return p.waterContainer
}
