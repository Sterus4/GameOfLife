package main

import (
	"GameOfLife/clicker"
	"GameOfLife/clicker/button"
	"GameOfLife/clicker/notification"
	"GameOfLife/clicker/slider"
	"GameOfLife/game"
	"errors"
	"fmt"
	_ "fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"time"
)

const LineWidth = 1
const SquareEdgeSize = 15
const BlockSize = LineWidth + SquareEdgeSize
const GameScreenWidth = 640
const GameScreenHeight = 480
const RealScreenSWidth = 640
const RealScreenSHeight = 480
const minFrateRate = 3
const maxFrameRate = 60
const fpsShowerLabelString = "Current fps"

var countOfHorizontalBlocks = GameScreenWidth/BlockSize + 3
var countOfVerticalBlocks = GameScreenHeight/BlockSize + 3
var GameFrameRate = 10

var WhiteColor = color.Color(color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
var BlackColor = color.Color(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF})
var PurpleColor = color.Color(color.RGBA{R: 0x68, G: 0x28, B: 0x80, A: 0xFF})
var LightPurpleColor = color.Color(color.RGBA{R: 0xAF, G: 0x46, B: 0xD5, A: 0xFF})

var stopButton = &button.GameButton{
	Name: "Stop",
	Rect: clicker.MyRect{
		LeftX:          GameScreenWidth - 80,
		TopY:           GameScreenHeight - 60,
		Width:          60,
		Height:         40,
		MainColor:      BlackColor,
		SecondaryColor: LightPurpleColor,
	},
	Handle: clicker.HandleStopRenderButton,
}

var clearButton = &button.GameButton{
	Name: "Clear",
	Rect: clicker.MyRect{
		LeftX:          GameScreenWidth - 150,
		TopY:           GameScreenHeight - 60,
		Width:          60,
		Height:         40,
		MainColor:      BlackColor,
		SecondaryColor: LightPurpleColor,
	},
	Handle: clicker.HandleClearButton,
}

var randomizeButton = &button.GameButton{
	Name: "Randomize",
	Rect: clicker.MyRect{
		LeftX:          GameScreenWidth - 240,
		TopY:           GameScreenHeight - 60,
		Width:          80,
		Height:         40,
		MainColor:      BlackColor,
		SecondaryColor: LightPurpleColor,
	},
	Handle: clicker.HandleRandomizeButton,
}

var exitButton = &button.GameButton{
	Name: "Exit",
	Rect: clicker.MyRect{
		LeftX:          GameScreenWidth - 80,
		TopY:           20,
		Width:          60,
		Height:         40,
		MainColor:      BlackColor,
		SecondaryColor: LightPurpleColor,
	},
	Handle: clicker.HandleExitButton,
}

var fpsShowerLabel = &button.GameButton{
	Name: fpsShowerLabelString,
	Rect: clicker.MyRect{
		LeftX:          20,
		TopY:           20,
		Width:          130,
		Height:         40,
		MainColor:      PurpleColor,
		SecondaryColor: LightPurpleColor,
	},
	Handle: func(g *game.State) {},
}

var speedSlider = &slider.GameSlider{
	Rect: clicker.MyRect{
		LeftX:          GameScreenWidth - 330,
		TopY:           GameScreenHeight - 60,
		Width:          80,
		Height:         40,
		MainColor:      BlackColor,
		SecondaryColor: LightPurpleColor,
	},
	CurrentValue: &GameFrameRate,
	MinValue:     minFrateRate,
	MaxValue:     maxFrameRate,
}

// TODO разобраться с выводом текста в несколько строке
var needToPauseNotification = &notification.NotificationImpl{
	Text:     "You should press stop button",
	Duration: 3 * time.Second,
	IsOnTop:  false,
}

var notifications = []clicker.Notification{
	needToPauseNotification,
}

var drawables = []clicker.Drawable{
	stopButton,
	clearButton,
	randomizeButton,
	exitButton,
	fpsShowerLabel,
	speedSlider,
}

var singleClickables = []clicker.Clickable{
	stopButton,
	clearButton,
	randomizeButton,
	exitButton,
	fpsShowerLabel,
}

var longClickables = []clicker.Clickable{
	speedSlider,
}

var state = game.State{
	StopUpdateFlag:          false,
	NeedToClearMatrix:       false,
	CountOfBlocksVertical:   countOfVerticalBlocks,
	CountOfBlocksHorizontal: countOfHorizontalBlocks,
	BlockSize:               BlockSize,
}

type Game struct {
	isFullScreenMode bool
}

func calculateLeftTopDotOfSquare(xPosition, yPosition int) (topLeftX, topLeftY int) {
	return xPosition * BlockSize, yPosition * BlockSize
}

func changeButtonName(button *button.GameButton, newName string) {
	button.Name = newName
}

func (g *Game) Update() error {
	if state.NeedToExit {
		return errors.New("exiting")
	}
	changeButtonName(fpsShowerLabel, fmt.Sprintf("%s: %d", fpsShowerLabelString, GameFrameRate))
	clicker.ProcessDrawingDot(state, needToPauseNotification)
	clicker.ProcessSingleMouseClick(singleClickables, &state)
	clicker.ProcessLongMouseClick(longClickables, &state)
	ebiten.SetTPS(GameFrameRate)
	if state.NeedToClearMatrix {
		for i := 1; i < len(state.GameMatrix)-1; i++ {
			for j := 1; j < len(state.GameMatrix[i])-1; j++ {
				state.GameMatrix[i][j].IsFilledOld = false
				state.GameMatrix[i][j].IsFilledNew = false
			}
		}
		state.NeedToClearMatrix = false
	}
	if !state.StopUpdateFlag {
		changeButtonName(stopButton, "Stop")
		for i := 1; i < len(state.GameMatrix)-1; i++ {
			for j := 1; j < len(state.GameMatrix[i])-1; j++ {
				neighbours := game.CountOfNeighbours(&state, i, j)
				if state.GameMatrix[i][j].IsFilledOld {
					if neighbours != 2 && neighbours != 3 {
						state.GameMatrix[i][j].IsFilledNew = false
					}
				} else {
					if neighbours == 3 {
						state.GameMatrix[i][j].IsFilledNew = true
					}
				}
			}
		}
		for i := 1; i < len(state.GameMatrix)-1; i++ {
			for j := 1; j < len(state.GameMatrix[i])-1; j++ {
				state.GameMatrix[i][j].IsFilledOld = state.GameMatrix[i][j].IsFilledNew
			}
		}
	} else {
		changeButtonName(stopButton, "Start")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//TODO не надо рисовать это каждый раз
	for i := 0; i < countOfHorizontalBlocks-3; i++ {
		vector.DrawFilledRect(screen, float32(i*BlockSize+SquareEdgeSize), 0, LineWidth, GameScreenHeight, WhiteColor, false)
	}
	for i := 0; i < countOfVerticalBlocks-3; i++ {
		vector.DrawFilledRect(screen, 0, float32(i*BlockSize+SquareEdgeSize), GameScreenWidth, LineWidth, WhiteColor, false)
	}

	for i := 1; i < len(state.GameMatrix)-2; i++ {
		for j := 1; j < len(state.GameMatrix[i])-2; j++ {

			gameSquare := state.GameMatrix[i][j]
			var squareColor color.Color
			if gameSquare.IsFilledOld {
				squareColor = WhiteColor
			} else {
				squareColor = BlackColor
			}
			var x, y = calculateLeftTopDotOfSquare(j-1, i-1)
			vector.DrawFilledRect(screen, float32(x), float32(y), SquareEdgeSize, SquareEdgeSize, squareColor, false)
		}
	}
	clicker.DrawDrawables(drawables, screen)
	clicker.DrawNotifications(notifications, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameScreenWidth, GameScreenHeight
}

func main() {
	g := &Game{}

	ebiten.SetWindowSize(RealScreenSWidth, RealScreenSHeight)
	ebiten.SetFullscreen(true)

	state.GameMatrix = make([][]game.Square, countOfVerticalBlocks)
	for i := 0; i < countOfVerticalBlocks; i++ {
		state.GameMatrix[i] = make([]game.Square, countOfHorizontalBlocks)
	}
	game.RandomizeMatrix(&state)

	ebiten.SetWindowTitle("Game of life!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
