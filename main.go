package main

import (
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
const GameFrameRate = 10

var countOfHorizontalBlocks = GlobalScreenWidth/BlockSize + 3
var countOfVerticalBlocks = GlobalScreenHeight/BlockSize + 3

var WhiteColor = color.Color(color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
var BlackColor = color.Color(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF})

type GameSquare struct {
	isFilledOld bool
	isFilledNew bool
}

var GameMatrix [][]GameSquare

type Game struct{}

func countOfNeighbours(y, x int) int {
	var result int
	if GameMatrix[y][x-1].isFilledOld {
		result += 1
	}
	if GameMatrix[y][x+1].isFilledOld {
		result += 1
	}
	if GameMatrix[y-1][x-1].isFilledOld {
		result += 1
	}
	if GameMatrix[y-1][x+1].isFilledOld {
		result += 1
	}
	if GameMatrix[y+1][x-1].isFilledOld {
		result += 1
	}
	if GameMatrix[y+1][x+1].isFilledOld {
		result += 1
	}
	if GameMatrix[y-1][x].isFilledOld {
		result += 1
	}
	if GameMatrix[y+1][x].isFilledOld {
		result += 1
	}
	return result
}

func (g *Game) Update() error {
	for i := 1; i < len(GameMatrix)-1; i++ {
		for j := 1; j < len(GameMatrix[i])-1; j++ {
			neighbours := countOfNeighbours(i, j)
			if GameMatrix[i][j].isFilledOld {
				if neighbours != 2 && neighbours != 3 {
					GameMatrix[i][j].isFilledNew = false
				}
			} else {
				if neighbours == 3 {
					GameMatrix[i][j].isFilledNew = true
				}
			}
		}
	}
	for i := 1; i < len(GameMatrix)-1; i++ {
		for j := 1; j < len(GameMatrix[i])-1; j++ {
			GameMatrix[i][j].isFilledOld = GameMatrix[i][j].isFilledNew
		}
	}
	return nil
}

func randomizeMatrix() {
	for i := 1; i < len(GameMatrix)-1; i++ {
		for j := 1; j < len(GameMatrix[i])-1; j++ {
			myRand := rand.Float32()
			if myRand < 0.5 {
				GameMatrix[i][j].isFilledOld = true
				GameMatrix[i][j].isFilledNew = true
			} else {
				GameMatrix[i][j].isFilledOld = false
				GameMatrix[i][j].isFilledNew = false
			}
		}
	}

}

func calculateLeftTopDotOfSquare(xPosition, yPosition int) (topLeftX, topLeftY int) {
	return xPosition * BlockSize, yPosition * BlockSize
}

func (g *Game) Draw(screen *ebiten.Image) {
	//TODO не надо рисовать это каждый раз
	for i := 0; i < countOfHorizontalBlocks-3; i++ {
		vector.DrawFilledRect(screen, float32(i*BlockSize+SquareEdgeSize), 0, LineWidth, GlobalScreenHeight, WhiteColor, false)
	}
	for i := 0; i < countOfVerticalBlocks-3; i++ {
		vector.DrawFilledRect(screen, 0, float32(i*BlockSize+SquareEdgeSize), GlobalScreenWidth, LineWidth, WhiteColor, false)
	}

	for i := 1; i < len(GameMatrix)-2; i++ {
		for j := 1; j < len(GameMatrix[i])-2; j++ {

			gameSquare := GameMatrix[i][j]
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(GlobalScreenWidth, GlobalScreenHeight)
	ebiten.SetFullscreen(true)
	ebiten.SetTPS(GameFrameRate)

	game := &Game{}

	GameMatrix = make([][]GameSquare, countOfVerticalBlocks)
	for i := 0; i < countOfVerticalBlocks; i++ {
		GameMatrix[i] = make([]GameSquare, countOfHorizontalBlocks)
	}
	randomizeMatrix()

	ebiten.SetWindowTitle("Game of life!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
