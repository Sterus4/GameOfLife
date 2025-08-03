package button

import (
	"GameOfLife/clicker/plot"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

const buttonBorderWidth = 2

type GameButton struct {
	Rect   plot.MyRect
	Name   string
	Handle func()
}

func (button *GameButton) Draw(screen *ebiten.Image) {
	button.Rect.DrawRectWithBorder(screen, buttonBorderWidth)
	fontFace := basicfont.Face7x13

	bounds, _ := font.BoundString(fontFace, button.Name)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
	textHeight := (bounds.Max.Y - bounds.Min.Y).Ceil()

	textX := button.Rect.LeftX + (button.Rect.Width-textWidth)/2
	textY := button.Rect.TopY + (button.Rect.Height-textHeight)/2 + textHeight

	text.Draw(screen, button.Name, fontFace, textX, textY, color.White)
}

func (button *GameButton) IsHit(x, y int) bool {
	return x >= button.Rect.LeftX && x <= button.Rect.LeftX+button.Rect.Width && y >= button.Rect.TopY && y <= button.Rect.TopY+button.Rect.Height
}

func (button *GameButton) ProcessClick(x, y int) bool {
	if !button.IsHit(x, y) {
		return false
	}
	fmt.Printf("button '%s' pressed\n", button.Name)
	button.Handle()
	return true
}
