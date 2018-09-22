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

	numPoints := 100
	radius := 0.5
	numSteps := 100
	stepSize := math.Pi / float64(2*(numSteps-1))
	diffMat := &m.DiffuseMaterial{m.NewColor(50, 150, 80)}
	complexObject := parametricObject(unitcircle, unitcircleDeriv, unitcircle2ndDeriv, numPoints, radius, numSteps, stepSize, diffMat)

	translation := m.Translate(m.Vector{1, 0, 2})
	rotation := m.RotateY(math.Pi)
	boom := m.NewSharedObject(complexObject, translation.Mul(rotation))
	scene.Add(boom)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
