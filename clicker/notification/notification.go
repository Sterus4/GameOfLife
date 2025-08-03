package notification

import (
	"GameOfLife/clicker/plot"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"time"
)

type Notification struct {
	Text     string
	Duration time.Duration
	IsOnTop  bool
	Rect     plot.MyRect
}

var notificationBorder = 2

func (notification *Notification) Popup() {
	go func() {
		notification.IsOnTop = true
		time.Sleep(notification.Duration)
		notification.IsOnTop = false
	}()
}

func (notification Notification) IsActive() bool {
	return notification.IsOnTop
}

func DrawNotifications(notifications []*Notification, screen *ebiten.Image) {
	for _, elem := range notifications {
		if elem.IsActive() {
			elem.Rect.DrawRectWithBorder(screen, notificationBorder)

			fontFace := basicfont.Face7x13
			bounds, _ := font.BoundString(fontFace, elem.GetText())
			textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
			textHeight := (bounds.Max.Y - bounds.Min.Y).Ceil()

			textX := elem.Rect.LeftX + (elem.Rect.Width-textWidth)/2
			textY := elem.Rect.TopY + (elem.Rect.Height-textHeight)/2 + textHeight

			text.Draw(screen, elem.GetText(), fontFace, textX, textY, color.White)
		}
	}
}

func (notification Notification) GetText() string {
	return notification.Text
}
