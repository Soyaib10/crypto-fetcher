package main

import (
	"encoding/json"
	"os"
)

type ExportData struct {
	Symbol  string    `json:"symbol"`
	Prices  []float64 `json:"prices"`
	Average float64   `json:"average"`
	Highest float64   `json:"highest"`
	Lowest  float64   `json:"lowest"`
}

func SaveToJSON(data ExportData) error {
	file, err := os.Create(data.Symbol + "_prices.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}