package clicker

import (
	"GameOfLife/game"
	"fmt"
)

func HandleStopRenderButton(state *game.State) {
	state.StopUpdateFlag = !state.StopUpdateFlag
	fmt.Println("Stop Render Button pressed")
}

func HandleClearButton(state *game.State) {
	state.NeedToClearMatrix = true
	state.StopUpdateFlag = true
	fmt.Println("Clear matrix")
}

func HandleRandomizeButton(state *game.State) {
	game.RandomizeMatrix(state)
}

func HandleExitButton(state *game.State) {
	state.NeedToExit = true
}
