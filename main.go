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

	m.SetBackgroundColor(m.NewColor(200, 200, 200))

	translation := m.Translate(m.Vector{0, 0.5, 1.5})
	rotation := m.RotateY(-math.Pi / 16.0)
	transform := translation.Mul(rotation)

	diffMat := &m.DiffuseMaterial{m.NewColor(250, 0, 0)}

	// points := gen.QuadraticKochIsland(4)
	// points := gen.DragonCurve(10)
	// points := gen.HexagonalGosperCurve(3)
	// points := gen.PeanoCurve(1)
	// points := gen.HilbertCurve3D(3)
	segments := gen.Branch2D(7)
	objects := []m.Object{}
	for _, points := range segments {
		if len(points) == 1 {
			continue
		}

		radial := gen.NewRadialCircle(func(t float64) float64 { return 0.005 }, 10)
		o := gen.BuildFromPoints(radial, points, diffMat)
		objects = append(objects, o)
		shared := m.NewSharedObject(o, transform)
		scene.Add(shared)
	}

	//points := gen.CenterPointsOnOrigin(s)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
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
