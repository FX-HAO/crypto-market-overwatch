package handler

import (
	"encoding/json"
	"net/http"
)

func TickerHanlder(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode({});
}
