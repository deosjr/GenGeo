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

	p0 := m.Vector{0.00, 0.00, 0.0}
	p1 := m.Vector{0.10, 1.5, 0.25}
	p2 := m.Vector{2.20, 2.6, 0.75}
	p3 := m.Vector{2.90, 1.23, 1.0}
	bezier := NewCubicBezier(p0, p1, p2, p3)
	bezierDeriv := bezier.derivative()
	bezier2nd := bezierDeriv.derivative()

	po := parametricObject{
		function:         bezier,
		derivative:       bezierDeriv,
		secondDerivative: bezier2nd,
		numPoints:        20,
		radius:           func(t float64) float64 { return (1 - t) * 0.01 },
		numSteps:         101,
		stepSize:         1.0 / 101.0,
		mat:              &m.DiffuseMaterial{m.NewColor(50, 150, 80)},
	}
	complexObject := po.build()

	translation := m.Translate(m.Vector{-1, 0, 2})
	// rotation := m.RotateY(math.Pi)
	boom := m.NewSharedObject(complexObject, translation) //.Mul(rotation))
	scene.Add(boom)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
