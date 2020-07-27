package gen

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

// bezier curve, t in [0-1]
type bezierCurve struct {
	ParametricFunction
	controlPoints []m.Vector
}

func NewBezierCurve(points []m.Vector) bezierCurve {
	return bezierCurve{
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

type cubicBezierCurve struct {
	bezierCurve
}

func NewCubicBezierCurve(p0, p1, p2, p3 m.Vector) cubicBezierCurve {
	points := []m.Vector{p0, p1, p2, p3}
	return cubicBezierCurve{
		bezierCurve{
			ParametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return cubicBezierFunc(t, p0, p1, p2, p3)},
			},
			controlPoints: points,
		},
	}
}

func cubicBezierFunc(t float64, p0, p1, p2, p3 m.Vector) m.Vector {
	k0 := float32(math.Pow((1-t), 3))
	k1 := float32(3*t*math.Pow((1-t), 2))
	k2 := float32(3*(1-t)*math.Pow(t, 2))
	k3 := float32(math.Pow(t, 3))
	return p0.Times(k0).Add(p1.Times(k1)).Add(p2.Times(k2)).Add(p3.Times(k3))
}

func (b cubicBezierCurve) Derivative() ParametricFunction {
	points := b.bezierCurve.controlPoints
	p0, p1, p2, p3 := points[0], points[1], points[2], points[3]
	newPoint := func(u, v m.Vector) m.Vector {
		return m.Vector{3 * (v.X - u.X), 3 * (v.Y - u.Y), 3 * (v.Z - u.Z)}
	}
	return NewQuadraticBezierCurve(newPoint(p0, p1), newPoint(p1, p2), newPoint(p2, p3))
}

func (b cubicBezierCurve) SecondDerivative() ParametricFunction {
	return b.Derivative().(quadraticBezierCurve).Derivative()
}

type quadraticBezierCurve struct {
	bezierCurve
}

func NewQuadraticBezierCurve(p0, p1, p2 m.Vector) quadraticBezierCurve {
	points := []m.Vector{p0, p1, p2}
	return quadraticBezierCurve{
		bezierCurve{
			ParametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return quadraticBezierFunc(t, p0, p1, p2)},
			},
			controlPoints: points,
		},
	}
}

func quadraticBezierFunc(t float64, p0, p1, p2 m.Vector) m.Vector {
	k0 := float32(math.Pow((1-t), 2))
	k1 := float32(2*(1-t)*t)
	k2 := float32(math.Pow(t, 2))
	return p0.Times(k0).Add(p1.Times(k1)).Add(p2.Times(k2))
}

func (b quadraticBezierCurve) Derivative() ParametricFunction {
	points := b.bezierCurve.controlPoints
	p0, p1, p2 := points[0], points[1], points[2]
	newPoint := func(u, v m.Vector) m.Vector {
		return m.Vector{2 * (v.X - u.X), 2 * (v.Y - u.Y), 2 * (v.Z - u.Z)}
	}
	return NewLinearBezierCurve(newPoint(p0, p1), newPoint(p1, p2))
}

type linearBezierCurve struct {
	bezierCurve
}

func NewLinearBezierCurve(p0, p1 m.Vector) linearBezierCurve {
	points := []m.Vector{p0, p1}
	return linearBezierCurve{
		bezierCurve{
			ParametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return linearBezierFunc(t, p0, p1)},
			},
			controlPoints: points,
		},
	}
}

func linearBezierFunc(t float64, p0, p1 m.Vector) m.Vector {
	return p0.Times(float32(1-t)).Add(p1.Times(float32(t)))
}


type bicubicBezierPatch struct {
	// ParametricSurface
	controlPoints []m.Vector // len 16
}

func NewBicubicBezierPatch(points []m.Vector) bicubicBezierPatch {
	return bicubicBezierPatch{
		controlPoints: points,
	}
}

func (b bicubicBezierPatch) Evaluate(u, v float64) m.Vector {
	p := make([]m.Vector, 4)
	for i:=0; i < 4; i++ {
		p0 := b.controlPoints[i*4]
		p1 := b.controlPoints[i*4 + 1]
		p2 := b.controlPoints[i*4 + 2]
		p3 := b.controlPoints[i*4 + 3]
		p[i] = cubicBezierFunc(u, p0, p1, p2, p3)
	}
	return cubicBezierFunc(v, p[0], p[1], p[2], p[3])
}