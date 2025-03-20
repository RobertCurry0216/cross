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
	cellNumberString     = "\u2080\u2081\u2082\u2083\u2084\u2085\u2086\u2087\u2088\u2089"
	boxString            = "┏┓┗┛━┃┣┫┳┻╋ .*"
	cellWidth        int = 3
	gridMinWidth     int = 60
	clueMaxWidth     int = 80
	clueMinWidth     int = 50
)

var (
	colorHighlightBG,
	colorHighlightFG,
	colorError,
	colorCorrect,
	colorStatusBar,
	colorFocusedBorder,
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
	emptySelected
)

var cellNumberRunes []string
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

	cellNumberRunes = make([]string, 0, len(cellNumberString))
	for _, r := range cellNumberString {
		cellNumberRunes = append(cellNumberRunes, string(r))
	}

	// colors
	colorError = lipgloss.AdaptiveColor{Light: "9", Dark: "1"}
	colorCorrect = lipgloss.AdaptiveColor{Light: "2", Dark: "10"}
	colorGridLine = lipgloss.AdaptiveColor{Light: "241", Dark: "7"}
	colorStatusBar = lipgloss.AdaptiveColor{Light: "4", Dark: "12"}
	colorHighlightBG = lipgloss.AdaptiveColor{Light: "8", Dark: "250"}
	colorHighlightFG = lipgloss.AdaptiveColor{Light: "15", Dark: "0"}
	colorFocusedBorder = lipgloss.AdaptiveColor{Light: "2", Dark: "10"}

	// styles
	styleBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	styleGridLine = lipgloss.NewStyle().Foreground(colorGridLine)
	styleTitle = lipgloss.NewStyle().Bold(true)
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

	return renderPuzzleView(layout, s.puzzle, clue)
}

// Helpers

func titledBorder(title string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.Top = fmt.Sprintf("─%s%s", title, strings.Repeat("─", 500))
	return border
}

type puzzleViewLayout struct {
	layout common.LayoutType
	puzzle common.LayoutBox
	clues  common.LayoutBox
	status common.LayoutBox
}

func calculateLayout(state *common.State) puzzleViewLayout {
	layout := puzzleViewLayout{}
	layout.layout = state.PuzzleView.Layout

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
	cluesHeight := state.Height - statusHeight

	// boxes
	layout.puzzle = common.LayoutBox{W: leftCol, H: state.Height - statusHeight}
	layout.status = common.LayoutBox{W: state.Width, H: statusHeight}
	layout.clues = common.LayoutBox{W: rightCol, H: cluesHeight}

	return layout
}

// Drawing functions
//
// |  __ \             | |       / ____|    (_)   | |
// | |__) |   _ ________ | ___  | |  __ _ __ _  __| |
// |  ___/ | | |_  /_  / |/ _ \ | | |_ | '__| |/ _` |
// | |   | |_| |/ / / /| |  __/ | |__| | |  | | (_| |
// |_|    \__,_/___/___|_|\___|  \_____|_|  |_|\__,_|

func renderPuzzleView(layout puzzleViewLayout, puzzle *puzzle.Puzzle, clue *puzzle.Clue) string {
	// render boxes
	grid := renderPuzzle(layout.puzzle, puzzle, clue, layout.layout == common.LayoutPuzzleFocus)
	if lipgloss.Width(grid) < gridMinWidth {
		grid = lipgloss.PlaceHorizontal(gridMinWidth, lipgloss.Center, grid)
	}

	// status bar
	status := renderStatusBar(layout.status)

	// combine
	rightColumn := renderClues(layout.clues, puzzle, clue, layout.layout == common.LayoutClueFocus)
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, grid)

	screen := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
	screen = lipgloss.JoinVertical(lipgloss.Center, screen, status)

	return screen
}

func renderPuzzle(box common.LayoutBox, puz *puzzle.Puzzle, selectedClue *puzzle.Clue, focus bool) string {
	buffer := NewBuffer(puz.Width*2+1, puz.Height*2+1)
	style := lipgloss.NewStyle().Border(titledBorder(puz.Title)).Height(box.H-2).Width(box.W-2).Align(lipgloss.Center, lipgloss.Center)

	if focus {
		style = style.BorderForeground(colorFocusedBorder)
	}

	insertCorners(puz, buffer)
	insertEdges(puz, buffer)
	insertCells(puz, buffer, selectedClue)

	return style.Render(buffer.String())
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
	var cell *puzzle.Cell
	var emptyC, emptyT, emptyL bool
	for y := 0; y < puz.Height+1; y++ {
		for x := 0; x < puz.Width+1; x++ {
			cell = puz.CellAt(x, y)
			emptyC = puzzle.IsCellBlankOrNil(cell)
			emptyT = puzzle.IsCellBlankOrNil(puz.CellAt(x, y-1))
			emptyL = puzzle.IsCellBlankOrNil(puz.CellAt(x-1, y))

			// vert lines
			if !emptyC || !emptyL {
				buffer.Set(y*2+1, x*2, styleGridLine.Render(boxRunes[vertLine]))
			} else {
				buffer.Set(y*2+1, x*2, boxRunes[blank])
			}

			// horiz lines
			if !emptyC && cell.Number() > 0 {
				runes := make([]string, cellWidth)
				for i := range cellWidth {
					runes[i] = boxRunes[horizLine]
				}
				for i, nr := range strconv.Itoa(cell.Number()) {
					runes[i] = string(nr)
				}
				buffer.Set(y*2, x*2+1, styleGridLine.Render(strings.Join(runes, "")))
			} else if !emptyC || !emptyT {
				buffer.Set(y*2, x*2+1, styleGridLine.Render(strings.Repeat(boxRunes[horizLine], cellWidth)))
			} else {
				buffer.Set(y*2, x*2+1, styleGridLine.Render(strings.Repeat(boxRunes[blank], cellWidth)))
			}
		}
	}
}

func insertCells(puz *puzzle.Puzzle, buffer *Buffer, selectedClue *puzzle.Clue) {
	for y := 0; y < puz.Height; y++ {
		for x := 0; x < puz.Width; x++ {
			cell := puz.CellAt(x, y)
			if !puzzle.IsCellBlankOrNil(cell) {
				text := boxRunes[empty]
				if !cell.IsEmpty() {
					text = string(*cell.Input)
				}

				style := lipgloss.NewStyle().Inherit(styleCellPadding)
				isSelected := cell.IsSelected
				isHighlighted := selectedClue == cell.ClueHoriz || selectedClue == cell.ClueVert

				if isSelected {
					style = style.Inherit(styleHighlightCell)
				} else if isHighlighted {
					style = style.Inherit(styleHighlightClue)
				}

				if cell.IsEmpty() && (isHighlighted || isSelected) {
					text = boxRunes[emptySelected]
				}

				if !cell.IsEmpty() && cell.ShowChecked {
					if cell.IsCorrect() {
						if isSelected {
							style = style.Background(colorCorrect)
						} else {
							style = style.Foreground(colorCorrect)
						}
					} else {
						if isSelected {
							style = style.Background(colorError)
						} else {
							style = style.Foreground(colorError)
						}
					}
				}

				if cell.IsCircled {
					text = fmt.Sprintf("(%s)", text)
				}

				buffer.Set(y*2+1, x*2+1, style.Render(text))
			} else {
				buffer.Set(y*2+1, x*2+1, strings.Repeat(boxRunes[blank], cellWidth))
			}
		}
	}
}

//   _____ _
//  / ____| |
// | |    | |_   _  ___ ___
// | |    | | | | |/ _ \ __|
// | |____| | |_| |  __\__ \
//  \_____|_|\__,_|\___|___/

func renderClues(box common.LayoutBox, puzzle *puzzle.Puzzle, selectedClue *puzzle.Clue, focus bool) string {
	acrossTitle := styleTitle.Render("Across:")
	acrossText, lnAcross := renderClueSet(box.W, puzzle.AcrossClues, selectedClue, focus)
	downTitle := styleTitle.Render("Down:")
	downText, lnDown := renderClueSet(box.W, puzzle.DownClues, selectedClue, focus)

	allClues := lipgloss.JoinVertical(lipgloss.Left, acrossTitle, acrossText, downTitle, downText)

	// ensure selected clue is visible
	lnSelected := lnAcross + 1
	if lnAcross < 0 {
		lnSelected = lipgloss.Height(acrossText) + 2 + lnDown
	}

	allClues = centerLine(allClues, box.H-2, lnSelected)

	style := styleBorder
	if focus {
		style = style.BorderForeground(colorFocusedBorder)
	}
	boxedClues := style.Border(titledBorder("Clues")).Height(box.H - 2).Width(box.W - 2).Render(allClues)

	return boxedClues
}

func renderClueSet(W int, clues []*puzzle.Clue, selectedClue *puzzle.Clue, focus bool) (string, int) {
	var out string
	var lineNum = -1

	for i, clue := range clues {
		num := fmt.Sprintf("%2d. ", clue.Number)
		clueText := common.WrapString(clue.Text, uint(W-4-lipgloss.Width(num)))
		clueText = fmt.Sprintf("%s (%d)", clueText, len(clue.Cells))
		if selectedClue == clue {
			clueText = styleHighlightClue.Render(clueText)
			if focus {
				focusText := renderFocusedClue(clue)
				clueText = lipgloss.JoinVertical(lipgloss.Left, clueText, focusText)
			}
			lineNum = lipgloss.Height(out) + lipgloss.Height(clueText)
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

	return out, lineNum
}

func renderFocusedClue(clue *puzzle.Clue) string {
	if clue == nil || len(clue.Cells) == 0 {
		return ""
	}

	w := len(clue.Cells)*2 + 1
	buffer := NewBuffer(w, 3)

	for i, cell := range clue.Cells {
		buffer.Set(0, i*2, boxRunes[topEdge])
		buffer.Set(1, i*2, boxRunes[vertLine])
		buffer.Set(2, i*2, boxRunes[bottomEdge])

		buffer.Set(0, i*2+1, strings.Repeat(boxRunes[horizLine], cellWidth))
		cellText := boxRunes[emptySelected]
		if !cell.IsEmpty() {
			cellText = string(*cell.Input)
		}
		cellText = styleCellPadding.Render(cellText)
		if cell.IsSelected {
			cellText = styleHighlightCell.Render(cellText)
		}
		buffer.Set(1, i*2+1, cellText)
		buffer.Set(2, i*2+1, strings.Repeat(boxRunes[horizLine], cellWidth))
	}

	buffer.Set(1, w-1, boxRunes[vertLine])
	buffer.Set(0, w-1, boxRunes[topRight])
	buffer.Set(2, w-1, boxRunes[bottomRight])
	buffer.Set(0, 0, boxRunes[topLeft])
	buffer.Set(2, 0, boxRunes[bottomLeft])

	return buffer.String()
}

func centerLine(str string, maxH, n int) string {
	lines := strings.Split(str, "\n")
	lineCount := len(lines)
	if lineCount < maxH || maxH < 0 || n < 0 {
		return str
	}
	startIdx := n - (maxH / 2)
	if startIdx < 0 {
		startIdx = 0
	}

	endIdx := startIdx + maxH
	if endIdx > lineCount {
		endIdx = lineCount
		startIdx = lineCount - maxH
	}

	return strings.Join(lines[startIdx:endIdx], "\n")
}

//   _____ _        _               ____
//  / ____| |      | |             |  _ \
// | (___ | |_ __ _| |_ _   _ ___  | |_) | __ _ _ __
//  \___ \| __/ _` | __| | | / __| |  _ < / _` | '__|
//  ____) | |_ (_| | |_| |_| \__ \ | |_) | (_| | |
// |_____/ \__\__,_|\__|\__,_|___/ |____/ \__,_|_|
//

func renderStatusBar(box common.LayoutBox) string {
	shortcuts := "esc: Exit | ctrl+l: Check letter | crtl+w: Check word | ctrl+a: Check puzzle | ctrl+r: Reveal word | ctrl+p: Reveal puzzle"
	version := "Cross-cli version 0.1"
	shortcuts = lipgloss.NewStyle().Foreground(colorStatusBar).Render(shortcuts)

	scLen := lipgloss.Width(shortcuts)
	vLen := lipgloss.Width(version)
	spacer := box.W - scLen - vLen
	if spacer < 0 {
		spacer = 0
	}

	return fmt.Sprintf("%s%s%s", shortcuts, strings.Repeat(" ", spacer), version)
}
