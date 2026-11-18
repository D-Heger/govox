// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/D-Heger/govox/internal/block"
	"github.com/D-Heger/govox/internal/chunk"
)

const blockSize = chunk.ChunkSize * chunk.ChunkDepth * chunk.ChunkSize

func path(dir string, pos chunk.Pos) string {
	return filepath.Join(dir, fmt.Sprintf("%d_%d_%d.chunk", pos.X, pos.Y, pos.Z))
}

// Save writes a chunk's raw block data. Only dirty chunks need saving, pristine chunks regenerate from the terrain seed.
func Save(dir string, pos chunk.Pos, blocks *chunk.Blocks) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	buf := make([]byte, blockSize)
	i := 0
	for x := range chunk.ChunkSize {
		for y := range chunk.ChunkDepth {
			for z := range chunk.ChunkSize {
				buf[i] = byte(blocks[x][y][z])
				i++
			}
		}
	}
	return os.WriteFile(path(dir, pos), buf, 0o644)
}

// Load reads a saved chunk. The bool reports whether a save existed.
func Load(dir string, pos chunk.Pos) (chunk.Blocks, bool, error) {
	var blocks chunk.Blocks
	data, err := os.ReadFile(path(dir, pos))
	if os.IsNotExist(err) {
		return blocks, false, nil
	}
	if err != nil {
		return blocks, false, err
	}
	if len(data) != blockSize {
		return blocks, false, fmt.Errorf("chunk %v: invalid size %d", pos, len(data))
	}
	i := 0
	for x := range chunk.ChunkSize {
		for y := range chunk.ChunkDepth {
			for z := range chunk.ChunkSize {
				blocks[x][y][z] = block.BlockID(data[i])
				i++
			}
		}
	}
	return blocks, true, nil
}
