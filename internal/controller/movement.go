// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package controller

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	movementSpeed float32 = 4.0
	sprintMult    float32 = 2.0
)

func getAxisMovement(positiveKey, negativeKey int32) float32 {
	var dir float32
	if rl.IsKeyDown(positiveKey) {
		dir += 1.0
	}
	if rl.IsKeyDown(negativeKey) {
		dir -= 1.0
	}
	return dir
}

func GetForwardsMovement() float32 {
	return getAxisMovement(rl.KeyW, rl.KeyS)
}

func GetSidewaysMovement() float32 {
	return getAxisMovement(rl.KeyD, rl.KeyA)
}

func GetUpwardsMovement() float32 {
	return getAxisMovement(rl.KeySpace, rl.KeyLeftShift)
}

// TODO fix if this just multiplies movement speed in the get movement method, it gets multiplied every frame without ever being reset
func isSprinting() bool {
	return rl.IsKeyDown(rl.KeyLeftControl)
}

func GetMovement() rl.Vector3 {
	dir := rl.Vector3{X: GetForwardsMovement(), Y: GetSidewaysMovement(), Z: GetUpwardsMovement()}
	if rl.Vector3Length(dir) > 0 {
		dir = rl.Vector3Normalize(dir)
	}
	currentSpeed := movementSpeed
	if isSprinting() {
		currentSpeed *= sprintMult
	}
	return rl.Vector3Scale(dir, currentSpeed*rl.GetFrameTime())
}
