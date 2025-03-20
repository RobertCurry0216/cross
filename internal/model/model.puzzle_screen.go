package model

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertcurry0216/cross/common"
	"github.com/robertcurry0216/cross/internal/puzzle"
	puz "github.com/robertcurry0216/cross/internal/puzzle"
)

func PuzzleScreenUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.state.Debug = msg.String()
		switch msg.String() {
		case "up":
			if m.state.PuzzleView.Layout == common.LayoutPuzzleFocus {
				SelectNextCell(&m, -1, 0)
			} else {
				SelectNextClue(&m, false)
			}
			return m, nil
		case "down":
			if m.state.PuzzleView.Layout == common.LayoutPuzzleFocus {
				SelectNextCell(&m, 1, 0)
			} else {
				SelectNextClue(&m, true)
			}
			return m, nil
		case "left":
			if m.state.PuzzleView.Layout == common.LayoutPuzzleFocus || !m.state.PuzzleView.IsVert {
				SelectNextCell(&m, 0, -1)
			} else {
				SelectNextCell(&m, -1, 0)
			}
			return m, nil
		case "right":
			if m.state.PuzzleView.Layout == common.LayoutPuzzleFocus || !m.state.PuzzleView.IsVert {
				SelectNextCell(&m, 0, 1)
			} else {
				SelectNextCell(&m, 1, 0)
			}
			return m, nil
		case " ":
			m.state.PuzzleView.IsVert = !m.state.PuzzleView.IsVert
			return m, nil
		case "backspace":
			if cell, ok := GetSelectedCell(&m); ok {
				if cell.IsEmpty() {
					if m.state.PuzzleView.IsVert {
						SelectNextCell(&m, -1, 0)
					} else {
						SelectNextCell(&m, 0, -1)
					}
					if cell, ok = GetSelectedCell(&m); !ok {
						return m, nil
					}
				}
				*cell.Input = '-'
			}
		case "ctrl+l":
			// check letter
			if cell, ok := GetSelectedCell(&m); ok {
				cell.ShowChecked = true
			}
		case "ctrl+w":
			// check word
			if cell, ok := GetSelectedCell(&m); ok {
				var clue *puzzle.Clue
				if m.state.PuzzleView.IsVert {
					clue = cell.ClueVert
				} else {
					clue = cell.ClueHoriz
				}

				for _, c := range clue.Cells {
					c.ShowChecked = true
				}
			}
		case "ctrl+a":
			// check puzzle
			for _, cell := range m.state.Puzzle.Grid {
				if !cell.IsBlank() {
					cell.ShowChecked = true
				}
			}
		case "ctrl+r":
			// reveal word
			if cell, ok := GetSelectedCell(&m); ok {
				var clue *puzzle.Clue
				if m.state.PuzzleView.IsVert {
					clue = cell.ClueVert
				} else {
					clue = cell.ClueHoriz
				}

				for _, c := range clue.Cells {
					*c.Input = c.Solution
				}
			}
		case "tab":
			// toggle focus
			if m.state.PuzzleView.Layout == common.LayoutPuzzleFocus {
				m.state.PuzzleView.Layout = common.LayoutClueFocus
			} else {
				m.state.PuzzleView.Layout = common.LayoutPuzzleFocus
			}
		default:
			pattern := `^[a-zA-Z]$`
			re := regexp.MustCompile(pattern)
			if re.MatchString(msg.String()) {
				SetLetter(&m, msg.String())
				if m.state.PuzzleView.IsVert {
					SelectNextCell(&m, 1, 0)
				} else {
					SelectNextCell(&m, 0, 1)
				}
			}
			m.state.Puzzle.Save()
			return m, nil
		}
	}
	return m, nil
}

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
	if cur := puz.CellAt(x, y); cur != nil {
		cur.IsSelected = false
	}

	for {
		x += xDir
		y += yDir
		if next := puz.CellAt(x, y); next == nil {
			break
		} else if !next.IsBlank() {
			view.X = x
			view.Y = y
			if view.IsVert && next.ClueVert != nil {
				next.ClueVert.Selected = true
			} else if next.ClueHoriz != nil {
				next.ClueHoriz.Selected = true
			}
			break
		}
	}

	x, y = view.X, view.Y
	if cur := puz.CellAt(x, y); cur != nil {
		cur.IsSelected = true
	}
}

func SetLetter(m *Model, letter string) {
	if cell, ok := GetSelectedCell(m); ok && len(letter) == 1 {
		*cell.Input = strings.ToUpper(letter)[0]
		cell.ShowChecked = false
	}
}

func GetSelectedCell(m *Model) (*puzzle.Cell, bool) {
	cell := m.state.Puzzle.CellAt(m.state.PuzzleView.X, m.state.PuzzleView.Y)
	return cell, cell != nil
}

func SelectNextClue(m *Model, forward bool) {
	view := &m.state.PuzzleView
	puz := m.state.Puzzle

	// Get the current clue
	cell, ok := GetSelectedCell(m)
	if !ok {
		return
	}
	cell.IsSelected = false
	var currentClue *puzzle.Clue
	var clues []*puzzle.Clue

	if view.IsVert {
		currentClue = cell.ClueVert
		clues = puz.DownClues
	} else {
		currentClue = cell.ClueHoriz
		clues = puz.AcrossClues
	}

	if currentClue == nil || len(clues) == 0 {
		return
	}

	// Find the index of the current clue
	currentIndex := -1
	for i, clue := range clues {
		if clue == currentClue {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return
	}

	// Calculate the next index
	nextIndex := currentIndex
	if forward {
		nextIndex = (currentIndex + 1) % len(clues)
	} else {
		nextIndex = (currentIndex - 1 + len(clues)) % len(clues)
	}

	// Select the first cell of the next clue
	nextClue := clues[nextIndex]
	if firstCell := nextClue.FirstCell(); firstCell != nil {
		// Update the position
		for y := 0; y < puz.Height; y++ {
			for x := 0; x < puz.Width; x++ {
				if puz.CellAt(x, y) == firstCell {
					view.X = x
					view.Y = y
					firstCell.IsSelected = true
					nextClue.Selected = true
					return
				}
			}
		}
	}
}
