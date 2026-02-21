package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)



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