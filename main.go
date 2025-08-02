package main

import (
	"GameOfLife/clicker"
	"GameOfLife/game"
	_ "fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math/rand"
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
}

var state = game.State{
	StopUpdateFlag:    false,
	NeedToClearMatrix: false,
}

type GameSquare struct {
	isFilledOld bool
	isFilledNew bool
}

type Game struct {
	GameMatrix       [][]GameSquare
	isFullScreenMode bool
}

func countOfNeighbours(g *Game, y, x int) int {
	var result int
	if g.GameMatrix[y][x-1].isFilledOld {
		result += 1
	}
	if g.GameMatrix[y][x+1].isFilledOld {
		result += 1
	}
	if g.GameMatrix[y-1][x-1].isFilledOld {
		result += 1
	}
	if g.GameMatrix[y-1][x+1].isFilledOld {
		result += 1
	}
	if g.GameMatrix[y+1][x-1].isFilledOld {
		result += 1
	}
	if g.GameMatrix[y+1][x+1].isFilledOld {
		result += 1
	}
	if g.GameMatrix[y-1][x].isFilledOld {
		result += 1
	}
	if g.GameMatrix[y+1][x].isFilledOld {
		result += 1
	}
	return result
}

func randomizeMatrix(g *Game) {
	for i := 1; i < len(g.GameMatrix)-1; i++ {
		for j := 1; j < len(g.GameMatrix[i])-1; j++ {
			myRand := rand.Float32()
			if myRand < 0.5 {
				g.GameMatrix[i][j].isFilledOld = true
				g.GameMatrix[i][j].isFilledNew = true
			} else {
				g.GameMatrix[i][j].isFilledOld = false
				g.GameMatrix[i][j].isFilledNew = false
			}
		}
	}

}

func calculateLeftTopDotOfSquare(xPosition, yPosition int) (topLeftX, topLeftY int) {
	return xPosition * BlockSize, yPosition * BlockSize
}

func (g *Game) Update() error {
	clicker.ProcessMouseClick(buttons, &state)
	if state.NeedToClearMatrix {
		for i := 1; i < len(g.GameMatrix)-1; i++ {
			for j := 1; j < len(g.GameMatrix[i])-1; j++ {
				g.GameMatrix[i][j].isFilledOld = false
				g.GameMatrix[i][j].isFilledNew = false
			}
		}
		state.NeedToClearMatrix = false
	}
	if !state.StopUpdateFlag {
		for i := 1; i < len(g.GameMatrix)-1; i++ {
			for j := 1; j < len(g.GameMatrix[i])-1; j++ {
				neighbours := countOfNeighbours(g, i, j)
				if g.GameMatrix[i][j].isFilledOld {
					if neighbours != 2 && neighbours != 3 {
						g.GameMatrix[i][j].isFilledNew = false
					}
				} else {
					if neighbours == 3 {
						g.GameMatrix[i][j].isFilledNew = true
					}
				}
			}
		}
		for i := 1; i < len(g.GameMatrix)-1; i++ {
			for j := 1; j < len(g.GameMatrix[i])-1; j++ {
				g.GameMatrix[i][j].isFilledOld = g.GameMatrix[i][j].isFilledNew
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

	for i := 1; i < len(g.GameMatrix)-2; i++ {
		for j := 1; j < len(g.GameMatrix[i])-2; j++ {

			gameSquare := g.GameMatrix[i][j]
			var squareColor color.Color
			if gameSquare.isFilledOld {
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

	g.GameMatrix = make([][]GameSquare, countOfVerticalBlocks)
	for i := 0; i < countOfVerticalBlocks; i++ {
		g.GameMatrix[i] = make([]GameSquare, countOfHorizontalBlocks)
	}
	randomizeMatrix(g)

	ebiten.SetWindowTitle("Game of life!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
