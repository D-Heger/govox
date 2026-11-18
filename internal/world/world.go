// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package world

import (
	"fmt"
	"log"

	"github.com/D-Heger/govox/internal/block"
	"github.com/D-Heger/govox/internal/chunk"
	"github.com/D-Heger/govox/internal/storage"
	"github.com/D-Heger/govox/internal/terrain"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MinWorldY         = -64
	MaxWorldY         = 64
	WorldHeightChunks = (MaxWorldY - MinWorldY) / chunk.ChunkDepth
	MinChunkY         = MinWorldY / chunk.ChunkDepth
	MaxChunkY         = MaxWorldY / chunk.ChunkDepth

	ViewRadius = 32
	loadBudget = 128
	meshBudget = 64
)

type World struct {
	Atlas   block.Atlas
	gen     *terrain.Generator
	saveDir string
	chunks  map[chunk.Pos]*chunk.Chunk
}

func Create() *World {
	return &World{
		Atlas:   block.LoadAtlas("assets/atlas.png"),
		gen:     terrain.New(MinWorldY),
		saveDir: fmt.Sprintf("saves/%d", terrain.Seed),
		chunks:  make(map[chunk.Pos]*chunk.Chunk),
	}
}

// Stream loads chunks around the player and unloads distant ones. Block data is
// kept one chunk beyond the meshed radius so interior meshes always have real
// neighbors. Work is budgeted per frame and done nearest-first.
func (w *World) Stream(center rl.Vector3) {
	c := chunk.PosAt(int(center.X), int(center.Y), int(center.Z))
	cx0, cz0 := c.X, c.Z

	for pos, ch := range w.chunks {
		d := chebyshev(pos.X-cx0, pos.Z-cz0)
		if d > ViewRadius+1 {
			w.release(pos, ch)
			delete(w.chunks, pos)
		} else if d > ViewRadius && ch.Meshed {
			ch.Unload()
		}
	}

	loads, meshes := 0, 0
	for d := 0; d <= ViewRadius+1; d++ {
		for _, off := range ringOffsets[d] {
			cx, cz := cx0+off[0], cz0+off[1]
			for cy := MinChunkY; cy < MaxChunkY; cy++ {
				pos := chunk.Pos{X: cx, Y: cy, Z: cz}
				ch, ok := w.chunks[pos]
				if !ok {
					if loads >= loadBudget {
						continue
					}
					ch = w.load(pos)
					w.chunks[pos] = ch
					loads++
				}
				if d <= ViewRadius && !ch.Meshed && meshes < meshBudget && w.neighborsLoaded(pos) {
					ch.BuildMesh(w.Atlas, w.IsSolid)
					meshes++
				}
			}
		}
	}
}

func (w *World) load(pos chunk.Pos) *chunk.Chunk {
	blocks, ok, err := storage.Load(w.saveDir, pos)
	if err != nil {
		log.Printf("load chunk %v: %v", pos, err)
	}
	if ok {
		return chunk.FromBlocks(pos, blocks)
	}
	return chunk.Create(pos, w.gen)
}

func (w *World) release(pos chunk.Pos, ch *chunk.Chunk) {
	if ch.Dirty {
		if err := storage.Save(w.saveDir, pos, &ch.BlockID); err != nil {
			log.Printf("save chunk %v: %v", pos, err)
		} else {
			ch.Dirty = false
		}
	}
	ch.Unload()
}

func (w *World) neighborsLoaded(pos chunk.Pos) bool {
	for _, o := range neighborOffsets {
		np := chunk.Pos{X: pos.X + o[0], Y: pos.Y + o[1], Z: pos.Z + o[2]}
		if np.Y < MinChunkY || np.Y >= MaxChunkY {
			continue
		}
		if _, ok := w.chunks[np]; !ok {
			return false
		}
	}
	return true
}

func (w *World) IsSolid(gx, gy, gz int) bool {
	if gy >= MaxWorldY {
		return false
	}
	if gy < MinWorldY {
		return true
	}

	pos := chunk.PosAt(gx, gy, gz)
	if ch, ok := w.chunks[pos]; ok {
		lx := gx - pos.X*chunk.ChunkSize
		ly := gy - pos.Y*chunk.ChunkDepth
		lz := gz - pos.Z*chunk.ChunkSize
		return ch.IsSolidLocal(lx, ly, lz)
	}

	return w.gen.BlockAt(gy, w.gen.SurfaceAt(gx, gz)) != block.Air
}

func (w *World) Draw() {
	for _, ch := range w.chunks {
		ch.Draw()
	}
}

func (w *World) Unload() {
	for pos, ch := range w.chunks {
		w.release(pos, ch)
	}
	w.Atlas.Unload()
}

var neighborOffsets = [6][3]int{{1, 0, 0}, {-1, 0, 0}, {0, 0, 1}, {0, 0, -1}, {0, 1, 0}, {0, -1, 0}}

// ringOffsets[d] lists the column offsets on the square perimeter of radius d.
var ringOffsets = func() [][][2]int {
	rings := make([][][2]int, ViewRadius+2)
	for d := range rings {
		if d == 0 {
			rings[d] = [][2]int{{0, 0}}
			continue
		}
		cols := make([][2]int, 0, 8*d)
		for i := -d; i < d; i++ {
			cols = append(cols,
				[2]int{i, -d},
				[2]int{d, i},
				[2]int{-i, d},
				[2]int{-d, -i},
			)
		}
		rings[d] = cols
	}
	return rings
}()

func chebyshev(dx, dz int) int {
	if dx < 0 {
		dx = -dx
	}
	if dz < 0 {
		dz = -dz
	}
	if dx > dz {
		return dx
	}
	return dz
}
