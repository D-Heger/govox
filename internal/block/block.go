// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package block

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const tileSize = 16

type Block struct {
	Model   rl.Model
	Texture rl.Texture2D
}

func Create() Block {
	var model = rl.LoadModelFromMesh(rl.GenMeshCube(1.0, 1.0, 1.0))

	var atlas = rl.LoadImage("assets/atlas.png")
	defer rl.UnloadImage(atlas)

	var rect = rl.Rectangle{X: 0, Y: 0, Width: tileSize, Height: tileSize}
	var crop = rl.ImageFromImage(*atlas, rect)
	rl.ImageFlipVertical(&crop)
	defer rl.UnloadImage(&crop)

	var texture = rl.LoadTextureFromImage(&crop)
	rl.SetTextureFilter(texture, rl.FilterPoint)
	rl.SetMaterialTexture(&model.GetMaterials()[0], rl.MapDiffuse, texture)

	return Block{Model: model, Texture: texture}
}

func (b *Block) Unload() {
	rl.UnloadTexture(b.Texture)
	rl.UnloadModel(b.Model)
}
