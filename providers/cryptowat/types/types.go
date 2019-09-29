package types

import (
	"time"
)

type OHLCRespCommon map[string]RespPairs
type RespPairs map[string]OHLCResp
type OHLCResp struct {
	Result    map[string][][]float64 `json:"result"`
	QueryTime time.Time
	Allowance Allowance                  `json:"allowance"`
}

type Allowance struct {
	Cost      int   `json:"cost"`
	Remaining int64 `json:"remaining"`
}
