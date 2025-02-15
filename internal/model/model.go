package model

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertcurry0216/cross/common"
	puz "github.com/robertcurry0216/cross/internal/puzzle"
)

type Model struct {
	state common.State
}

// Constructor

func NewModel() Model {
	views := make([]common.Viewable, 0, 10)
	return Model{state: common.State{Views: views}}
}

// bubble tea functions

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// global actions
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.PopView()
			if len(m.state.Views) == 0 {
				return m, tea.Quit
			} else {
				return m, nil
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.state.Width = msg.Width
		m.state.Height = msg.Height
		return m, nil
	}

	// puzzle view actions

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.SelectNextCell(-1, 0)
			return m, nil
		case "down":
			m.SelectNextCell(1, 0)
			return m, nil
		case "left":
			m.SelectNextCell(0, -1)
			return m, nil
		case "right":
			m.SelectNextCell(0, 1)
			return m, nil
		case " ":
			m.state.PuzzleView.IsVert = !m.state.PuzzleView.IsVert
		default:
			pattern := `^[a-zA-Z]$`
			re := regexp.MustCompile(pattern)
			if re.MatchString(msg.String()) {
				m.SetLetter(msg.String())
				if m.state.PuzzleView.IsVert {
					m.SelectNextCell(1, 0)
				} else {
					m.SelectNextCell(0, 1)
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.state.Views) > 0 {
		view := m.state.Views[len(m.state.Views)-1]
		return view.View(m.state)
	} else {
		return ""
	}
}

// model functions

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

func (m *Model) SetPuzzle(puzzle *puz.Puzzle) {
	m.state.Puzzle = puzzle
	m.SelectNextCell(0, 1)
}

func (m *Model) SelectNextCell(yDir, xDir int) {
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

func (m *Model) SetLetter(letter string) {
	cell := m.state.Puzzle.CellAt(m.state.PuzzleView.X, m.state.PuzzleView.Y)
	if cell == nil || len(letter) != 1 {
		return
	} else {
		cell.Input = strings.ToUpper(letter)[0]
	}
}
