package main

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestMockExchanger_AllMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEx := NewMockExchanger(ctrl)

	mockEx.EXPECT().GetTicker().Return(nil, nil)
	mockEx.EXPECT().GetTrades("BTC_USD").Return(nil, nil)
	mockEx.EXPECT().GetOrderBook(10, "BTC_USD").Return(nil, nil)
	mockEx.EXPECT().GetCurrencies().Return(nil, nil)
	mockEx.EXPECT().GetCandlesHistory("BTC_USD", 30, gomock.Any(), gomock.Any()).Return(CandlesHistory{}, nil)
	mockEx.EXPECT().GetClosePrice("BTC_USD", 30, gomock.Any(), gomock.Any()).Return([]float64{1, 2, 3}, nil)

	_, _ = mockEx.GetTicker()
	_, _ = mockEx.GetTrades("BTC_USD")
	_, _ = mockEx.GetOrderBook(10, "BTC_USD")
	_, _ = mockEx.GetCurrencies()
	_, _ = mockEx.GetCandlesHistory("BTC_USD", 30, time.Now(), time.Now())
	_, _ = mockEx.GetClosePrice("BTC_USD", 30, time.Now(), time.Now())
}

func TestIndicator_SMA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEx := NewMockExchanger(ctrl)
	mockEx.
		EXPECT().
		GetClosePrice("BTC_USD", 60, gomock.Any(), gomock.Any()).
		Return([]float64{1, 2, 3, 4, 5}, nil)

	ind := NewIndicator(mockEx, WithCalculateSMA(calculateSMA))
	sma, err := ind.SMA("BTC_USD", 60, 3, time.Now().Add(-time.Hour), time.Now())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(sma) == 0 {
		t.Error("expected SMA result, got empty")
	}
}

func TestIndicator_EMA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEx := NewMockExchanger(ctrl)
	mockEx.
		EXPECT().
		GetClosePrice("BTC_USD", 60, gomock.Any(), gomock.Any()).
		Return([]float64{1, 2, 3, 4, 5}, nil)

	ind := NewIndicator(mockEx, WithCalculateEMA(calculateEMA))
	ema, err := ind.EMA("BTC_USD", 60, 3, time.Now().Add(-time.Hour), time.Now())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(ema) == 0 {
		t.Error("expected EMA result, got empty")
	}
}

func TestIndicator_ErrorFromExchange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEx := NewMockExchanger(ctrl)
	mockEx.
		EXPECT().
		GetClosePrice("BTC_USD", 60, gomock.Any(), gomock.Any()).
		Return(nil, errors.New("exchange error"))

	ind := NewIndicator(mockEx, WithCalculateSMA(calculateSMA))
	_, err := ind.SMA("BTC_USD", 60, 3, time.Now().Add(-time.Hour), time.Now())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestNewIndicator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEx := NewMockExchanger(ctrl)

	ind := NewIndicator(mockEx, WithCalculateSMA(calculateSMA), WithCalculateEMA(calculateEMA))
	if ind == nil {
		t.Error("expected indicator to be created, got nil")
	}
}

func TestNewIndicatorWithoutOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEx := NewMockExchanger(ctrl)

	ind := NewIndicator(mockEx)
	if ind == nil {
		t.Error("expected indicator to be created, got nil")
	}
}

func TestIndicator_SMA_WithEmptyData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEx := NewMockExchanger(ctrl)
	mockEx.
		EXPECT().
		GetClosePrice("BTC_USD", 30, gomock.Any(), gomock.Any()).
		Return([]float64{}, nil)

	ind := NewIndicator(mockEx, WithCalculateSMA(calculateSMA))
	sma, err := ind.SMA("BTC_USD", 30, 5, time.Now().Add(-time.Hour), time.Now())

	if err == nil {
		t.Error("expected error for empty data, got nil")
	}
	if sma != nil {
		t.Error("expected nil SMA result for empty data")
	}
}

func TestIndicator_EMA_WithEmptyData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEx := NewMockExchanger(ctrl)
	mockEx.
		EXPECT().
		GetClosePrice("BTC_USD", 30, gomock.Any(), gomock.Any()).
		Return([]float64{}, nil)

	ind := NewIndicator(mockEx, WithCalculateEMA(calculateEMA))
	ema, err := ind.EMA("BTC_USD", 30, 5, time.Now().Add(-time.Hour), time.Now())

	if err == nil {
		t.Error("expected error for empty data, got nil")
	}
	if ema != nil {
		t.Error("expected nil EMA result for empty data")
	}
}

func TestCalculateSMA(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	period := 3

	result := calculateSMA(data, period)
	if result == nil {
		t.Error("expected SMA calculation result, got nil")
	}
}

func TestCalculateEMA(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	period := 3

	result := calculateEMA(data, period)
	if result == nil {
		t.Error("expected EMA calculation result, got nil")
	}
}

func TestNewExmo(t *testing.T) {
	exmo := NewExmo()
	if exmo == nil {
		t.Error("expected Exmo instance, got nil")
	}
}

func TestExmo_GetTicker(t *testing.T) {
	exmo := NewExmo()
	ticker, err := exmo.GetTicker()
	if err != nil {
		t.Errorf("unexpected error for GetTicker: %v", err)
	}
	if ticker == nil {
		t.Error("expected ticker data, got nil")
	}
}

func TestExmo_GetTrades(t *testing.T) {
	exmo := NewExmo()
	trades, err := exmo.GetTrades("BTC_USD")
	if err != nil {
		t.Errorf("unexpected error for GetTrades: %v", err)
	}
	if trades == nil {
		t.Error("expected trades data, got nil")
	}
}

func TestExmo_GetOrderBook(t *testing.T) {
	exmo := NewExmo()
	orderBook, err := exmo.GetOrderBook(10, "BTC_USD")
	if err != nil {
		t.Errorf("unexpected error for GetOrderBook: %v", err)
	}
	if orderBook == nil {
		t.Error("expected order book data, got nil")
	}
}

func TestExmo_GetCurrencies(t *testing.T) {
	exmo := NewExmo()
	currencies, err := exmo.GetCurrencies()
	if err != nil {
		t.Errorf("unexpected error for GetCurrencies: %v", err)
	}
	if currencies == nil {
		t.Error("expected currencies data, got nil")
	}
}

func TestExmo_GetCandlesHistory(t *testing.T) {
	exmo := NewExmo()
	candles, err := exmo.GetCandlesHistory("BTC_USD", 30, time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		t.Errorf("unexpected error for GetCandlesHistory: %v", err)
	}
	if candles.Candles == nil {
		t.Error("expected currencies data, got nil")
	}
}

func TestExmo_GetClosePrice(t *testing.T) {
	exmo := NewExmo()
	prices, err := exmo.GetClosePrice("BTC_USD", 30, time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		t.Errorf("unexpected error for GetClosePrice: %v", err)
	}
	if prices == nil {
		t.Error("expected prices data, got nil")
	}
}

func TestNewExmo_WithCustomOptions(t *testing.T) {
	exmo := NewExmo(func(e *Exmo) {})
	if exmo == nil {
		t.Error("expected Exmo instance with options, got nil")
	}
}

func TestCandlesHistory_Unmarshal(t *testing.T) {
	jsonData := `{
		"candles": [
			{
				"t": 1640995200,
				"o": 100.0,
				"c": 105.0,
				"h": 110.0,
				"l": 95.0,
				"v": 1000.0
			}
		]
	}`

	ch, err := UnmarshalCandlesHistory([]byte(jsonData))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ch.Candles) != 1 {
		t.Errorf("expected 1 candle, got %d", len(ch.Candles))
	}

	candle := ch.Candles[0]
	if candle.T != 1640995200 {
		t.Errorf("expected timestamp 1640995200, got %d", candle.T)
	}
	if candle.O != 100.0 {
		t.Errorf("expected open 100.0, got %f", candle.O)
	}
}

func TestCandlesHistory_Marshal(t *testing.T) {
	ch := CandlesHistory{
		Candles: []Candle{
			{
				T: 1640995200,
				O: 100.0,
				C: 105.0,
				H: 110.0,
				L: 95.0,
				V: 1000.0,
			},
		},
	}

	data, err := ch.Marshal()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected marshaled data, got empty")
	}
}

func TestOrderBook_Unmarshal(t *testing.T) {
	jsonData := `{
		"BTC_USD": {
			"ask_quantity": "1.5",
			"ask_amount": "45000",
			"ask_top": "30000",
			"bid_quantity": "2.0",
			"bid_amount": "60000",
			"bid_top": "30000",
			"ask": [["30000", "1.5"]],
			"bid": [["30000", "2.0"]]
		}
	}`

	ob, err := UnmarshalOrderBook([]byte(jsonData))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ob) != 1 {
		t.Errorf("expected 1 pair, got %d", len(ob))
	}

	pair, exists := ob["BTC_USD"]
	if !exists {
		t.Error("expected BTC_USD pair, not found")
	}

	if pair.AskQuantity != "1.5" {
		t.Errorf("expected ask_quantity 1.5, got %s", pair.AskQuantity)
	}
}

func TestOrderBook_Marshal(t *testing.T) {
	ob := OrderBook{
		"BTC_USD": OrderBookPair{
			AskQuantity: "1.5",
			AskAmount:   "45000",
			AskTop:      "30000",
			BidQuantity: "2.0",
			BidAmount:   "60000",
			BidTop:      "30000",
			Ask:         [][]string{{"30000", "1.5"}},
			Bid:         [][]string{{"30000", "2.0"}},
		},
	}

	data, err := ob.Marshal()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected marshaled data, got empty")
	}
}

func TestTicker_Unmarshal(t *testing.T) {
	jsonData := `{
		"BTC_USD": {
			"buy_price": "30000",
			"sell_price": "30001",
			"last_trade": "30000",
			"high": "31000",
			"low": "29000",
			"avg": "30000",
			"vol": "100",
			"vol_curr": "3000000",
			"updated": 1640995200
		}
	}`

	ticker, err := UnmarshalTicker([]byte(jsonData))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ticker) != 1 {
		t.Errorf("expected 1 pair, got %d", len(ticker))
	}

	value, exists := ticker["BTC_USD"]
	if !exists {
		t.Error("expected BTC_USD pair, not found")
	}

	if value.BuyPrice != "30000" {
		t.Errorf("expected buy_price 30000, got %s", value.BuyPrice)
	}
}

func TestTicker_Marshal(t *testing.T) {
	ticker := Ticker{
		"BTC_USD": TickerValue{
			BuyPrice:  "30000",
			SellPrice: "30001",
			LastTrade: "30000",
			High:      "31000",
			Low:       "29000",
			Avg:       "30000",
			Vol:       "100",
			VolCurr:   "3000000",
			Updated:   1640995200,
		},
	}

	data, err := ticker.Marshal()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected marshaled data, got empty")
	}
}

func TestTrades_Unmarshal(t *testing.T) {
	jsonData := `{
		"BTC_USD": [
			{
				"trade_id": 12345,
				"date": 1640995200,
				"type": "buy",
				"quantity": "1.5",
				"price": "30000",
				"amount": "45000"
			}
		]
	}`

	trades, err := UnmarshalTrades([]byte(jsonData))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(trades) != 1 {
		t.Errorf("expected 1 pair, got %d", len(trades))
	}

	pairs, exists := trades["BTC_USD"]
	if !exists {
		t.Error("expected BTC_USD pair, not found")
	}

	if len(pairs) != 1 {
		t.Errorf("expected 1 trade, got %d", len(pairs))
	}

	trade := pairs[0]
	if trade.TradeID != 12345 {
		t.Errorf("expected trade_id 12345, got %d", trade.TradeID)
	}
	if trade.Type != Buy {
		t.Errorf("expected type buy, got %s", trade.Type)
	}
}

func TestTrades_Marshal(t *testing.T) {
	trades := Trades{
		"BTC_USD": []Pair{
			{
				TradeID:  12345,
				Date:     1640995200,
				Type:     Buy,
				Quantity: "1.5",
				Price:    "30000",
				Amount:   "45000",
			},
		},
	}

	data, err := trades.Marshal()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected marshaled data, got empty")
	}
}

func TestType_Constants(t *testing.T) {
	if Buy != "buy" {
		t.Errorf("expected Buy to be 'buy', got %s", Buy)
	}
	if Sell != "sell" {
		t.Errorf("expected Sell to be 'sell', got %s", Sell)
	}
}

func TestCandlesHistory_UnmarshalInvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`

	_, err := UnmarshalCandlesHistory([]byte(invalidJSON))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestOrderBook_UnmarshalInvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`

	_, err := UnmarshalOrderBook([]byte(invalidJSON))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestTicker_UnmarshalInvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`

	_, err := UnmarshalTicker([]byte(invalidJSON))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestTrades_UnmarshalInvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`

	_, err := UnmarshalTrades([]byte(invalidJSON))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestMain_Integration(t *testing.T) {
	originalStdout := os.Stdout

	tmpFile, err := os.CreateTemp("", "test_output")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	os.Stdout = tmpFile
	main()
	os.Stdout = originalStdout

	info, err := tmpFile.Stat()
	if err != nil {
		t.Fatalf("failed to get file info: %v", err)
	}

	if info.Size() == 0 {
		t.Error("expected main function to produce output, got empty")
	}
}

// func TestIndicator_WithRealExchange(t *testing.T) {
// 	exchange := NewExmo()
// 	indicator := NewIndicator(exchange, WithCalculateSMA(calculateSMA), WithCalculateEMA(calculateEMA))

// 	sma, err := indicator.SMA("BTC_USD", 30, 5, time.Now().AddDate(0, 0, -2), time.Now())
// 	if err != nil {
// 		t.Errorf("unexpected error from SMA: %v", err)
// 	}
// 	if sma == nil {
// 		t.Error("expected SMA data, got nil")
// 	}

// 	ema, err := indicator.EMA("BTC_USD", 30, 5, time.Now().AddDate(0, 0, -2), time.Now())
// 	if err != nil {
// 		t.Errorf("unexpected error from EMA: %v", err)
// 	}
// 	if ema == nil {
// 		t.Error("expected EMA data, got nil")
// 	}
// }