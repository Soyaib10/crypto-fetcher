package main

import (
	"sync"
	"time"

	"crypto-fetcher/logger"
	"crypto-fetcher/store"
	"crypto-fetcher/fetch"
	"crypto-fetcher/export"
)

func main() {
	symbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}
	priceStore := store.NewPriceStore()

	var wg sync.WaitGroup
	wg.Add(len(symbols))

	for _, sym := range symbols {
		go func(s string) {
			defer wg.Done()
			prices, err := fetch.FetchLast10PricesWithTimeout(s, 10*time.Second)
			if err != nil {
				logger.LogError("Failed to fetch %s: %v", s, err)
				return
			}
			priceStore.SetPrices(s, prices)
			logger.LogInfo("Fetched %s prices: %v", s, prices)
		}(sym)
	}

	wg.Wait()
	logger.LogInfo("All symbols fetched")

	exporters := []struct {
		e export.Exporter
		f string
	}{
		{export.CSVExporter{}, "prices.csv"},
		{export.JSONExporter{}, "prices.json"},
	}

	for _, exp := range exporters {
		if err := exp.e.Export(priceStore, exp.f); err != nil {
			logger.LogError("Export failed for %s: %v", exp.f, err)
		} else {
			logger.LogInfo("Exported successfully: %s", exp.f)
		}
	}
}