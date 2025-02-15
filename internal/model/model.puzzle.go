package model

// import (
// 	"regexp"
// 	"strings"

// 	tea "github.com/charmbracelet/bubbletea"
// 	puz "github.com/robertcurry0216/cross/internal/puzzle"
// )

// func PuzzleUpdate(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "up":
// 			SelectNextCell(m, -1, 0)
// 			return m, nil
// 		case "down":
// 			SelectNextCell(m, 1, 0)
// 			return m, nil
// 		case "left":
// 			SelectNextCell(m, 0, -1)
// 			return m, nil
// 		case "right":
// 			SelectNextCell(m, 0, 1)
// 			return m, nil
// 		case " ":
// 			m.state.PuzzleView.IsVert = !m.state.PuzzleView.IsVert
// 		default:
// 			pattern := `^[a-zA-Z]$`
// 			re := regexp.MustCompile(pattern)
// 			if re.MatchString(msg.String()) {
// 				SetLetter(m, msg.String())
// 				if m.state.PuzzleView.IsVert {
// 					SelectNextCell(m, 1, 0)
// 				} else {
// 					SelectNextCell(m, 0, 1)
// 				}
// 			}
// 		}
// 	}
// 	return m, nil
// }

// func SetPuzzle(m *Model, puzzle *puz.Puzzle) {
// 	m.state.Puzzle = puzzle
// 	SelectNextCell(m, 0, 1)
// }

// func SelectNextCell(m *Model, yDir, xDir int) {
// 	view := m.state.PuzzleView
// 	puz := m.state.Puzzle
// 	x, y := view.X, view.Y

// 	for {
// 		x += xDir
// 		y += yDir
// 		if next := puz.CellAt(x, y); next == nil {
// 			return
// 		} else if !next.IsBlank() {
// 			view.X = x
// 			view.Y = y
// 			next.Selected = true
// 			if view.IsVert && next.ClueVert != nil {
// 				next.ClueVert.Selected = true
// 			} else if next.ClueHoriz != nil {
// 				next.ClueHoriz.Selected = true
// 			}
// 			return
// 		}
// 	}
// }

// func SetLetter(m *Model, letter string) {
// 	cell := m.state.Puzzle.CellAt(m.state.PuzzleView.X, m.state.PuzzleView.Y)
// 	if cell == nil || len(letter) != 1 {
// 		return
// 	} else {
// 		cell.Input = strings.ToUpper(letter)[0]
// 	}
// }
