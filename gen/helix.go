package gen

import "math"

type helix struct {
	a func(t float64) float64 //radius
	b func(t float64) float64 //slope b/a, pitch 2pi*b
}

func NewHelix(a, b func(t float64) float64) helix {
	return helix{a: a, b: b}
}

func (h helix) Function() ParametricFunction {
	return parametricFunctionPiecewise{
		x: func(t float64) float64 { return h.a(t) * math.Cos(t) },
		y: func(t float64) float64 { return h.a(t) * math.Sin(t) },
		z: func(t float64) float64 { return h.b(t) * t },
	}
}

func (h helix) Derivative() ParametricFunction {
	return parametricFunctionPiecewise{
		x: func(t float64) float64 { return -h.a(t) * math.Sin(t) },
		y: func(t float64) float64 { return h.a(t) * math.Cos(t) },
		z: func(t float64) float64 { return h.b(t) },
	}
}

func (h helix) SecondDerivative() ParametricFunction {
	return parametricFunctionPiecewise{
		x: func(t float64) float64 { return -h.a(t) * math.Cos(t) },
		y: func(t float64) float64 { return -h.a(t) * math.Sin(t) },
		z: func(t float64) float64 { return 0.0 },
	}
}
