package game

import (
	"GameOfLife/clicker"
	"GameOfLife/clicker/button"
	"GameOfLife/clicker/notification"
	"GameOfLife/clicker/plot"
	"GameOfLife/clicker/slider"
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
	"time"
)

const (
	fpsShowerLabelString = "Current fps"
	minFrateRate         = 1
	maxFrameRate         = 60
	LineWidth            = 1
	SquareEdgeSize       = 15
	BlockSize            = LineWidth + SquareEdgeSize
)

var FrameRate = 10
var lastUpdate = time.Now()
var timeBetweenUpdate = time.Second * time.Duration(FrameRate) / 60

var stopButton = &button.GameButton{
	Name: "Stop",
	Rect: plot.MyRect{
		LeftX:          -1,
		TopY:           -1,
		Width:          60,
		Height:         40,
		MainColor:      plot.BlackColor,
		SecondaryColor: plot.LightPurpleColor,
	},
	Handle: HandleStopRenderButton,
}

var clearButton = &button.GameButton{
	Name: "Clear",
	Rect: plot.MyRect{
		LeftX:/*GameScreenWidth - 150*/ 0,
		TopY:/*GameScreenHeight - 60*/ 0,
		Width:          60,
		Height:         40,
		MainColor:      plot.BlackColor,
		SecondaryColor: plot.LightPurpleColor,
	},
	Handle: HandleClearButton,
}

var randomizeButton = &button.GameButton{
	Name: "Randomize",
	Rect: plot.MyRect{
		LeftX:          0,
		TopY:           0,
		Width:          80,
		Height:         40,
		MainColor:      plot.BlackColor,
		SecondaryColor: plot.LightPurpleColor,
	},
	Handle: HandleRandomizeButton,
}

var exitButton = &button.GameButton{
	Name: "Exit",
	Rect: plot.MyRect{
		LeftX:          0,
		TopY:           20,
		Width:          60,
		Height:         40,
		MainColor:      plot.BlackColor,
		SecondaryColor: plot.LightPurpleColor,
	},
	Handle: HandleExitButton,
}

var fpsShowerLabel = &button.GameButton{
	Name: fpsShowerLabelString,
	Rect: plot.MyRect{
		LeftX:          20,
		TopY:           20,
		Width:          130,
		Height:         40,
		MainColor:      plot.PurpleColor,
		SecondaryColor: plot.LightPurpleColor,
	},
	Handle: func() {},
}

var speedSlider = &slider.GameSlider{
	Rect: plot.MyRect{
		LeftX:          0,
		TopY:           0,
		Width:          80,
		Height:         40,
		MainColor:      plot.BlackColor,
		SecondaryColor: plot.LightPurpleColor,
	},
	CurrentValue: &FrameRate,
	MinValue:     minFrateRate,
	MaxValue:     maxFrameRate,
}

var needToPauseNotification = &notification.Notification{
	Text:     "You should press stop button",
	Duration: 3 * time.Second,
	IsOnTop:  false,
	Rect: plot.MyRect{
		LeftX:          0,
		TopY:           20,
		Width:          200,
		Height:         40,
		MainColor:      plot.BlackColor,
		SecondaryColor: plot.LightPurpleColor,
	},
}

var notifications = []*notification.Notification{
	needToPauseNotification,
}

var drawables = []plot.Drawable{
	stopButton,
	clearButton,
	randomizeButton,
	exitButton,
	fpsShowerLabel,
	speedSlider,
}

var singleClickables = []plot.Clickable{
	stopButton,
	clearButton,
	randomizeButton,
	exitButton,
	fpsShowerLabel,
}

var longClickables = []plot.Clickable{
	speedSlider,
}

type State struct {
	StopUpdateFlag          bool
	NeedToClearMatrix       bool
	GameMatrix              [][]Square
	CountOfBlocksVertical   int
	CountOfBlocksHorizontal int
	BlockSize               int
	NeedToExit              bool
	WindowGameWidth         int
	WindowGameHeight        int
}

var state = &State{
	StopUpdateFlag:          false,
	NeedToClearMatrix:       false,
	CountOfBlocksVertical:   -1,
	CountOfBlocksHorizontal: -1,
	BlockSize:               -1,
	WindowGameWidth:         0,
	WindowGameHeight:        0,
}

func HandleStopRenderButton() {
	state.StopUpdateFlag = !state.StopUpdateFlag
}

func HandleClearButton() {
	state.NeedToClearMatrix = true
	state.StopUpdateFlag = true
	fmt.Println("Clear matrix")
}

func HandleRandomizeButton() {
	RandomizeMatrix(state)
}

func HandleExitButton() {
	state.NeedToExit = true
}

func notValidClick(x, y int) bool {
	return x < 0 || y < 0 || x > state.WindowGameWidth || y > state.WindowGameHeight
}

func initButtonsPositions(state *State) {
	stopButton.Rect.LeftX = state.WindowGameWidth - 80
	stopButton.Rect.TopY = state.WindowGameHeight - 60

	clearButton.Rect.LeftX = state.WindowGameWidth - 150
	clearButton.Rect.TopY = state.WindowGameHeight - 60

	randomizeButton.Rect.LeftX = state.WindowGameWidth - 240
	randomizeButton.Rect.TopY = state.WindowGameHeight - 60

	exitButton.Rect.LeftX = state.WindowGameWidth - 80
	exitButton.Rect.TopY = 20

	speedSlider.Rect.LeftX = state.WindowGameWidth - 330
	speedSlider.Rect.TopY = state.WindowGameHeight - 60

	needToPauseNotification.Rect.LeftX = state.WindowGameWidth/2 - needToPauseNotification.Rect.Width/2
}

func ProcessDrawingDot() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if state.StopUpdateFlag {
			if notValidClick(x, y) {
				return
			}
			DrawDot(state, x, y)
		} else if !someButtonPressed(x, y) {
			needToPauseNotification.Popup()
		}
	}
}

func someButtonPressed(x, y int) bool {
	return stopButton.IsHit(x, y) ||
		clearButton.IsHit(x, y) ||
		randomizeButton.IsHit(x, y) ||
		exitButton.IsHit(x, y) ||
		fpsShowerLabel.IsHit(x, y) ||
		speedSlider.IsHit(x, y)
}

func calculateLeftTopDotOfSquare(xPosition, yPosition int) (topLeftX, topLeftY int) {
	return xPosition * BlockSize, yPosition * BlockSize
}

func CountOfNeighbours(state *State, y, x int) int {
	var result int
	if state.GameMatrix[y][x-1].IsFilledOld {
		result += 1
	}
	if state.GameMatrix[y][x+1].IsFilledOld {
		result += 1
	}
	if state.GameMatrix[y-1][x-1].IsFilledOld {
		result += 1
	}
	if state.GameMatrix[y-1][x+1].IsFilledOld {
		result += 1
	}
	if state.GameMatrix[y+1][x-1].IsFilledOld {
		result += 1
	}
	if state.GameMatrix[y+1][x+1].IsFilledOld {
		result += 1
	}
	if state.GameMatrix[y-1][x].IsFilledOld {
		result += 1
	}
	if state.GameMatrix[y+1][x].IsFilledOld {
		result += 1
	}
	return result
}

func RandomizeMatrix(state *State) {
	for i := 1; i < len(state.GameMatrix)-1; i++ {
		for j := 1; j < len(state.GameMatrix[i])-1; j++ {
			myRand := rand.Float32()
			if myRand < 0.5 {
				state.GameMatrix[i][j].IsFilledOld = true
				state.GameMatrix[i][j].IsFilledNew = true
			} else {
				state.GameMatrix[i][j].IsFilledOld = false
				state.GameMatrix[i][j].IsFilledNew = false
			}
		}
	}
}

func InitGame(GameScreenHeight, GameScreenWidth int) {
	state.WindowGameHeight = GameScreenHeight
	state.WindowGameWidth = GameScreenWidth

	state.CountOfBlocksHorizontal = GameScreenWidth/BlockSize + 3
	state.CountOfBlocksVertical = GameScreenHeight/BlockSize + 3
	state.BlockSize = BlockSize

	state.GameMatrix = make([][]Square, state.CountOfBlocksVertical)
	for i := 0; i < state.CountOfBlocksVertical; i++ {
		state.GameMatrix[i] = make([]Square, state.CountOfBlocksHorizontal)
	}
	initButtonsPositions(state)
	RandomizeMatrix(state)

}

type Square struct {
	IsFilledOld bool
	IsFilledNew bool
}

func DrawDot(state *State, x, y int) {
	dotX, dotY := x/state.BlockSize+1, y/state.BlockSize+1
	state.GameMatrix[dotY][dotX].IsFilledOld = true
	state.GameMatrix[dotY][dotX].IsFilledNew = true
}

func changeButtonName(button *button.GameButton, newName string) {
	button.Name = newName
}

func clearMatrix() {
	for i := 1; i < len(state.GameMatrix)-1; i++ {
		for j := 1; j < len(state.GameMatrix[i])-1; j++ {
			state.GameMatrix[i][j].IsFilledOld = false
			state.GameMatrix[i][j].IsFilledNew = false
		}
	}
}

func oneGameTick() {
	for i := 1; i < len(state.GameMatrix)-1; i++ {
		for j := 1; j < len(state.GameMatrix[i])-1; j++ {
			neighbours := CountOfNeighbours(state, i, j)
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
}

// TODO заложить это в интерфейс Screen
func DrawMainScreen(screen *ebiten.Image) {
	//TODO не надо рисовать это каждый раз (или надо?)
	for i := 0; i < state.CountOfBlocksHorizontal-3; i++ {
		vector.DrawFilledRect(screen, float32(i*BlockSize+SquareEdgeSize), 0, LineWidth, float32(state.WindowGameHeight), plot.WhiteColor, false)
	}
	for i := 0; i < state.CountOfBlocksVertical-3; i++ {
		vector.DrawFilledRect(screen, 0, float32(i*BlockSize+SquareEdgeSize), float32(state.WindowGameWidth), LineWidth, plot.WhiteColor, false)
	}

	for i := 1; i < len(state.GameMatrix)-2; i++ {
		for j := 1; j < len(state.GameMatrix[i])-2; j++ {

			gameSquare := state.GameMatrix[i][j]
			var squareColor color.Color
			if gameSquare.IsFilledOld {
				squareColor = plot.WhiteColor
			} else {
				squareColor = plot.BlackColor
			}
			var x, y = calculateLeftTopDotOfSquare(j-1, i-1)
			vector.DrawFilledRect(screen, float32(x), float32(y), SquareEdgeSize, SquareEdgeSize, squareColor, false)
		}
	}
	clicker.DrawDrawables(drawables, screen)
	notification.DrawNotifications(notifications, screen)
}

func UpdateMainScreen() error {
	if state.NeedToExit {
		return errors.New("exiting")
	}
	timeBetweenUpdate = time.Second / time.Duration(FrameRate)
	changeButtonName(fpsShowerLabel, fmt.Sprintf("%s: %d", fpsShowerLabelString, FrameRate))
	ProcessDrawingDot()
	clicker.ProcessSingleMouseClick(singleClickables)
	clicker.ProcessLongMouseClick(longClickables)

	if state.NeedToClearMatrix {
		clearMatrix()
		state.NeedToClearMatrix = false
	}
	if !state.StopUpdateFlag {
		changeButtonName(stopButton, "Stop")
		if time.Now().Sub(lastUpdate) > timeBetweenUpdate {
			oneGameTick()
			lastUpdate = time.Now()
		}
	} else {
		changeButtonName(stopButton, "Start")
	}
	return nil
}
