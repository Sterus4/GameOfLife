package main

import (
	"GameOfLife/clicker"
	"GameOfLife/game"
	_ "fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	GameScreenWidth   = 640
	GameScreenHeight  = 480
	RealScreenSWidth  = 1920
	RealScreenSHeight = 1080
)

// TODO Понять что делать с размерами кнопок

type GameStatus int

const (
	GameStatusGame GameStatus = iota
	GameStatusSettings
)

type Game struct{}

func (g *Game) Update() error {
	if err := game.UpdateMainScreen(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	game.DrawMainScreen(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameScreenWidth, GameScreenHeight
}

func main() {
	g := &Game{}
	game.InitGame(GameScreenHeight, GameScreenWidth)
	clicker.InitClicker(GameScreenWidth, GameScreenHeight)

	ebiten.SetWindowSize(RealScreenSWidth, RealScreenSHeight)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Game of life!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
