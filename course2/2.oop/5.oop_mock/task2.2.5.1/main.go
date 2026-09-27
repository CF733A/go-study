package main

import (
	"fmt"
	"time"
)

func main() {
	exchange := NewExmo()
	ind := NewIndicator(exchange, WithCalculateEMA(calculateEMA), WithCalculateSMA(calculateSMA))
	sma, err := ind.SMA("BTC_USD", 60, 5, time.Now().AddDate(0, 0, -1), time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("SMA:", sma)
	ema, err := ind.EMA("BTC_USD", 60, 5, time.Now().AddDate(0, 0, -1), time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("EMA:", ema)
}
