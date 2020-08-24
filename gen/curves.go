package gen

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

type ellipse struct {
	radiusx float32
	radiusy float32
	center  m.Vector
}

func NewEllipse(c m.Vector, a, b float32) ellipse {
	return ellipse{center: c, radiusx: a, radiusy: b}
}
func NewCircle(c m.Vector, r float32) ellipse {
	return ellipse{center: c, radiusx: r, radiusy: r}
}

func (e ellipse) Point(phase float64, xaxis, yaxis m.Vector) m.Vector {
	xVector := xaxis.Times(e.radiusx * float32(math.Cos(phase)))
	yVector := yaxis.Times(e.radiusy * float32(math.Sin(phase)))
	return e.center.Add(xVector).Add(yVector)
}

// PointsPhaseRange returns n points on an ellipse between phase min and max
func (e ellipse) PointsPhaseRange(minPhase, maxPhase float64, n int, xaxis, yaxis m.Vector) []m.Vector {
	diff := maxPhase - minPhase
	angle := (1 / float64(n-1)) * diff
	l := make([]m.Vector, n)
	for i := 0; i < n; i++ {
		phase := math.Mod(minPhase+float64(i)*angle, 2*math.Pi)
		l[i] = e.Point(phase, xaxis, yaxis)
	}
	return l
}

type sphere struct {
	center m.Vector
	radius float32
}

func NewSphere(c m.Vector, r float32) sphere {
	return sphere{center: c, radius: r}
}

// NOTE: sharing multiple of these spheres is not optimal since nmat cannot be shared nicely
// this function included mainly to document a use of normal mapping for now
func (s sphere) NormalMappedSphere(mat m.Material, n int) m.Object {
	nmat := &m.NormalMappingMaterial{
		WrappedMaterial: mat,
		NormalFunc: func(si *m.SurfaceInteraction) m.Vector {
            p := si.UntransformedPoint
			return m.VectorFromTo(s.center, p).Normalize()
		},
	}
	triangles := s.Triangulate(n, nmat)
	return m.NewTriangleComplexObject(triangles)
}

// triangulate starting from an octahedron approximating a sphere
// in n iterations of recursive subdivision
func (s sphere) Triangulate(n int, mat m.Material) []m.Triangle {
	triangles := AxisAlignedOctahedron(s.center, s.radius, mat)
	for i := 0; i < n; i++ {
		newT := []m.Triangle{}
		for _, t := range triangles {
			sub := subdivide(t)
			for _, st := range sub {
				newT = append(newT, m.NewTriangle(s.pushOut(st.P0), s.pushOut(st.P1), s.pushOut(st.P2), st.Material))
			}
		}
		triangles = newT
	}
	return triangles
}

func (s sphere) pushOut(p m.Vector) m.Vector {
	v := m.VectorFromTo(s.center, p)
	d := s.radius - v.Length()
	dv := v.Normalize().Times(d)
	return p.Add(dv)
}

// AxisAlignedOctahedron returns an octahedron centered on point p
// with vertices at distance r from p along the X, Y and Z axes
func AxisAlignedOctahedron(p m.Vector, r float32, mat m.Material) []m.Triangle {
	minX, maxX := m.Vector{p.X - r, p.Y, p.Z}, m.Vector{p.X + r, p.Y, p.Z}
	minY, maxY := m.Vector{p.X, p.Y - r, p.Z}, m.Vector{p.X, p.Y + r, p.Z}
	minZ, maxZ := m.Vector{p.X, p.Y, p.Z - r}, m.Vector{p.X, p.Y, p.Z + r}
	return []m.Triangle{
		m.NewTriangle(minX, maxZ, maxY, mat),
		m.NewTriangle(minZ, minX, maxY, mat),
		m.NewTriangle(maxX, minZ, maxY, mat),
		m.NewTriangle(maxZ, maxX, maxY, mat),
		m.NewTriangle(maxZ, minX, minY, mat),
		m.NewTriangle(minX, minZ, minY, mat),
		m.NewTriangle(minZ, maxX, minY, mat),
		m.NewTriangle(maxX, maxZ, minY, mat),
	}
}
