package coin

import (
	"strings"
)

// Coin represents each crypto coins
type Coin struct {
	ID                         string  `json:"id"`
	Name                       string  `json:"name"`
	Symbol                     string  `json:"symbol"`
	Rank                       float64 `json:"rank,string"`
	PriceUSD                   float64 `json:"price_usd,string"`
	PriceBTC                   float64 `json:"price_btc,string"`
	PastDayVolumeUSD           float64 `json:"24h_volume_usd,string"`
	MarketCapUSD               float64 `json:"market_cap_usd,string"`
	AvailableSupply            float64 `json:"available_supply,string"`
	TotalSupply                float64 `json:"total_supply,string"`
	MaxSupply                  float64 `json:"max_supply,string"`
	PercentChangeOneHour       float64 `json:"percent_change_1h,string"`
	PercentChangePastDay       float64 `json:"percent_change_24h,string"`
	PercentChangePastSevenDays float64 `json:"percent_change_7d,string"`
	LastUpdated                int64   `json:"-"`
}

func (coin *Coin) Metrics() map[string]float64 {
	metrics := make(map[string]float64)
	metrics[coin.concatenateMetricName("price_usd")] = coin.PriceUSD
	metrics[coin.concatenateMetricName("price_btc")] = coin.PriceBTC
	metrics[coin.concatenateMetricName("24h_volume_usd")] = coin.PastDayVolumeUSD
	metrics[coin.concatenateMetricName("market_cap_usd")] = coin.MarketCapUSD
	metrics[coin.concatenateMetricName("available_supply")] = coin.AvailableSupply
	metrics[coin.concatenateMetricName("total_supply")] = coin.TotalSupply
	metrics[coin.concatenateMetricName("percent_change_1h")] = coin.PercentChangeOneHour
	metrics[coin.concatenateMetricName("percent_change_24h")] = coin.PercentChangePastDay
	metrics[coin.concatenateMetricName("percent_change_7d")] = coin.PercentChangePastSevenDays

	return metrics
}

func (coin *Coin) concatenateMetricName(attr string) string {
	return strings.Join([]string{"", strings.Replace(strings.ToLower(coin.ID), "-", "_", -1), attr}, "_")
}

type Coins []*Coin
