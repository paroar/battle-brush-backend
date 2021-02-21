package utils

// Average returns the average of a slice of floats
func Average(collection []float64) float64 {

	if len(collection) <= 0 {
		return 0.0
	}

	sum := 0.0
	for _, value := range collection {
		sum += value
	}
	return (sum / float64(len(collection)))
}
