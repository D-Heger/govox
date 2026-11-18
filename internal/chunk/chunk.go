// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package chunk

import (
	"github.com/D-Heger/govox/internal/block"
	"github.com/D-Heger/govox/internal/terrain"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ChunkSize  = 16
	ChunkDepth = 16
)

// SolidFunc checks for a solid block at the world coordinate
type SolidFunc func(gx, gy, gz int) bool

// Pos is a chunk coordinate in chunk space.
type Pos struct {
	X, Y, Z int
}

// Blocks is a 3D array of BlockIDs representing the blocks in a chunk.
type Blocks [ChunkSize][ChunkDepth][ChunkSize]block.BlockID

// Chunk represents a chunk of blocks in the game world.
type Chunk struct {
	Pos        Pos
	BlockID    Blocks
	solidCount int
	Model      rl.Model
	Position   rl.Vector3
	Dirty      bool
	Meshed     bool
}

// PosAt returns the chunk containing the given block coordinate.
func PosAt(gx, gy, gz int) Pos {
	return Pos{X: floorDiv(gx, ChunkSize), Y: floorDiv(gy, ChunkDepth), Z: floorDiv(gz, ChunkSize)}
}

// floorDiv performs integer division that rounds towards negative infinity, useful for chunk coordinate calculations.
func floorDiv(a, b int) int {
	q := a / b
	if a%b != 0 && a < 0 {
		q--
	}
	return q
}

// origin returns the world-space origin of the chunk based on its chunk coordinates.
func (p Pos) origin() rl.Vector3 {
	return rl.Vector3{
		X: float32(p.X * ChunkSize),
		Y: float32(p.Y * ChunkDepth),
		Z: float32(p.Z * ChunkSize),
	}
}

// faceDev holds the data for a single face of a cube, including its neighbor offset, normal, brightness, vertices, and UV mapping.
type faceDev struct {
	neighborOffset [3]int //dx, dy, dz
	faceNormal     rl.Vector3
	faceBrightness uint8
	vertices       [6][3]float32
	uvIndex        [6][2]int // 0=uMin, 1=vMin, 2=uMax, 3=vMax
}

// cubeFaces defines the six faces of a cube, each with its neighbor offset, normal vector, brightness, vertex positions, and UV mapping indices.
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

// Create generates a new chunk at the given position using the provided terrain generator.
func Create(pos Pos, gen *terrain.Generator) *Chunk {
	chunk := &Chunk{Pos: pos, Position: pos.origin()}
	chunk.generate(gen)
	return chunk
}

// FromBlocks wraps already-loaded block data (e.g. from disk).
func FromBlocks(pos Pos, blocks Blocks) *Chunk {
	chunk := &Chunk{Pos: pos, Position: pos.origin(), BlockID: blocks}
	for x := range ChunkSize {
		for y := range ChunkDepth {
			for z := range ChunkSize {
				if blocks[x][y][z] != block.Air {
					chunk.solidCount++
				}
			}
		}
	}
	return chunk
}

// BuildMesh generates the mesh for the chunk based on its block data and the provided texture atlas. It uses the worldIsSolid function to determine if neighboring blocks are solid.
func (chunk *Chunk) BuildMesh(atlas block.Atlas, worldIsSolid SolidFunc) {
	chunk.Meshed = true
	if chunk.solidCount == 0 {
		return
	}

	mesh := chunk.createMesh(atlas, worldIsSolid)
	if mesh.VertexCount == 0 {
		return
	}

	chunk.Model = rl.LoadModelFromMesh(mesh)
	rl.SetMaterialTexture(&chunk.Model.GetMaterials()[0], rl.MapDiffuse, atlas.Texture)
}

// Draw renders the chunk's model at its world position.
func (chunk *Chunk) Draw() {
	if chunk.Model.MeshCount == 0 {
		return
	}
	rl.DrawModel(chunk.Model, chunk.Position, 1.0, rl.White)
}

// Unload releases the resources associated with the chunk's model.
func (chunk *Chunk) Unload() {
	chunk.Meshed = false
	if chunk.Model.MeshCount == 0 {
		return
	}
	rl.UnloadModel(chunk.Model)
	chunk.Model = rl.Model{}
}

// generate populates the chunk's block data using the provided terrain generator.
func (chunk *Chunk) generate(gen *terrain.Generator) {
	baseX := chunk.Pos.X * ChunkSize
	baseY := chunk.Pos.Y * ChunkDepth
	baseZ := chunk.Pos.Z * ChunkSize

	for x := range ChunkSize {
		for z := range ChunkSize {
			surface := gen.SurfaceAt(baseX+x, baseZ+z)
			for y := range ChunkDepth {
				b := gen.BlockAt(baseY+y, surface)
				chunk.BlockID[x][y][z] = b
				if b != block.Air {
					chunk.solidCount++
				}
			}
		}
	}
}

// IsSolidLocal checks if the block at the given chunk-space coordinates (x, y, z) is solid (not Air). It returns false if the coordinates are out of bounds.
func (c *Chunk) IsSolidLocal(x, y, z int) bool {
	if x < 0 || x >= ChunkSize || y < 0 || y >= ChunkDepth || z < 0 || z >= ChunkSize {
		return false
	}
	return c.BlockID[x][y][z] != block.Air
}

// isSolid checks if the block at the given chunk-space coordinates (x, y, z) is solid (not Air). If the coordinates are out of bounds, it uses the provided worldIsSolid function to check neighboring chunks.
func (chunk *Chunk) isSolid(x, y, z int, worldIsSolid SolidFunc) bool {
	if x >= 0 && x < ChunkSize &&
		y >= 0 && y < ChunkDepth &&
		z >= 0 && z < ChunkSize {
		return chunk.BlockID[x][y][z] != block.Air
	}

	if worldIsSolid != nil {
		gx := int(chunk.Position.X) + x
		gy := int(chunk.Position.Y) + y
		gz := int(chunk.Position.Z) + z
		return worldIsSolid(gx, gy, gz)
	}

	return false
}

// createMesh generates the mesh for the chunk based on its block data and the provided texture atlas. It uses the worldIsSolid function to determine if neighboring blocks are solid.
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

		rl.UploadMesh(&mesh, false)

		mesh.Vertices = nil
		mesh.Normals = nil
		mesh.Texcoords = nil
		mesh.Colors = nil
	}

	return mesh
}
