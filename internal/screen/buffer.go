package screen

import (
	"errors"
	"strings"
)

type Buffer struct {
	data []string
	w    int
	h    int
}

func NewBuffer(width, height int) *Buffer {
	data := make([]string, height*width)
	return &Buffer{data: data, w: width, h: height}
}

func (b *Buffer) getIndex(row, col int) (int, error) {
	idx := (row * b.w) + col
	if idx < 0 || idx >= len(b.data) || row < 0 || row >= b.h || col < 0 || col >= b.w {
		return -1, errors.New("index out of bounds")
	}
	return idx, nil
}

// Set safely sets a value at the specified row and column
func (b *Buffer) Set(row, col int, value string) error {
	if idx, err := b.getIndex(row, col); err == nil {
		b.data[idx] = value
		return nil
	} else {
		return err
	}
}

// Get safely retrieves a value at the specified row and column
func (b *Buffer) Get(row, col int) (string, error) {
	if idx, err := b.getIndex(row, col); err == nil {
		return b.data[idx], nil
	} else {
		return "", err
	}
}

func (b *Buffer) String() string {
	var sb strings.Builder

	for row := range b.h {
		if row > 0 {
			sb.WriteString("\n")
		}
		for col := range b.w {
			if cell, err := b.Get(row, col); err == nil {
				sb.WriteString(cell)
			}
		}
	}

	return sb.String()
}
