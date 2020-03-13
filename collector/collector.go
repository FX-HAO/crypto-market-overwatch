package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/FX-HAO/crypto-market-overwatch/coin"
	"github.com/FX-HAO/crypto-market-overwatch/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	resty.SetHTTPMode()

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

// Collector represents a collector that collects data from the server
type Collector struct {
	mu sync.RWMutex

	coins  map[string]*coin.Coin
	gauges map[string]prometheus.Gauge

	interval    int
	lastUpdated int64

	closed int32
}

// NewCollector creates a new Collector and returns it.
func newCollector(interval int) *Collector {
	collector := &Collector{
		coins:    make(map[string]*coin.Coin),
		gauges:   make(map[string]prometheus.Gauge),
		interval: interval,
	}
	return collector
}

func ListenAndServe(host string, port, interval int) {
	c := newCollector(interval)
	c.Start()

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	router.HandleFunc("/api/coins", c.coinsHandler).Methods("GET")
	router.HandleFunc("/api/coins/{coin}", c.coinHandler).Methods("GET")

	router.Use(restfulMiddleware)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}

// Start does some initial preparation and starts to work
func (c *Collector) Start() {
	c.collect()
}

func (c *Collector) fetch() (coin.Coins, error) {
	resp, err := resty.R().Get("https://web-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?convert=USD,CNY")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Cannot access api, status code: %d", resp.StatusCode())
	}
	var data struct {
		Data coin.Coins `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	coins := data.Data
	coins.Init()
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
			log.Infof("fetch %d coins", len(coins))

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
