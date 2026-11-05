// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package world

import (
	"github.com/D-Heger/govox/internal/block"
	"github.com/D-Heger/govox/internal/chunk"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const WorldSize = 16

type World struct {
	Atlas  block.Atlas
	Chunks [WorldSize][WorldSize]chunk.Chunk
}

func Create() World {
	world := World{
		Atlas: block.LoadAtlas("assets/atlas.png"),
	}

	// just create chunks on first pass
	for cx := range WorldSize {
		for cz := range WorldSize {
			position := rl.Vector3{
				X: float32(cx * chunk.ChunkSize),
				Y: 0,
				Z: float32(cz * chunk.ChunkSize),
			}
			world.Chunks[cx][cz] = chunk.Create(position)
		}
	}

	// second pass creates the meshes
	for cx := range WorldSize {
		for cz := range WorldSize {
			world.Chunks[cx][cz].BuildMesh(world.Atlas, world.IsSolid)
		}
	}

	return world
}

func (w *World) IsSolid(gx, gy, gz int) bool {
	if gx < 0 || gy < 0 || gy >= chunk.ChunkDepth || gz < 0 {
		return false
	}

	cx := gx / chunk.ChunkSize
	cz := gz / chunk.ChunkSize

	if cx >= WorldSize || cz >= WorldSize {
		return false
	}

	lx := gx % chunk.ChunkSize
	lz := gz % chunk.ChunkSize

	return w.Chunks[cx][cz].IsSolidLocal(lx, gy, lz)
}

func (world *World) Draw() {
	for cx := range WorldSize {
		for cz := range WorldSize {
			world.Chunks[cx][cz].Draw()
		}
	}
}

func (world *World) Unload() {
	for cx := range WorldSize {
		for cz := range WorldSize {
			world.Chunks[cx][cz].Unload()
		}
	}
	world.Atlas.Unload()
}
