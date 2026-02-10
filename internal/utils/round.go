package utils

import (
	"math"
)

func Round(v float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(v*pow) / pow
}
