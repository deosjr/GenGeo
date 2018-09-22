package main

import (
	"testing"

	m "github.com/deosjr/GRayT/src/model"
)

func TestCubicBezier(t *testing.T) {
	for i, tt := range []struct {
		p0   m.Vector
		p1   m.Vector
		p2   m.Vector
		p3   m.Vector
		t    float64
		want m.Vector
	}{
		{
			p0:   m.Vector{1, 2, 3},
			p1:   m.Vector{4, 5, 6},
			p2:   m.Vector{7, 8, 9},
			p3:   m.Vector{10, 11, 12},
			t:    0.0,
			want: m.Vector{1, 2, 3},
		},
		{
			p0:   m.Vector{1, 2, 3},
			p1:   m.Vector{4, 5, 6},
			p2:   m.Vector{7, 8, 9},
			p3:   m.Vector{10, 11, 12},
			t:    1.0,
			want: m.Vector{10, 11, 12},
		},
	} {
		bezier := NewCubicBezier(tt.p0, tt.p1, tt.p2, tt.p3)
		got := bezier.Vector(tt.t)
		if !compareVector(got, tt.want) {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}
