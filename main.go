package main

import (
	"GameOfLife/clicker"
	"GameOfLife/game"
	_ "fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
)

const LineWidth = 1
const SquareEdgeSize = 15
const BlockSize = LineWidth + SquareEdgeSize
const GlobalScreenWidth = 640
const GlobalScreenHeight = 480
const GameFrameRate = 20

var countOfHorizontalBlocks = GlobalScreenWidth/BlockSize + 3
var countOfVerticalBlocks = GlobalScreenHeight/BlockSize + 3

var WhiteColor = color.Color(color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
var BlackColor = color.Color(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF})

var buttons = []clicker.Button{
	{
		Name:   "Stop",
		LeftX:  GlobalScreenWidth - 80,
		TopY:   GlobalScreenHeight - 60,
		Width:  60,
		Height: 40,
		Handle: clicker.HandleStopRenderButton,
		Color:  BlackColor,
	},
	{
		Name:   "Clear",
		LeftX:  GlobalScreenWidth - 150,
		TopY:   GlobalScreenHeight - 60,
		Width:  60,
		Height: 40,
		Handle: clicker.HandleClearButton,
		Color:  BlackColor,
	},
	{
		Name:   "Randomize",
		LeftX:  GlobalScreenWidth - 240,
		TopY:   GlobalScreenHeight - 60,
		Width:  80,
		Height: 40,
		Handle: clicker.HandleRandomizeButton,
		Color:  BlackColor,
	},
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

func countOfNeighbours(y, x int) int {
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

func calculateLeftTopDotOfSquare(xPosition, yPosition int) (topLeftX, topLeftY int) {
	return xPosition * BlockSize, yPosition * BlockSize
}

func (g *Game) Update() error {
	clicker.ProcessDrawingDot(state)
	clicker.ProcessMouseClick(buttons, &state)
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
		for i := 1; i < len(state.GameMatrix)-1; i++ {
			for j := 1; j < len(state.GameMatrix[i])-1; j++ {
				neighbours := countOfNeighbours(i, j)
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//TODO не надо рисовать это каждый раз
	for i := 0; i < countOfHorizontalBlocks-3; i++ {
		vector.DrawFilledRect(screen, float32(i*BlockSize+SquareEdgeSize), 0, LineWidth, GlobalScreenHeight, WhiteColor, false)
	}
	for i := 0; i < countOfVerticalBlocks-3; i++ {
		vector.DrawFilledRect(screen, 0, float32(i*BlockSize+SquareEdgeSize), GlobalScreenWidth, LineWidth, WhiteColor, false)
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
	clicker.DrawButtons(buttons, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{}

	ebiten.SetWindowSize(GlobalScreenWidth, GlobalScreenHeight)
	ebiten.SetFullscreen(g.isFullScreenMode)
	ebiten.SetTPS(GameFrameRate)

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
