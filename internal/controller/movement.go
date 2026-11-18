// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package controller

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	movementSpeed float32 = 4.0
	sprintMult    float32 = 4.0
)

// getAxisMovement returns a float32 value representing the movement direction based on the provided positive and negative key codes. It returns 1.0 for positive movement, -1.0 for negative movement, and 0.0 if neither or both keys are pressed.
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

// GetForwardsMovement returns the forward/backward movement direction based on the W and S keys. It returns 1.0 for forward, -1.0 for backward, and 0.0 if neither or both keys are pressed.
func GetForwardsMovement() float32 {
	return getAxisMovement(rl.KeyW, rl.KeyS)
}

// GetSidewaysMovement returns the left/right movement direction based on the D and A keys. It returns 1.0 for right, -1.0 for left, and 0.0 if neither or both keys are pressed.
func GetSidewaysMovement() float32 {
	return getAxisMovement(rl.KeyD, rl.KeyA)
}

// GetUpwardsMovement returns the upward/downward movement direction based on the Space and Left Shift keys. It returns 1.0 for upward, -1.0 for downward, and 0.0 if neither or both keys are pressed.
func GetUpwardsMovement() float32 {
	return getAxisMovement(rl.KeySpace, rl.KeyLeftShift)
}

// isSprinting checks if the Left Control key is currently pressed, indicating that the player is sprinting.
func isSprinting() bool {
	return rl.IsKeyDown(rl.KeyLeftControl)
}

// GetMovement calculates and returns the movement vector based on the current input state. It combines forward/backward, left/right, and upward/downward movements, normalizes the direction, and scales it by the movement speed and frame time. If sprinting, the speed is multiplied by a sprint multiplier.
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
