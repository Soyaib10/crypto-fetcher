package main

func GetStats(prices []float64) (avg, high, low float64) {
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
