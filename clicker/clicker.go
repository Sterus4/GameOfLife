package clicker

import (
	"GameOfLife/clicker/notification"
	"GameOfLife/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

var mousePressed bool
var width, height = ebiten.WindowSize()
var notificationBorder = 2

var AllNotificationsRect = MyRect{
	LeftX:          width/2 - 100,
	TopY:           20,
	Width:          200,
	Height:         50,
	MainColor:      color.Color(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}),
	SecondaryColor: color.Color(color.RGBA{R: 0xAF, G: 0x46, B: 0xD5, A: 0xFF}),
}

type Hittable interface {
	IsHit(x, y int) bool
}

type Clickable interface {
	ProcessClick(x, y int, state *game.State) bool
}

type Drawable interface {
	Draw(screen *ebiten.Image)
}

type Notification interface {
	Popup()
	IsActive() bool
	GetText() string
}

type MyRect struct {
	LeftX, TopY, Width, Height int
	MainColor                  color.Color
	SecondaryColor             color.Color
}

func notValidClick(x, y int) bool {

	return x < 0 || y < 0 || x > width || y > height
}

func ProcessDrawingDot(state game.State, notification *notification.NotificationImpl) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if state.StopUpdateFlag {
			x, y := ebiten.CursorPosition()
			if notValidClick(x, y) {
				return
			}
			game.DrawDot(state, x, y)
		} else {
			notification.Popup()
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

func DrawDrawables(drawables []Drawable, screen *ebiten.Image) {
	for _, elem := range drawables {
		elem.Draw(screen)
	}
}

func DrawNotifications(notifications []Notification, screen *ebiten.Image) {
	for _, elem := range notifications {
		if elem.IsActive() {
			AllNotificationsRect.DrawRectWithBorder(screen, notificationBorder)

			fontFace := basicfont.Face7x13
			bounds, _ := font.BoundString(fontFace, elem.GetText())
			textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
			textHeight := (bounds.Max.Y - bounds.Min.Y).Ceil()

			textX := AllNotificationsRect.LeftX + (AllNotificationsRect.Width-textWidth)/2
			textY := AllNotificationsRect.TopY + (AllNotificationsRect.Height-textHeight)/2 + textHeight

			text.Draw(screen, elem.GetText(), fontFace, textX, textY, color.White)
		}
	}
}
