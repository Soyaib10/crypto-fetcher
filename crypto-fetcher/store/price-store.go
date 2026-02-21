package store

import "sync"

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

// Getter function: returns a copy of internal map (thread-safe)
func (ps *PriceStore) GetAllData() map[string][]float64 {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	dataCopy := make(map[string][]float64, len(ps.data))
	for k, v := range ps.data {
		prices := make([]float64, len(v))
		copy(prices, v)
		dataCopy[k] = prices
	}
	return dataCopy
}

// Deadlock-free stats calculation
func (ps *PriceStore) CalculateStats(prices []float64) (avg, high, low float64) {
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