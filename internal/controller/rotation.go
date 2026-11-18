// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package controller

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const rotationSpeed = 0.05

// GetRotation calculates and returns the rotation vector based on the current mouse movement. It scales the mouse delta by a predefined rotation speed to determine the rotation in the X and Y axes, while the Z axis remains unchanged.
func GetRotation() rl.Vector3 {
	mouseDelta := rl.GetMouseDelta()
	return rl.Vector3{X: mouseDelta.X * rotationSpeed, Y: mouseDelta.Y * rotationSpeed, Z: 0.0}
}
