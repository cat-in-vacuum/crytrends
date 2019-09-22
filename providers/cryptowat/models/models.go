package models

type OHLCRespCommon map[string]map[string]OHLCResp
type OHLCResp struct {
	Result    map[string][][]float64 `json:"result"`
	Allowance Allowance              `json:"allowance"`
}

type Allowance struct {
	Cost      int   `json:"cost"`
	Remaining int64 `json:"remaining"`
}
