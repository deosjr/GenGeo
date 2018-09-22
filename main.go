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

	m.SetBackgroundColor(m.NewColor(10, 10, 50))

	aFunc := func(t float64) float64 { return t * 0.3 }
	bFunc := func(t float64) float64 { return t * 0.05 }
	helix := NewHelix(aFunc, bFunc)
	po := parametricObject{
		function:         helix.function(),
		derivative:       helix.derivative(),
		secondDerivative: helix.secondDerivative(),
		numPoints:        20,
		radius:           func(t float64) float64 { return 0.1 },
		numSteps:         300,
		stepSize:         math.Pi / 32.0,
		mat:              &m.DiffuseMaterial{m.NewColor(200, 100, 0)},
	}
	complexObject := po.build()

	translation := m.Translate(m.Vector{0, -0.5, 2})
	rotation := m.RotateX(-math.Pi / 2.0)
	helix1 := m.NewSharedObject(complexObject, translation.Mul(rotation))
	scene.Add(helix1)
	rotation = m.RotateY(math.Pi).Mul(m.RotateX(-math.Pi / 2.0))
	helix2 := m.NewSharedObject(complexObject, translation.Mul(rotation))
	scene.Add(helix2)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
