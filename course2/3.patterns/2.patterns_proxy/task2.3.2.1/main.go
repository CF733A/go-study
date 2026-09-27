package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cinar/indicator"
)

type Indicator interface {
	StochPrice() ([]float64, []float64)
	RSI(period int) ([]float64, []float64)
	StochRSI(rsiPeriod int) ([]float64, []float64)
	SMA(period int) []float64
	MACD() ([]float64, []float64)
	EMA() []float64
}

func UnmarshalKLines(data []byte) (KLines, error) {
	var r KLines
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *KLines) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type KLines struct {
	Pair    string   `json:"pair"`
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

type Lines struct {
	high    []float64
	low     []float64
	closing []float64
}

func (t *Lines) StochPrice() ([]float64, []float64) {
	k, d := indicator.StochasticOscillator(t.high, t.low, t.closing)

	return k, d
}

func (t *Lines) RSI(period int) ([]float64, []float64) {
	rs, rsi := indicator.RsiPeriod(period, t.closing)

	return rs, rsi
}

func (t *Lines) StochRSI(rsiPeriod int) ([]float64, []float64) {
	_, rsi := t.RSI(rsiPeriod)
	k, d := indicator.StochasticOscillator(rsi, rsi, rsi)

	return k, d
}

func (t *Lines) SMA(period int) []float64 {
	return indicator.Sma(period, t.closing)
}

func (t *Lines) MACD() ([]float64, []float64) {
	return indicator.Macd(t.closing)
}

func (t *Lines) EMA() []float64 {
	return indicator.Ema(5, t.closing)
}

type LinesProxy struct {
	lines Indicator
	cache map[string][]float64
}

func LoadKlinesProxy(data []byte) *LinesProxy {
	return &LinesProxy{
		lines: LoadKlines(data),
		cache: make(map[string][]float64),
	}
}

func (lp *LinesProxy) StochPrice() ([]float64, []float64) {
	kKey := "k_stochprice"
	dKey := "d_stochprice"

	if kkey, ok := lp.cache[kKey]; ok {
		if dkey, ok := lp.cache[dKey]; ok {
			return kkey, dkey
		}
	}

	k, d := lp.lines.StochPrice()
	lp.cache[kKey] = k
	lp.cache[dKey] = d
	return k, d
}

func (lp *LinesProxy) RSI(period int) ([]float64, []float64) {
	rsKey := fmt.Sprintf("rs_%v", period)
	rsiKey := fmt.Sprintf("rsi_%v", period)

	if rs, ok := lp.cache[rsKey]; ok {
		if rsi, ok := lp.cache[rsiKey]; ok {
			return rs, rsi
		}
	}

	rs, rsi := lp.lines.RSI(period)
	lp.cache[rsKey] = rs
	lp.cache[rsiKey] = rsi
	return rs, rsi
}

func (lp *LinesProxy) StochRSI(rsiPeriod int) ([]float64, []float64) {
	kKey := fmt.Sprintf("k_stochrsi_%v", rsiPeriod)
	dKey := fmt.Sprintf("d_stochrsi_%v", rsiPeriod)

	if kkey, ok := lp.cache[kKey]; ok {
		if dkey, ok := lp.cache[dKey]; ok {
			return kkey, dkey
		}
	}

	k, d := lp.lines.StochRSI(rsiPeriod)
	lp.cache[kKey] = k
	lp.cache[dKey] = d
	return k, d
}

func (lp *LinesProxy) SMA(period int) []float64 {
	sKey := fmt.Sprintf("sma_%v", period)

	if skey, ok := lp.cache[sKey]; ok {
		return skey
	}

	s := lp.lines.SMA(period)
	lp.cache[sKey] = s
	return s
}

func (lp *LinesProxy) MACD() ([]float64, []float64) {
	const (
		macdKey   = "macd"
		signalKey = "signal"
	)

	if macd, ok := lp.cache[macdKey]; ok {
		if signal, ok := lp.cache[signalKey]; ok {
			return macd, signal
		}
	}

	macd, signal := lp.lines.MACD()
	lp.cache[macdKey] = macd
	lp.cache[signalKey] = signal
	return macd, signal
}

func (lp *LinesProxy) EMA() []float64 {
	const key = "ema"

	if val, ok := lp.cache[key]; ok {
		return val
	}

	val := lp.lines.EMA()
	lp.cache[key] = val
	return val
}

func LoadKlines(data []byte) *Lines {
	klines, err := UnmarshalKLines(data)
	if err != nil {
		log.Fatal(err)
	}

	t := &Lines{}
	for _, v := range klines.Candles {
		t.closing = append(t.closing, v.C)
		t.low = append(t.low, v.L)
		t.high = append(t.high, v.H)
	}

	return t
}

func LoadCandles(pair string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.exmo.com/v1.1/candles_history?symbol=%s&resolution=30&from=1703056979&to=1705476839", pair), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func main() {
	pair := "BTC_USD"
	candles := LoadCandles(pair)
	lines := LoadKlinesProxy(candles)
	lines.RSI(3)
}
