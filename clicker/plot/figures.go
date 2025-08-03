package plot

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type MyRect struct {
	LeftX, TopY, Width, Height int
	MainColor                  color.Color
	SecondaryColor             color.Color
}

func (Rect *MyRect) DrawRectWithBorder(screen *ebiten.Image, borderSpan int) {
	vector.DrawFilledRect(screen, float32(Rect.LeftX-borderSpan), float32(Rect.TopY-borderSpan), float32(Rect.Width+borderSpan*2), float32(Rect.Height+borderSpan*2), Rect.SecondaryColor, false)
	vector.DrawFilledRect(screen, float32(Rect.LeftX), float32(Rect.TopY), float32(Rect.Width), float32(Rect.Height), Rect.MainColor, false)
}
