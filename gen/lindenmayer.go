package gen

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

// l-systems using turtle graphics interpretation of symbols
// only considering deterministic connected structures for now
// so no moving without drawing and one applicable rule per symbol
// working from http://algorithmicbotany.org/papers/abop/abop.pdf

type turtle struct {
	pos     m.Vector
	heading m.Vector
}

type Lsystem struct {
	Axiom       string
	Productions map[rune]string
}

// rewrite axiom n times according to productions
// then draw points from the result string
// d is length of initial line drawn by F at iteration 0
// dFactor is the factor by which d shrinks every iteration
// delta is the size of angle change by orientation changes
func (l Lsystem) Evaluate(n int, d, dFactor, delta float64) [][]m.Vector {
	s := l.Axiom
	for i := 0; i < n; i++ {
		s = l.rewrite(s)
	}
	dNew := d * math.Pow(dFactor, float64(n))
	return l.draw(s, dNew, delta)
}

func (l Lsystem) Nonbranching(n int, d, dFactor, delta float64) []m.Vector {
	a := l.Evaluate(n, d, dFactor, delta)
	return a[0]
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

type savedPos struct {
	turtle  turtle
	H, L, U m.Vector
}

func (l Lsystem) draw(s string, d, delta float64) [][]m.Vector {
	// turtle starts in origin facing up
	origin := m.Vector{0, 0, 0}
	H, L, U := m.Vector{0, 1, 0}, m.Vector{1, 0, 0}, m.Vector{0, 0, 1}
	t := turtle{origin, H.Times(d)}
	stack := []savedPos{}

	segments := [][]m.Vector{}
	segment := []m.Vector{t.pos}
	for _, r := range s {
		switch r {
		case 'F', 'G', 'L', 'R':
			t.pos = t.pos.Add(t.heading)
			segment = append(segment, t.pos)
		case '+':
			t.heading, L, H = transformAxes(delta, U, t.heading, L, H)
		case '-':
			t.heading, L, H = transformAxes(-delta, U, t.heading, L, H)
		case '&':
			t.heading, U, H = transformAxes(delta, L, t.heading, U, H)
		case '^':
			t.heading, U, H = transformAxes(-delta, L, t.heading, U, H)
		case '\\':
			t.heading, U, L = transformAxes(delta, H, t.heading, U, L)
		case '/':
			t.heading, U, L = transformAxes(-delta, H, t.heading, U, L)
		case '|':
			t.heading, L, H = transformAxes(math.Pi, U, t.heading, L, H)
		case '[':
			lastPos := savedPos{t, H, L, U}
			stack = append(stack, lastPos)
		case ']':
			newPos := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			t = newPos.turtle
			H, L, U = newPos.H, newPos.L, newPos.U
			segments = append(segments, segment)
			segment = []m.Vector{t.pos}
		}
	}
	if len(segment) > 1 {
		segments = append(segments, segment)
	}
	return segments
}

// rotate turtle heading and the other axes delta degrees around the principal axis
func transformAxes(delta float64, rotationAxis, th, n, bn m.Vector) (m.Vector, m.Vector, m.Vector) {
	transform := m.Rotate(delta, rotationAxis)
	newHeading := transform.Vector(th)
	normal := transform.Vector(n)
	binormal := transform.Vector(bn)
	return newHeading, normal, binormal
}

// some famous 2D L-system examples from the book:
func QuadraticKochIsland(n int) []m.Vector {
	l := Lsystem{
		Axiom: "F-F-F-F",
		Productions: map[rune]string{
			'F': "F-F+F+FF-F-F+F",
		},
	}
	return l.Nonbranching(n, 1.0, 0.25, math.Pi/2.0)
}

func DragonCurve(n int) []m.Vector {
	l := Lsystem{
		Axiom: "F",
		Productions: map[rune]string{
			'F': "F+G+",
			'G': "-F-G",
		},
	}
	return l.Nonbranching(n, 1.0, 0.75, math.Pi/2.0)
}

func HexagonalGosperCurve(n int) []m.Vector {
	l := Lsystem{
		Axiom: "F",
		Productions: map[rune]string{
			'F': "F+G++G-F--FF-G+",
			'G': "-F+GG++G+F--F-G",
		},
	}
	return l.Nonbranching(n, 1.0, 0.5, math.Pi/3.0)
}

func PeanoCurve(n int) []m.Vector {
	l := Lsystem{
		Axiom: "L",
		Productions: map[rune]string{
			'L': "LFRFL-F-RFLFR+F+LFRFL",
			'R': "RFLFR+F+LFRFL-F-RFLFR",
		},
	}
	return l.Nonbranching(n, 1.0, 0.25, math.Pi/2.0)
}

// and some 3D examples:
func HilbertCurve3D(n int) []m.Vector {
	l := Lsystem{
		Axiom: "A",
		Productions: map[rune]string{
			'A': "B-F+CFC+F-D&F^D-F+&&CFC+F+B//",
			'B': "A&F^CFB^F^D^^-F-D^|F^B|FC^F^A//",
			'C': "|D^|F^B-F+C^F^A&&FA&F^C+F+B^F^D//",
			'D': "|CFB-F+B|FA&F^A&&FB-F+B|FC//",
		},
	}
	return l.Nonbranching(n, 1.0, 0.5, math.Pi/2.0)
}

// branching 2D
func Branch2D(n int) [][]m.Vector {
	l := Lsystem{
		Axiom: "X",
		Productions: map[rune]string{
			'X': "F[+X]F[-X]+X",
			'F': "FF",
		},
	}
	return l.Evaluate(n, 5.0, 0.4, (20.0/360.0)*(2.0*math.Pi))
}
