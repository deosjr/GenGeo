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
	rotation := m.RotateY(math.Pi / 4.0)
	transform := translation.Mul(rotation)

	// TODO: normal mapping as an option for parametric objects in general?
	// solving the function in t is hard though..

	// Normal mapping example:
	// Given a point m on the geometry, we can find the actual normal
	// of the full geometry (a ring torus in this case)
	// by finding the point p(t) on the parametric function closest to
	// p and drawing the normal vector from p(t) to m

	// This example takes the unit circle in x-y as a function
	// and a ring torus as the parametric object.
	// We find t by solving v.p'(t) = 0, the dot product of v and p'(t),
	// where v is the vector m to p(t) (or vice versa, doesn't matter)
	// and p'(t) is the tangent in p(t) given by its first derivative.
	// The closest point p(t) to m is where v and p'(t) are perpendicular, hence the 0

	// v.p'(t) = 0
	// p'x(t)*(mx - px(t)) + p'y(t)*(my - py(t)) = 0
	// -sin(t)*(mx - cos(t)) + cos(t)*(my - sin(t)) = 0
	// -mx*sin(t) + sin(t)cos(t) + my*cost(t) - sin(t)cos(t) = 0
	// my*cos(t) = mx*sin(t)
	// my/mx = sin(t)/cos(t) = tan(t)
	// t = arctan(my/mx)

	diffMat := &m.DiffuseMaterial{m.NewColor(250, 200, 40)}
	nmat := &m.NormalMappingMaterial{
		WrappedMaterial: diffMat,
		NormalFunc: func(si *m.SurfaceInteraction) m.Vector {
			p := si.Point
			// Note: without reversing translation this calculation is incorrect
			p = transform.Inverse().Point(p)
			// arctan range is -pi/2 , pi/2
			// therefore only half the circle is shaded correctly
			// using Atan2 with range -pi, pi instead
			t := math.Atan2(p.Y, p.X)
			pt := m.Vector{math.Cos(t), math.Sin(t), 0.0}
			return m.VectorFromTo(pt, p).Normalize()
		},
	}

	unitcircle := gen.NewParametricFunction(
		func(t float64) float64 { return math.Cos(t) },
		func(t float64) float64 { return math.Sin(t) },
		func(t float64) float64 { return 0.0 },
	)
	unitcircleDeriv := gen.NewParametricFunction(
		func(t float64) float64 { return -math.Sin(t) },
		func(t float64) float64 { return math.Cos(t) },
		func(t float64) float64 { return 0.0 },
	)
	unitcircle2ndDeriv := gen.NewParametricFunction(
		func(t float64) float64 { return -math.Cos(t) },
		func(t float64) float64 { return -math.Sin(t) },
		func(t float64) float64 { return 0.0 },
	)

	fUnitCircle := gen.NewC2Differentiable(unitcircle, unitcircleDeriv, unitcircle2ndDeriv)

	radial := gen.NewRadialCircle(func(t float64) float64 { return 0.1 }, 10)
	numSteps := 33
	stepSize := math.Pi / 16.0
	complexObject := gen.NewParametricObject(fUnitCircle, radial, numSteps, stepSize, nmat).Build()

	c := m.NewSharedObject(complexObject, transform)
	scene.Add(c)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
