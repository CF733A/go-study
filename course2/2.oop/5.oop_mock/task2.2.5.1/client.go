package main

import (
	"errors"
	"fmt"
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

//go:generate mockgen -source=client.go -destination=mock_client.go -package=main

type Exmo struct {
	client *http.Client
	url    string
}

func NewExmo(opts ...func(exmo *Exmo)) Exchanger {
	exmo := &Exmo{
		client: &http.Client{},
		url:    "https://api.exmo.com/v1",
	}
	for _, opt := range opts {
		opt(exmo)
	}
	return exmo
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

type Exchanger interface {
	GetTicker() (Ticker, error)
	GetTrades(pairs ...string) (Trades, error)
	GetOrderBook(limit int, pairs ...string) (OrderBook, error)
	GetCurrencies() (Currencies, error)
	GetCandlesHistory(pair string, period int, start, end time.Time) (CandlesHistory, error)
	GetClosePrice(pair string, resolution int, start, end time.Time) ([]float64, error)
}

func (e *Exmo) GetTicker() (Ticker, error) {
	result, err := e.DoRequest(ticker, nil)
	if err != nil {
		return nil, errors.New("не смог получить тикер getticker")
	}
	ticker, err := UnmarshalTicker(result)
	return ticker, nil
}

func (e *Exmo) GetTrades(pairs ...string) (Trades, error) {
	params := url.Values{}
	if len(pairs) > 0 {
		params.Add("pair", strings.Join(pairs, ","))
	}
	result, err := e.DoRequest(trades, params)
	trades, err := UnmarshalTrades(result)
	return trades, err
}

func (e *Exmo) GetOrderBook(limit int, pairs ...string) (OrderBook, error) {
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}
	if len(pairs) > 0 {
		params.Add("pair", strings.Join(pairs, ","))
	}
	result, err := e.DoRequest(orderBook, params)
	orderBook, err := UnmarshalOrderBook(result)
	return orderBook, err
}

func (e *Exmo) GetCurrencies() (Currencies, error) {
	result, err := e.DoRequest(currency, nil)
	currency, err := UnmarshalCurrency(result)
	return currency, err
}

func (e *Exmo) GetCandlesHistory(pair string, period int, start, end time.Time) (CandlesHistory, error) {

	params := url.Values{}
	params.Add("symbol", pair)
	params.Add("resolution", strconv.Itoa(period))
	params.Add("from", fmt.Sprintf("%d", start.Unix()))
	params.Add("to", fmt.Sprintf("%d", end.Unix()))

	result, err := e.DoRequest(candlesHistory, params)
	candlesHistory, err := UnmarshalCandlesHistory(result)
	return candlesHistory, err
}

func (e *Exmo) GetClosePrice(pair string, limit int, start, end time.Time) ([]float64, error) {
	candleHistory, err := e.GetCandlesHistory(pair, limit, start, end)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения свечей: %v", err)
	}

	var closePrice []float64
	for _, cancandle := range candleHistory.Candles {
		closePrice = append(closePrice, cancandle.C)
	}

	return closePrice, err
}
