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

var (
	priceUSD = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_price_usd",
			Help: "Current price of the coin",
		},
		[]string{"currency"},
	)
	priceBTC = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_price_btc",
			Help: "Current btc price of the coin",
		},
		[]string{"currency"},
	)
	pastDayVolumeUSD = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_past_day_volume_usd",
			Help: "Volume of the coin in past 24 hours",
		},
		[]string{"currency"},
	)
	marketCapUSD = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_market_cap_usd",
			Help: "Market capitalization of the coin",
		},
		[]string{"currency"},
	)
	availableSupply = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_available_supply",
			Help: "Available supply",
		},
		[]string{"currency"},
	)
	totalSupply = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_total_supply",
			Help: "Total supply",
		},
		[]string{"currency"},
	)
	maxSupply = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_max_supply",
			Help: "Maximum supply",
		},
		[]string{"currency"},
	)
	percentChangeOneHour = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_percent_change_1h",
			Help: "Percent change in 1 hour",
		},
		[]string{"currency"},
	)
	percentChangePastDay = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_percent_change_24h",
			Help: "Percent change in 24 hours",
		},
		[]string{"currency"},
	)
	percentChangePastSevenDays = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_percent_change_7d",
			Help: "Percent change in 7 days",
		},
		[]string{"currency"},
	)
	priceCNY = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_price_cny",
			Help: "Current CNY price of the coin",
		},
		[]string{"currency"},
	)
	pastDayVolumeCNY = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_past_day_volume_cny",
			Help: "CNY volume of the coin in past 24 hours",
		},
		[]string{"currency"},
	)
	marketCapCNY = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coin_market_cap_cny",
			Help: "CNY market capitalization of the coin",
		},
		[]string{"currency"},
	)
)

func init() {
	prometheus.MustRegister(priceUSD)
	prometheus.MustRegister(priceBTC)
	prometheus.MustRegister(pastDayVolumeUSD)
	prometheus.MustRegister(marketCapUSD)
	prometheus.MustRegister(availableSupply)
	prometheus.MustRegister(totalSupply)
	prometheus.MustRegister(maxSupply)
	prometheus.MustRegister(percentChangeOneHour)
	prometheus.MustRegister(percentChangePastDay)
	prometheus.MustRegister(percentChangePastSevenDays)
	prometheus.MustRegister(priceCNY)
	prometheus.MustRegister(pastDayVolumeCNY)
	prometheus.MustRegister(marketCapCNY)
}

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
	resp, err := resty.R().Get("https://api.coinmarketcap.com/v1/ticker/?convert=CNY")
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
				(&gaugeCoin{priceUSD}).setCoin(coin, coin.PriceUSD)
				(&gaugeCoin{priceBTC}).setCoin(coin, coin.PriceBTC)
				(&gaugeCoin{pastDayVolumeUSD}).setCoin(coin, coin.PastDayVolumeUSD)
				(&gaugeCoin{marketCapUSD}).setCoin(coin, coin.MarketCapUSD)
				(&gaugeCoin{availableSupply}).setCoin(coin, coin.AvailableSupply)
				(&gaugeCoin{totalSupply}).setCoin(coin, coin.TotalSupply)
				(&gaugeCoin{maxSupply}).setCoin(coin, coin.MaxSupply)
				(&gaugeCoin{percentChangeOneHour}).setCoin(coin, coin.PercentChangeOneHour)
				(&gaugeCoin{percentChangePastDay}).setCoin(coin, coin.PercentChangePastDay)
				(&gaugeCoin{percentChangePastSevenDays}).setCoin(coin, coin.PercentChangePastSevenDays)
				(&gaugeCoin{priceCNY}).setCoin(coin, coin.PriceCNY)
				(&gaugeCoin{pastDayVolumeCNY}).setCoin(coin, coin.PastDayVolumeCNY)
				(&gaugeCoin{marketCapCNY}).setCoin(coin, coin.MarketCapCNY)

				c.coins[coin.ID] = coin
			}
			c.lastUpdated = time.Now().Unix()
			c.mu.Unlock()
		}
	}()
}

type gaugeCoin struct {
	*prometheus.GaugeVec
}

func (gauge *gaugeCoin) setCoin(coin *coin.Coin, v float64) {
	gauge.WithLabelValues(coin.ID).Set(v)
}