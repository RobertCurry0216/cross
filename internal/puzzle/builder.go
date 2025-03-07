package puzzle

import (
	"fmt"
	"os"
	"path/filepath"
)

type IBuilder interface {
	Build() (*Puzzle, error)
	Validate() error
	Write()
}

func NewBuilderFromFile(path string) (IBuilder, error) {
	ext := filepath.Ext(path)
	switch ext {
	case ".puz":
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		return NewPuzBuilder(raw, path), nil
	}

	return nil, fmt.Errorf("failed to parse file with ext: %v", ext)
}

func InitPuzzle(puz *Puzzle) error {
	size := puz.Width * puz.Height
	puz.Grid = make([]*Cell, size)

	for i := 0; i < size; i++ {
		puz.Grid[i] = NewCell()
		if puz.Solution[i] != '.' {
			puz.Grid[i].Solution = puz.Solution[i]
		}
		puz.Grid[i].Input = &puz.Input[i]
	}

	AssignClues(puz)

	return nil
}

func NeedsAcrossClue(puz *Puzzle, row, col int) bool {
	cell := puz.CellAt(col, row)
	if cell == nil || cell.IsBlank() {
		return false
	}
	left := puz.CellAt(col-1, row)
	right := puz.CellAt(col+1, row)

	if (left == nil || left.IsBlank()) && (right != nil && !right.IsBlank()) {
		return true
	}

	return false
}

func NeedsDownClue(puz *Puzzle, row, col int) bool {
	cell := puz.CellAt(col, row)
	if cell == nil || cell.IsBlank() {
		return false
	}
	above := puz.CellAt(col, row-1)
	below := puz.CellAt(col, row+1)

	if (above == nil || above.IsBlank()) && (below != nil && !below.IsBlank()) {
		return true
	}

	return false
}

func AssignClues(puz *Puzzle) error {
	puz.DownClues = make([]*Clue, 0, len(puz.Clues))
	puz.AcrossClues = make([]*Clue, 0, len(puz.Clues))
	clueIdx := 0
	clueNum := 1
	maxClues := len(puz.Clues)

	if maxClues == 0 {
		return nil
	}

	for row := range puz.Height {
		for col := range puz.Width {
			cell := puz.CellAt(col, row)
			needsAcross := NeedsAcrossClue(puz, row, col)
			needsDown := NeedsDownClue(puz, row, col)

			if needsAcross {
				clue := puz.Clues[clueIdx]
				clue.Number = clueNum
				puz.AcrossClues = append(puz.AcrossClues, clue)
				cell.ClueHoriz = clue
				clueIdx++

				// assign rest of cells to the clue
				clue.Cells = append(clue.Cells, cell)
				var next *Cell
				for i := 1; ; i++ {
					next = puz.CellAt(col+i, row)
					if next == nil || next.IsBlank() {
						break
					}
					next.ClueHoriz = clue
					clue.Cells = append(clue.Cells, next)
				}
			}

			if needsDown {
				clue := puz.Clues[clueIdx]
				clue.Number = clueNum
				puz.DownClues = append(puz.DownClues, clue)
				cell.ClueVert = clue
				clueIdx++

				// assign rest of cells to the clue
				clue.Cells = append(clue.Cells, cell)
				var next *Cell
				for i := 1; ; i++ {
					next = puz.CellAt(col, row+i)
					if next == nil || next.IsBlank() {
						break
					}
					next.ClueVert = clue
					clue.Cells = append(clue.Cells, next)
				}
			}

			if needsAcross || needsDown {
				clueNum++
			}

			if clueIdx > maxClues {
				return fmt.Errorf("Tried to assign too many clue indexes")
			}
		}
	}

	return nil
}
