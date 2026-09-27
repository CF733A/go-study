package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIndicator struct {
	mock.Mock
}

func (m *MockIndicator) StochPrice() ([]float64, []float64) {
	args := m.Called()
	return args.Get(0).([]float64), args.Get(1).([]float64)
}

func (m *MockIndicator) RSI(period int) ([]float64, []float64) {
	args := m.Called(period)
	return args.Get(0).([]float64), args.Get(1).([]float64)
}

func (m *MockIndicator) StochRSI(rsiPeriod int) ([]float64, []float64) {
	args := m.Called(rsiPeriod)
	return args.Get(0).([]float64), args.Get(1).([]float64)
}

func (m *MockIndicator) SMA(period int) []float64 {
	args := m.Called(period)
	return args.Get(0).([]float64)
}

func (m *MockIndicator) MACD() ([]float64, []float64) {
	args := m.Called()
	return args.Get(0).([]float64), args.Get(1).([]float64)
}

func (m *MockIndicator) EMA() []float64 {
	args := m.Called()
	return args.Get(0).([]float64)
}

func TestLinesProxy_StochPrice(t *testing.T) {
	mockIndicator := new(MockIndicator)
	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	mockIndicator.On("StochPrice").Return([]float64{1, 2}, []float64{3, 4}).Once()
	k, d := proxy.StochPrice()
	assert.Equal(t, []float64{1, 2}, k)
	assert.Equal(t, []float64{3, 4}, d)

	k, d = proxy.StochPrice()
	assert.Equal(t, []float64{1, 2}, k)
	assert.Equal(t, []float64{3, 4}, d)

	mockIndicator.AssertExpectations(t)
}

func TestLinesProxy_RSI(t *testing.T) {
	mockIndicator := new(MockIndicator)
	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	mockIndicator.On("RSI", 14).Return([]float64{1}, []float64{2}).Once()
	rs, rsi := proxy.RSI(14)
	assert.Equal(t, []float64{1}, rs)
	assert.Equal(t, []float64{2}, rsi)

	rs, rsi = proxy.RSI(14)
	assert.Equal(t, []float64{1}, rs)
	assert.Equal(t, []float64{2}, rsi)

	mockIndicator.On("RSI", 7).Return([]float64{3}, []float64{4}).Once()
	rs, rsi = proxy.RSI(7)
	assert.Equal(t, []float64{3}, rs)
	assert.Equal(t, []float64{4}, rsi)

	mockIndicator.AssertExpectations(t)
}

func TestLinesProxy_StochRSI(t *testing.T) {
	mockIndicator := new(MockIndicator)
	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	mockIndicator.On("StochRSI", 14).Return([]float64{1}, []float64{2}).Once()
	k, d := proxy.StochRSI(14)
	assert.Equal(t, []float64{1}, k)
	assert.Equal(t, []float64{2}, d)

	k, d = proxy.StochRSI(14)
	assert.Equal(t, []float64{1}, k)
	assert.Equal(t, []float64{2}, d)

	mockIndicator.AssertExpectations(t)
}

func TestLinesProxy_SMA(t *testing.T) {
	mockIndicator := new(MockIndicator)
	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	mockIndicator.On("SMA", 10).Return([]float64{1, 2, 3}).Once()
	sma := proxy.SMA(10)
	assert.Equal(t, []float64{1, 2, 3}, sma)

	sma = proxy.SMA(10)
	assert.Equal(t, []float64{1, 2, 3}, sma)

	mockIndicator.On("SMA", 20).Return([]float64{4, 5}).Once()
	sma = proxy.SMA(20)
	assert.Equal(t, []float64{4, 5}, sma)

	mockIndicator.AssertExpectations(t)
}

func TestLinesProxy_MACD(t *testing.T) {
	mockIndicator := new(MockIndicator)
	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	mockIndicator.On("MACD").Return([]float64{1}, []float64{2}).Once()
	macd, signal := proxy.MACD()
	assert.Equal(t, []float64{1}, macd)
	assert.Equal(t, []float64{2}, signal)

	macd, signal = proxy.MACD()
	assert.Equal(t, []float64{1}, macd)
	assert.Equal(t, []float64{2}, signal)

	mockIndicator.AssertExpectations(t)
}

func TestLinesProxy_EMA(t *testing.T) {
	mockIndicator := new(MockIndicator)
	proxy := &LinesProxy{
		lines: mockIndicator,
		cache: make(map[string][]float64),
	}

	mockIndicator.On("EMA").Return([]float64{1, 2, 3}).Once()
	ema := proxy.EMA()
	assert.Equal(t, []float64{1, 2, 3}, ema)

	ema = proxy.EMA()
	assert.Equal(t, []float64{1, 2, 3}, ema)

	mockIndicator.AssertExpectations(t)
}

func TestLoadKlines(t *testing.T) {
	testData := []byte(`{
		"pair": "BTC_USD",
		"candles": [
			{"t": 123, "o": 1.1, "c": 1.2, "h": 1.3, "l": 1.0, "v": 100},
			{"t": 124, "o": 1.2, "c": 1.3, "h": 1.4, "l": 1.1, "v": 200}
		]
	}`)

	lines := LoadKlines(testData)
	assert.Equal(t, []float64{1.3, 1.4}, lines.high)
	assert.Equal(t, []float64{1.0, 1.1}, lines.low)
	assert.Equal(t, []float64{1.2, 1.3}, lines.closing)
}

func TestUnmarshalKLines(t *testing.T) {
	t.Run("Valid data", func(t *testing.T) {
		testData := []byte(`{"pair": "TEST", "candles": [{"t": 123, "o": 1.1, "c": 1.2, "h": 1.3, "l": 1.0, "v": 100}]}`)
		klines, err := UnmarshalKLines(testData)
		assert.NoError(t, err)
		assert.Equal(t, "TEST", klines.Pair)
		assert.Len(t, klines.Candles, 1)
	})

	t.Run("Invalid data", func(t *testing.T) {
		testData := []byte(`{"invalid": "data"`)
		_, err := UnmarshalKLines(testData)
		assert.Error(t, err)
	})
}

func TestLoadKlinesProxy(t *testing.T) {
	testData := []byte(`{
		"pair": "BTC_USD",
		"candles": [{"t": 123, "o": 1.1, "c": 1.2, "h": 1.3, "l": 1.0, "v": 100}]
	}`)

	proxy := LoadKlinesProxy(testData)
	assert.NotNil(t, proxy)
	assert.NotNil(t, proxy.cache)
	assert.NotNil(t, proxy.lines)
}

func TestLoadKlines_EmptyData(t *testing.T) {
	lines := LoadKlines([]byte(`{"pair": "BTC_USD", "candles": []}`))
	assert.Empty(t, lines.high)
	assert.Empty(t, lines.low)
	assert.Empty(t, lines.closing)
}

func TestKLines_Marshal(t *testing.T) {
	klines := KLines{
		Pair: "TEST",
		Candles: []Candle{
			{T: 123, O: 1.1, C: 1.2, H: 1.3, L: 1.0, V: 100},
		},
	}
	data, err := klines.Marshal()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestKLines_MarshalError(t *testing.T) {
	klines := KLines{
		Pair: "TEST",
		Candles: []Candle{
			{T: 123, O: 1.1, C: 1.2, H: 1.3, L: 1.0, V: math.Inf(1)},
		},
	}
	_, err := klines.Marshal()
	assert.Error(t, err)
}
