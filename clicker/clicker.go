package clicker

import (
	"GameOfLife/clicker/plot"
	"github.com/hajimehoshi/ebiten/v2"
)

var mousePressed bool

var window = &MyWindow{}

type MyWindow struct {
	WindowWidth  int
	WindowHeight int
}

func InitClicker(windowWidth, windowHeight int) {
	window.WindowWidth = windowWidth
	window.WindowHeight = windowHeight
}

func notValidClick(x, y int) bool {
	return x < 0 || y < 0 || x > window.WindowWidth || y > window.WindowHeight
}

func ProcessSingleMouseClick(clickables []plot.Clickable) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !mousePressed {
			x, y := ebiten.CursorPosition()
			if notValidClick(x, y) {
				return
			}
			mousePressed = true
			for _, clickable := range clickables {
				if clickable.ProcessClick(x, y) {
					break
				}
			}
		}
	} else {
		mousePressed = false
	}
}

func ProcessLongMouseClick(clickables []plot.Clickable) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if notValidClick(x, y) {
			return
		}
		for _, clickable := range clickables {
			if clickable.ProcessClick(x, y) {
				break
			}
		}
	}
}

func DrawDrawables(drawables []plot.Drawable, screen *ebiten.Image) {
	for _, elem := range drawables {
		elem.Draw(screen)
	}
}
