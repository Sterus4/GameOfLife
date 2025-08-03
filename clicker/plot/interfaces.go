package plot

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Hittable interface {
	IsHit(x, y int) bool
}

type Clickable interface {
	ProcessClick(x, y int) bool
	IsClickable() bool
}

type Drawable interface {
	Draw(screen *ebiten.Image)
	IsVisible() bool
}
