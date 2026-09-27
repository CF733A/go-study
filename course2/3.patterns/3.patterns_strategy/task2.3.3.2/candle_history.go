package main

import "encoding/json"

func UnmarshalCandlesHistory(data []byte) (CandlesHistory, error) {
	var r CandlesHistory
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CandlesHistory) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CandlesHistory struct {
	Candles []Candle `json:"candles"`
}

type Candle struct {
	T int64   `json:"t"`
	O float64 `json:"o"`
	C float64 `json:"c"`
	H float64 `json:"h"`
	L float64 `json:"l"`
	V float64 `json:"v"`
}