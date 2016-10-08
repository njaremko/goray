package main

import (
	"fmt"
	"os"
	"strconv"
)

// Median returns the median element in a slice of floats
func Median(items []float64) float64 {
	n := len(items)
	switch {
	case n == 0:
		return 0
	case n%2 == 1:
		return items[n/2]
	default:
		a := items[n/2-1]
		b := items[n/2]
		return (a + b) / 2
	}
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

// ParseFloats parses a string and returns a slice of float64s
func ParseFloats(items []string) []float64 {
	result := make([]float64, len(items))
	for i, item := range items {
		f, err := strconv.ParseFloat(item, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		result[i] = f
	}
	return result
}
