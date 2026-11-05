// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package game

import (
	"github.com/D-Heger/govox/internal/controller"
	"github.com/D-Heger/govox/internal/window"
	"github.com/D-Heger/govox/internal/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Run(width, height int32, title string) {
	window.InitWindow(width, height, title)
	defer window.Close()

	camera := controller.CreateCamera()

	gameWorld := world.Create()
	defer gameWorld.Unload()

	for !window.ShouldClose() {
		update(&camera)
		draw(camera, &gameWorld)
	}
}

func update(camera *rl.Camera3D) {
	rl.UpdateCameraPro(camera, controller.GetMovement(), controller.GetRotation(), 0.0)
}

func draw(camera rl.Camera3D, gameWorld *world.World) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.SkyBlue)
	rl.DrawFPS(10, 10)
	rl.BeginMode3D(camera)
	defer rl.EndMode3D()

	gameWorld.Draw()
}
