package gen

import (
	m "github.com/deosjr/GRayT/src/model"
)

// extrusion of face with holes is a more involved problem
// lets start simple with a solid face
// assumption: points form a convex 2d object when projected on
// the plane perpendicular to extrusionVector
func ExtrudeSolidFace(points []m.Vector, extrusionVector m.Vector, mat m.Material) m.Object {
	triangles := []m.Object{}
	// dumb algorithm for front face: join triangles radiating from one point
	// note: facing is counterclockwise
	p0 := points[0]
	for i, p1 := range points[1 : len(points)-1] {
		p2 := points[i+2]
		t := m.NewTriangle(p0, p1, p2, mat)
		triangles = append(triangles, t)
	}
	// side faces
	ex := make([]m.Vector, len(points))
	for i, p := range points {
		ex[i] = p.Add(extrusionVector)
	}
	triangles = append(triangles, JoinPoints([][]m.Vector{points, ex}, mat)...)
	// and back face similar to front
	// note: reverse point order
	p0 = ex[0]
	for i, p1 := range ex[1 : len(ex)-1] {
		p2 := ex[i+2]
		t := m.NewTriangle(p2, p1, p0, mat)
		triangles = append(triangles, t)
	}
	return m.NewComplexObject(triangles)
}
