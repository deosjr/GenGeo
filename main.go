package main

import (
	"fmt"
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
	"github.com/deosjr/GenGeo/gen"
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

	m.SetBackgroundColor(m.NewColor(20, 20, 20))

	translation := m.Translate(m.Vector{0, 2, 1.5})
	rotation := m.RotateY(-math.Pi / 16.0)
	transform := translation.Mul(rotation)

	diffMat := &m.DiffuseMaterial{m.NewColor(250, 0, 0)}

	// points := gen.QuadraticKochIsland(4)
	// points := gen.DragonCurve(10)
	// points := gen.HexagonalGosperCurve(3)
	// points := gen.PeanoCurve(2)
	points := gen.HilbertCurve3D(3)
	points = gen.CenterPointsOnOrigin(points)

	radial := gen.NewRadialCircle(func(t float64) float64 { return 0.02 }, 20)
	o := gen.BuildFromPoints(radial, points, diffMat)
	c := m.NewSharedObject(o, transform)
	scene.Add(c)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
