package coin

import (
	"encoding/json"
	"strings"
)

type Currency struct {
	Price                      float64 `json:"price"`
	PastDayVolume              float64 `json:"volume_24h"`
	PercentChangeOneHour       float64 `json:"percent_change_1h"`
	PercentChangePastDay       float64 `json:"percent_change_24h"`
	PercentChangePastSevenDays float64 `json:"percent_change_7d"`
	MarketCap                  float64 `json:"market_cap"`
}

type Quote struct {
	CNY *Currency `json:"CNY"`
	USD *Currency `json:"USD"`
}

// Coin represents each crypto coins
type Coin struct {
	// ID                         int64   `json:"id,int"`
	ID                         string  `json:"-"`
	Name                       string  `json:"name"`
	Symbol                     string  `json:"symbol"`
	Rank                       float64 `json:"rank,string"`
	PriceUSD                   float64 `json:"price_usd,string"`
	PriceBTC                   float64 `json:"price_btc,string"`
	PastDayVolumeUSD           float64 `json:"24h_volume_usd,string"`
	MarketCapUSD               float64 `json:"market_cap_usd,string"`
	AvailableSupply            float64 `json:"available_supply,string"`
	TotalSupply                float64 `json:"total_supply"`
	MaxSupply                  float64 `json:"max_supply"`
	PercentChangeOneHour       float64 `json:"percent_change_1h,string"`
	PercentChangePastDay       float64 `json:"percent_change_24h,string"`
	PercentChangePastSevenDays float64 `json:"percent_change_7d,string"`
	LastUpdated                int64   `json:"-"`
	PriceCNY                   float64 `json:"price_cny,string"`
	PastDayVolumeCNY           float64 `json:"24h_volume_cny,string"`
	MarketCapCNY               float64 `json:"market_cap_cny,string"`
	Quote                      *Quote  `json:"quote"`
}

type coinEncodedWithID struct {
	// ID                         int64   `json:"id,int"`
	ID                         string  `json:"id"`
	Name                       string  `json:"name"`
	Symbol                     string  `json:"symbol"`
	Rank                       float64 `json:"rank,string"`
	PriceUSD                   float64 `json:"price_usd,string"`
	PriceBTC                   float64 `json:"price_btc,string"`
	PastDayVolumeUSD           float64 `json:"24h_volume_usd,string"`
	MarketCapUSD               float64 `json:"market_cap_usd,string"`
	AvailableSupply            float64 `json:"available_supply,string"`
	TotalSupply                float64 `json:"total_supply"`
	MaxSupply                  float64 `json:"max_supply"`
	PercentChangeOneHour       float64 `json:"percent_change_1h,string"`
	PercentChangePastDay       float64 `json:"percent_change_24h,string"`
	PercentChangePastSevenDays float64 `json:"percent_change_7d,string"`
	LastUpdated                int64   `json:"-"`
	PriceCNY                   float64 `json:"price_cny,string"`
	PastDayVolumeCNY           float64 `json:"24h_volume_cny,string"`
	MarketCapCNY               float64 `json:"market_cap_cny,string"`
	Quote                      *Quote  `json:"quote"`
}

func (coin *Coin) MarshalJSON() ([]byte, error) {
	if bytes, err := json.Marshal(coinEncodedWithID(*coin)); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func (coin *Coin) metrics() map[string]float64 {
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
	return strings.Join([]string{strings.Replace(strings.ToLower(coin.ID), "-", "_", -1), attr}, "_")
}

// Coins represents an array of coins
type Coins []*Coin

func (coins Coins) Init() {
	for _, coin := range coins {
		coin.ID = strings.Replace(strings.ToLower(coin.Name), " ", "-", -1)

		coin.PriceUSD = coin.Quote.USD.Price
		coin.PastDayVolumeUSD = coin.Quote.USD.PastDayVolume
		coin.MarketCapUSD = coin.Quote.USD.MarketCap
		coin.PriceCNY = coin.Quote.CNY.Price
		coin.PastDayVolumeCNY = coin.Quote.CNY.PastDayVolume
		coin.MarketCapCNY = coin.Quote.CNY.MarketCap
	}
}
