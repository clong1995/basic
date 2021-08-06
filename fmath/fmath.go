package fmath

import (
	"github.com/shopspring/decimal"
)

// Add 加
func Add(args ...float64) float64 {
	var newFloat64 decimal.Decimal
	for i, arg := range args {
		if i == 0 {
			newFloat64 = decimal.NewFromFloat(arg)
		} else {
			newFloat64 = newFloat64.Add(decimal.NewFromFloat(arg))
		}
	}
	result, _ := newFloat64.Float64()
	return result
}

// Sub 减
func Sub(args ...float64) float64 {
	var newFloat64 decimal.Decimal
	for i, arg := range args {
		if i == 0 {
			newFloat64 = decimal.NewFromFloat(arg)
		} else {
			newFloat64 = newFloat64.Sub(decimal.NewFromFloat(arg))
		}
	}
	result, _ := newFloat64.Float64()
	return result
}

//Mul 乘
func Mul(args ...float64) float64 {
	var newFloat64 decimal.Decimal
	for i, arg := range args {
		if i == 0 {
			newFloat64 = decimal.NewFromFloat(arg)
		} else {
			newFloat64 = newFloat64.Mul(decimal.NewFromFloat(arg))
		}
	}
	result, _ := newFloat64.Float64()
	return result
}

//Div 除
func Div(args ...float64) float64 {
	var newFloat64 decimal.Decimal
	for i, arg := range args {
		if i == 0 {
			newFloat64 = decimal.NewFromFloat(arg)
		} else {
			newFloat64 = newFloat64.Div(decimal.NewFromFloat(arg))
		}
	}
	result, _ := newFloat64.Float64()
	return result
}

//Round 四舍五入
func Round(number float64, places int32) float64 {
	result, _ := decimal.NewFromFloat(number).Round(places).Float64()

	return result
}
