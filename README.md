# crypto-market-overwatch

crypto-market-overwatch is a exportor for prometheus to track crypto coins.

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

Then start your prometheus and configure the target to the server, see [prometheus.yml](https://github.com/FX-HAO/crypto-market-overwatch/blob/master/prometheus/prometheus.yml).

## Options

## Roadmap

- [ ] Alert support
- [ ] Docker support
- [ ] Grafana configuration
