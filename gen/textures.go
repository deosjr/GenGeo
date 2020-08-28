package gen

import (
    "image"
    "image/color"
    "math/rand"
    "time"
)

type GrayScottInput struct {
    Width int
    Height int
    Iterations int
    FeedRate float64
    KillRate float64
    DiffRateA float64
    DiffRateB float64
}

type coord struct {
    X int
    Y int
}

type ab struct {
    A float64
    B float64
}

const (
    centerWeight = -1.0
    adjacentWeight = 0.2
    diagonalWeight = 0.05
    seedRate = 1.0 / 50.0
)

// grid based approximation to reaction-diffusion using Gray-Scott model
func GrayScott(input GrayScottInput) image.Image {
    w := input.Width
    h := input.Height
    f := input.FeedRate
    k := input.KillRate
    da := input.DiffRateA
    db := input.DiffRateB
    m := map[coord]ab{}

    // initial seeding of grid with A and B values
    rand.Seed(time.Now().Unix())
    for y:=0; y < h; y++ {
        for x:=0; x < w; x++ {
            b := 0.0
            if rand.Float64() < seedRate {
                b = 1.0
            }
            m[coord{X:x, Y:y}] = ab{A:1.0, B:b}
        }
    }

    // do n iterations of grid update
    for n:=0; n < input.Iterations; n++ {
        m = grayScottLoop(m, w, h, f, k, da, db)
    }

    // collect image and return
	img := image.NewRGBA(image.Rect(0, 0, w, h))
    for c, ab := range m {
        bw := uint8((1 - ab.B) * 255)
        color := color.RGBA{bw, bw, bw, 255}
        img.Set(c.X, c.Y, color)
    }
    return img
}


// A' = A + (D_a * La - AB^2 + f(1-A))dt
// B' = B + (D_b * Lb + AB^2 - (k+f)B)dt
// where D is diffusion rate and L the 2D laplacian,
// which we will approximate with a laplacian matrix
// dt is delta time which we will set to 1 (ignore)
// TODO: variable f, k based on location
func grayScottLoop(m map[coord]ab, w, h int, f, k, da, db float64) map[coord]ab {
    newM := map[coord]ab{}
    for c, cab := range m {
        x, y := c.X, c.Y
        a, b := cab.A, cab.B
        la := a * centerWeight
        lb := b * centerWeight
        adjacents, diagonals := neighbours(m, x, y, w, h)
        for _, n := range adjacents {
            la += n.A * adjacentWeight
            lb += n.B * adjacentWeight
        }
        for _, n := range diagonals {
            la += n.A * diagonalWeight
            lb += n.B * diagonalWeight
        }
        newa := a + ((da * la) - (a * b * b) + (f * (1.0-a)))
        newb := b + ((db * lb) + (a * b * b) - ((k + f) * b))
        newM[c] = ab{A:newa, B:newb}
    }
    return newM
}

func neighbours(m map[coord]ab, x, y, w, h int) ([]ab, []ab) {
    xn := x+1
    if xn == w {
        xn = 0
    }
    xp := x-1
    if x == 0 {
        xp = w-1
    }
    yn := y+1
    if yn == h {
        yn = 0
    }
    yp := y-1
    if y == 0 {
        yp = h-1
    }
    return []ab{
        m[coord{X:x, Y:yp}],
        m[coord{X:x, Y:yn}],
        m[coord{X:xp, Y:y}],
        m[coord{X:xn, Y:y}],
    }, []ab {
        m[coord{X:xp, Y:yp}],
        m[coord{X:xp, Y:yn}],
        m[coord{X:xn, Y:yp}],
        m[coord{X:xn, Y:yn}],
    }
}
