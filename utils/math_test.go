package utils

import (
	"testing"
)

func TestAverage(t *testing.T) {
	var averageTests = []struct {
		collection []float64
		expected   float64
	}{
		{[]float64{1, 1, 1}, 1},
		{[]float64{1, 2, 3}, 2},
		{[]float64{}, 0},
	}

	for _, test := range averageTests {
		if output := Average(test.collection); output != test.expected {
			t.Errorf("Test Failed: input %v, expected %v, recieved %v", test.collection, test.expected, output)
		}
	}
}
