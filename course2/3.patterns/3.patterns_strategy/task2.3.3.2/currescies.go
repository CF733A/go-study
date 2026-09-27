package main

import "encoding/json"

type Currencies []string

func UnmarshalCurrencies(data []byte) (Currencies, error) {
	var r Currencies
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Currencies) Marshal() ([]byte, error) {
	return json.Marshal(r)
}