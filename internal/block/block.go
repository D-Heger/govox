// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package block

type Face byte

const (
	FaceNorth Face = iota // z-
	FaceSouth             // z+
	FaceEast              // x+
	FaceWest              // x-
	FaceUp                // y+
	FaceDown              // y-
)

// TileID represents the ID of a tile (16x16 pixel area, left-to-right, top-to-bottom) in the texture atlas.
type TileID int8

// Tile constants for the texture atlas. These IDs correspond to the order of tiles in the atlas image.
const (
	TileDirt      TileID = 0
	TileGrassSide TileID = 1
	TileGrassTop  TileID = 2
	TileStone     TileID = 3
	TileDiabase   TileID = 4
)

// BlockID represents the ID of a block type in the game world.
type BlockID int8

// Block constants representing different types of blocks in the game world.
const (
	Air     BlockID = -1
	Dirt    BlockID = 0
	Grass   BlockID = 1
	Stone   BlockID = 2
	Diabase BlockID = 3
)

// blockFaces maps each BlockID to its TileID per Face, in the order of the face enum
var blockFaces = map[BlockID][6]TileID{
	Dirt: {
		TileDirt, TileDirt, TileDirt, TileDirt, TileDirt, TileDirt,
	},
	Grass: {
		TileGrassSide, TileGrassSide, TileGrassSide, TileGrassSide, TileGrassTop, TileDirt,
	},
	Stone: {
		TileStone, TileStone, TileStone, TileStone, TileStone, TileStone,
	},
	Diabase: {
		TileDiabase, TileDiabase, TileDiabase, TileDiabase, TileDiabase, TileDiabase,
	},
}

// GetTileForFace returns the TileID for a given BlockID and face index (0-5).
func GetTileForFace(blockID BlockID, face int) TileID {
	tiles, ok := blockFaces[blockID]
	if !ok || face < 0 || face >= len(tiles) {
		return TileDirt
	}
	return tiles[face]
}
