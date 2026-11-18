// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package terrain

import (
	"math"

	"github.com/D-Heger/govox/internal/block"
)

const (
	Seed        = 1337
	baseHeight  = 0
	amplitude   = 28
	scale       = 1.0 / 80
	octaves     = 4
	persistence = 0.5
	lacunarity  = 2.0
	dirtDepth   = 3
)

type Generator struct {
	noise   noise2D
	bottomY int
}

func New(bottomY int) *Generator {
	return &Generator{noise: newNoise2D(Seed), bottomY: bottomY}
}

// SurfaceAt returns the grass-layer height for a column.
func (g *Generator) SurfaceAt(gx, gz int) int {
	n := g.noise.fbm(float64(gx)*scale, float64(gz)*scale, octaves, persistence, lacunarity)
	return baseHeight + int(math.Round(n*amplitude))
}

// BlockAt returns the block at height gy in a column whose surface is surface.
func (g *Generator) BlockAt(gy, surface int) block.BlockID {
	switch {
	case gy <= g.bottomY:
		return block.Diabase
	case gy > surface:
		return block.Air
	case gy == surface:
		return block.Grass
	case gy >= surface-dirtDepth:
		return block.Dirt
	default:
		return block.Stone
	}
}
