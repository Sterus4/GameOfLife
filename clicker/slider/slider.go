package slider

import (
	"GameOfLife/clicker"
	"GameOfLife/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const sliderBorderWidth = 2
const sliderSpan = 10

func (slider *GameSlider) ProcessClick(x, y int, state *game.State) bool {
	return true
}

type GameSlider struct {
	Rect         clicker.MyRect
	CurrentValue *int
	MinValue     int
	MaxValue     int
}

func (slider *GameSlider) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(slider.Rect.LeftX), float32(slider.Rect.TopY+slider.Rect.Height/2-sliderSpan/2), float32(slider.Rect.Width), float32(sliderSpan), slider.Rect.SecondaryColor, false)
	vector.DrawFilledRect(screen, float32(slider.Rect.LeftX+sliderBorderWidth), float32(slider.Rect.TopY+slider.Rect.Height/2-sliderSpan/2+sliderBorderWidth), float32(slider.Rect.Width-2*sliderBorderWidth), float32(sliderSpan-2*sliderBorderWidth), slider.Rect.MainColor, false)
	centerByX := slider.Rect.Width * *slider.CurrentValue / (slider.MaxValue - slider.MinValue)
	vector.DrawFilledRect(screen, float32(slider.Rect.LeftX+centerByX-sliderSpan/2), float32(slider.Rect.TopY), float32(sliderSpan), float32(slider.Rect.Height), slider.Rect.SecondaryColor, false)
	vector.DrawFilledRect(screen, float32(slider.Rect.LeftX+centerByX-sliderSpan/2+sliderBorderWidth), float32(slider.Rect.TopY+sliderBorderWidth), float32(sliderSpan-2*sliderBorderWidth), float32(slider.Rect.Height-2*sliderBorderWidth), slider.Rect.MainColor, false)
}

func DrawSliders(sliders []*GameSlider, screen *ebiten.Image) {
	for _, elem := range sliders {
		elem.Draw(screen)
	}
}
