// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package controller

import rl "github.com/gen2brain/raylib-go/raylib"

// CreateCamera initializes and returns a new 3D camera with a predefined position, target, up vector, and field of view.
func CreateCamera() rl.Camera3D {
	var position = rl.Vector3{X: 48.0, Y: 56.0, Z: 48.0}
	var target = rl.Vector3{X: 0.0, Y: 0.0, Z: 0.0}
	var up = rl.Vector3{X: 0.0, Y: 1.0, Z: 0.0}
	var fovY float32 = 65
	return rl.NewCamera3D(position, target, up, fovY, rl.CameraPerspective)
}
