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
	return m.Vector{float32(f.x(t)), float32(f.y(t)), float32(f.z(t))}
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
	return m.NewTriangleComplexObject(triangles)
}

// a radial2d defines points around a center according to a pattern
// simplest example is a circle, drawing n points with radius r
// this captures all regular convex polygonals by adjusting n
type radial2d interface {
	Points(p, normal, binormal m.Vector, t float64) []m.Vector
}

type radialEllipse struct {
	radiusx   func(t float64) float32
	radiusy   func(t float64) float32
	numPoints int
}

func NewRadialEllipse(a, b func(t float64) float32, n int) radialEllipse {
	return radialEllipse{radiusx: a, radiusy: b, numPoints: n}
}
func NewRadialCircle(radius func(t float64) float32, n int) radialEllipse {
	return radialEllipse{radiusx: radius, radiusy: radius, numPoints: n}
}

func (re radialEllipse) Points(p, normal, binormal m.Vector, t float64) []m.Vector {
	e := NewEllipse(p, re.radiusx(t), re.radiusy(t))
	// PointsPhaseRange is inclusive on both ends, dont want phase 0 (=2pi) twice
	inclusive := e.PointsPhaseRange(0, 2*math.Pi, re.numPoints+1, normal, binormal)
	return inclusive[:len(inclusive)-1]
}

// assumption: radial does not depend on t
// ez taken as random vector to calculate normal/binormal
// TODO: random normal leads to crossed bindings
// new assumption: radial is always a circle with the same radius
func BuildFromPoints(radial radialEllipse, points []m.Vector, mat m.Material) m.Object {
	normal := m.Vector{0, 0, -1}
	binormal := m.Vector{0, 1, 0}
	radialPoints := make([][]m.Vector, len(points))
	radius := radial.radiusx(0)

	radialPoints[0] = radial.Points(points[0], normal, binormal, 0)
	for i := 1; i < len(points)-1; i++ {
		p := points[i]
		prev := points[i-1]
		next := points[i+1]
		prevnext := m.VectorFromTo(prev, next).Normalize()
		n_temp := m.Vector{0, 0, 1}.Cross(prevnext).Normalize()
		b_temp := prevnext.Cross(n_temp).Normalize()

		//temp normal/binormal are now semirandom vectors in the circle slice between points
		//we can translate the previous normal/binormal on that circle using a plane intersection test

		planenormal := n_temp.Cross(b_temp).Normalize()
		rdirection := m.VectorFromTo(prev, p).Normalize()
		rorigin := prev.Add(normal.Times(radius))
		ln := rdirection.Dot(planenormal)
		// TODO: figure out why ln can be 0, leads to NaN triangles
		// set to 1 for now but that is totally arbitrary
		if ln == 0 {
			ln = 1
		}
		d := m.VectorFromTo(rorigin, p).Dot(planenormal) / ln
		npoint := rorigin.Add(rdirection.Times(d))
		normal = m.VectorFromTo(p, npoint).Normalize()

		rorigin = prev.Add(binormal.Times(radius))
		d = m.VectorFromTo(rorigin, p).Dot(planenormal) / ln
		bpoint := rorigin.Add(rdirection.Times(d))
		binormal = m.VectorFromTo(p, bpoint).Normalize()

		radialPoints[i] = radial.Points(p, normal, binormal, 0)
	}
	radialPoints[len(points)-1] = radial.Points(points[len(points)-1], normal, binormal, 0)

	triangles := JoinPoints(radialPoints, mat)
	return m.NewTriangleComplexObject(triangles)
}

// build path of points: node object at each point, radial around vertices
// assumption: points are on a line
func BuildNodesVertices(node m.Object, radial radial2d, points []m.Vector, mat m.Material) m.Object {
	objects := []m.Object{}
	ez := m.Vector{0, 0, 1}

	for i := 0; i < len(points)-1; i++ {
		p := points[i]
		next := points[i+1]
		objects = append(objects, m.NewSharedObject(node, m.Translate(p)))
		heading := m.VectorFromTo(p, next).Normalize()
		normal := ez.Cross(heading).Normalize()
		binormal := normal.Cross(heading)
		radialP := radial.Points(p, normal, binormal, 0)
		radialNext := radial.Points(next, normal, binormal, 0)
		pointsList := [][]m.Vector{radialP, radialNext}
		for _, t := range JoinPoints(pointsList, mat) {
			objects = append(objects, t)
		}
	}
	translation := m.Translate(points[len(points)-1])
	objects = append(objects, m.NewSharedObject(node, translation))
	return m.NewComplexObject(objects)
}
