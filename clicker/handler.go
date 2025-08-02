package clicker

import (
	"GameOfLife/game"
	"fmt"
)

func HandleStopRenderButton(state *game.State) {
	state.StopUpdateFlag = !state.StopUpdateFlag
	fmt.Println("Stop Render Button", state.StopUpdateFlag)
}

func HandleClearButton(state *game.State) {
	state.NeedToClearMatrix = true
	fmt.Println("Clear matrix")
}

func HandleRandomizeButton(state *game.State) {
	game.RandomizeMatrix(state)
}
