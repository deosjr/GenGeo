package main

import (
	"fmt"
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
)

var (
	width      uint = 1600
	height     uint = 1200
	numWorkers      = 10

	ex = m.Vector{1, 0, 0}
	ey = m.Vector{0, 1, 0}
	ez = m.Vector{0, 0, 1}
)

func main() {
	fmt.Println("Creating scene...")
	camera := m.NewPerspectiveCamera(width, height, 0.5*math.Pi)
	scene := m.NewScene(camera)

	l1 := m.NewDistantLight(m.Vector{1, -1, 1}, m.NewColor(255, 255, 255), 50)
	scene.AddLights(l1)

	m.SetBackgroundColor(m.NewColor(50, 100, 150))

	// diffMat := &m.DiffuseMaterial{m.NewColor(50, 10, 100)}
	// reflMat := &m.ReflectiveMaterial{scene}
	translation := m.Translate(m.Vector{1, 0, 2})
	rotation := m.RotateY(math.Pi)

	numPoints := 20
	radius := 0.5

	unitcircle := parametricFunction{
		x: func(t float64) float64 { return math.Cos(t) },
		y: func(t float64) float64 { return math.Sin(t) },
		z: func(t float64) float64 { return 0.0 },
	}

	unitcircleDeriv := parametricFunction{
		x: func(t float64) float64 { return -math.Sin(t) },
		y: func(t float64) float64 { return math.Cos(t) },
		z: func(t float64) float64 { return 0.0 },
	}

	unitcircle2ndDeriv := parametricFunction{
		x: func(t float64) float64 { return -math.Cos(t) },
		y: func(t float64) float64 { return -math.Sin(t) },
		z: func(t float64) float64 { return 0.0 },
	}

	f := frenetFrame(unitcircle, unitcircleDeriv, unitcircle2ndDeriv)
	p0, _, n0, b0 := f(0)
	p1, _, n1, b1 := f(math.Pi / 8.0)
	p2, _, n2, b2 := f(math.Pi / 4.0)
	p3, _, n3, b3 := f(3 * math.Pi / 8.0)
	p4, _, n4, b4 := f(math.Pi / 2.0)

	c0 := pointsOnCircle(p0, n0, b0, numPoints, radius)
	c1 := pointsOnCircle(p1, n1, b1, numPoints, radius)
	c2 := pointsOnCircle(p2, n2, b2, numPoints, radius)
	c3 := pointsOnCircle(p3, n3, b3, numPoints, radius)
	c4 := pointsOnCircle(p4, n4, b4, numPoints, radius)

	triangles := joinCirclePoints([][]m.Vector{c4, c3, c2, c1, c0})
	complexObject := m.NewComplexObject(triangles)
	boom := m.NewSharedObject(complexObject, translation.Mul(rotation))
	scene.Add(boom)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
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
func joinCirclePoints(pointLists [][]m.Vector) []m.Object {
	diffMat := &m.DiffuseMaterial{m.NewColor(50, 200, 100)}

	numLists := len(pointLists)
	numPoints := len(pointLists[0])
	triangles := make([]m.Object, 2*numPoints*(numLists-1))

	for i := 0; i < numLists-1; i++ {
		offset := 2 * numPoints * i
		c1 := pointLists[i]
		c2 := pointLists[i+1]

		triangles[offset] = m.NewTriangle(c1[numPoints-1], c1[0], c2[numPoints-1], diffMat)
		triangles[offset+1] = m.NewTriangle(c2[numPoints-1], c1[0], c2[0], diffMat)
		for j := 0; j < numPoints-1; j++ {
			triangles[offset+(j+1)*2] = m.NewTriangle(c1[j], c1[j+1], c2[j], diffMat)
			triangles[offset+(j+1)*2+1] = m.NewTriangle(c2[j], c1[j+1], c2[j+1], diffMat)
		}
	}

	return triangles
}

type parametricFunction struct {
	x func(t float64) float64
	y func(t float64) float64
	z func(t float64) float64
}

func (f parametricFunction) Vector(t float64) m.Vector {
	v := m.Vector{f.x(t), f.y(t), f.z(t)}
	return v.Normalize()
}

func frenetFrame(f1, f2, f3 parametricFunction) func(t float64) (p, tangent, normal, binormal m.Vector) {
	frenet1 := f2.Vector
	frenet2 := func(t float64) m.Vector {
		secondDeriv := f3.Vector(t)
		e1 := frenet1(t)
		return secondDeriv.Sub(e1.Times(secondDeriv.Dot(e1)))
	}
	frenet3 := func(t float64) m.Vector {
		return frenet1(t).Cross(frenet2(t))
	}
	return func(t float64) (p, tangent, normal, binormal m.Vector) {
		p = f1.Vector(t)
		tangent = frenet1(t)
		normal = frenet2(t)
		binormal = frenet3(t)
		return p, tangent, normal, binormal
	}
}
