package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ticker         = "/ticker"
	trades         = "/trades"
	orderBook      = "/order_book"
	currency       = "/currency"
	candlesHistory = "/candles_history"
)

type Ticker map[string]TickerValue

type TickerValue struct {
	BuyPrice  string `json:"buy_price"`
	SellPrice string `json:"sell_price"`
	LastTrade string `json:"last_trade"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Avg       string `json:"avg"`
	Vol       string `json:"vol"`
	VolCurr   string `json:"vol_curr"`
	Updated   int64  `json:"updated"`
}

type Trades map[string][]Pair

type Pair struct {
	TradeID  int64   `json:"trade_id"`
	Type     string  `json:"type"`
	Price    float64 `json:"price,string"`
	Quantity float64 `json:"quantity,string"`
	Amount   float64 `json:"amount,string"`
	Date     int64   `json:"date"`
}

type OrderBook map[string]OrderBookPair

type OrderBookPair struct {
	AskQuantity string     `json:"ask_quantity"`
	AskAmount   string     `json:"ask_amount"`
	AskTop      string     `json:"ask_top"`
	BidQuantity string     `json:"bid_quantity"`
	BidAmount   string     `json:"bid_amount"`
	BidTop      string     `json:"bid_top"`
	Ask         [][]string `json:"ask"`
	Bid         [][]string `json:"bid"`
}

type Currencies []string

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

type Exmo struct {
	url    string
	client *http.Client
}

func NewExmo(opts ...func(*Exmo)) *Exmo {
    ex := &Exmo{
        url:    "https://api.exmo.com/v1",
        client: &http.Client{Timeout: 30 * time.Second},
    }
    
    for _, opt := range opts {
        opt(ex)
    }
    
    return ex
}

func WithClient(client *http.Client) func(*Exmo) {
	return func(e *Exmo) {
		e.client = client
	}
}

func WithURL(url string) func(*Exmo) {
	return func(e *Exmo) {
		e.url = url
	}
}

func (e *Exmo) doRequest(endpoint string, params url.Values, result interface{}) error {
	reqURL := e.url + endpoint

	if params != nil {
		reqURL += "?" + params.Encode()
	}

	response, err := e.client.Get(reqURL)
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа %s", err)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("ошибка десериализации ответа %s", err)
	}

	return nil
}

type Exchanger interface {
	GetTicker() (Ticker, error)
	GetTrades(pairs ...string) (Trades, error)
	GetOrderBook(limit int, pairs ...string) (OrderBook, error)
	GetCurrencies() (Currencies, error)
	GetCandlesHistory(pair string, limit int, start, end time.Time) (CandlesHistory, error)
	GetClosePrice(pair string, limit int, start, end time.Time) ([]float64, error)
}

func (e *Exmo) GetTicker() (Ticker, error) {
	var result Ticker
	err := e.doRequest(ticker, nil, &result)
	return result, err
}

func (e *Exmo) GetTrades(pairs ...string) (Trades, error) {
	var result Trades
	params := url.Values{}
	if len(pairs) > 0 {
		params.Add("pair", strings.Join(pairs, ","))
	}
	err := e.doRequest(trades, params, &result)
	return result, err
}

func (e *Exmo) GetOrderBook(limit int, pairs ...string) (OrderBook, error) {
	var result OrderBook
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}
	if len(pairs) > 0 {
		params.Add("pair", strings.Join(pairs, ","))
	}
	err := e.doRequest(orderBook, params, &result)
	return result, err
}

func (e *Exmo) GetCurrencies() (Currencies, error) {
	var result Currencies
	err := e.doRequest(currency, nil, &result)
	return result, err
}

func (e *Exmo) GetCandlesHistory(pair string, limit int, start, end time.Time) (CandlesHistory, error) {
	var result CandlesHistory

	params := url.Values{}
	params.Add("symbol", pair)
	params.Add("resolution", strconv.Itoa(limit))
	params.Add("from", fmt.Sprintf("%d", start.Unix()))
	params.Add("to", fmt.Sprintf("%d", end.Unix()))

	err := e.doRequest(candlesHistory, params, &result)
	return result, err
}

func (e *Exmo) GetClosePrice(pair string, limit int, start, end time.Time) ([]float64, error) {

	candles, err := e.GetCandlesHistory(pair, limit, start, end)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения свечей: %v", err)
	}

	var closePrice []float64
	for _, cancandle := range candles.Candles {
		closePrice = append(closePrice, cancandle.C)
	}

	return closePrice, err
}

func main() {
	var exchange Exchanger
	exchange = NewExmo()
	ticker, err := exchange.GetCandlesHistory("BTC_USD", 30, time.Now().Add(-time.Hour*24), time.Now())
	if err != nil {
		return
	}
	fmt.Println(ticker)
}