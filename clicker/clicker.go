package clicker

import (
	"GameOfLife/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

var mousePressed bool

type Hittable interface {
	IsHit(x, y int) bool
}

type Clickable interface {
	ProcessClick(x, y int, state *game.State) bool
}

type MyRect struct {
	LeftX, TopY, Width, Height int
	MainColor                  color.Color
	SecondaryColor             color.Color
}

func notValidClick(x, y int) bool {
	width, height := ebiten.WindowSize()
	return x < 0 || y < 0 || x > width || y > height
}

func ProcessDrawingDot(state game.State) {
	if state.StopUpdateFlag {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if notValidClick(x, y) {
				return
			}
			fmt.Printf("Matrix pressed: X=%d, Y=%d\n", x, y)
			game.DrawDot(state, x, y)
		}
	}
}

func ProcessMouseClick(clickables []Clickable, state *game.State) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !mousePressed {
			x, y := ebiten.CursorPosition()
			if notValidClick(x, y) {
				return
			}
			mousePressed = true
			for _, clickable := range clickables {
				if clickable.ProcessClick(x, y, state) {
					break
				}
			}
		}
	} else {
		mousePressed = false
	}
}
