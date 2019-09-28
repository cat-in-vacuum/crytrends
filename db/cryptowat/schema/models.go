package schema

import "time"

type OHLCSchema struct {
	ID        int
	QueryTime time.Time
	Market    string
	Pair      string
	Period    string
	OHLC      [6]float64
}
