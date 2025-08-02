package notification

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type NotificationImpl struct {
	Text     string
	Duration time.Duration
	IsOnTop  bool
}

func (notification *NotificationImpl) Popup() {
	go func() {
		notification.IsOnTop = true
		time.Sleep(notification.Duration)
		notification.IsOnTop = false
	}()
}

func (notification NotificationImpl) IsActive() bool {
	return notification.IsOnTop
}

func (notification NotificationImpl) Draw(screen *ebiten.Image) {

}

func (notification NotificationImpl) GetText() string {
	return notification.Text
}
