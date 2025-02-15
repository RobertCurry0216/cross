package model

import (
	"strings"

	"github.com/robertcurry0216/cross/common"
	puz "github.com/robertcurry0216/cross/internal/puzzle"
)

// model methods
func (m *Model) PushView(view common.Viewable) {
	view.Init(m.state)
	m.state.Views = append(m.state.Views, view)
}

func (m *Model) PopView() {
	if len(m.state.Views) > 0 {
		m.state.Views = m.state.Views[:len(m.state.Views)-1]
	}
}

// Puzzle view
func SetPuzzle(m *Model, puzzle *puz.Puzzle) {
	m.state.Puzzle = puzzle
	// lazy way to ensure the initial cell isn't blank
	SelectNextCell(m, 0, 1)
	SelectNextCell(m, 0, -1)
}

func SelectNextCell(m *Model, yDir, xDir int) {
	view := &m.state.PuzzleView
	puz := m.state.Puzzle
	x, y := view.X, view.Y

	for {
		x += xDir
		y += yDir
		if next := puz.CellAt(x, y); next == nil {
			return
		} else if !next.IsBlank() {
			view.X = x
			view.Y = y
			next.Selected = true
			if view.IsVert && next.ClueVert != nil {
				next.ClueVert.Selected = true
			} else if next.ClueHoriz != nil {
				next.ClueHoriz.Selected = true
			}
			return
		}
	}
}

func SetLetter(m *Model, letter string) {
	cell := m.state.Puzzle.CellAt(m.state.PuzzleView.X, m.state.PuzzleView.Y)
	if cell == nil || len(letter) != 1 {
		return
	} else {
		cell.Input = strings.ToUpper(letter)[0]
	}
}
