package main

import (
	"fmt"
	"math"

    //"image/png"
    //"os"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
    "github.com/deosjr/GenGeo/gen"
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

	//diffMat := m.NewDiffuseMaterial(m.ConstantTexture{Color: m.NewColor(250, 0, 0)})
    //texture := m.NewUVTexture(m.TriangleMeshUVFunc)
	//diffMat := m.NewDiffuseMaterial(texture)
    //texture := m.NewCheckerboardTexture(4, m.TriangleMeshUVFunc)
	//diffMat := m.NewDiffuseMaterial(texture)
/*    
    existingImageFile, err := os.Open("out.png")
	if err != nil {
        panic(err)
	}
	defer existingImageFile.Close()
    loadedImage, err := png.Decode(existingImageFile)
	if err != nil {
        panic(err)
	}
    texture := m.NewImageTexture(loadedImage, m.TriangleMeshUVFunc)
	diffMat := m.NewDiffuseMaterial(texture)
 */   

    input := gen.GrayScottInput{
        Width: 100,
        Height: 100,
        Iterations: 10000,
        FeedRate: 0.055,
        KillRate: 0.062,
        DiffRateA: 1.0,
        DiffRateB: 0.5,
    }
    img := gen.GrayScott(input)
    texture := m.NewImageTexture(img, m.TriangleMeshUVFunc)
	diffMat := m.NewDiffuseMaterial(texture)

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
		//diffMat := &m.DiffuseMaterial{Color: m.NewColor(uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)))}
		complexObject := patch.TriangulateWithNormalMapping(32, diffMat)
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

