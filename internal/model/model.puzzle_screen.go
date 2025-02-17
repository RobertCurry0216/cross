package model

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertcurry0216/cross/internal/puzzle"
	puz "github.com/robertcurry0216/cross/internal/puzzle"
)

func PuzzleScreenUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.debug = msg.String()
		switch msg.String() {
		case "up":
			SelectNextCell(&m, -1, 0)
			return m, nil
		case "down":
			SelectNextCell(&m, 1, 0)
			return m, nil
		case "left":
			SelectNextCell(&m, 0, -1)
			return m, nil
		case "right":
			SelectNextCell(&m, 0, 1)
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
				cell.Input = ' '
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
	if cell, ok := GetSelectedCell(m); ok && len(letter) == 1 {
		cell.Input = strings.ToUpper(letter)[0]
	}
}

func GetSelectedCell(m *Model) (*puzzle.Cell, bool) {
	cell := m.state.Puzzle.CellAt(m.state.PuzzleView.X, m.state.PuzzleView.Y)
	return cell, cell != nil
}
