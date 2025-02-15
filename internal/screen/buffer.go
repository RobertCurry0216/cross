package screen

import (
	"errors"
)

// Buffer is a wrapper around a 2D string slice that provides safe access
type Buffer struct {
	data [][]string
}

// NewSafe2DArray creates a new Safe2DArray from a 2D string slice
func NewBuffer(width, height int) *Buffer {
	data := make([][]string, height)
	for i := 0; i < len(data); i++ {
		data[i] = make([]string, width)
	}

	return &Buffer{data: data}
}

// Set safely sets a value at the specified row and column
func (s *Buffer) Set(row, col int, value string) error {
	if row < 0 || row >= len(s.data) {
		return errors.New("row index out of bounds")
	}
	if col < 0 || col >= len(s.data[row]) {
		return errors.New("column index out of bounds")
	}
	s.data[row][col] = value
	return nil
}

// Get safely retrieves a value at the specified row and column
func (s *Buffer) Get(row, col int) (string, error) {
	if row < 0 || row >= len(s.data) {
		return "", errors.New("row index out of bounds")
	}
	if col < 0 || col >= len(s.data[row]) {
		return "", errors.New("column index out of bounds")
	}
	return s.data[row][col], nil
}

func (b *Buffer) Size() (int, int) {
	if len(b.data) == 0 {
		return 0, 0
	}
	return len(b.data), len(b.data[0])
}
