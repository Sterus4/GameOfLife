package game

import "math/rand"

type State struct {
	StopUpdateFlag          bool
	NeedToClearMatrix       bool
	GameMatrix              [][]Square
	CountOfBlocksVertical   int
	CountOfBlocksHorizontal int
	BlockSize               int
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

type Square struct {
	IsFilledOld bool
	IsFilledNew bool
}

func DrawDot(state State, x, y int) {
	dotX, dotY := x/state.BlockSize+1, y/state.BlockSize+1
	state.GameMatrix[dotY][dotX].IsFilledOld = true
	state.GameMatrix[dotY][dotX].IsFilledNew = true
}
