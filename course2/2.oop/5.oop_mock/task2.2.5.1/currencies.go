package main

import "encoding/json"

type Currencies []string

func UnmarshalCurrency(data []byte) (Currencies, error) {
	var r Currencies
	err := json.Unmarshal(data, &r)
	return r, err
}

func (c *Currencies) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
