package clicker

import (
	"GameOfLife/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
			game.DrawDot(state, x, y)
		}
	}
}

func (Rect *MyRect) DrawRectWithBorder(screen *ebiten.Image, borderSpan int) {
	vector.DrawFilledRect(screen, float32(Rect.LeftX-borderSpan), float32(Rect.TopY-borderSpan), float32(Rect.Width+borderSpan*2), float32(Rect.Height+borderSpan*2), Rect.SecondaryColor, false)
	vector.DrawFilledRect(screen, float32(Rect.LeftX), float32(Rect.TopY), float32(Rect.Width), float32(Rect.Height), Rect.MainColor, false)
}

func ProcessSingleMouseClick(clickables []Clickable, state *game.State) {
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

func ProcessLongMouseClick(clickables []Clickable, state *game.State) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if notValidClick(x, y) {
			return
		}
		for _, clickable := range clickables {
			if clickable.ProcessClick(x, y, state) {
				break
			}
		}
	}
}
