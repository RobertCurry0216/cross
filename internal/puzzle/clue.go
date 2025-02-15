package puzzle

type Clue struct {
	Text     string
	Number   int
	Selected bool
	Cells    []*Cell
}

func NewClue(text string) *Clue {
	cells := make([]*Cell, 0, 15)
	return &Clue{Text: text, Cells: cells}
}

func (clue *Clue) FirstCell() *Cell {
	if len(clue.Cells) > 0 {
		return clue.Cells[0]
	}
	return nil
}
