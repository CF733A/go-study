package main

import (
	"errors"
	"time"

	"github.com/cinar/indicator"
)

type Indicator struct {
	exchange     Exchanger
	calculateSMA func(data []float64, period int) []float64
	calculateEMA func(data []float64, period int) []float64
}

type IndicatorOption func(*Indicator)

type Indicatorer interface {
	SMA(pair string, resolution, period int, from, to time.Time) ([]float64, error)
	EMA(pair string, resolution, period int, from, to time.Time) ([]float64, error)
}

func (i *Indicator) SMA(pair string, resolution, period int, from, to time.Time) ([]float64, error) {
	if period <= 0 {
		return nil, errors.New("period must be positive")
	}

	data, err := i.exchange.GetClosePrice(pair, resolution, from, to)
	if err != nil {
		return nil, err
	}

	if len(data) < period {
		return nil, errors.New("not enough data points for the specified period")
	}

	return i.calculateSMA(data, period), nil
}

func (i *Indicator) EMA(pair string, resolution, period int, from, to time.Time) ([]float64, error) {
	if period <= 0 {
		return nil, errors.New("period must be positive")
	}

	data, err := i.exchange.GetClosePrice(pair, resolution, from, to)
	if err != nil {
		return nil, err
	}

	if len(data) < period {
		return nil, errors.New("not enough data points for the specified period")
	}

	return i.calculateEMA(data, period), nil
}

func NewIndicator(exchange Exchanger, opts ...IndicatorOption) Indicatorer {
	ind := &Indicator{
		exchange: exchange,
	}

	for _, opt := range opts {
		opt(ind)
	}

	return ind
}

func calculateSMA(data []float64, period int) []float64 {
	return indicator.Sma(period, data)
}

func calculateEMA(data []float64, period int) []float64 {
	return indicator.Ema(period, data)
}

func WithCalculateSMA(SMAfunc func(data []float64, period int) []float64) IndicatorOption {
	return func(ind *Indicator) {
		ind.calculateSMA = SMAfunc
	}
}

func WithCalculateEMA(EMAfunc func(data []float64, period int) []float64) IndicatorOption {
	return func(ind *Indicator) {
		ind.calculateEMA = EMAfunc
	}
}
