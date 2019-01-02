package gen

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

type ParametricFunction interface {
	X(t float64) float64
	Y(t float64) float64
	Z(t float64) float64
	Vector(t float64) m.Vector
}

// class C2 functions have a continuous first and second derivative
type c2Differentiable interface {
	Function() ParametricFunction
	Derivative() ParametricFunction
	SecondDerivative() ParametricFunction
}

// simple c2 class differential function with fixed derivative functions
type c2 struct {
	function         ParametricFunction
	derivative       ParametricFunction
	secondDerivative ParametricFunction
}

func NewC2Differentiable(f1, f2, f3 ParametricFunction) c2Differentiable {
	return c2{
		function:         f1,
		derivative:       f2,
		secondDerivative: f3,
	}
}

func (c c2) Function() ParametricFunction {
	return c.function
}

func (c c2) Derivative() ParametricFunction {
	return c.derivative
}

func (c c2) SecondDerivative() ParametricFunction {
	return c.secondDerivative
}

type parametricFunction struct {
	x func(t float64) float64
	y func(t float64) float64
	z func(t float64) float64
}

func NewParametricFunction(x, y, z func(t float64) float64) parametricFunction {
	return parametricFunction{
		x: x,
		y: y,
		z: z,
	}
}

func (f parametricFunction) X(t float64) float64 {
	return f.x(t)
}
func (f parametricFunction) Y(t float64) float64 {
	return f.y(t)
}
func (f parametricFunction) Z(t float64) float64 {
	return f.z(t)
}

func (f parametricFunction) Vector(t float64) m.Vector {
	return m.Vector{f.x(t), f.y(t), f.z(t)}
}

func (f parametricFunction) Function() ParametricFunction {
	return f
}

// return a function that evaluates both the point p at time t
// and the Frenet-Serret frame in this point
func frenetSerret(f c2Differentiable) func(t float64) (p, tangent, normal, binormal m.Vector) {
	frenet1 := func(t float64) m.Vector {
		return f.Derivative().Vector(t).Normalize()
	}
	frenet2 := func(t float64) m.Vector {
		secondDeriv := f.SecondDerivative().Vector(t)
		e1 := frenet1(t)
		return secondDeriv.Sub(e1.Times(secondDeriv.Dot(e1)))
	}
	frenet3 := func(t float64) m.Vector {
		return frenet1(t).Cross(frenet2(t))
	}
	return func(t float64) (p, tangent, normal, binormal m.Vector) {
		p = f.Function().Vector(t)
		tangent = frenet1(t)
		normal = frenet2(t)
		binormal = frenet3(t)
		return p, tangent, normal, binormal
	}
}

// a parametric object builds a 3d model by repeatedly drawing
// a 2d radial around a parametric function in the perpendicular plane
// and joining the resulting points in a mesh
type parametricObject struct {
	function c2Differentiable
	radial   radial2d
	numSteps int
	stepSize float64
	mat      m.Material
}

func NewParametricObject(f c2Differentiable, r radial2d, n int, s float64, m m.Material) parametricObject {
	return parametricObject{
		function: f,
		radial:   r,
		numSteps: n,
		stepSize: s,
		mat:      m,
	}
}

func (po parametricObject) Build() m.Object {
	f := frenetSerret(po.function)
	points := make([][]m.Vector, po.numSteps)

	for i := 0; i < po.numSteps; i++ {
		t := float64(i) * po.stepSize
		p, _, normal, binormal := f(t)
		points[i] = po.radial.Points(p, normal, binormal, t)
	}

	triangles := JoinPoints(points, po.mat)
	return m.NewComplexObject(triangles)
}

// a radial2d defines points around a center according to a pattern
// simplest example is a circle, drawing n points with radius r
// this captures all regular convex polygonals by adjusting n
type radial2d interface {
	Points(p, normal, binormal m.Vector, t float64) []m.Vector
}

type radialEllipse struct {
	radiusx   func(t float64) float64
	radiusy   func(t float64) float64
	numPoints int
}

func NewRadialEllipse(a, b func(t float64) float64, n int) radialEllipse {
	return radialEllipse{radiusx: a, radiusy: b, numPoints: n}
}
func NewRadialCircle(radius func(t float64) float64, n int) radialEllipse {
	return radialEllipse{radiusx: radius, radiusy: radius, numPoints: n}
}

func (re radialEllipse) Points(p, normal, binormal m.Vector, t float64) []m.Vector {
	e := NewEllipse(p, re.radiusx(t), re.radiusy(t))
	return e.PointsPhaseRange(0, 2*math.Pi, re.numPoints, normal, binormal)
}
