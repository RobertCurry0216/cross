package model

import (
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
)

func PuzzleScreenUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return m, nil
}
