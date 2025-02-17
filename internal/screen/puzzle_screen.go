package screen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/robertcurry0216/cross/common"
	"github.com/robertcurry0216/cross/internal/puzzle"
)

const (
	boxString        = "┏┓┗┛━┃┣┫┳┻╋ *"
	blankCell        = ' '
	emptyCell        = '░'
	cellWidth    int = 3
	gridMinWidth int = 60
	clueMaxWidth int = 80
	clueMinWidth int = 50
)

var (
	colorHighlightBG,
	colorHighlightFG,
	colorError,
	colorCorrect,
	colorGridLine lipgloss.AdaptiveColor
)

const (
	topLeft = iota
	topRight
	bottomLeft
	bottomRight
	horizLine
	vertLine
	leftEdge
	rightEdge
	topEdge
	bottomEdge
	cross
	blank
	empty
)

type layoutType int

const (
	layoutCluesRight layoutType = iota
)

var boxRunes []string
var (
	styleTitle,
	styleCellText,
	styleBorder,
	styleGridLine,
	styleHighlightClue,
	styleHighlightCell,
	styleCellPadding lipgloss.Style
)

// init function to initialize boxRunes
func init() {
	boxRunes = make([]string, 0, len(boxString))
	for _, r := range boxString {
		boxRunes = append(boxRunes, string(r))
	}

	// colors
	colorError = lipgloss.AdaptiveColor{Light: "13", Dark: "13"}
	colorCorrect = lipgloss.AdaptiveColor{Light: "10", Dark: "10"}
	colorGridLine = lipgloss.AdaptiveColor{Light: "241", Dark: "7"}
	colorHighlightBG = lipgloss.AdaptiveColor{Light: "8", Dark: "250"}
	colorHighlightFG = lipgloss.AdaptiveColor{Light: "15", Dark: "0"}

	// styles
	styleBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 2)
	styleGridLine = lipgloss.NewStyle().Foreground(colorGridLine)
	styleTitle = lipgloss.NewStyle().Underline(true).Bold(true)
	styleCellText = lipgloss.NewStyle().Faint(true)
	styleHighlightClue = lipgloss.NewStyle().Faint(false).Bold(true)
	styleHighlightCell = lipgloss.NewStyle().Faint(false).Bold(true).Background(colorHighlightBG).Foreground(colorHighlightFG)
	styleCellPadding = lipgloss.NewStyle().Width(cellWidth).Align(lipgloss.Center)
}

type PuzzleScreen struct {
	puzzle *puzzle.Puzzle
}

func (s *PuzzleScreen) Init(state common.State) {
	s.puzzle = state.Puzzle
}

func (s *PuzzleScreen) View(state common.State) string {
	puzState := state.PuzzleView
	cell := s.puzzle.CellAt(puzState.X, puzState.Y)
	var clue *puzzle.Clue
	if cell != nil {
		if puzState.IsVert {
			clue = cell.ClueVert
		} else {
			clue = cell.ClueHoriz
		}
	}

	layout := calculateLayout(&state)

	// Apply the style to the text
	grid := renderPuzzle(layout.puzzle, s.puzzle, clue, cell)
	if lipgloss.Width(grid) < gridMinWidth {
		grid = lipgloss.PlaceHorizontal(gridMinWidth, lipgloss.Center, grid)
	}

	// combine
	hClues := renderClues(layout.acrossClues, s.puzzle.HorizClues, "Across", clue)
	vClues := renderClues(layout.downClues, s.puzzle.VertClues, "Down", clue)

	rightColumn := lipgloss.JoinVertical(lipgloss.Left, hClues, vClues)
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, grid)

	screen := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)

	// Print the styled box
	return screen
}

// Helpers

func titledBorder(title string) lipgloss.Style {
	border := lipgloss.RoundedBorder()
	border.Top = fmt.Sprintf("─%s%s", title, strings.Repeat("─", 500))
	return styleBorder.Border(border)
}

type puzzleViewLayout struct {
	layout      layoutType
	puzzle      common.LayoutBox
	downClues   common.LayoutBox
	acrossClues common.LayoutBox
	title       common.LayoutBox
	status      common.LayoutBox
}

func calculateLayout(state *common.State) puzzleViewLayout {
	layout := puzzleViewLayout{}

	layout.layout = layoutCluesRight

	// column widths
	leftColMin := state.Puzzle.Width * 4
	leftCol := int(float64(state.Width) * 0.6)
	if leftCol < leftColMin {
		leftCol = leftColMin
	}
	rightCol := state.Width - leftCol
	if rightCol < 20 {
		rightCol = 20
	}

	// heights
	statusHeight := 1
	titleHeight := 0
	acrossHeight := int(float64(state.Height-statusHeight) / 2)
	downHeight := state.Height - acrossHeight - statusHeight

	// boxes
	layout.puzzle = common.LayoutBox{
		W: leftCol,
		H: state.Height - statusHeight - titleHeight,
	}

	layout.title = common.LayoutBox{
		W: leftCol,
		H: titleHeight,
	}

	layout.status = common.LayoutBox{
		W: state.Width,
		H: statusHeight,
	}

	layout.acrossClues = common.LayoutBox{
		W: rightCol,
		H: acrossHeight,
	}

	layout.downClues = common.LayoutBox{
		W: rightCol,
		H: downHeight,
	}

	return layout
}

// Drawing functions

func renderPuzzle(box common.LayoutBox, puz *puzzle.Puzzle, selectedClue *puzzle.Clue, selectedCell *puzzle.Cell) string {
	buffer := NewBuffer(puz.Width*2+1, puz.Height*2+1)
	style := titledBorder("Puzzle").Height(box.H-2).Width(box.W-2).Align(lipgloss.Center, lipgloss.Center)

	insertCorners(puz, buffer)
	insertEdges(puz, buffer)
	insertCells(puz, buffer, selectedClue, selectedCell)

	var sb strings.Builder

	h, w := buffer.Size()
	for row := range h {
		if row > 0 {
			sb.WriteString("\n")
		}
		for col := range w {
			if cell, err := buffer.Get(row, col); err == nil {
				sb.WriteString(cell)
			}
		}
	}

	return style.Render(sb.String())
}

func insertCorners(puz *puzzle.Puzzle, buffer *Buffer) {
	var emptyBR, emptyTR, emptyBL, emptyTL bool
	for y := 0; y < puz.Height+1; y++ {
		rIdx := y * 2
		for x := 0; x < puz.Width+1; x++ {
			cIdx := x * 2
			emptyBR = puzzle.IsCellBlankOrNil(puz.CellAt(x, y))
			emptyTR = puzzle.IsCellBlankOrNil(puz.CellAt(x, y-1))
			emptyBL = puzzle.IsCellBlankOrNil(puz.CellAt(x-1, y))
			emptyTL = puzzle.IsCellBlankOrNil(puz.CellAt(x-1, y-1))

			var cell string

			switch {
			case (!emptyTL && !emptyBR) || (!emptyTR && !emptyBL):
				cell = boxRunes[cross]
			case emptyTL && emptyTR && !emptyBR && emptyBL:
				cell = boxRunes[topLeft]
			case emptyTL && emptyTR && emptyBR && !emptyBL:
				cell = boxRunes[topRight]
			case emptyTL && !emptyTR && emptyBR && emptyBL:
				cell = boxRunes[bottomLeft]
			case !emptyTL && emptyTR && emptyBR && emptyBL:
				cell = boxRunes[bottomRight]
			case !emptyTL && emptyTR && emptyBR && !emptyBL:
				cell = boxRunes[rightEdge]
			case emptyTL && !emptyTR && !emptyBR && emptyBL:
				cell = boxRunes[leftEdge]
			case !emptyTL && !emptyTR && emptyBR && emptyBL:
				cell = boxRunes[bottomEdge]
			case emptyTL && emptyTR && !emptyBR && !emptyBL:
				cell = boxRunes[topEdge]
			default:
				cell = boxRunes[blank]
			}

			buffer.Set(rIdx, cIdx, styleGridLine.Render(cell))

		}
	}
}

func insertEdges(puz *puzzle.Puzzle, buffer *Buffer) {
	var emptyC, emptyT, emptyL bool
	for y := 0; y < puz.Height+1; y++ {
		for x := 0; x < puz.Width+1; x++ {
			emptyC = puzzle.IsCellBlankOrNil(puz.CellAt(x, y))
			emptyT = puzzle.IsCellBlankOrNil(puz.CellAt(x, y-1))
			emptyL = puzzle.IsCellBlankOrNil(puz.CellAt(x-1, y))

			if !emptyC || !emptyL {
				buffer.Set(y*2+1, x*2, styleGridLine.Render(boxRunes[vertLine]))
			} else {
				buffer.Set(y*2+1, x*2, boxRunes[blank])
			}

			if !emptyC || !emptyT {
				buffer.Set(y*2, x*2+1, styleGridLine.Render(strings.Repeat(boxRunes[horizLine], cellWidth)))
			} else {
				buffer.Set(y*2, x*2+1, styleGridLine.Render(strings.Repeat(boxRunes[blank], cellWidth)))
			}
		}
	}
}

func insertCells(puz *puzzle.Puzzle, buffer *Buffer, selectedClue *puzzle.Clue, selectedCell *puzzle.Cell) {
	styleEmpty := styleCellText.Inherit(styleCellPadding)
	styleHighlight := styleHighlightClue.Inherit(styleCellPadding).Inherit(styleCellText)
	styleHighlightCell := styleHighlightCell.Inherit(styleCellPadding).Inherit(styleCellText)

	for y := 0; y < puz.Height; y++ {
		for x := 0; x < puz.Width; x++ {
			cell := puz.CellAt(x, y)
			if !puzzle.IsCellBlankOrNil(cell) {
				text := boxRunes[empty]
				if !cell.IsEmpty() {
					text = string(cell.Input)
				} else if n := cell.Number(); n > 0 {
					text = strconv.Itoa(n)
				}

				if selectedCell == cell {
					buffer.Set(y*2+1, x*2+1, styleHighlightCell.Render(text))
				} else if selectedClue == cell.ClueHoriz || selectedClue == cell.ClueVert {
					buffer.Set(y*2+1, x*2+1, styleHighlight.Render(text))
				} else {
					if cell.IsEmpty() {
						buffer.Set(y*2+1, x*2+1, styleEmpty.Render(text))
					} else {
						buffer.Set(y*2+1, x*2+1, styleCellPadding.Render(text))
					}
				}
			} else {
				buffer.Set(y*2+1, x*2+1, strings.Repeat(boxRunes[blank], cellWidth))
			}
		}
	}
}

func renderClues(box common.LayoutBox, clues []*puzzle.Clue, title string, selectedClue *puzzle.Clue) string {
	var out string

	for i, clue := range clues {
		num := fmt.Sprintf("%2d. ", clue.Number)
		clueText := common.WrapString(clue.Text, uint(box.W-4-lipgloss.Width(num)))
		clueText = fmt.Sprintf("%s (%d)", clueText, len(clue.Cells))
		if selectedClue == clue {
			clueText = styleHighlightClue.Render(clueText)
		} else {
			clueText = styleCellText.Render(clueText)
		}
		fullClue := lipgloss.JoinHorizontal(lipgloss.Top, num, clueText)
		if i == 0 {
			out = fullClue
		} else {
			out = lipgloss.JoinVertical(lipgloss.Left, out, fullClue)
		}
	}

	border := lipgloss.RoundedBorder()
	border.Top = fmt.Sprintf("─%s%s", title, strings.Repeat("─", 100))
	style := styleBorder.Border(border).Height(box.H - 2).Width(box.W - 2)

	out = style.Render(out)

	return out
}
