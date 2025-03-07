package puzzle_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/robertcurry0216/cross/internal/puzzle"
)

func TestPuzzleWrite(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "puzzle_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy a test puzzle file to the temp directory
	testPuzPath := filepath.Join("testdata", "test.puz")
	tempPuzPath := filepath.Join(tempDir, "test.puz")

	// Read the original test file
	originalData, err := os.ReadFile(testPuzPath)
	if err != nil {
		t.Fatalf("Failed to read test puzzle file: %v", err)
	}

	// Write it to the temp location
	err = os.WriteFile(tempPuzPath, originalData, 0644)
	if err != nil {
		t.Fatalf("Failed to write temp puzzle file: %v", err)
	}

	// Load the puzzle from the temp file
	builder, err := puzzle.NewBuilderFromFile(tempPuzPath)
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}

	puz, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build puzzle: %v", err)
	}

	// Make some changes to the puzzle
	// Find a non-blank cell to modify
	var cellIndex int
	for i, cell := range puz.Grid {
		if !cell.IsBlank() {
			cellIndex = i
			break
		}
	}

	// Change the input for that cell
	originalInput := puz.Input[cellIndex]
	puz.Input[cellIndex] = 'X' // Change to 'X' regardless of what it was

	// Save the puzzle
	puz.Save()

	// Reload the puzzle to verify changes were saved
	builder2, err := puzzle.NewBuilderFromFile(tempPuzPath)
	if err != nil {
		t.Fatalf("Failed to create second builder: %v", err)
	}

	puz2, err := builder2.Build()
	if err != nil {
		t.Fatalf("Failed to build second puzzle: %v", err)
	}

	// Verify the change was saved
	if puz2.Input[cellIndex] != 'X' {
		t.Errorf("Expected input at index %d to be 'X', got '%c'", cellIndex, puz2.Input[cellIndex])
	}

	// Verify the file was actually written
	newData, err := os.ReadFile(tempPuzPath)
	if err != nil {
		t.Fatalf("Failed to read modified puzzle file: %v", err)
	}

	if bytes.Equal(originalData, newData) {
		t.Error("File data was not modified")
	}

	// Restore the original input for cleanup
	puz2.Input[cellIndex] = originalInput
	puz2.Save()
}

// TestPuzzleWriteValidation tests that the puzzle file remains valid after writing
func TestPuzzleWriteValidation(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "puzzle_validation_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy a test puzzle file to the temp directory
	testPuzPath := filepath.Join("testdata", "test.puz")
	tempPuzPath := filepath.Join(tempDir, "test.puz")

	// Read the original test file
	originalData, err := os.ReadFile(testPuzPath)
	if err != nil {
		t.Fatalf("Failed to read test puzzle file: %v", err)
	}

	// Write it to the temp location
	err = os.WriteFile(tempPuzPath, originalData, 0644)
	if err != nil {
		t.Fatalf("Failed to write temp puzzle file: %v", err)
	}

	// Load the puzzle from the temp file
	builder, err := puzzle.NewBuilderFromFile(tempPuzPath)
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}

	puz, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build puzzle: %v", err)
	}

	// Make changes to multiple cells
	for i, cell := range puz.Grid {
		if !cell.IsBlank() && i%5 == 0 { // Change every 5th non-blank cell
			puz.Input[i] = 'Z'
		}
	}

	// Save the puzzle
	puz.Save()

	// Validate the saved puzzle
	builder2, err := puzzle.NewBuilderFromFile(tempPuzPath)
	if err != nil {
		t.Fatalf("Failed to create second builder: %v", err)
	}

	// This should succeed if the file format is valid
	_, err = builder2.Build()
	if err != nil {
		t.Fatalf("Failed to build puzzle after saving: %v", err)
	}

	// Explicitly validate the puzzle
	err = builder2.Validate()
	if err != nil {
		t.Fatalf("Puzzle validation failed after saving: %v", err)
	}
}
