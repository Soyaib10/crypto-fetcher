package main

import (
	"encoding/json"
	"os"
)

func ExportData(store *PriceStore, filename string) error {
	store.mu.Lock()
	dataCopy := make(map[string][]float64, len(store.data))
	for symbol, prices := range store.data {
		pricesCopy := make([]float64, len(prices))
		copy(pricesCopy, prices)
		dataCopy[symbol] = pricesCopy
	}
	store.mu.Unlock() // Unlock early

	// Step 2: Prepare final data with stats
	finalData := map[string]interface{}{}
	for symbol, prices := range dataCopy {
		avg, high, low := GetStats(prices) // Deadlock-free
		finalData[symbol] = map[string]interface{}{
			"prices": prices,
			"avg":    avg,
			"high":   high,
			"low":    low,
		}
	}

	// Step 3: Write JSON to file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(finalData)
}
