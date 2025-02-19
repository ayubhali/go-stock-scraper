package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Stock struct {
	company, price, change string
}

type AlphaVantageResponse struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Price           string `json:"05. price"`
		ChangePercent   string `json:"10. change percent"`
	} `json:"Global Quote"`
}

func main() {
	fmt.Println("Starting...")
	apiKey := os.Getenv("ALPHA_VANTAGE_KEY")
	if apiKey == "" {
		apiKey = "YOUR_API_KEY"
	}

	symbols := []string{"AAPL", "MSFT", "GOOGL"}
	var stocks []Stock

	for _, symbol := range symbols {
		fmt.Printf("Getting %s...\n", symbol)
		url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, apiKey)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var result AlphaVantageResponse
		json.Unmarshal(body, &result)

		stock := Stock{
			company: symbol,
			price:   result.GlobalQuote.Price,
			change:  result.GlobalQuote.ChangePercent,
		}

		fmt.Printf("%s: $%s (%s)\n", stock.company, stock.price, stock.change)
		stocks = append(stocks, stock)
		time.Sleep(1 * time.Second)
	}

	file, _ := os.Create("stocks.csv")
	writer := csv.NewWriter(file)
	writer.Write([]string{"Company", "Price", "Change"})
	for _, stock := range stocks {
		writer.Write([]string{stock.company, stock.price, stock.change})
	}
	writer.Flush()
	file.Close()

	fmt.Println("Done - check stocks.csv")
}