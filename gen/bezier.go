package gen

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

// bezier curve, t in [0-1]
type bezier struct {
	parametricFunction
	controlPoints []m.Vector
}

func NewBezier(points []m.Vector) bezier {
	return bezier{
		parametricFunction: parametricFunction{
			x: parametrizeBezierFunc(bezierFunc, points, m.X),
			y: parametrizeBezierFunc(bezierFunc, points, m.Y),
			z: parametrizeBezierFunc(bezierFunc, points, m.Z),
		},
		controlPoints: points,
	}
}

func parametrizeBezierFunc(f func(float64, []m.Vector, m.Dimension) float64, points []m.Vector, d m.Dimension) func(float64) float64 {
	return func(t float64) float64 {
		return f(t, points, d)
	}
}

func bezierFunc(t float64, points []m.Vector, d m.Dimension) float64 {
	//TODO, dont feel like implementing binomials etc right now
	return 0.0
}

type cubicBezier struct {
	bezier
}

func NewCubicBezier(p0, p1, p2, p3 m.Vector) cubicBezier {
	points := []m.Vector{p0, p1, p2, p3}
	return cubicBezier{
		bezier{
			parametricFunction: parametricFunction{
				x: parametrizeBezierFunc(cubicBezierFunc, points, m.X),
				y: parametrizeBezierFunc(cubicBezierFunc, points, m.Y),
				z: parametrizeBezierFunc(cubicBezierFunc, points, m.Z),
			},
			controlPoints: points,
		},
	}
}

func cubicBezierFunc(t float64, points []m.Vector, d m.Dimension) float64 {
	p0, p1, p2, p3 := points[0], points[1], points[2], points[3]
	return p0.Get(d)*math.Pow((1-t), 3) + p1.Get(d)*3*math.Pow((1-t), 2)*t +
		p2.Get(d)*3*(1-t)*math.Pow(t, 2) + p3.Get(d)*math.Pow(t, 3)
}

func (b cubicBezier) Derivative() ParametricFunction {
	points := b.bezier.controlPoints
	p0, p1, p2, p3 := points[0], points[1], points[2], points[3]
	newPoint := func(u, v m.Vector) m.Vector {
		return m.Vector{3 * (v.X - u.X), 3 * (v.Y - u.Y), 3 * (v.Z - u.Z)}
	}
	return NewQuadraticBezier(newPoint(p0, p1), newPoint(p1, p2), newPoint(p2, p3))
}

func (b cubicBezier) SecondDerivative() ParametricFunction {
	return b.Derivative().(quadraticBezier).Derivative()
}

type quadraticBezier struct {
	bezier
}

func NewQuadraticBezier(p0, p1, p2 m.Vector) quadraticBezier {
	points := []m.Vector{p0, p1, p2}
	return quadraticBezier{
		bezier{
			parametricFunction: parametricFunction{
				x: parametrizeBezierFunc(quadraticBezierFunc, points, m.X),
				y: parametrizeBezierFunc(quadraticBezierFunc, points, m.Y),
				z: parametrizeBezierFunc(quadraticBezierFunc, points, m.Z),
			},
			controlPoints: points,
		},
	}
}

func quadraticBezierFunc(t float64, points []m.Vector, d m.Dimension) float64 {
	p0, p1, p2 := points[0], points[1], points[2]
	return p0.Get(d)*math.Pow((1-t), 2) + p1.Get(d)*2*(1-t)*t + p2.Get(d)*math.Pow(t, 2)
}

func (b quadraticBezier) Derivative() ParametricFunction {
	points := b.bezier.controlPoints
	p0, p1, p2 := points[0], points[1], points[2]
	newPoint := func(u, v m.Vector) m.Vector {
		return m.Vector{2 * (v.X - u.X), 2 * (v.Y - u.Y), 2 * (v.Z - u.Z)}
	}
	return NewLinearBezier(newPoint(p0, p1), newPoint(p1, p2))
}

type linearBezier struct {
	bezier
}

func NewLinearBezier(p0, p1 m.Vector) linearBezier {
	points := []m.Vector{p0, p1}
	return linearBezier{
		bezier{
			parametricFunction: parametricFunction{
				x: parametrizeBezierFunc(linearBezierFunc, points, m.X),
				y: parametrizeBezierFunc(linearBezierFunc, points, m.Y),
				z: parametrizeBezierFunc(linearBezierFunc, points, m.Z),
			},
			controlPoints: points,
		},
	}
}

func linearBezierFunc(t float64, points []m.Vector, d m.Dimension) float64 {
	p0, p1 := points[0], points[1]
	return p0.Get(d)*(1-t) + p1.Get(d)*t
}
