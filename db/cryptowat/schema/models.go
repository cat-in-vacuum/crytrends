package schema

import "time"

type OHLCSchema struct {
	QueryTime time.Time
	Market    string
	Pair      string
	Period    string
	OHLC      [][]float64
}
