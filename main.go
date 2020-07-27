package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

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

	diffMat := &m.DiffuseMaterial{Color: m.NewColor(250, 0, 0)}

	rand.Seed(time.Now().UTC().UnixNano())
	controlP := make([]m.Vector, 16)
	for i:=0; i < 4; i++ {
		controlP[i*4] = m.Vector{float32(i), rand.Float32(), 0}
		controlP[i*4+1] = m.Vector{float32(i), rand.Float32(), 1}
		controlP[i*4+2] = m.Vector{float32(i), rand.Float32(), 2}
		controlP[i*4+3] = m.Vector{float32(i), rand.Float32(), 3}
	}
	patch := gen.NewBicubicBezierPatch(controlP)
	triangles := []m.Triangle{}
	samples := 16
	for u:=0; u<samples-1; u++ {
		for v:=0; v<samples-1; v++ {
			llhc := patch.Evaluate(float64(u)/float64(samples), float64(v)/float64(samples))
			lrhc := patch.Evaluate(float64(u+1)/float64(samples), float64(v)/float64(samples))
			ulhc := patch.Evaluate(float64(u)/float64(samples), float64(v+1)/float64(samples))
			urhc := patch.Evaluate(float64(u+1)/float64(samples), float64(v+1)/float64(samples))
			triangles = append(triangles, m.NewTriangle(lrhc, llhc, ulhc, diffMat))
			triangles = append(triangles, m.NewTriangle(lrhc, ulhc, urhc, diffMat))
		}
	}

	translation := m.Translate(m.Vector{0, 1, 2})
	rotation := m.RotateY(math.Pi).Mul(m.RotateX(math.Pi/8.0))
	complexObject := m.NewTriangleComplexObject(triangles)
	boom := m.NewSharedObject(complexObject, translation.Mul(rotation))
	scene.Add(boom)

	//complex := m.NewComplexObject(objects)
	//fmt.Println(SaveObj(complex))

	//points := gen.CenterPointsOnOrigin(s)
	scene.Precompute()

	from, to := m.Vector{0, 2, -5}, m.Vector{0, 0, 10}
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

func SaveObj(o m.Object) string {
	triangles := trianglesFromObject(o)
	vertices := []m.Vector{}
	vertexMap := map[m.Vector]int64{}
	faces := make([]m.Face, len(triangles))
	for i, t := range triangles {
		v0, ok := vertexMap[t.P0]
		if !ok {
			v0 = int64(len(vertexMap)) + 1
			vertexMap[t.P0] = v0
			vertices = append(vertices, t.P0)
		}
		v1, ok := vertexMap[t.P1]
		if !ok {
			v1 = int64(len(vertexMap)) + 1
			vertexMap[t.P1] = v1
			vertices = append(vertices, t.P1)
		}
		v2, ok := vertexMap[t.P2]
		if !ok {
			v2 = int64(len(vertexMap)) + 1
			vertexMap[t.P2] = v2
			vertices = append(vertices, t.P2)
		}
		// TODO: coordinate handedness!
		faces[i] = m.Face{v2, v1, v0}
	}

	s := ""
	for _, v := range vertices {
		s += fmt.Sprintf("v %f %f %f\n", v.X, v.Y, v.Z)
	}
	for _, f := range faces {
		s += fmt.Sprintf("f %d %d %d\n", f.V0, f.V1, f.V2)
	}
	return s
}

func trianglesFromObject(objects ...m.Object) []m.Triangle {
	triangles := []m.Triangle{}
	for _, o := range objects {
		switch t := o.(type) {
		case m.Triangle:
			triangles = append(triangles, t)
		case *m.ComplexObject:
			triangles = append(triangles, trianglesFromObject(t.Objects()...)...)
		case *m.SharedObject:
			trs := trianglesFromObject(t.Object)
			for _, tr := range trs {
				p0 := t.ObjectToWorld.Point(tr.P0)
				p1 := t.ObjectToWorld.Point(tr.P1)
				p2 := t.ObjectToWorld.Point(tr.P2)
				newTr := m.NewTriangle(p0, p1, p2, tr.Material)
				triangles = append(triangles, newTr)
			}
		}
	}
	return triangles
}
