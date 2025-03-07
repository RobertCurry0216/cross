package puzzle

type Cell struct {
	ClueVert    *Clue
	ClueHoriz   *Clue
	Solution    byte
	Input       *byte
	Selected    bool
	ShowChecked bool
	IsCircled   bool
}

func NewCell() *Cell {
	return &Cell{}
}

func (cell *Cell) IsBlank() bool {
	return cell.Solution == 0
}

func (cell *Cell) IsEmpty() bool {
	return *cell.Input < 'A' || *cell.Input > 'Z'
}

func (cell *Cell) Number() int {
	if cell.ClueHoriz != nil && cell.ClueHoriz.FirstCell() == cell {
		return cell.ClueHoriz.Number
	} else if cell.ClueVert != nil && cell.ClueVert.FirstCell() == cell {
		return cell.ClueVert.Number
	}
	return -1
}

func (cell *Cell) IsCorrect() bool {
	return *cell.Input == cell.Solution
}

// helpers

func IsCellBlankOrNil(cell *Cell) bool {
	return cell == nil || cell.IsBlank()
}
