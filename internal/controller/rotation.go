// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package controller

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const rotationSpeed = 0.05

func GetRotation() rl.Vector3 {
	mouseDelta := rl.GetMouseDelta()
	return rl.Vector3{X: mouseDelta.X * rotationSpeed, Y: mouseDelta.Y * rotationSpeed, Z: 0.0}
}
