// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package game

import (
	"github.com/D-Heger/govox/internal/block"
	"github.com/D-Heger/govox/internal/controller"
	"github.com/D-Heger/govox/internal/window"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Run(width, height int32, title string) {
	window.InitWindow(width, height, title)
	defer window.Close()

	camera := controller.CreateCamera()
	model := block.Create()
	defer model.Unload()

	for !window.ShouldClose() {
		update(&camera)
		draw(camera, model.Model)
	}
}

func update(camera *rl.Camera3D) {
	rl.UpdateCameraPro(camera, controller.GetMovement(), controller.GetRotation(), 0.0)
}

func draw(camera rl.Camera3D, model rl.Model) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.SkyBlue)
	rl.DrawFPS(10, 10)
	rl.BeginMode3D(camera)
	defer rl.EndMode3D()

	rl.DrawGrid(10, 1.0)
	rl.DrawModel(model, rl.Vector3{X: 0.5, Y: 0.5, Z: 0.5}, 1.0, rl.White)
}
