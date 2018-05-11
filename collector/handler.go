package collector

import (
	"encoding/json"
	"net/http"

	"github.com/FX-HAO/crypto-market-overwatch/coin"
	"github.com/gorilla/mux"
)

func restfulMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (c *Collector) coinsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(c.coins)
}

func (c *Collector) coinHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var coin *coin.Coin
	var ok bool

	if coin, ok = c.coins[vars["coin"]]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(coin)
}
