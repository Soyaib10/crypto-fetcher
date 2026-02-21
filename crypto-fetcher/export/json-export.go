package export

import (
	"crypto-fetcher/store"
	"encoding/json"
	"os"
)

type JSONExporter struct{}

func (j JSONExporter) Export(store *store.PriceStore, filename string) error {
	// Use getter
	dataCopy := store.GetAllData()

	finalData := map[string]interface{}{}
	for symbol, prices := range dataCopy {
		avg, high, low := store.CalculateStats(prices)
		finalData[symbol] = map[string]interface{}{
			"prices": prices,
			"avg":    avg,
			"high":   high,
			"low":    low,
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(finalData)
}