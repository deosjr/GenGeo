package gen

import (
	m "github.com/deosjr/GRayT/src/model"
)

// extrusion of simple solid shape
// assumption: points form a convex 2d object when projected on
// the plane perpendicular to extrusionVector
// the back face is the front face mirrored
// the side is made by pairwise joining the points of front and back
func ExtrudeSolidFace(points []m.Vector, extrusionVector m.Vector, mat m.Material) m.Object {
	triangles := []m.Triangle{}
	// dumb algorithm for front face: join triangles radiating from one point
	// note: facing is counterclockwise
	p0 := points[0]
	for i, p1 := range points[1 : len(points)-1] {
		p2 := points[i+2]
		t := m.NewTriangle(p0, p1, p2, mat)
		triangles = append(triangles, t)
		exp0, exp1, exp2 := p0.Add(extrusionVector), p1.Add(extrusionVector), p2.Add(extrusionVector)
		backT := m.NewTriangle(exp2, exp1, exp0, mat)
		triangles = append(triangles, backT)
	}

	// side faces
	ex := make([]m.Vector, len(points))
	for i, p := range points {
		ex[i] = p.Add(extrusionVector)
	}
	triangles = append(triangles, JoinPoints([][]m.Vector{points, ex}, mat)...)
	return m.NewTriangleComplexObject(triangles)
}

// more complex / less naive:
// what if the front face was given as a list of triangles
// the back face is the front face translated and with mirrored facing
// we can indicate holes/outlines with separate lists of points
// for which order matters but JoinPoints will work just fine

type ExtrusionFace struct {
	Front    []m.Triangle
	Outer    [][]m.Vector
	Inner    [][]m.Vector
	Material m.Material
}

func (ef ExtrusionFace) Extrude(extrusionVector m.Vector) m.Object {
	triangles := []m.Triangle{}
	for _, t := range ef.Front {
		triangles = append(triangles, m.NewTriangle(t.P0, t.P1, t.P2, ef.Material))
		p0, p1, p2 := t.P0.Add(extrusionVector), t.P1.Add(extrusionVector), t.P2.Add(extrusionVector)
		backT := m.NewTriangle(p2, p1, p0, ef.Material)
		triangles = append(triangles, backT)
	}
	for _, list := range ef.Outer {
		ex := make([]m.Vector, len(list))
		for i, p := range list {
			ex[i] = p.Add(extrusionVector)
		}
		triangles = append(triangles, JoinPoints([][]m.Vector{list, ex}, ef.Material)...)
	}
	for _, list := range ef.Inner {
		ex := make([]m.Vector, len(list))
		for i, p := range list {
			ex[i] = p.Add(extrusionVector)
		}
		triangles = append(triangles, JoinPoints([][]m.Vector{list, ex}, ef.Material)...)
	}
	return m.NewTriangleComplexObject(triangles)
}

func (ef ExtrusionFace) ExtrudeNonCircular(extrusionVector m.Vector) m.Object {
	triangles := []m.Triangle{}
	for _, t := range ef.Front {
		triangles = append(triangles, m.NewTriangle(t.P0, t.P1, t.P2, ef.Material))
		p0, p1, p2 := t.P0.Add(extrusionVector), t.P1.Add(extrusionVector), t.P2.Add(extrusionVector)
		backT := m.NewTriangle(p2, p1, p0, ef.Material)
		triangles = append(triangles, backT)
	}
	for _, list := range ef.Outer {
		ex := make([]m.Vector, len(list))
		for i, p := range list {
			ex[i] = p.Add(extrusionVector)
		}
		triangles = append(triangles, JoinPointsNonCircular([][]m.Vector{list, ex}, ef.Material)...)
	}
	for _, list := range ef.Inner {
		ex := make([]m.Vector, len(list))
		for i, p := range list {
			ex[i] = p.Add(extrusionVector)
		}
		triangles = append(triangles, JoinPointsNonCircular([][]m.Vector{list, ex}, ef.Material)...)
	}
	return m.NewTriangleComplexObject(triangles)
}
