package model

import "math"

func (m ScaledNumberType) GetValue() float64 {
	if m.Number == nil {
		return 0
	}
	var scale float64 = 0
	if m.Scale != nil {
		scale = float64(*m.Scale)
	}
	return float64(*m.Number) * math.Pow(10, scale)
}

func NewScaledNumberType(value float64) *ScaledNumberType {
	m := &ScaledNumberType{}

	var maxDecimals float64 = 3

	numberValue := NumberType(value * math.Pow(10, maxDecimals))
	m.Number = &numberValue

	if numberValue != 0 {
		scaleValue := ScaleType(-maxDecimals)
		m.Scale = &scaleValue
	}

	return m
}
