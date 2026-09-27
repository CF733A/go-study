package main

import (
	"errors"
	"fmt"
	"io"
	"log"
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

//go:generate mockgen -source=main.go -destination=mock_client.go -package=main

type GeneralIndicatorer interface {
	GetData(pair string, period int, from, to time.Time, indicator Indicatorer) ([]float64, error)
}

type GeneralIndicator struct{}

func (g *GeneralIndicator) GetData(pair string, period int, from, to time.Time, indicator Indicatorer) ([]float64, error) {
	return indicator.GetData(pair, period, from, to)
}

type Indicatorer interface {
	GetData(pair string, period int, from, to time.Time) ([]float64, error)
}

type IndicatorSMA struct {
	exchange Exchanger
}

func NewIndicatorSMA(exchange Exchanger) *IndicatorSMA {
	return &IndicatorSMA{
		exchange: exchange,
	}
}

func (i *IndicatorSMA) GetData(pair string, period int, from, to time.Time) ([]float64, error) {
	closedPrice, err := i.exchange.GetClosePrice(pair, period, from, to)
	if err != nil {
		return closedPrice, fmt.Errorf("failed get data SMA: %w", err)
	}

	sma := calculateSMA(closedPrice, 3)
	return sma, nil
}

type IndicatorEMA struct {
	exchange Exchanger
}

func NewIndicatorEMA(exchange Exchanger) *IndicatorEMA {
	return &IndicatorEMA{
		exchange: exchange,
	}
}

func (i *IndicatorEMA) GetData(pair string, period int, from, to time.Time) ([]float64, error) {
	closedPrice, err := i.exchange.GetClosePrice(pair, period, from, to)
	if err != nil {
		return closedPrice, fmt.Errorf("failed get data EMA: %w", err)
	}

	sma := calculateSMA(closedPrice, 3)
	ema := calculateEMA(sma, 3)
	return ema, nil
}

// Функция для расчета простого скользящего среднего (SMA)
func calculateSMA(data []float64, period int) []float64 {
	var sma = make([]float64, len(data)/period)
	for i := range sma {
		sum := 0.0
		for _, d := range data[i*period : i*period+period] {
			sum += d
		}
		sma[i] = sum / float64(period)
	}

	return sma
}

// Функция для расчета экспоненциального скользящего среднего (EMA)
func calculateEMA(data []float64, period int) []float64 {
	if len(data) == 0 || period <= 0 {
		return nil
	}

	alpha := 2.0 / (float64(period) + 1.0)
	ema := make([]float64, len(data))

	ema[0] = data[0]

	for i := 1; i < len(data); i++ {
		ema[i] = alpha*data[i] + (1-alpha)*ema[i-1]
	}

	return ema
}

type Exchanger interface {
	GetTicker() (Ticker, error)
	GetTrades(pairs ...string) (Trades, error)
	GetOrderBook(limit int, pairs ...string) (OrderBook, error)
	GetCurrencies() (Currencies, error)
	GetCandlesHistory(pair string, resolution int, start, end time.Time) (CandlesHistory, error)
	GetClosePrice(pair string, resolution int, start, end time.Time) ([]float64, error)
}

type Exmo struct {
	client *http.Client
	url    string
}

func NewExmo() Exchanger {
	return &Exmo{
		client: &http.Client{},
		url:    "https://api.exmo.com/v1",
	}
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
	currency, err := UnmarshalCurrencies(result)
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

func (e *Exmo) DoRequest(endpoint string, params url.Values) ([]byte, error) {
	reqURL := e.url + endpoint

	if params != nil {
		reqURL += "?" + params.Encode()
	}

	response, err := e.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа %s", err)
	}

	return body, nil
}

func main() {
	var exchange Exchanger
	exchange = NewExmo()

	indicatorSMA := NewIndicatorSMA(exchange)

	generalIndicator := &GeneralIndicator{}
	sma, err := generalIndicator.GetData("BTC_USD", 30, time.Now().Add(-time.Hour*24*5), time.Now(), indicatorSMA)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sma)

	indicatorEMA := NewIndicatorEMA(exchange)
	ema, err := generalIndicator.GetData("BTC_USD", 30, time.Now().Add(-time.Hour*24*5), time.Now(), indicatorEMA)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ema)
}
