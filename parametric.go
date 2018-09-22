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
	v := m.Vector{f.x(t), f.y(t), f.z(t)}
	return v.Normalize()
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

func pointsOnCircle(p, normal, binormal m.Vector, numPoints int, radius float64) []m.Vector {
	angle := (1 / (float64(numPoints))) * (2 * math.Pi)
	l := make([]m.Vector, numPoints)
	for i := 0; i < numPoints; i++ {
		xVector := normal.Times(radius * math.Cos(float64(i)*angle))
		yVector := binormal.Times(radius * math.Sin(float64(i)*angle))
		newP := p.Add(xVector).Add(yVector)
		l[i] = newP
	}
	return l
}

// assumes each list has the same number of points
func joinCirclePoints(pointLists [][]m.Vector, mat m.Material) []m.Object {
	numLists := len(pointLists)
	numPoints := len(pointLists[0])
	triangles := make([]m.Object, 2*numPoints*(numLists-1))

	for i := 0; i < numLists-1; i++ {
		offset := 2 * numPoints * i
		c1 := pointLists[i]
		c2 := pointLists[i+1]

		triangles[offset] = m.NewTriangle(c1[numPoints-1], c1[0], c2[numPoints-1], mat)
		triangles[offset+1] = m.NewTriangle(c2[numPoints-1], c1[0], c2[0], mat)
		for j := 0; j < numPoints-1; j++ {
			triangles[offset+(j+1)*2] = m.NewTriangle(c1[j], c1[j+1], c2[j], mat)
			triangles[offset+(j+1)*2+1] = m.NewTriangle(c2[j], c1[j+1], c2[j+1], mat)
		}
	}

	return triangles
}

// TODO: cleanup parameters
func parametricObject(f1, f2, f3 ParametricFunction, numPoints int, radius float64, numSteps int, stepSize float64, mat m.Material) m.Object {
	f := frenetFunctions(f1, f2, f3)
	points := make([][]m.Vector, numSteps)

	for i := 0; i < numSteps; i++ {
		t := float64(i) * stepSize
		p, _, normal, binormal := f(t)
		// reverse index because of triangle normals?
		points[numSteps-1-i] = pointsOnCircle(p, normal, binormal, numPoints, radius)
	}

	triangles := joinCirclePoints(points, mat)
	return m.NewComplexObject(triangles)
}
