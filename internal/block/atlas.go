// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package block

import rl "github.com/gen2brain/raylib-go/raylib"

const tileSize = 16

type Atlas struct {
	Texture rl.Texture2D
}

func LoadAtlas(path string) Atlas {
	image := rl.LoadImage(path)
	defer rl.UnloadImage(image)

	texture := rl.LoadTextureFromImage(image)
	rl.SetTextureFilter(texture, rl.FilterPoint)

	return Atlas{Texture: texture}
}

// Returns [4]float32{uMin, vMin, uMax, vMax} for a given index tile
func (a Atlas) GetUV(blockID int8) [4]float32 {
	columns := int(a.Texture.Width) / tileSize
	if columns == 0 {
		return [4]float32{0, 0, 1, 1}
	}

	tileX := int(blockID) % columns
	tileY := int(blockID) / columns

	uMin := float32(tileX*tileSize) / float32(a.Texture.Width)
	vMin := float32(tileY*tileSize) / float32(a.Texture.Height)
	uMax := float32((tileX+1)*tileSize) / float32(a.Texture.Width)
	vMax := float32((tileY+1)*tileSize) / float32(a.Texture.Height)

	return [4]float32{uMin, vMin, uMax, vMax}
}

func (a Atlas) Unload() {
	rl.UnloadTexture(a.Texture)
}
