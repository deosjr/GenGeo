package main

import (
	"math"
	"testing"

	m "github.com/deosjr/GRayT/src/model"
)

func TestPointsOnCircle(t *testing.T) {
	for i, tt := range []struct {
		point     m.Vector
		normal    m.Vector
		binormal  m.Vector
		numPoints int
		radius    float64
		want      []m.Vector
	}{
		{
			point:     m.Vector{0, 0, 0},
			normal:    m.Vector{1, 0, 0},
			binormal:  m.Vector{0, 0, 1},
			numPoints: 4,
			radius:    1.0,
			want:      []m.Vector{{1, 0, 0}, {0, 0, 1}, {-1, 0, 0}, {0, 0, -1}},
		},
		{
			point:     m.Vector{0, 0, 0},
			normal:    m.Vector{1, 0, 0},
			binormal:  m.Vector{0, 0, 1},
			numPoints: 8,
			radius:    1.0,
			want: []m.Vector{
				{1, 0, 0},
				{math.Cos(math.Pi / 4.0), 0, math.Sin(math.Pi / 4.0)},
				{0, 0, 1},
				{-math.Cos(math.Pi / 4.0), 0, math.Sin(math.Pi / 4.0)},
				{-1, 0, 0},
				{-math.Cos(math.Pi / 4.0), 0, -math.Sin(math.Pi / 4.0)},
				{0, 0, -1},
				{math.Cos(math.Pi / 4.0), 0, -math.Sin(math.Pi / 4.0)},
			},
		},
	} {
		got := pointsOnCircle(tt.point, tt.normal, tt.binormal, tt.numPoints, tt.radius)
		if !compareVectors(got, tt.want) {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}

func compareVector(u, v m.Vector) bool {
	imprecision := 1000.0
	for d := 0; d < 3; d++ {
		dim := m.Dimension(d)
		if int(u.Get(dim)*imprecision) != int(v.Get(dim)*imprecision) {
			return false
		}
	}
	return true
}

func compareVectors(vl1, vl2 []m.Vector) bool {
	if len(vl1) != len(vl2) {
		return false
	}
	for i := 0; i < len(vl1); i++ {
		if !compareVector(vl1[i], vl2[i]) {
			return false
		}
	}
	return true
}
