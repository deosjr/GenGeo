package main

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

type parametricFunction struct {
	x func(t float64) float64
	y func(t float64) float64
	z func(t float64) float64
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

func frenetFunctions(f, deriv, secDeriv ParametricFunction) func(t float64) (p, tangent, normal, binormal m.Vector) {
	frenet1 := func(t float64) m.Vector {
		return deriv.Vector(t).Normalize()
	}
	frenet2 := func(t float64) m.Vector {
		secondDeriv := secDeriv.Vector(t)
		e1 := frenet1(t)
		return secondDeriv.Sub(e1.Times(secondDeriv.Dot(e1)))
	}
	frenet3 := func(t float64) m.Vector {
		return frenet1(t).Cross(frenet2(t))
	}
	return func(t float64) (p, tangent, normal, binormal m.Vector) {
		p = f.Vector(t)
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
	function         ParametricFunction
	derivative       ParametricFunction
	secondDerivative ParametricFunction
	radial           radial2d
	numSteps         int
	stepSize         float64
	mat              m.Material
}

func (po parametricObject) build() m.Object {
	f := frenetFunctions(po.function, po.derivative, po.secondDerivative)
	points := make([][]m.Vector, po.numSteps)

	for i := 0; i < po.numSteps; i++ {
		t := float64(i) * po.stepSize
		p, _, normal, binormal := f(t)
		points[i] = po.radial.points(p, normal, binormal, t)
	}

	triangles := joinPoints(points, po.mat)
	return m.NewComplexObject(triangles)
}

// a radial2d defines points around a center according to a pattern
// simplest example is a circle, drawing n points with radius r
// this captures all regular convex polygonals by adjusting n
type radial2d interface {
	points(p, normal, binormal m.Vector, t float64) []m.Vector
}

type ellipse struct {
	radiusx   func(t float64) float64
	radiusy   func(t float64) float64
	numPoints int
}

func newEllipse(a, b func(t float64) float64, n int) ellipse {
	return ellipse{radiusx: a, radiusy: b, numPoints: n}
}
func newCircle(radius func(t float64) float64, n int) ellipse {
	return ellipse{radiusx: radius, radiusy: radius, numPoints: n}
}

func (e ellipse) points(p, normal, binormal m.Vector, t float64) []m.Vector {
	angle := (1 / (float64(e.numPoints))) * (2 * math.Pi)
	radiusx := e.radiusx(t)
	radiusy := e.radiusy(t)
	l := make([]m.Vector, e.numPoints)
	for i := 0; i < e.numPoints; i++ {
		xVector := normal.Times(radiusx * math.Cos(float64(i)*angle))
		yVector := binormal.Times(radiusy * math.Sin(float64(i)*angle))
		newP := p.Add(xVector).Add(yVector)
		l[i] = newP
	}
	return l
}

// assumes each list has the same number of points
func joinPoints(pointLists [][]m.Vector, mat m.Material) []m.Object {
	numLists := len(pointLists)
	numPoints := len(pointLists[0])
	triangles := make([]m.Object, 2*numPoints*(numLists-1))

	for i := 0; i < numLists-1; i++ {
		offset := 2 * numPoints * i
		c1 := pointLists[i]
		c2 := pointLists[i+1]

		triangles[offset] = m.NewTriangle(c1[numPoints-1], c2[numPoints-1], c1[0], mat)
		triangles[offset+1] = m.NewTriangle(c2[numPoints-1], c2[0], c1[0], mat)
		for j := 0; j < numPoints-1; j++ {
			triangles[offset+(j+1)*2] = m.NewTriangle(c1[j], c2[j], c1[j+1], mat)
			triangles[offset+(j+1)*2+1] = m.NewTriangle(c2[j], c2[j+1], c1[j+1], mat)
		}
	}
	return triangles
}
