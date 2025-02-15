package model

import (
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertcurry0216/cross/common"
	"github.com/robertcurry0216/cross/internal/screen"
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

func (m Model) View() string {
	if len(m.state.Views) > 0 {
		view := m.state.Views[len(m.state.Views)-1]
		return view.View(m.state)
	} else {
		return ""
	}
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

	if len(m.state.Views) == 0 {
		return m, nil
	}

	// view specific actions
	switch m.state.Views[len(m.state.Views)-1].(type) {
	case *screen.PuzzleScreen:
		switch msg := msg.(type) {
		case tea.KeyMsg:
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
	}

	// catch all return
	return m, nil
}
