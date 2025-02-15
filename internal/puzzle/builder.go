package puzzle

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type Builder struct {
	raw []byte
}

type iBuilder interface {
	Build() (*Puzzle, error)
	Validate() error
}

func NewBuilderFromFile(path string) (iBuilder, error) {
	ext := filepath.Ext(path)
	switch ext {
	case ".puz":
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		return &PuzBuilder{raw: raw}, nil
	}

	return nil, fmt.Errorf("failed to parse file with ext: %v", ext)
}

func readString(data []byte) (string, int) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return "", -1
	}
	str := string(data[:nullIndex])
	offset := nullIndex + 1 // Move past the null terminator

	return str, offset
}

func InitPuzzle(puz *Puzzle) error {
	puz.Grid = make([][]Cell, puz.Height)

	for y := range int(puz.Height) {
		puz.Grid[y] = make([]Cell, puz.Width)
		for x := range int(puz.Width) {
			puz.Grid[y][x] = NewCell()
			cell := puz.CellAt(x, y)

			cell.Solution = puz.SolutionAt(x, y)
			cell.Input = puz.InputAt(x, y)
		}
	}

	assignClues(puz)

	return nil
}

func needsAcrossClue(puz *Puzzle, row, col int) bool {
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

func needsDownClue(puz *Puzzle, row, col int) bool {
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

func assignClues(puz *Puzzle) error {
	puz.VertClues = make([]*Clue, 0, len(puz.Clues))
	puz.HorizClues = make([]*Clue, 0, len(puz.Clues))
	clueIdx := 0
	clueNum := 1
	maxClues := len(puz.Clues)

	for row := range puz.Height {
		for col := range puz.Width {
			cell := puz.CellAt(col, row)
			needsAcross := needsAcrossClue(puz, row, col)
			needsDown := needsDownClue(puz, row, col)

			if needsAcross {
				clue := puz.Clues[clueIdx]
				clue.Number = clueNum
				puz.HorizClues = append(puz.HorizClues, clue)
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
				puz.VertClues = append(puz.VertClues, clue)
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
