package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"crypto-fetcher/logger"
)

type PriceStore struct {
	data map[string][]float64
}

func NewPriceStore() *PriceStore {
	return &PriceStore{
		data: make(map[string][]float64),
	}
}

func (ps *PriceStore) SetPrices(symbol string, prices []float64) {
	ps.data[symbol] = prices
}

func (ps *PriceStore) GetPrices(symbol string) []float64 {
	return ps.data[symbol]
}

func (ps *PriceStore) GetStats(symbol string) (avg, high, low float64) {
	prices := ps.data[symbol]
	if len(prices) == 0 {
		return 0, 0, 0
	}

	sum := 0.0
	high = prices[0]
	low = prices[0]

	for _, p := range prices {
		sum += p
		if p > high {
			high = p
		}
		if p < low {
			low = p
		}
	}

	avg = sum / float64(len(prices))
	return
}

func FetchLast10Prices(symbol string) ([]float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=1m&limit=10", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var klines [][]interface{}
	if err := json.Unmarshal(body, &klines); err != nil {
		return nil, err
	}

	prices := make([]float64, 0, 10)
	for _, k := range klines {
		closeStr, ok := k[4].(string)
		if !ok {
			continue
		}
		price, err := strconv.ParseFloat(closeStr, 64)
		if err != nil {
			continue
		}
		prices = append(prices, price)
	}

	return prices, nil
}

func main() {
	logger.Init()

	symbol := "BTCUSDT"

	logger.Info("Fetching last 10 prices for %s", symbol)

	prices, err := FetchLast10Prices(symbol)
	if err != nil {
		logger.Error("Failed to fetch prices: %v", err)
		return
	}

	store := NewPriceStore()
	store.SetPrices(symbol, prices)

	avg, high, low := store.GetStats(symbol)

	// CSV save
	err = SaveToCSV(symbol, store.GetPrices(symbol))
	if err != nil {
		logger.Error("Failed to save CSV: %v", err)
	}

	// JSON save
	exportData := ExportData{
		Symbol:  symbol,
		Prices:  store.GetPrices(symbol),
		Average: avg,
		Highest: high,
		Lowest:  low,
	}

	err = SaveToJSON(exportData)
	if err != nil {
		logger.Error("Failed to save JSON: %v", err)
	}

	logger.Info("Data successfully exported to CSV and JSON")

	fmt.Println("====================================")
	fmt.Printf("Symbol: %s\n", symbol)
	fmt.Printf("Last 10 Close Prices: %v\n", store.GetPrices(symbol))
	fmt.Printf("Average Price: %.2f\n", avg)
	fmt.Printf("Highest Price: %.2f\n", high)
	fmt.Printf("Lowest Price: %.2f\n", low)
	fmt.Println("====================================")
}
