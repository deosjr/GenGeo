package main

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

// bezier curve, t in [0-1]
type bezier struct {
	controlPoints []m.Vector
}

func bezierFunc(t float64, b bezier, d m.Dimension) float64 {
	//TODO, dont feel like implementing binomials etc right now
	return 0.0
}

func (b bezier) X(t float64) float64 {
	return bezierFunc(t, b, m.X)
}

func (b bezier) Y(t float64) float64 {
	return bezierFunc(t, b, m.Y)
}

func (b bezier) Z(t float64) float64 {
	return bezierFunc(t, b, m.Z)
}

func (b bezier) Vector(t float64) m.Vector {
	return m.Vector{b.X(t), b.Y(t), b.Z(t)}
}

type cubicBezier struct {
	bezier
}

func NewCubicBezier(p0, p1, p2, p3 m.Vector) cubicBezier {
	return cubicBezier{bezier{[]m.Vector{p0, p1, p2, p3}}}
}

func cubicBezierFunc(t float64, b bezier, d m.Dimension) float64 {
	p0, p1, p2, p3 := b.controlPoints[0], b.controlPoints[1], b.controlPoints[2], b.controlPoints[3]
	return p0.Get(d)*math.Pow((1-t), 3) + p1.Get(d)*3*math.Pow((1-t), 2)*t +
		p2.Get(d)*3*(1-t)*math.Pow(t, 2) + p3.Get(d)*math.Pow(t, 3)
}

func (b cubicBezier) X(t float64) float64 {
	return cubicBezierFunc(t, b.bezier, m.X)
}

func (b cubicBezier) Y(t float64) float64 {
	return cubicBezierFunc(t, b.bezier, m.Y)
}

func (b cubicBezier) Z(t float64) float64 {
	return cubicBezierFunc(t, b.bezier, m.Z)
}

func (b cubicBezier) Vector(t float64) m.Vector {
	return m.Vector{b.X(t), b.Y(t), b.Z(t)}
}

func (b cubicBezier) derivative() quadraticBezier {
	p0, p1, p2, p3 := b.controlPoints[0], b.controlPoints[1], b.controlPoints[2], b.controlPoints[3]
	newPoint := func(u, v m.Vector) m.Vector {
		return m.Vector{3 * (v.X - u.X), 3 * (v.Y - u.Y), 3 * (v.Z - u.Z)}
	}
	return NewQuadraticBezier(newPoint(p0, p1), newPoint(p1, p2), newPoint(p2, p3))
}

type quadraticBezier struct {
	bezier
}

func NewQuadraticBezier(p0, p1, p2 m.Vector) quadraticBezier {
	return quadraticBezier{bezier{[]m.Vector{p0, p1, p2}}}
}

func quadraticBezierFunc(t float64, b bezier, d m.Dimension) float64 {
	p0, p1, p2 := b.controlPoints[0], b.controlPoints[1], b.controlPoints[2]
	return p0.Get(d)*math.Pow((1-t), 2) + p1.Get(d)*2*(1-t)*t + p2.Get(d)*math.Pow(t, 2)
}

func (b quadraticBezier) X(t float64) float64 {
	return quadraticBezierFunc(t, b.bezier, m.X)
}

func (b quadraticBezier) Y(t float64) float64 {
	return quadraticBezierFunc(t, b.bezier, m.Y)
}

func (b quadraticBezier) Z(t float64) float64 {
	return quadraticBezierFunc(t, b.bezier, m.Z)
}

func (b quadraticBezier) Vector(t float64) m.Vector {
	return m.Vector{b.X(t), b.Y(t), b.Z(t)}
}

func (b quadraticBezier) derivative() linearBezier {
	p0, p1, p2 := b.controlPoints[0], b.controlPoints[1], b.controlPoints[2]
	newPoint := func(u, v m.Vector) m.Vector {
		return m.Vector{2 * (v.X - u.X), 2 * (v.Y - u.Y), 2 * (v.Z - u.Z)}
	}
	return NewLinearBezier(newPoint(p0, p1), newPoint(p1, p2))
}

type linearBezier struct {
	bezier
}

func NewLinearBezier(p0, p1 m.Vector) linearBezier {
	return linearBezier{bezier{[]m.Vector{p0, p1}}}
}

func linearBezierFunc(t float64, b bezier, d m.Dimension) float64 {
	p0, p1 := b.controlPoints[0], b.controlPoints[1]
	return p0.Get(d)*(1-t) + p1.Get(d)*t
}

func (b linearBezier) X(t float64) float64 {
	return linearBezierFunc(t, b.bezier, m.X)
}

func (b linearBezier) Y(t float64) float64 {
	return linearBezierFunc(t, b.bezier, m.Y)
}

func (b linearBezier) Z(t float64) float64 {
	return linearBezierFunc(t, b.bezier, m.Z)
}

func (b linearBezier) Vector(t float64) m.Vector {
	return m.Vector{b.X(t), b.Y(t), b.Z(t)}
}
