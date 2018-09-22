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
	rotation := m.RotateY(math.Pi / 4)

	numPoints := 20
	radius := 0.5
	c0 := pointsOnCircle(m.Vector{0, 0, 0}, ex, ez, numPoints, radius)
	c1 := pointsOnCircle(m.Vector{0, 2, 0}, ex, ez, numPoints, radius*1.5)
	c2 := pointsOnCircle(m.Vector{0, 4, 0}, ex, ez, numPoints, radius)

	triangles := joinCirclePoints([][]m.Vector{c0, c1, c2})
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
