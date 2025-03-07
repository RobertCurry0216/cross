package puzzle_test

import (
	"os"
	"testing"

	"github.com/robertcurry0216/cross/internal/puzzle"
)

func TestNewBuilderFromFile(t *testing.T) {
	// Test with non-existent file
	_, err := puzzle.NewBuilderFromFile("non_existent_file.puz")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test with unsupported extension
	tempFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	_, err = puzzle.NewBuilderFromFile(tempFile.Name())
	if err == nil {
		t.Error("Expected error for unsupported file extension, got nil")
	}
}

func TestInitPuzzle(t *testing.T) {
	// Create a simple test puzzle
	puz := puzzle.NewPuzzle()
	puz.Width = 3
	puz.Height = 3
	puz.Solution = []byte("ABC.DEFGH")
	puz.Input = make([]byte, 9)

	err := puzzle.InitPuzzle(puz)
	if err != nil {
		t.Fatalf("InitPuzzle failed: %v", err)
	}

	// Verify grid was created correctly
	if len(puz.Grid) != 9 {
		t.Errorf("Expected grid size 9, got %d", len(puz.Grid))
	}

	// Check that cells were initialized properly
	for i, cell := range puz.Grid {
		if i == 3 { // This is the '.' in the solution
			if !cell.IsBlank() {
				t.Errorf("Expected blank cell at index 3, got non-blank")
			}
		} else {
			if cell.Solution != puz.Solution[i] {
				t.Errorf("Expected solution %c at index %d, got %c", puz.Solution[i], i, cell.Solution)
			}
			if cell.Input != &puz.Input[i] {
				t.Errorf("Input pointer not set correctly at index %d", i)
			}
		}
	}
}

func TestNeedsAcrossClue(t *testing.T) {
	// Create a test puzzle with a simple pattern
	// . A B
	// C D E
	// F . G
	puz := puzzle.NewPuzzle()
	puz.Width = 3
	puz.Height = 3
	puz.Solution = []byte(".ABCDEF.G")
	puz.Input = make([]byte, 9)
	puzzle.InitPuzzle(puz)

	testCases := []struct {
		row, col int
		expected bool
	}{
		{0, 0, false}, // blank cell
		{0, 1, true},  // start of across clue
		{0, 2, false}, // middle of across clue
		{1, 0, true},  // start of across clue
		{1, 1, false}, // middle of across clue
		{1, 2, false}, // end of across clue
		{2, 0, false}, // single cell across clue
		{2, 1, false}, // blank cell
		{2, 2, false}, // single cell across clue
	}

	for _, tc := range testCases {
		result := puzzle.NeedsAcrossClue(puz, tc.row, tc.col)
		if result != tc.expected {
			t.Errorf("needsAcrossClue(%d, %d) = %v, expected %v", tc.row, tc.col, result, tc.expected)
		}
	}
}

func TestNeedsDownClue(t *testing.T) {
	// Create a test puzzle with a simple pattern
	// . A B
	// C D E
	// F . G
	puz := puzzle.NewPuzzle()
	puz.Width = 3
	puz.Height = 3
	puz.Solution = []byte(".ABCDEF.G")
	puz.Input = make([]byte, 9)
	puzzle.InitPuzzle(puz)

	testCases := []struct {
		row, col int
		expected bool
	}{
		{0, 0, false}, // blank cell
		{0, 1, true},  // start of down clue
		{0, 2, true},  // start of down clue
		{1, 0, true},  // start of down clue
		{1, 1, false}, // middle of down clue
		{1, 2, false}, // middle of down clue
		{2, 0, false}, // end of down clue
		{2, 1, false}, // blank cell
		{2, 2, false}, // end of down clue
	}

	for _, tc := range testCases {
		result := puzzle.NeedsDownClue(puz, tc.row, tc.col)
		if result != tc.expected {
			t.Errorf("needsDownClue(%d, %d) = %v, expected %v", tc.row, tc.col, result, tc.expected)
		}
	}
}

func TestAssignClues(t *testing.T) {
	// Create a test puzzle with a simple pattern
	// . A B
	// C D E
	// F . G
	puz := puzzle.NewPuzzle()
	puz.Width = 3
	puz.Height = 3
	puz.Solution = []byte(".ABCDEF.G")
	puz.Input = make([]byte, 9)
	puzzle.InitPuzzle(puz)

	// Create some test clues
	puz.Clues = []*puzzle.Clue{
		puzzle.NewClue("Across 1"),
		puzzle.NewClue("Across 4"),
		puzzle.NewClue("Across 6"),
		puzzle.NewClue("Across 8"),
		puzzle.NewClue("Down 1"),
		puzzle.NewClue("Down 2"),
		puzzle.NewClue("Down 3"),
	}

	err := puzzle.AssignClues(puz)
	if err != nil {
		t.Fatalf("assignClues failed: %v", err)
	}

	// Check that we have the right number of across and down clues
	if len(puz.AcrossClues) != 2 {
		t.Errorf("Expected 2 across clues, got %d", len(puz.AcrossClues))
	}
	if len(puz.DownClues) != 3 {
		t.Errorf("Expected 3 down clues, got %d", len(puz.DownClues))
	}

	// Check that clue numbers are assigned correctly
	expectedAcrossNumbers := []int{1, 3}
	for i, clue := range puz.AcrossClues {
		if clue.Number != expectedAcrossNumbers[i] {
			t.Errorf("Expected across clue %d to have number %d, got %d",
				i, expectedAcrossNumbers[i], clue.Number)
		}
	}

	expectedDownNumbers := []int{1, 2, 3}
	for i, clue := range puz.DownClues {
		if clue.Number != expectedDownNumbers[i] {
			t.Errorf("Expected down clue %d to have number %d, got %d",
				i, expectedDownNumbers[i], clue.Number)
		}
	}

	// Check that cells have the correct clue references
	// Cell (0,1) should have horizontal clue #1
	if puz.CellAt(1, 0).ClueHoriz == nil || puz.CellAt(1, 0).ClueHoriz.Number != 1 {
		t.Errorf("Cell (1,0) should have horizontal clue #1")
	}

	// Cell (1,0) should have horizontal clue #3
	if puz.CellAt(0, 1).ClueHoriz == nil || puz.CellAt(0, 1).ClueHoriz.Number != 3 {
		t.Errorf("Cell (0,1) should have horizontal clue #3")
	}
}

func TestPuzBuilderChecksumRegion(t *testing.T) {
	testData := []byte("ABCDEFG")
	result := puzzle.ChecksumRegion(testData, 0)

	// The checksum algorithm is deterministic, so we can test for a specific value
	// This expected value would need to be calculated manually or verified once
	if result == 0 {
		t.Error("checksumRegion returned 0, which is likely incorrect")
	}

	// Test that the same input always produces the same output
	result2 := puzzle.ChecksumRegion(testData, 0)
	if result != result2 {
		t.Errorf("checksumRegion not deterministic: %d != %d", result, result2)
	}

	// Test that different inputs produce different outputs
	testData2 := []byte("ABCDEFGH")
	result3 := puzzle.ChecksumRegion(testData2, 0)
	if result == result3 {
		t.Error("checksumRegion should produce different results for different inputs")
	}
}
