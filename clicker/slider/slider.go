package slider

import (
	"GameOfLife/clicker/plot"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const sliderBorderWidth = 2
const sliderSpan = 10

func (slider *GameSlider) ProcessClick(x, y int) bool {
	if !slider.IsHit(x, y) {
		return false
	}
	var newVal = (slider.MaxValue - slider.MinValue) * (x - slider.Rect.LeftX) / slider.Rect.Width
	newVal = max(newVal, slider.MinValue)
	newVal = min(newVal, slider.MaxValue)
	*slider.CurrentValue = newVal
	return true
}

type GameSlider struct {
	Rect         plot.MyRect
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

func (slider *GameSlider) IsHit(x, y int) bool {
	return x > slider.Rect.LeftX && x < slider.Rect.LeftX+slider.Rect.Width && y > slider.Rect.TopY+slider.Rect.Height/2-sliderSpan/2 && y < slider.Rect.TopY+slider.Rect.Height/2-sliderSpan/2+sliderSpan
}
