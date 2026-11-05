// Copyright (c) 2026 David Heger. All rights reserved.
// Use of this source code is governed by the GNU Affero General Public License v3.0 that can be found in the LICENSE file.

package window

import rl "github.com/gen2brain/raylib-go/raylib"

func InitWindow(width, height int32, title string) {
	rl.InitWindow(width, height, title)
	rl.SetTargetFPS(60)
	rl.DisableCursor()
}

func ShouldClose() bool {
	return rl.WindowShouldClose()
}

func Close() {
	rl.CloseWindow()
}
