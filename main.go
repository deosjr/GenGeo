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

	m.SetBackgroundColor(m.NewColor(10, 30, 100))

	var decayFunc = func(init, decay float64) func(float64) float64 {
		return func(t float64) float64 {
			i := init
			if t/decay > i {
				return 0.0
			}
			return i - t/decay
		}
	}

	aFunc := decayFunc(0.5, 100.0)
	bFunc := decayFunc(0.07, 1000.0)
	helix := NewHelix(aFunc, bFunc)
	po := parametricObject{
		function:         helix.function(),
		derivative:       helix.derivative(),
		secondDerivative: helix.secondDerivative(),
		radial:           newEllipse(decayFunc(0.5, 100.0), decayFunc(0.4, 100.0), 50),
		numSteps:         100,
		stepSize:         math.Pi / 8.0,
		mat:              &m.DiffuseMaterial{m.NewColor(150, 100, 30)},
	}
	complexObject := po.build()

	translation := m.Translate(m.Vector{0, 1, 1.5})
	rotation := m.RotateY(-3 * math.Pi / 4.0).Mul(m.RotateX(-math.Pi / 2.0))
	helix1 := m.NewSharedObject(complexObject, translation.Mul(rotation))
	scene.Add(helix1)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
