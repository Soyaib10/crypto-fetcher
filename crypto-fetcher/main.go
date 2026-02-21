package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// --- Logger (simple) ---
func logInfo(format string, a ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", a...)
}
func logError(format string, a ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", a...)
}

// --- PriceStore ---
type PriceStore struct {
	data map[string][]float64
	mu   sync.Mutex
}

func NewPriceStore() *PriceStore {
	return &PriceStore{
		data: make(map[string][]float64),
	}
}

func (ps *PriceStore) SetPrices(symbol string, prices []float64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.data[symbol] = prices
}

func (ps *PriceStore) GetPrices(symbol string) []float64 {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.data[symbol]
}

func (ps *PriceStore) GetStats(symbol string) (avg, high, low float64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	prices := ps.data[symbol]
	if len(prices) == 0 {
		return 0, 0, 0
	}
	high, low = prices[0], prices[0]
	sum := 0.0
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

// --- Fetch Function with Timeout ---
func FetchLast10PricesWithTimeout(symbol string, timeout time.Duration) ([]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=1m&limit=10", symbol)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
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
		p, _ := strconv.ParseFloat(closeStr, 64)
		prices = append(prices, p)
	}

	return prices, nil
}

// --- Main ---
func main() {
	symbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}
	store := NewPriceStore()

	var wg sync.WaitGroup
	wg.Add(len(symbols))

	for _, sym := range symbols {
		go func(s string) {
			defer wg.Done()

			prices, err := FetchLast10PricesWithTimeout(s, 5*time.Second)
			if err != nil {
				logError("Failed to fetch %s: %v", s, err)
				return
			}

			store.SetPrices(s, prices)
			logInfo("Fetched %s prices: %v", s, prices)

			avg, high, low := store.GetStats(s)
			logInfo("%s → Avg: %.2f, High: %.2f, Low: %.2f", s, avg, high, low)

		}(sym)
	}

	wg.Wait()
	logInfo("All symbols fetched")

	err := SaveToCSV(store, "prices.csv")
	if err != nil {
		logError("CSV export failed: %v", err)
	} else {
		logInfo("Prices saved to prices.csv")
	}

	err = ExportData(store, "prices.json")
	if err != nil {
		logError("JSON export failed: %v", err)
	} else {
		logInfo("Prices saved to prices.json")
	}
}
