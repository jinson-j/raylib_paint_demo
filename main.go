package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DrawMode int

const (
	ModeFree DrawMode = iota
	ModeLine
	ModeCircle
)

var (
	currentMode  DrawMode = ModeFree
	currentColor rl.Color = rl.Black
	brushSize    float32  = 5.0
	canvas       rl.RenderTexture2D
	startPos     rl.Vector2
	endPos       rl.Vector2
	isDrawing    bool
)

const window_width = 500
const window_height = 500

func calculate_distance(a, b rl.Vector2) float32 {
	return float32(math.Sqrt(float64((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y))))
}

func main() {
	rl.InitWindow(int32(window_width), int32(window_height), "Paint Program")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	canvas = rl.LoadRenderTexture(window_width, window_height)

	rl.BeginTextureMode(canvas)
	rl.ClearBackground(rl.White)
	rl.EndTextureMode()

	for !rl.WindowShouldClose() {
		mousePos := rl.GetMousePosition()

		// mode
		if rl.IsKeyPressed(rl.KeyF) {
			currentMode = ModeFree
		} else if rl.IsKeyPressed(rl.KeyL) {
			currentMode = ModeLine
		} else if rl.IsKeyPressed(rl.KeyC) {
			currentMode = ModeCircle
		}
		// color
		if rl.IsKeyPressed(rl.KeyOne) {
			currentColor = rl.Black
		} else if rl.IsKeyPressed(rl.KeyTwo) {
			currentColor = rl.Red
		} else if rl.IsKeyPressed(rl.KeyThree) {
			currentColor = rl.Blue
		} else if rl.IsKeyPressed(rl.KeyFour) {
			currentColor = rl.Green
		}

		// brush
		if rl.IsKeyPressed(rl.KeyEqual) && brushSize < 50 {
			brushSize++
		} else if rl.IsKeyPressed(rl.KeyMinus) && brushSize > 1 {
			brushSize--
		}

		// Handle drawing
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			startPos = mousePos
			if currentMode == ModeFree {
				rl.BeginTextureMode(canvas)
				rl.DrawCircleV(mousePos, brushSize, currentColor)
				rl.EndTextureMode()
			}
			isDrawing = true
		}

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			endPos = mousePos
			rl.BeginTextureMode(canvas)

			switch currentMode {
			case ModeLine:
				rl.DrawLineEx(startPos, endPos, brushSize*2, currentColor)
			case ModeCircle:
				radius := calculate_distance(startPos, endPos)
				rl.DrawCircleV(startPos, radius, currentColor)
			}

			rl.EndTextureMode()
			isDrawing = false
		}

		// Free drawing mode (continuous)
		if isDrawing && currentMode == ModeFree && rl.IsMouseButtonDown(rl.MouseLeftButton) {
			rl.BeginTextureMode(canvas)
			rl.DrawLineEx(startPos, mousePos, brushSize*2, currentColor)
			rl.EndTextureMode()
			startPos = mousePos
		}

		// Drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Canvas
		rl.DrawTextureRec(
			canvas.Texture,
			rl.NewRectangle(0, 0, float32(canvas.Texture.Width), float32(-canvas.Texture.Height)),
			rl.NewVector2(0, 0),
			rl.White,
		)

		rl.DrawRectangle(0, 0, window_width, 30, rl.LightGray)
		rl.DrawText("Paint Time!", 5, 5, 20, rl.Black)

		modeText := "Free"
		if currentMode == ModeLine {
			modeText = "Line"
		} else if currentMode == ModeCircle {
			modeText = "Circle"
		}
		rl.DrawText(modeText, 135, 5, 20, rl.Black)
		rl.DrawText("Size: "+fmt.Sprintf("%.0f", brushSize), 220, 5, 20, rl.Black)

		// Preview for shapes
		if isDrawing {
			switch currentMode {
			case ModeLine:
				rl.DrawLineEx(startPos, mousePos, brushSize*2, rl.Fade(currentColor, 0.6))
			case ModeCircle:
				radius := calculate_distance(startPos, mousePos)
				rl.DrawCircleLines(int32(startPos.X), int32(startPos.Y), radius, rl.Fade(currentColor, 0.6))
			}
		}

		rl.DrawText("F:Free    L:Line    C:Circle    1-4:Colors    +/-:Size", 10, window_height-20, 10, rl.DarkGray)
		rl.DrawFPS(window_width-85, window_height-30)

		rl.EndDrawing()
	}

	rl.UnloadRenderTexture(canvas)
}
