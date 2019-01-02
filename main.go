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

	translation := m.Translate(m.Vector{0.5, 1.5, 2})
	rotation := m.RotateY(math.Pi)
	transform := translation.Mul(rotation)

	diffMat := &m.DiffuseMaterial{m.NewColor(250, 200, 40)}
	nmat := &m.NormalMappingMaterial{
		WrappedMaterial: diffMat,
		NormalFunc: func(si *m.SurfaceInteraction) m.Vector {
			p := si.Point
			// Note: without reversing translation this calculation is incorrect
			p = transform.Inverse().Point(p)
			return m.VectorFromTo(m.Vector{0, 0, 0}, p).Normalize()
		},
	}

	s := gen.NewSphere(m.Vector{0, 0, 0}, 1.0)
	triangles := s.Triangulate(5, nmat)
	objs := make([]m.Object, len(triangles))
	for i, t := range triangles {
		objs[i] = t
	}
	complexObject := m.NewComplexObject(objs)
	c := m.NewSharedObject(complexObject, transform)
	scene.Add(c)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
