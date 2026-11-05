// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package block

const (
	FaceNorth = iota // z-
	FaceSouth        // z+
	FaceEast         // x+
	FaceWest         // x-
	FaceUp           // y+
	FaceDown         // y-
)

// TileID
const (
	TileDirt      int8 = 0
	TileGrassSide int8 = 1
	TileGrassTop  int8 = 2
	TileStone     int8 = 3
)

// BlockID
const (
	Air   int8 = -1
	Dirt  int8 = 0
	Grass int8 = 1
	Stone int8 = 2
)

// Maps each BlockID to its TileID per Face, in the order of the face enum
var blockFaces = map[int8][6]int8{
	Dirt: {
		TileDirt, TileDirt, TileDirt, TileDirt, TileDirt, TileDirt,
	},
	Grass: {
		TileGrassSide, TileGrassSide, TileGrassSide, TileGrassSide, TileGrassTop, TileDirt,
	},
	Stone: {
		TileStone, TileStone, TileStone, TileStone, TileStone, TileStone,
	},
}

func GetTileForFace(blockID int8, face int) int8 {
	if tiles, ok := blockFaces[blockID]; ok {
		return tiles[face]
	}
	return TileDirt
}
