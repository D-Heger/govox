// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package chunk

import (
	"github.com/D-Heger/govox/internal/block"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ChunkSize  = 16
	ChunkDepth = 16
)

// Checks for a solid block at the world coordinate
type SolidFunc func(gx, gy, gz int) bool

type Chunk struct {
	BlockID  [ChunkSize][ChunkDepth][ChunkSize]int8
	Model    rl.Model
	Position rl.Vector3
}

type faceDev struct {
	neighborOffset [3]int //dx, dy, dz
	faceNormal     rl.Vector3
	faceBrightness uint8
	vertices       [6][3]float32
	uvIndex        [6][2]int // 0=uMin, 1=vMin, 2=uMax, 3=vMax
}

var cubeFaces = [6]faceDev{
	// North z-
	{
		neighborOffset: [3]int{0, 0, -1}, faceNormal: rl.Vector3{X: 0, Y: 0, Z: -1}, faceBrightness: 200,
		vertices: [6][3]float32{{0, 0, 0}, {1, 1, 0}, {1, 0, 0}, {0, 0, 0}, {0, 1, 0}, {1, 1, 0}},
		uvIndex:  [6][2]int{{2, 3}, {0, 1}, {0, 3}, {2, 3}, {2, 1}, {0, 1}},
	},
	// South z+
	{
		neighborOffset: [3]int{0, 0, 1}, faceNormal: rl.Vector3{X: 0, Y: 0, Z: 1}, faceBrightness: 220,
		vertices: [6][3]float32{{0, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 0, 1}, {1, 1, 1}, {0, 1, 1}},
		uvIndex:  [6][2]int{{0, 3}, {2, 3}, {2, 1}, {0, 3}, {2, 1}, {0, 1}},
	},
	// East x+
	{
		neighborOffset: [3]int{1, 0, 0}, faceNormal: rl.Vector3{X: 1, Y: 0, Z: 0}, faceBrightness: 175,
		vertices: [6][3]float32{{1, 0, 1}, {1, 0, 0}, {1, 1, 0}, {1, 0, 1}, {1, 1, 0}, {1, 1, 1}},
		uvIndex:  [6][2]int{{0, 3}, {2, 3}, {2, 1}, {0, 3}, {2, 1}, {0, 1}},
	},
	// West x-
	{
		neighborOffset: [3]int{-1, 0, 0}, faceNormal: rl.Vector3{X: -1, Y: 0, Z: 0}, faceBrightness: 160,
		vertices: [6][3]float32{{0, 0, 1}, {0, 1, 0}, {0, 0, 0}, {0, 0, 1}, {0, 1, 1}, {0, 1, 0}},
		uvIndex:  [6][2]int{{2, 3}, {0, 1}, {0, 3}, {2, 3}, {2, 1}, {0, 1}},
	},
	// Up y+
	{
		neighborOffset: [3]int{0, 1, 0}, faceNormal: rl.Vector3{X: 0, Y: 1, Z: 0}, faceBrightness: 255,
		vertices: [6][3]float32{{0, 1, 0}, {1, 1, 1}, {1, 1, 0}, {0, 1, 0}, {0, 1, 1}, {1, 1, 1}},
		uvIndex:  [6][2]int{{0, 1}, {2, 3}, {2, 1}, {0, 1}, {0, 3}, {2, 3}},
	},
	// Down y-
	{
		neighborOffset: [3]int{0, -1, 0}, faceNormal: rl.Vector3{X: 0, Y: -1, Z: 0}, faceBrightness: 130,
		vertices: [6][3]float32{{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 0}, {1, 0, 1}, {0, 0, 1}},
		uvIndex:  [6][2]int{{0, 1}, {2, 1}, {2, 3}, {0, 1}, {2, 3}, {0, 3}},
	},
}

func Create(position rl.Vector3) Chunk {
	chunk := Chunk{Position: position}
	chunk.generate()
	return chunk
}

func (chunk *Chunk) BuildMesh(atlas block.Atlas, worldIsSolid SolidFunc) {
	mesh := chunk.createMesh(atlas, worldIsSolid)
	chunk.Model = rl.LoadModelFromMesh(mesh)
	rl.SetMaterialTexture(&chunk.Model.GetMaterials()[0], rl.MapDiffuse, atlas.Texture)
}

func (chunk *Chunk) Draw() {
	rl.DrawModel(chunk.Model, chunk.Position, 1.0, rl.White)
}

func (chunk *Chunk) Unload() {
	rl.UnloadModel(chunk.Model)
}

func (chunk *Chunk) generate() {
	for y := range ChunkDepth {
		var b int8 = block.Stone
		if y > 10 {
			b = block.Air
		}
		if y == 10 {
			b = block.Grass
		}
		if y <= 9 && y >= 7 {
			b = block.Dirt
		}

		for x := range ChunkSize {
			for z := range ChunkSize {
				chunk.BlockID[x][y][z] = b
			}
		}
	}

	// Random Holes
	//for range 600 {
	//	chunk.BlockID[rand.Intn(ChunkSize)][rand.Intn(11)][rand.Intn(ChunkSize)] = block.Air
	//}

	for x := range ChunkSize {
		for y := range ChunkDepth {
			for z := range ChunkSize {
				if chunk.BlockID[x][y][z] == block.Dirt && !chunk.IsSolidLocal(x, y+1, z) {
					chunk.BlockID[x][y][z] = block.Grass
				}
			}
		}
	}
}

// chunk-cords lookup
func (c *Chunk) IsSolidLocal(x, y, z int) bool {
	if x < 0 || x >= ChunkSize || y < 0 || y >= ChunkDepth || z < 0 || z >= ChunkSize {
		return false
	}
	return c.BlockID[x][y][z] != block.Air
}

// Checks local first, then world
func (chunk *Chunk) isSolid(x, y, z int, worldIsSolid SolidFunc) bool {
	if y < 0 || y >= ChunkDepth {
		return false
	}
	if x < 0 || x >= ChunkSize || z < 0 || z >= ChunkSize {
		if worldIsSolid != nil {
			gx := int(chunk.Position.X) + x
			gz := int(chunk.Position.Z) + z
			return worldIsSolid(gx, y, gz)
		}
		return false
	}
	return chunk.BlockID[x][y][z] != block.Air
}

func (chunk *Chunk) createMesh(atlas block.Atlas, worldIsSolid SolidFunc) rl.Mesh {
	var vertices []float32
	var normals []float32
	var texcoords []float32
	var colors []uint8

	for x := range ChunkSize {
		for y := range ChunkDepth {
			for z := range ChunkSize {
				if !chunk.IsSolidLocal(x, y, z) {
					continue
				}

				id := chunk.BlockID[x][y][z]

				for faceIdx, face := range cubeFaces {
					// Only generate face if there is Air next to it
					if !chunk.isSolid(x+face.neighborOffset[0], y+face.neighborOffset[1], z+face.neighborOffset[2], worldIsSolid) {
						// Get the specific tile texture for this block's face
						tileID := block.GetTileForFace(id, faceIdx)
						uv := atlas.GetUV(tileID)

						for i := 0; i < 6; i++ {
							vertices = append(vertices,
								float32(x)+face.vertices[i][0],
								float32(y)+face.vertices[i][1],
								float32(z)+face.vertices[i][2],
							)
							normals = append(normals, face.faceNormal.X, face.faceNormal.Y, face.faceNormal.Z)
							texcoords = append(texcoords, uv[face.uvIndex[i][0]], uv[face.uvIndex[i][1]])
							colors = append(colors, face.faceBrightness, face.faceBrightness, face.faceBrightness, 255)
						}
					}
				}
			}
		}
	}

	var mesh rl.Mesh
	mesh.VertexCount = int32(len(vertices) / 3)
	mesh.TriangleCount = mesh.VertexCount / 3

	if mesh.VertexCount > 0 {
		mesh.Vertices = &vertices[0]
		mesh.Normals = &normals[0]
		mesh.Texcoords = &texcoords[0]
		mesh.Colors = &colors[0]
	}

	rl.UploadMesh(&mesh, false)

	mesh.Vertices = nil
	mesh.Normals = nil
	mesh.Texcoords = nil
	mesh.Colors = nil

	return mesh
}
