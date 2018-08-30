package zxcvbnmath

import (
	"math"
	"math/big"
)


// NChoseK returns the binomial co-efficient taking and returning float64
// It is simply a type adjusting wrapper for big.Binomial()
func NChoseK(n, k float64) float64 {
	if k > n {
		return 0
	} else if k == 0 {
		return 1
	}
	coef := new(big.Int).Binomial(int64(n), int64(k))

	f := new(big.Float).SetInt(coef)
	r, _ := f.Float64()

	return r
}

// Round a number
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
