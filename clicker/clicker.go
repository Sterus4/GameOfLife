package clicker

import (
	"GameOfLife/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

const buttonBorderWidth = 2

var mousePressed bool

type Button struct {
	LeftX, TopY, Width, Height int
	Name                       string
	Handle                     func(*game.State)
	Color                      color.Color
}

func checkHit(r Button, x, y int) bool {
	return x >= r.LeftX && x <= r.LeftX+r.Width && y >= r.TopY && y <= r.TopY+r.Height
}

func processButton(x, y int, buttons []Button, state *game.State) {
	for _, elem := range buttons {
		if checkHit(elem, x, y) {
			fmt.Printf("button '%s' pressed\n", elem.Name)
			elem.Handle(state)
			return
		}
	}
	fmt.Println("No button pressed")
}

func (button *Button) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(button.LeftX)-buttonBorderWidth, float32(button.TopY)-buttonBorderWidth, float32(button.Width)+buttonBorderWidth*2, float32(button.Height)+buttonBorderWidth*2, color.White, false)
	vector.DrawFilledRect(screen, float32(button.LeftX), float32(button.TopY), float32(button.Width), float32(button.Height), button.Color, false)
	fontFace := basicfont.Face7x13

	bounds, _ := font.BoundString(fontFace, button.Name)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
	textHeight := (bounds.Max.Y - bounds.Min.Y).Ceil()

	textX := button.LeftX + (button.Width-textWidth)/2
	textY := button.TopY + (button.Height-textHeight)/2 + textHeight

	text.Draw(screen, button.Name, fontFace, textX, textY, color.White)
}

func notValidClick(x, y int) bool {
	width, height := ebiten.WindowSize()
	return x < 0 || y < 0 || x > width || y > height
}

func DrawButtons(buttons []Button, screen *ebiten.Image) {
	for _, elem := range buttons {
		elem.Draw(screen)
	}
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

func ProcessMouseClick(buttons []Button, state *game.State) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !mousePressed {
			x, y := ebiten.CursorPosition()
			if notValidClick(x, y) {
				return
			}
			mousePressed = true
			processButton(x, y, buttons, state)
		}
	} else {
		mousePressed = false
	}
}
