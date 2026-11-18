// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package terrain

import (
	"math"
	"math/rand"
)

// Seeded 2D gradient noise with a fractal (multi-octave) sum.
type noise2D struct {
	perm [512]int
}

func newNoise2D(seed int64) noise2D {
	var n noise2D
	var p [256]int
	for i := range p {
		p[i] = i
	}
	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
	for i := range n.perm {
		n.perm[i] = p[i&255]
	}
	return n
}

func (n noise2D) at(x, y float64) float64 {
	fx, fy := math.Floor(x), math.Floor(y)
	X, Y := int(fx)&255, int(fy)&255
	xf, yf := x-fx, y-fy
	u, v := fade(xf), fade(yf)

	p := n.perm
	aa, ba := p[p[X]+Y], p[p[X+1]+Y]
	ab, bb := p[p[X]+Y+1], p[p[X+1]+Y+1]

	x1 := lerp(grad(aa, xf, yf), grad(ba, xf-1, yf), u)
	x2 := lerp(grad(ab, xf, yf-1), grad(bb, xf-1, yf-1), u)
	return lerp(x1, x2, v)
}

// fbm layers octaves of noise: each octave doubles frequency and halves amplitude.
func (n noise2D) fbm(x, y float64, octaves int, persistence, lacunarity float64) float64 {
	var total, amp, freq, norm float64 = 0, 1, 1, 0
	for i := 0; i < octaves; i++ {
		total += amp * n.at(x*freq, y*freq)
		norm += amp
		amp *= persistence
		freq *= lacunarity
	}
	return total / norm
}

func grad(hash int, x, y float64) float64 {
	switch hash & 7 {
	case 0:
		return x + y
	case 1:
		return -x + y
	case 2:
		return x - y
	case 3:
		return -x - y
	case 4:
		return x
	case 5:
		return -x
	case 6:
		return y
	default:
		return -y
	}
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}
