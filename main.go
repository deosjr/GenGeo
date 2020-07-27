package main

import (
	"fmt"
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
	//"github.com/deosjr/GenGeo/gen"
)

var (
	width      uint = 1600
	height     uint = 1200
	numWorkers      = 10
	numSamples      = 10

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

	m.SetBackgroundColor(m.NewColor(200, 200, 200))

	diffMat := &m.DiffuseMaterial{Color: m.NewColor(250, 0, 0)}

	patches, err := LoadPatches("teapot")
	if err != nil {
		fmt.Println(err)
		return
	}
	translation := m.Translate(m.Vector{0, -2, 2})
	rotation := m.RotateX(-math.Pi/2.0).Mul(m.RotateX(math.Pi/8.0))
	// NOTE: enable scale to render the original teapot
	//scale := m.Scale(1, 4.0/3.0, 1)
	transformation := translation.Mul(rotation)//.Mul(scale)
	for _, patch := range patches {
		complexObject := patch.Triangulate(32, diffMat)
		// not really a mesh but I guess thats WIP
		patchMesh := m.NewSharedObject(complexObject, transformation)
		scene.Add(patchMesh)
	}

	scene.Precompute()

	from, to := m.Vector{0, 2, -3}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	params := render.Params{
		Scene:        scene,
		NumWorkers:   numWorkers,
		NumSamples:   numSamples,
		AntiAliasing: true,
		TracerType:   m.WhittedStyle,
	}
	fmt.Println("Rendering...")
	film := render.Render(params)
	film.SaveAsPNG("out.png")
}