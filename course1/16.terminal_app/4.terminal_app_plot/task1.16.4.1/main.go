package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
	"github.com/guptarohit/asciigraph"
)

type CurrencyData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func main() {
	selectCurrensy := make(chan string)
	defer close(selectCurrensy)

	resetChan := make(chan bool)
	defer close(resetChan)

	go asciiGraph(selectCurrensy, resetChan)

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	getMenu()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("ошибка создания считывателя кейборда")
			return
		}

		switch {
		case char == '1':
			selectCurrensy <- "BTCUSDT"
		case char == '2':
			selectCurrensy <- "LTCUSDT"
		case char == '3':
			selectCurrensy <- "ETHUSDT"
		case char == 'q':
			return
		case key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2:
			resetChan <- true
			getMenu()
		default:
			fmt.Println("нажата неизвестная программе клавиша")
		}
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getMenu() {
	clear()
	fmt.Println("1. BTC_USD\n2. LTC_USD\n3. ETH_USD\nPress 1-3 to change symbol, press q to exit")
}

func getHttpCource(currency string) (CurrencyData, error) {
	url := "https://api.binance.com/api/v3/ticker/price?symbol=" + currency

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return CurrencyData{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return CurrencyData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return CurrencyData{}, err
	}

	var cource CurrencyData
	err = json.Unmarshal(body, &cource)
	if err != nil {
		fmt.Println("ошибка десериализации")
	}

	return cource, nil
}

func getPrice(val string) (float64, error) {
	cource, err := getHttpCource(val)
	if err != nil {
		fmt.Println("ошибка получения валюты")
		return 0, err
	}

	transformedCource, err := strconv.ParseFloat(cource.Price, 64)
	if err != nil {
		fmt.Println("ошибка конвертации:", err)
		return 0, err
	}
	return transformedCource, nil
}

func asciiGraph(cur chan string, reset chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	courseArr := []float64{}
	var currentCurrency string

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	for {
		select {
		case currency := <-cur:
			currentCurrency = currency
			courseArr = []float64{}
			clear()
		case <-reset:
			currentCurrency = ""
			courseArr = []float64{}
			clear()
		case <-ticker.C:
			if currentCurrency != "" {

				price, err := getPrice(currentCurrency)
				if err != nil {
					fmt.Println("ошибка получения валюты в asciigraph:", err)
					break
				}

				if len(courseArr) < 100 {
					for i := 0; i < 100; i++ {
						courseArr = append(courseArr, price)
					}

				} else {
					courseArr = append(courseArr[1:], price)
				}

				
				fmt.Fprintf(writer, "%s, %.2f\n", currentCurrency, price)
				graph := makeGraph(courseArr)
				fmt.Fprintln(writer, graph)
				currentTime := time.Now()
				fmt.Fprintf(writer, "%s\n", currentTime.Format("15:04:05"))
				fmt.Fprintf(writer, "%s\n", currentTime.Format("2006-01-02"))
				

			}
		}
	}
}

func makeGraph(data []float64) string {
	graph := asciigraph.Plot(
		data,
		asciigraph.Height(10),
		asciigraph.Width(100),
	)

	return graph
}
