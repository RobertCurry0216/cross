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


func buildAndValidatePuzzle(path string) (*puz.Puzzle, error) {
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
	flag.Parse()
	
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Please provide a crossword file path")
		fmt.Println("Usage: cross [crossword_file.puz]")
		os.Exit(1)
	}
	
	path := args[0]
	p, err := buildAndValidatePuzzle(path)
	if err != nil {
		fmt.Printf("Error loading puzzle: %v\n", err)
		os.Exit(1)
	}

	// return
	m := model.NewModel()

	model.SetPuzzle(&m, p)
	m.PushView(&screen.PuzzleScreen{})

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
