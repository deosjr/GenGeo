package gen

import (
	"math"
	"reflect"
	"testing"

	m "github.com/deosjr/GRayT/src/model"
)

func TestPointsOnCircle(t *testing.T) {
	for i, tt := range []struct {
		point     m.Vector
		normal    m.Vector
		binormal  m.Vector
		numPoints int
		radius    func(t float64) float64
		t         float64
		want      []m.Vector
	}{
		{
			point:     m.Vector{0, 0, 0},
			normal:    m.Vector{1, 0, 0},
			binormal:  m.Vector{0, 0, 1},
			numPoints: 4,
			radius:    func(t float64) float64 { return 1.0 },
			want:      []m.Vector{{1, 0, 0}, {0, 0, 1}, {-1, 0, 0}, {0, 0, -1}},
		},
		{
			point:     m.Vector{0, 0, 0},
			normal:    m.Vector{1, 0, 0},
			binormal:  m.Vector{0, 0, 1},
			numPoints: 8,
			radius:    func(t float64) float64 { return 1.0 },
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
		c := NewRadialCircle(tt.radius, tt.numPoints)
		got := c.Points(tt.point, tt.normal, tt.binormal, tt.t)
		if !compareVectors(got, tt.want) {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}

func TestJoinCirclePoints(t *testing.T) {
	for i, tt := range []struct {
		points [][]m.Vector
		want   []m.Object
	}{
		{
			points: [][]m.Vector{{}},
			want:   []m.Object{},
		},
		{
			points: [][]m.Vector{
				{{1, 0, 0}, {0, 0, 1}, {-1, 0, 0}, {0, 0, -1}},
				{{1, 1, 0}, {0, 1, 1}, {-1, 1, 0}, {0, 1, -1}},
			},
			want: []m.Object{
				m.Triangle{P0: v(0, 0, -1), P1: v(0, 1, -1), P2: v(1, 0, 0)},
				m.Triangle{P0: v(0, 1, -1), P1: v(1, 1, 0), P2: v(1, 0, 0)},
				m.Triangle{P0: v(1, 0, 0), P1: v(1, 1, 0), P2: v(0, 0, 1)},
				m.Triangle{P0: v(1, 1, 0), P1: v(0, 1, 1), P2: v(0, 0, 1)},
				m.Triangle{P0: v(0, 0, 1), P1: v(0, 1, 1), P2: v(-1, 0, 0)},
				m.Triangle{P0: v(0, 1, 1), P1: v(-1, 1, 0), P2: v(-1, 0, 0)},
				m.Triangle{P0: v(-1, 0, 0), P1: v(-1, 1, 0), P2: v(0, 0, -1)},
				m.Triangle{P0: v(-1, 1, 0), P1: v(0, 1, -1), P2: v(0, 0, -1)},
			},
		},
	} {
		got := JoinPoints(tt.points, nil)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}

func v(x, y, z float64) m.Vector {
	return m.Vector{x, y, z}
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
