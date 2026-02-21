package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func SaveToCSV(symbol string, prices []float64) error {
	file, err := os.Create(symbol + "_prices.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"symbol", "price"})

	for _, p := range prices {
		writer.Write([]string{
			symbol,
			strconv.FormatFloat(p, 'f', -1, 64),
		})
	}

	return nil
}