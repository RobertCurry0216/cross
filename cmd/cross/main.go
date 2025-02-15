package main

import (
	"flag"
	"fmt"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertcurry0216/cross/internal/model"
	puz "github.com/robertcurry0216/cross/internal/puzzle"
	"github.com/robertcurry0216/cross/internal/screen"
)

var path string

func init() {
	const (
		pathDefault = ""
		usage       = "the file path to the crossword"
	)

	flag.StringVar(&path, "filepath", pathDefault, usage)
	flag.StringVar(&path, "f", pathDefault, usage+" (short)")
}

func buildAndValidatePuzzle() (*puz.Puzzle, error) {
	builder, err := puz.NewBuilderFromFile(path)
	if err != nil {
		return &puz.Puzzle{}, err
	}

	p, err := builder.Build()
	if err != nil {
		return &puz.Puzzle{}, err
	}

	if err := builder.Validate(); err != nil {
		return &puz.Puzzle{}, err
	}
	return p, nil
}

func main() {
	// for i := 0; i < 16; i++ {
	// 	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color(strconv.Itoa(i))).Render("Hi: " + strconv.Itoa(i)))
	// }

	// return

	flag.Parse()
	p, _ := buildAndValidatePuzzle()

	m := model.NewModel()

	model.SetPuzzle(&m, p)
	m.PushView(&screen.PuzzleScreen{})

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
