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
func (l Lsystem) Evaluate(n int, d float32, dFactor, delta float64) []Lsegment {
	instrs := l.rewriteN(l.Axiom, 0, n)
	dNew := d * float32(math.Pow(dFactor, float64(n)))
	return draw(instrs, dNew, delta)
}

func (l Lsystem) rewriteN(s string, depth, n int) []turtleInstruction {
	instrs := []turtleInstruction{}
	if n == 0 {
		for _, r := range s {
			instrs = append(instrs, lookup(r))
		}
		return instrs
	}
	for _, r := range s {
		p, ok := l.Productions[r]
		if !ok {
			instrs = append(instrs, lookup(r))
			continue
		}
		recInstrs := l.rewriteN(p, depth+1, n-1)
		instrs = append(instrs, recInstrs...)
	}
	return instrs
}

type savedPos struct {
	turtle  turtle
	H, L, U m.Vector
}

type turtleInstruction struct {
	operation turtleOperation
	// so far unused, but can support depth-related info
	// such as varying trunk width etc?
	param float32
}

type turtleOperation uint8

const (
	none turtleOperation = iota
	forward
	turnLeft
	turnRight
	pitchDown
	pitchUp
	rollLeft
	rollRight
	turnAround
	addStack
	popStack
	startLeaf
	endLeaf
)

func lookup(r rune) turtleInstruction {
	i := turtleInstruction{}
	switch r {
	case 'F', 'f', 'G', 'L', 'R':
		i.operation = forward
	case '+':
		i.operation = turnLeft
	case '-':
		i.operation = turnRight
	case '&':
		i.operation = pitchDown
	case '^':
		i.operation = pitchUp
	case '\\':
		i.operation = rollLeft
	case '/':
		i.operation = rollRight
	case '|':
		i.operation = turnAround
	case '[':
		i.operation = addStack
	case ']':
		i.operation = popStack
	case '{':
		i.operation = startLeaf
	case '}':
		i.operation = endLeaf
	default:
		i.operation = none
	}
	return i
}

type Lsegment interface {
	GetPoints() []m.Vector
}

type Lleaf struct {
	points []m.Vector
}

func (l Lleaf) GetPoints() []m.Vector {
	return l.points
}

type Lbranch struct {
	points []m.Vector
}

func (b Lbranch) GetPoints() []m.Vector {
	return b.points
}

func draw(instrs []turtleInstruction, d float32, delta float64) []Lsegment {
	// turtle starts in origin facing up
	origin := m.Vector{0, 0, 0}
	H, L, U := m.Vector{0, 1, 0}, m.Vector{1, 0, 0}, m.Vector{0, 0, 1}
	t := turtle{origin, H.Times(d)}
	stack := []savedPos{}

	segments := []Lsegment{}
	seg := []m.Vector{t.pos}
	var leafSeg []m.Vector
	var leafMaking bool
	for _, instr := range instrs {
		switch instr.operation {
		case forward:
			t.pos = t.pos.Add(t.heading)
			if leafMaking {
				leafSeg = append(leafSeg, t.pos)
			} else {
				seg = append(seg, t.pos)
			}
		case turnLeft:
			t.heading, L, H = transformAxes(delta, U, t.heading, L, H)
		case turnRight:
			t.heading, L, H = transformAxes(-delta, U, t.heading, L, H)
		case pitchDown:
			t.heading, U, H = transformAxes(delta, L, t.heading, U, H)
		case pitchUp:
			t.heading, U, H = transformAxes(-delta, L, t.heading, U, H)
		case rollLeft:
			t.heading, U, L = transformAxes(delta, H, t.heading, U, L)
		case rollRight:
			t.heading, U, L = transformAxes(-delta, H, t.heading, U, L)
		case turnAround:
			t.heading, L, H = transformAxes(math.Pi, U, t.heading, L, H)
		case addStack:
			lastPos := savedPos{t, H, L, U}
			stack = append(stack, lastPos)
		case popStack:
			newPos := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			t = newPos.turtle
			H, L, U = newPos.H, newPos.L, newPos.U
			if len(seg) > 1 {
				segments = append(segments, Lbranch{points: seg})
			}
			seg = []m.Vector{t.pos}
		case startLeaf:
			leafMaking = true
			leafSeg = []m.Vector{t.pos}
		case endLeaf:
			leafMaking = false
			segments = append(segments, Lleaf{points: leafSeg})
		}
	}
	if len(seg) > 1 {
		segments = append(segments, Lbranch{points: seg})
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
func QuadraticKochIsland(n int) []Lsegment {
	l := Lsystem{
		Axiom: "F-F-F-F",
		Productions: map[rune]string{
			'F': "F-F+F+FF-F-F+F",
		},
	}
	return l.Evaluate(n, 1.0, 0.25, math.Pi/2.0)
}

func DragonCurve(n int) []Lsegment {
	l := Lsystem{
		Axiom: "F",
		Productions: map[rune]string{
			'F': "F+G+",
			'G': "-F-G",
		},
	}
	return l.Evaluate(n, 1.0, 0.75, math.Pi/2.0)
}

func HexagonalGosperCurve(n int) []Lsegment {
	l := Lsystem{
		Axiom: "F",
		Productions: map[rune]string{
			'F': "F+G++G-F--FF-G+",
			'G': "-F+GG++G+F--F-G",
		},
	}
	return l.Evaluate(n, 1.0, 0.5, math.Pi/3.0)
}

func PeanoCurve(n int) []Lsegment {
	l := Lsystem{
		Axiom: "L",
		Productions: map[rune]string{
			'L': "LFRFL-F-RFLFR+F+LFRFL",
			'R': "RFLFR+F+LFRFL-F-RFLFR",
		},
	}
	return l.Evaluate(n, 1.0, 0.25, math.Pi/2.0)
}

// and some 3D examples:
func HilbertCurve3D(n int) []Lsegment {
	l := Lsystem{
		Axiom: "A",
		Productions: map[rune]string{
			'A': "B-F+CFC+F-D&F^D-F+&&CFC+F+B//",
			'B': "A&F^CFB^F^D^^-F-D^|F^B|FC^F^A//",
			'C': "|D^|F^B-F+C^F^A&&FA&F^C+F+B^F^D//",
			'D': "|CFB-F+B|FA&F^A&&FB-F+B|FC//",
		},
	}
	return l.Evaluate(n, 1.0, 0.5, math.Pi/2.0)
}

// branching 2D
func Branch2D_a(n int) []Lsegment {
	l := Lsystem{
		Axiom: "F",
		Productions: map[rune]string{
			'F': "F[+F]F[-F]F",
		},
	}
	return l.Evaluate(n, 5.0, 0.4, (25.7/360.0)*(2.0*math.Pi))
}

func Branch2D_b(n int) []Lsegment {
	l := Lsystem{
		Axiom: "F",
		Productions: map[rune]string{
			'F': "F[+F]F[-F][F]",
		},
	}
	return l.Evaluate(n, 5.0, 0.4, (20.0/360.0)*(2.0*math.Pi))
}

func Branch2D_d(n int) []Lsegment {
	l := Lsystem{
		Axiom: "X",
		Productions: map[rune]string{
			'X': "F[+X]F[-X]+X",
			'F': "FF",
		},
	}
	return l.Evaluate(n, 5.0, 0.4, (20.0/360.0)*(2.0*math.Pi))
}

// branching 3D: Fig 1.25
// not yet respecting the following operators:
// ! means decrement the diameter of segment
// ' means increment the index on color table (?)
// note: L is used for forward, so leaf is X
func Branch3D(n int) []Lsegment {
	l := Lsystem{
		Axiom: "A",
		Productions: map[rune]string{
			'A': "[&FX!A]/////'[&FX!A]/////'[&FX!A]",
			'F': "S/////F",
			'S': "FX",
			'X': "['''^^{-f+f+f-|-f+f+f}]",
		},
	}
	return l.Evaluate(n, 5.0, 0.5, (22.5/360.0)*(2.0*math.Pi))
}

// Fig 1.26
// will look better once diameter segment changes are implemented?
// no use of ! operator anywhere though...
func Branch3D_2(n int) []Lsegment {
	l := Lsystem{
		Axiom: "P",
		Productions: map[rune]string{
			// plant
			'P': "I+[P+O]--//[--X]I[++X]-[PO]++PO",
			// internode
			'I': "FS[//&&X][//^^X]FS",
			// seg
			'S': "SFS",
			// leaf
			'X': "['{+f-ff-f+|+f-ff-f}]",
			// flower
			'O': "[&&&C'/W////W////W////W////W]",
			// pedicel
			'C': "FF",
			// wedge
			'W': "['^F][{&&&&-f+f|-f+f}]",
		},
	}
	return l.Evaluate(n, 1.0, 0.5, (18.0/360.0)*(2.0*math.Pi))
}
