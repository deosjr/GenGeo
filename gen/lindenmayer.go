package gen

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

// starting with 2d l-systems using turtle graphics interpretation of symbols
// only considering deterministic connected structures for now
// so no moving without drawing and one applicable rule per symbol
// working from http://algorithmicbotany.org/papers/abop/abop.pdf

type turtle struct {
	x, y  float64
	angle float64
}

type Lsystem struct {
	Axiom       string
	Productions map[rune]string
}

// rewrite axiom n times according to productions
// then draw points from the result string
// d is length of initial line drawn by F at iteration 0
// dFactor is the factor by which d shrinks every iteration
// delta is the size of angle change by +/- operations
func (l Lsystem) Evaluate(n int, d, dFactor, delta float64) []m.Vector {
	s := l.Axiom
	for i := 0; i < n; i++ {
		s = l.rewrite(s)
	}
	dNew := d * math.Pow(dFactor, float64(n))
	return l.draw(s, dNew, delta)
}

func (l Lsystem) rewrite(s string) string {
	newS := ""
	for _, r := range s {
		p, ok := l.Productions[r]
		if !ok {
			newS = newS + string(r)
			continue
		}
		newS = newS + p
	}
	return newS
}

func (l Lsystem) draw(s string, d, delta float64) []m.Vector {
	// turtle starts in origin facing up
	t := turtle{0, 0, math.Pi / 2.0}

	points := []m.Vector{m.Vector{0, 0, 0}}
	for _, r := range s {
		switch r {
		case 'F', 'G':
			t.x = t.x + d*math.Cos(t.angle)
			t.y = t.y + d*math.Sin(t.angle)
			points = append(points, m.Vector{t.x, t.y, 0})
		case '+':
			t.angle = t.angle + delta
		case '-':
			t.angle = t.angle - delta
		}
	}
	return points
}

// some famous L-system examples from the book:
func QuadraticKochIsland(n int) []m.Vector {
	l := Lsystem{
		Axiom: "F-F-F-F",
		Productions: map[rune]string{
			'F': "F-F+F+FF-F-F+F",
		},
	}
	return l.Evaluate(n, 1.0, 0.25, math.Pi/2.0)
}

func DragonCurve(n int) []m.Vector {
	l := Lsystem{
		Axiom: "F",
		Productions: map[rune]string{
			'F': "F+G+",
			'G': "-F-G",
		},
	}
	return l.Evaluate(n, 1.0, 0.75, math.Pi/2.0)
}

func HexagonalGosperCurve(n int) []m.Vector {
	l := Lsystem{
		Axiom: "F",
		Productions: map[rune]string{
			'F': "F+G++G-F--FF-G+",
			'G': "-F+GG++G+F--F-G",
		},
	}
	return l.Evaluate(n, 1.0, 0.5, math.Pi/3.0)
}
