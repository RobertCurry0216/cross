package common

import "github.com/robertcurry0216/cross/internal/puzzle"

// the State for the cli
type State struct {
	Puzzle *puzzle.Puzzle
	Views  []Viewable
	Width  int
	Height int

	PuzzleView PuzzleView
}

type PuzzleView struct {
	X      int
	Y      int
	IsVert bool
}
