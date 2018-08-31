package php

import (
	"math"
	"math/cmplx"
	"math/rand"
	"strconv"
	"time"
)

func Abs(n float64) float64 {
	return math.Abs(n)
}

func Acos(val complex128) complex128 {
	return cmplx.Acos(val)
}

func Acosh(val complex128) complex128 {
	return cmplx.Acosh(val)
}

func Asin(val complex128) complex128 {
	return cmplx.Asin(val)
}

func Asinh(val complex128) complex128 {
	return cmplx.Asinh(val)
}

func Atan(val complex128) complex128 {
	return cmplx.Atan(val)
}

func Atanh(val complex128) complex128 {
	return cmplx.Atanh(val)
}

func Base_convert(num string, frombase, tobase int) (string, error) {
	i, err := strconv.ParseInt(num, frombase, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, tobase), nil
}

func Bindec(bin string) (int64, error) {
	i, err := strconv.ParseInt(bin, 2, 0)
	if err != nil {
		err = err.(*strconv.NumError).Err
	}
	return i, err
}

func Ceil(f float64) float64 {
	return math.Ceil(f)
}

// func Ceil(f float64) float64 {
// 	return float64(math.Ceil(float64(f)))
// }

func Max(x, y float64) float64 {
	return math.Max(float64(x), float64(y))
}

func Min(x, y float64) float64 {
	return math.Min(float64(x), float64(y))
}

var initRand bool

/* ============================================================================================ */
func Rand(m ...int) int {
	if initRand == false {
		rand.Seed(time.Now().UnixNano())
		initRand = true
	}

	argc := len(m)
	min := 0
	max := 32768
	if argc == 0 {
		return rand.Intn(max)
	} else if argc == 1 {
		if m[0] > 0 {
			max = m[0]
		}
		return rand.Intn(max)
	} else {
		min = m[0]
		max = m[1]

		if min < 0 {
			min = 0
		}
		if max < 0 {
			return rand.Intn(32768)
		}
		if min > max {
			return rand.Intn(32768)
		}
		return rand.Intn(max-min+1) + min
	}
}

func Mt_rand(m ...int) int {
	return Rand(m...)
}

func Pow(x, y float64) float64 {
	return math.Pow(x, y)
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func FloatPrecision(num float64, precision int) float64 {

	multiplier := math.Pow(10, float64(precision))
	return float64(Round(num*multiplier)) / multiplier
}
