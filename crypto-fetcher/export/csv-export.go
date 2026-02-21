package export

import (
	"crypto-fetcher/store"
	"encoding/csv"
	"os"
	"strconv"
)

type CSVExporter struct{}

func (c CSVExporter) Export(store *store.PriceStore, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Symbol", "Prices", "Avg", "High", "Low"})

	// Use getter
	dataCopy := store.GetAllData()

	for symbol, prices := range dataCopy {
		avg, high, low := store.CalculateStats(prices)
		pricesStr := ""
		for i, p := range prices {
			if i != 0 {
				pricesStr += ";"
			}
			pricesStr += strconv.FormatFloat(p, 'f', 2, 64)
		}
		writer.Write([]string{
			symbol,
			pricesStr,
			strconv.FormatFloat(avg, 'f', 2, 64),
			strconv.FormatFloat(high, 'f', 2, 64),
			strconv.FormatFloat(low, 'f', 2, 64),
		})
	}

	return nil
}