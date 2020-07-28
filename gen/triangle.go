package gen

import (
	m "github.com/deosjr/GRayT/src/model"
)

func subdivide(t m.Triangle) []m.Triangle {
	p01 := t.P0.Add(m.VectorFromTo(t.P0, t.P1).Times(0.5))
	p02 := t.P0.Add(m.VectorFromTo(t.P0, t.P2).Times(0.5))
	p12 := t.P1.Add(m.VectorFromTo(t.P1, t.P2).Times(0.5))
	return []m.Triangle{
		m.NewTriangle(t.P0, p01, p02, t.Material),
		m.NewTriangle(t.P1, p12, p01, t.Material),
		m.NewTriangle(t.P2, p02, p12, t.Material),
		m.NewTriangle(p01, p12, p02, t.Material),
	}
}

// assumes each list has the same number of points
func JoinPoints(pointLists [][]m.Vector, mat m.Material) []m.Triangle {
	numLists := len(pointLists)
	numPoints := len(pointLists[0])
	triangles := make([]m.Triangle, 2*numPoints*(numLists-1))

	for i := 0; i < numLists-1; i++ {
		offset := 2 * numPoints * i
		c1 := pointLists[i]
		c2 := pointLists[i+1]

		triangles[offset] = m.NewTriangle(c1[0], c1[numPoints-1], c2[numPoints-1], mat)
		triangles[offset+1] = m.NewTriangle(c1[0], c2[numPoints-1], c2[0], mat)
		for j := 0; j < numPoints-1; j++ {
			triangles[offset+(j+1)*2] = m.NewTriangle(c1[j], c1[j+1], c2[j], mat)
			triangles[offset+(j+1)*2+1] = m.NewTriangle(c2[j], c1[j+1], c2[j+1], mat)
		}
	}
	return triangles
}

// as JoinPoints, but the points are not considered to form a loop
// therefore the first and last points will not be linked up
// assumes each list has the same number of points
func JoinPointsNonCircular(pointLists [][]m.Vector, mat m.Material) []m.Triangle {
	numLists := len(pointLists)
	numPoints := len(pointLists[0])
	triangles := make([]m.Triangle, 2*(numPoints-1)*(numLists-1))

	for i := 0; i < numLists-1; i++ {
		offset := 2 * numPoints * i
		c1 := pointLists[i]
		c2 := pointLists[i+1]

		for j := 0; j < numPoints-1; j++ {
			triangles[offset+j*2] = m.NewTriangle(c1[j], c1[j+1], c2[j], mat)
			triangles[offset+j*2+1] = m.NewTriangle(c2[j], c1[j+1], c2[j+1], mat)
		}
	}
	return triangles
}

func CenterPointsOnOrigin(points []m.Vector) []m.Vector {
	centroid := m.Vector{0, 0, 0}
	for _, p := range points {
		centroid = centroid.Add(p)
	}
	centroid = centroid.Times(1.0 / float32(len(points)))
	objectToOrigin := m.Translate(centroid).Inverse()

	centered := make([]m.Vector, len(points))
	for i, p := range points {
		centered[i] = objectToOrigin.Point(p)
	}
	return centered
}
