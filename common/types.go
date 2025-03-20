package common

import "github.com/robertcurry0216/cross/internal/puzzle"

// the State for the cli
type State struct {
	Debug  string
	Puzzle *puzzle.Puzzle
	Views  []Viewable
	Width  int
	Height int

	PuzzleView PuzzleView
}

type LayoutType int

const (
	LayoutPuzzleFocus LayoutType = iota
	LayoutClueFocus
)

type PuzzleView struct {
	X      int
	Y      int
	IsVert bool
	Layout LayoutType
}

type LayoutBox struct {
	W int
	H int
}
