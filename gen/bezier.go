package gen

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
)

// bezier curve, t in [0-1]
type bezierCurve struct {
	parametricFunction
	controlPoints []m.Vector
}

func NewBezierCurve(points []m.Vector) bezierCurve {
	return bezierCurve{
		parametricFunction: parametricFunction{
			f: func(t float64) m.Vector{ return bezierFunc(t, points)},
		},
		controlPoints: points,
	}
}

func bezierFunc(t float64, points []m.Vector) m.Vector {
	//TODO, dont feel like implementing binomials etc right now
	return m.Vector{}
}

type cubicBezierCurve struct {
	bezierCurve
}

func NewCubicBezierCurve(p0, p1, p2, p3 m.Vector) cubicBezierCurve {
	points := []m.Vector{p0, p1, p2, p3}
	return cubicBezierCurve{
		bezierCurve{
			parametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return cubicBezierFunc(t, p0, p1, p2, p3)},
			},
			controlPoints: points,
		},
	}
}

func cubicBezierFunc(t float64, p0, p1, p2, p3 m.Vector) m.Vector {
	k0 := float32(math.Pow((1-t), 3))
	k1 := float32(3*t*math.Pow((1-t), 2))
	k2 := float32(3*(1-t)*math.Pow(t, 2))
	k3 := float32(math.Pow(t, 3))
	return p0.Times(k0).Add(p1.Times(k1)).Add(p2.Times(k2)).Add(p3.Times(k3))
}

func (b cubicBezierCurve) Derivative() ParametricFunction {
	points := b.bezierCurve.controlPoints
	p0, p1, p2, p3 := points[0], points[1], points[2], points[3]
	newPoint := func(u, v m.Vector) m.Vector {
		return m.Vector{3 * (v.X - u.X), 3 * (v.Y - u.Y), 3 * (v.Z - u.Z)}
	}
	return NewQuadraticBezierCurve(newPoint(p0, p1), newPoint(p1, p2), newPoint(p2, p3))
}

func (b cubicBezierCurve) SecondDerivative() ParametricFunction {
	return b.Derivative().(quadraticBezierCurve).Derivative()
}

type quadraticBezierCurve struct {
	bezierCurve
}

func NewQuadraticBezierCurve(p0, p1, p2 m.Vector) quadraticBezierCurve {
	points := []m.Vector{p0, p1, p2}
	return quadraticBezierCurve{
		bezierCurve{
			parametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return quadraticBezierFunc(t, p0, p1, p2)},
			},
			controlPoints: points,
		},
	}
}

func quadraticBezierFunc(t float64, p0, p1, p2 m.Vector) m.Vector {
	k0 := float32(math.Pow((1-t), 2))
	k1 := float32(2*(1-t)*t)
	k2 := float32(math.Pow(t, 2))
	return p0.Times(k0).Add(p1.Times(k1)).Add(p2.Times(k2))
}

func (b quadraticBezierCurve) Derivative() ParametricFunction {
	points := b.bezierCurve.controlPoints
	p0, p1, p2 := points[0], points[1], points[2]
	newPoint := func(u, v m.Vector) m.Vector {
		return m.Vector{2 * (v.X - u.X), 2 * (v.Y - u.Y), 2 * (v.Z - u.Z)}
	}
	return NewLinearBezierCurve(newPoint(p0, p1), newPoint(p1, p2))
}

type linearBezierCurve struct {
	bezierCurve
}

func NewLinearBezierCurve(p0, p1 m.Vector) linearBezierCurve {
	points := []m.Vector{p0, p1}
	return linearBezierCurve{
		bezierCurve{
			parametricFunction: parametricFunction{
				f: func(t float64) m.Vector{ return linearBezierFunc(t, p0, p1)},
			},
			controlPoints: points,
		},
	}
}

func linearBezierFunc(t float64, p0, p1 m.Vector) m.Vector {
	return p0.Times(float32(1-t)).Add(p1.Times(float32(t)))
}

type ParametricSurface interface {
	Evaluate(u, v float64) m.Vector
	Triangulate(samples int, mat m.Material) m.Object
	TriangulateWithNormalMapping(samples int, mat m.Material, transform m.Transform) m.Object
}

type bicubicBezierPatch struct {
	controlPoints []m.Vector // len 16
}

func NewBicubicBezierPatch(points []m.Vector) bicubicBezierPatch {
	return bicubicBezierPatch{
		controlPoints: points,
	}
}

func (b bicubicBezierPatch) Evaluate(u, v float64) m.Vector {
	p := make([]m.Vector, 4)
	for i:=0; i < 4; i++ {
		p0 := b.controlPoints[i*4]
		p1 := b.controlPoints[i*4 + 1]
		p2 := b.controlPoints[i*4 + 2]
		p3 := b.controlPoints[i*4 + 3]
		p[i] = cubicBezierFunc(u, p0, p1, p2, p3)
	}
	return cubicBezierFunc(v, p[0], p[1], p[2], p[3])
}

// samples is number of triangles in one dimension of the surface,
// for samples*samples amount of triangles in total
func (b bicubicBezierPatch) Triangulate(samples int, mat m.Material) m.Object {
	triangles := []m.Triangle{}
	for u:=0; u<samples; u++ {
		for v:=0; v<samples; v++ {
			f64s := float64(samples)
			llhc := b.Evaluate(float64(u)/f64s, float64(v)/f64s)
			lrhc := b.Evaluate(float64(u+1)/f64s, float64(v)/f64s)
			ulhc := b.Evaluate(float64(u)/f64s, float64(v+1)/f64s)
			urhc := b.Evaluate(float64(u+1)/f64s, float64(v+1)/f64s)
			triangles = append(triangles, m.NewTriangle(llhc, lrhc, ulhc, mat))
			triangles = append(triangles, m.NewTriangle(lrhc, urhc, ulhc, mat))
		}
	}
	return m.NewTriangleComplexObject(triangles)
}

func dUBezier(controlPoints []m.Vector, u64, v float64) m.Vector {
    p := make([]m.Vector, 4)
    for i:=0;i<4;i++ {
        p0 := controlPoints[i]
        p1 := controlPoints[i+4]
        p2 := controlPoints[i+8]
        p3 := controlPoints[i+12]
        p[i] = cubicBezierFunc(v, p0, p1, p2, p3)
    }
    u := float32(u64)
    v1 := p[0].Times(-3 * (1-u) * (1-u))
    v2 := p[1].Times((3 * (1-u) * (1-u)) - (6 * u * (1-u)))
    v3 := p[2].Times((6 * u * (1-u)) - (3 * u * u))
    v4 := p[3].Times(3 * u * u)
    return v1.Add(v2).Add(v3).Add(v4)
}

func dVBezier(controlPoints []m.Vector, u, v64 float64) m.Vector {
    p := make([]m.Vector, 4)
    for i:=0;i<4;i++ {
        p0 := controlPoints[i*4]
        p1 := controlPoints[i*4+1]
        p2 := controlPoints[i*4+2]
        p3 := controlPoints[i*4+3]
        p[i] = cubicBezierFunc(u, p0, p1, p2, p3)
    }
    v := float32(v64)
    v1 := p[0].Times(-3 * (1-v) * (1-v))
    v2 := p[1].Times((3 * (1-v) * (1-v)) - (6 * v * (1-v)))
    v3 := p[2].Times((6 * v * (1-v)) - (3 * v * v))
    v4 := p[3].Times(3 * v * v)
    return v1.Add(v2).Add(v3).Add(v4)
}

func (b bicubicBezierPatch) normal(u, v float64) m.Vector {
    dU := dUBezier(b.controlPoints, u, v)
    dV := dVBezier(b.controlPoints, u, v)
    return dU.Cross(dV).Normalize()
}

// TODO: should be stored on the triangle meshes
var tempMap = map[m.Vector]m.Vector{}

// TODO: store vertex normals and use them instead of relying on closures 
// if nothing else, will make debugging easier
func (b bicubicBezierPatch) TriangulateWithNormalMapping(samples int, baseMat m.Material, transform m.Transform) m.Object {
	triangles := []m.Triangle{}
	for u:=0; u<samples; u++ {
		for v:=0; v<samples; v++ {
			f64s := float64(samples)
			llhc := b.Evaluate(float64(u)/f64s, float64(v)/f64s)
			nllhc := b.normal(float64(u)/f64s, float64(v)/f64s)
			lrhc := b.Evaluate(float64(u+1)/f64s, float64(v)/f64s)
			nlrhc := b.normal(float64(u+1)/f64s, float64(v)/f64s)
			ulhc := b.Evaluate(float64(u)/f64s, float64(v+1)/f64s)
			nulhc := b.normal(float64(u)/f64s, float64(v+1)/f64s)
			urhc := b.Evaluate(float64(u+1)/f64s, float64(v+1)/f64s)
			nurhc := b.normal(float64(u+1)/f64s, float64(v+1)/f64s)

            tempMap[llhc] = nllhc
            tempMap[lrhc] = nlrhc
            tempMap[ulhc] = nulhc
            tempMap[urhc] = nurhc

            mat := &m.NormalMappingMaterial{
                WrappedMaterial: baseMat,
                NormalFunc: func(si *m.SurfaceInteraction) m.Vector {
                    tr := si.GetObject().(m.Triangle)
			        p := transform.Inverse().Point(si.Point)
                    l0, l1, l2 := tr.Barycentric(p)
                    nl0 := tempMap[tr.P0]
                    nl1 := tempMap[tr.P1]
                    nl2 := tempMap[tr.P2]
                    return nl0.Times(l0).Add(nl1.Times(l1)).Add(nl2.Times(l2))
                },
            }
			triangles = append(triangles, m.NewTriangle(llhc, lrhc, ulhc, mat))
			triangles = append(triangles, m.NewTriangle(lrhc, urhc, ulhc, mat))
		}
	}
	return m.NewTriangleComplexObject(triangles)
}
