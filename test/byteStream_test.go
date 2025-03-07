package puzzle

import (
	"bytes"
	"testing"

	"github.com/robertcurry0216/cross/internal/puzzle"
)

func TestNewByteStream(t *testing.T) {
	testData := []byte{1, 2, 3, 4, 5}
	stream := puzzle.NewByteStream(testData)

	if stream.Size != len(testData) {
		t.Errorf("Expected size %d, got %d", len(testData), stream.Size)
	}

	if !bytes.Equal(stream.Raw, testData) {
		t.Errorf("Raw data not stored correctly")
	}

	if stream.Pointer != 0 {
		t.Errorf("Expected initial pointer 0, got %d", stream.Pointer)
	}
}

func TestByteStreamIncPointer(t *testing.T) {
	stream := puzzle.NewByteStream([]byte{1, 2, 3, 4, 5})

	stream.IncPointer(3)
	if stream.Pointer != 3 {
		t.Errorf("Expected pointer 3, got %d", stream.Pointer)
	}

	stream.IncPointer(1)
	if stream.Pointer != 4 {
		t.Errorf("Expected pointer 4, got %d", stream.Pointer)
	}
}

func TestByteStreamChompN(t *testing.T) {
	testData := []byte{1, 2, 3, 4, 5}
	stream := puzzle.NewByteStream(testData)

	// Test normal case
	data, n := stream.ChompN(3)
	if n != 3 {
		t.Errorf("Expected to read 3 bytes, got %d", n)
	}
	if !bytes.Equal(data, []byte{1, 2, 3}) {
		t.Errorf("Expected [1,2,3], got %v", data)
	}
	if stream.Pointer != 3 {
		t.Errorf("Expected pointer 3, got %d", stream.Pointer)
	}

	// Test reading past end
	data, n = stream.ChompN(3)
	if n != 2 {
		t.Errorf("Expected to read 2 bytes, got %d", n)
	}
	if !bytes.Equal(data, []byte{4, 5}) {
		t.Errorf("Expected [4,5], got %v", data)
	}
	if stream.Pointer != 5 {
		t.Errorf("Expected pointer 5, got %d", stream.Pointer)
	}

	// Test reading at end
	data, n = stream.ChompN(1)
	if n != 0 {
		t.Errorf("Expected to read 0 bytes, got %d", n)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty slice, got %v", data)
	}

	// Test negative count
	stream.Pointer = 0
	data, n = stream.ChompN(-1)
	if n != 0 {
		t.Errorf("Expected to read 0 bytes with negative count, got %d", n)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty slice with negative count, got %v", data)
	}
}

func TestByteStreamChomp(t *testing.T) {
	testData := []byte{1, 2, 3}
	stream := puzzle.NewByteStream(testData)

	data, n := stream.Chomp()
	if n != 1 {
		t.Errorf("Expected to read 1 byte, got %d", n)
	}
	if data[0] != 1 {
		t.Errorf("Expected 1, got %d", data[0])
	}

	data, n = stream.Chomp()
	if n != 1 {
		t.Errorf("Expected to read 1 byte, got %d", n)
	}
	if data[0] != 2 {
		t.Errorf("Expected 2, got %d", data[0])
	}

	data, n = stream.Chomp()
	if n != 1 {
		t.Errorf("Expected to read 1 byte, got %d", n)
	}
	if data[0] != 3 {
		t.Errorf("Expected 3, got %d", data[0])
	}

	data, n = stream.Chomp()
	if n != 0 {
		t.Errorf("Expected to read 0 bytes at end, got %d", n)
	}
}

func TestByteStreamReadString(t *testing.T) {
	// Test normal case
	testData := []byte{'H', 'e', 'l', 'l', 'o', 0, 'W', 'o', 'r', 'l', 'd', 0}
	stream := puzzle.NewByteStream(testData)

	str, n := stream.ReadString()
	if str != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", str)
	}
	if n != 6 {
		t.Errorf("Expected offset 6, got %d", n)
	}
	if stream.Pointer != 6 {
		t.Errorf("Expected pointer 6, got %d", stream.Pointer)
	}

	str, n = stream.ReadString()
	if str != "World" {
		t.Errorf("Expected 'World', got '%s'", str)
	}
	if n != 6 {
		t.Errorf("Expected offset 6, got %d", n)
	}
	if stream.Pointer != 12 {
		t.Errorf("Expected pointer 12, got %d", stream.Pointer)
	}

	// Test missing null terminator
	testData = []byte{'T', 'e', 's', 't'}
	stream = puzzle.NewByteStream(testData)

	str, n = stream.ReadString()
	if n != -1 {
		t.Errorf("Expected error code -1 for missing null terminator, got %d", n)
	}
	if stream.Pointer != 0 {
		t.Errorf("Pointer should not move on error, got %d", stream.Pointer)
	}

	// Test empty string
	testData = []byte{0, 'A', 'B', 'C'}
	stream = puzzle.NewByteStream(testData)

	str, n = stream.ReadString()
	if str != "" {
		t.Errorf("Expected empty string, got '%s'", str)
	}
	if n != 1 {
		t.Errorf("Expected offset 1 for empty string, got %d", n)
	}
	if stream.Pointer != 1 {
		t.Errorf("Expected pointer 1, got %d", stream.Pointer)
	}
}
