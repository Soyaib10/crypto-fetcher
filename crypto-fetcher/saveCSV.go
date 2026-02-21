package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

// Save prices to CSV
func SaveToCSV(store *PriceStore, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	writer.Write([]string{"Symbol", "Prices", "Avg", "High", "Low"})

	// Data
	for symbol := range store.data {
		prices := store.GetPrices(symbol)
		avg, high, low := store.GetStats(symbol)
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