package puzzle

import (
	"fmt"
)

type Puzzle struct {
	input    []byte
	solution []byte
	builder  iBuilder

	Width       int
	Height      int
	Clues       []*Clue
	DownClues   []*Clue
	AcrossClues []*Clue
	Grid        []*Cell

	Title     string
	Author    string
	Copyright string
	Notes     string
}

func NewPuzzle() *Puzzle {
	return &Puzzle{}
}

func (puz *Puzzle) String() string {
	return fmt.Sprintf("<Puzzle size:%vx%v clues: %v title: %v>", puz.Width, puz.Height, len(puz.Clues), puz.Title)
}

func (puz *Puzzle) SolutionAt(x, y int) byte {
	idx := (y * puz.Width) + x
	if idx > len(puz.solution) {
		return 0
	}
	return puz.solution[idx]
}

func (puz *Puzzle) InputAt(x, y int) byte {
	idx := (y * puz.Width) + x
	if idx > len(puz.input) {
		return 0
	}
	return puz.input[idx]
}

func (puz *Puzzle) CellAt(x, y int) *Cell {
	if x < 0 || y < 0 || x >= puz.Width || y >= puz.Height {
		return nil
	}
	return puz.Grid[y*puz.Width+x]
}

func (puz *Puzzle) Save() {
	puz.builder.Write()
}
