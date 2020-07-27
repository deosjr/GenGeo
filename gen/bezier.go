package gen

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

// bezier curve, t in [0-1]
type bezier struct {
	ParametricFunction
	controlPoints []m.Vector
}

func NewBezier(points []m.Vector) bezier {
	return bezier{
		ParametricFunction: parametricFunction{
			f: func(t float64) m.Vector{ return bezierFunc(t, points)},
		},
		controlPoints: points,
	}
}

func bezierFunc(t float64, points []m.Vector) m.Vector {
	//TODO, dont feel like implementing binomials etc right now
	return m.Vector{} 
}

type cubicBezier struct {
	bezier
}

func NewCubicBezier(p0, p1, p2, p3 m.Vector) cubicBezier {
	points := []m.Vector{p0, p1, p2, p3}
	return cubicBezier{
		bezier{
			ParametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return cubicBezierFunc(t, points)},
			},
			controlPoints: points,
		},
	}
}

func cubicBezierFunc(t float64, points []m.Vector) m.Vector {
	p0, p1, p2, p3 := points[0], points[1], points[2], points[3]
	k0 := float32(math.Pow((1-t), 3))
	k1 := float32(3*t*math.Pow((1-t), 2))
	k2 := float32(3*(1-t)*math.Pow(t, 2))
	k3 := float32(math.Pow(t, 3))
	return p0.Times(k0).Add(p1.Times(k1)).Add(p2.Times(k2)).Add(p3.Times(k3))
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
			ParametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return quadraticBezierFunc(t, points)},
			},
			controlPoints: points,
		},
	}
}

func quadraticBezierFunc(t float64, points []m.Vector) m.Vector {
	p0, p1, p2 := points[0], points[1], points[2]
	k0 := float32(math.Pow((1-t), 2))
	k1 := float32(2*(1-t)*t)
	k2 := float32(math.Pow(t, 2))
	return p0.Times(k0).Add(p1.Times(k1)).Add(p2.Times(k2))
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
			ParametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return linearBezierFunc(t, points)},
			},
			controlPoints: points,
		},
	}
}

func linearBezierFunc(t float64, points []m.Vector) m.Vector {
	p0, p1 := points[0], points[1]
	return p0.Times(float32(1-t)).Add(p1.Times(float32(t)))
}
