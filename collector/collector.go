package collector

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/FX-HAO/crypto-market-overwatch/coin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	resty "gopkg.in/resty.v0"
)

type Collector struct {
	mu sync.RWMutex

	coins  map[string]*coin.Coin
	gauges map[string]prometheus.Gauge

	interval    int
	lastUpdated int64

	closed int32
}

func NewCollector(interval int) *Collector {
	collector := &Collector{
		coins:    make(map[string]*coin.Coin),
		gauges:   make(map[string]prometheus.Gauge),
		interval: interval,
	}
	return collector
}

// Start does some initial preparation and starts to work
func (c *Collector) Start() {
	c.collect()
}

func (c *Collector) fetch() (coin.Coins, error) {
	resp, err := resty.R().Get("https://api.coinmarketcap.com/v1/ticker/")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp.StatusCode() != 200 {
		log.Error("Cannot access routes api")
	}
	coins := []*coin.Coin{}
	if err := json.Unmarshal(resp.Body(), &coins); err != nil {
		return nil, err
	}
	return coins, nil
}

// collect starts a goroutine to collect the data
func (c *Collector) collect() {
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(c.interval))
			c.mu.Lock()
			coins, err := c.fetch()
			if err != nil {
				c.mu.Unlock()
				log.Error(err)
				continue
			}
			for _, coin := range coins {
				id := coin.ID
				metrics := coin.Metrics()
				if _, ok := c.coins[id]; !ok {
					for k, _ := range metrics {
						gauge := prometheus.NewGaugeVec(
							prometheus.GaugeOpts{
								Name: k,
								Help: k,
							},
							[]string{},
						).WithLabelValues()
						prometheus.MustRegister(gauge)
						c.gauges[k] = gauge
						log.Infof("Register %s", k)
					}
				}

				for k, v := range metrics {
					gauge := c.gauges[k]
					gauge.Set(v)
				}
				c.coins[id] = coin
			}
			c.lastUpdated = time.Now().Unix()
			c.mu.Unlock()
		}
	}()
}
