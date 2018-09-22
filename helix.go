package main

import "math"

type helix struct {
	a func(t float64) float64 //radius
	b func(t float64) float64 //slope b/a
}

func NewHelix(a, b func(t float64) float64) helix {
	return helix{a: a, b: b}
}

func (h helix) function() parametricFunction {
	return parametricFunction{
		x: func(t float64) float64 { return h.a(t) * math.Cos(t) },
		y: func(t float64) float64 { return h.a(t) * math.Sin(t) },
		z: func(t float64) float64 { return h.b(t) * t },
	}
}

func (h helix) derivative() parametricFunction {
	return parametricFunction{
		x: func(t float64) float64 { return -h.a(t) * math.Sin(t) },
		y: func(t float64) float64 { return h.a(t) * math.Cos(t) },
		z: func(t float64) float64 { return h.b(t) },
	}
}

func (h helix) secondDerivative() parametricFunction {
	return parametricFunction{
		x: func(t float64) float64 { return -h.a(t) * math.Cos(t) },
		y: func(t float64) float64 { return -h.a(t) * math.Sin(t) },
		z: func(t float64) float64 { return 0.0 },
	}
}
