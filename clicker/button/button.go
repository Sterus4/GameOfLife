package button

import (
	"GameOfLife/clicker"
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

type GameButton struct {
	Rect   clicker.MyRect
	Name   string
	Handle func(*game.State)
}

func (button *GameButton) draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(button.Rect.LeftX)-buttonBorderWidth, float32(button.Rect.TopY)-buttonBorderWidth, float32(button.Rect.Width)+buttonBorderWidth*2, float32(button.Rect.Height)+buttonBorderWidth*2, button.Rect.SecondaryColor, false)
	vector.DrawFilledRect(screen, float32(button.Rect.LeftX), float32(button.Rect.TopY), float32(button.Rect.Width), float32(button.Rect.Height), button.Rect.MainColor, false)
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

func (button *GameButton) ProcessClick(x, y int, state *game.State) bool {
	if button.IsHit(x, y) {
		fmt.Printf("button '%s' pressed\n", button.Name)
		button.Handle(state)
		return true
	}
	return false
}

func DrawButtons(buttons []*GameButton, screen *ebiten.Image) {
	for _, elem := range buttons {
		elem.draw(screen)
	}
}
