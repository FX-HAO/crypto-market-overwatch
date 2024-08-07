# crypto-market-overwatch

[![](https://images.microbadger.com/badges/image/fuxin/crypto-market-overwatch.svg)](https://microbadger.com/images/fuxin/crypto-market-overwatch "Get your own image badge on microbadger.com")
[![Go Report Card](https://goreportcard.com/badge/github.com/FX-HAO/crypto-market-overwatch)](https://goreportcard.com/report/github.com/FX-HAO/crypto-market-overwatch)

crypto-market-overwatch is an exporter for prometheus to track crypto coins. Grafana demo: http://45.32.68.4/, both of username and password are `guest`.

## Installation

```go
go get github.com/FX-HAO/crypto-market-overwatch
cd $GOPATH/src/github.com/FX-HAO/crypto-market-overwatch
go build && ./crypto-market-overwatch
```

## Usage

See detail: 

```go
./crypto-market-overwatch --help
```

It's recommended to use docker to keep your environment clean and much easier to use.

```
docker run -it -d -p 80:80 -e X_CMC_PRO_API_KEY=$X_CMC_PRO_API_KEY --name crypto-market-overwatch fuxin/crypto-market-overwatch
```

Then start your prometheus and configure the target to the server, see [prometheus.yml](https://github.com/FX-HAO/crypto-market-overwatch/blob/master/prometheus/prometheus.yml).

## APIs

```
/api/coins       // return all coins information
/api/coins/{coin} // return specific coin information
```

## Options

You can set up a Grafana dashboard for Prometheus. More info: [Grafana supports querying Prometheus](https://prometheus.io/docs/visualization/grafana/)

## Roadmap

- [ ] Alert support
- [x] Docker support
- [ ] Grafana configuration
