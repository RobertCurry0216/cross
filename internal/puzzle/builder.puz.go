package puzzle

import (
	"encoding/binary"
	"fmt"
)

type PuzBuilder struct {
	raw    []byte
	Puzzle *Puzzle
}

func (b *PuzBuilder) Build() (*Puzzle, error) {
	puz := NewPuzzle()

	// extract basic data
	// checksum := binary.LittleEndian.Uint16(b.raw[0:2])
	puz.Width = int(b.raw[0x2c])
	puz.Height = int(b.raw[0x2d])

	// extract grid information
	gridSize := puz.Width * puz.Height
	offset := int(0x34)

	puz.solution = make([]byte, gridSize)
	copy(puz.solution, b.raw[offset:offset+gridSize])
	offset += gridSize

	puz.input = make([]byte, gridSize)
	copy(puz.input, b.raw[offset:offset+gridSize])
	offset += gridSize

	// extract title
	if title, n := readString(b.raw[offset:]); n == -1 {
		return nil, fmt.Errorf("Malformed .puz file: missing null terminator in title")
	} else {
		puz.Title = title
		offset += n
	}

	// extract author
	if author, n := readString(b.raw[offset:]); n == -1 {
		return nil, fmt.Errorf("Malformed .puz file: missing null terminator in author")
	} else {
		puz.Author = author
		offset += n
	}

	// extract copyright
	if copyright, n := readString(b.raw[offset:]); n == -1 {
		return nil, fmt.Errorf("Malformed .puz file: missing null terminator in copyright")
	} else {
		puz.Copyright = copyright
		offset += n
	}

	// extract clues
	clueCount := binary.LittleEndian.Uint16(b.raw[0x2e : 0x2e+2])
	puz.Clues = make([]*Clue, clueCount) // Prepare slice to hold clues

	for i := range clueCount {
		// Extract the clue string
		if clueText, n := readString(b.raw[offset:]); n == -1 {
			return nil, fmt.Errorf("Malformed .puz file: missing null terminator in clues")
		} else {
			puz.Clues[i] = NewClue(clueText)
			offset += n
		}
	}

	puz.Notes, _ = readString(b.raw[offset:])

	// cross reference
	b.Puzzle = puz
	puz.builder = b

	InitPuzzle(puz)

	return puz, nil
}

func (b *PuzBuilder) Validate() error {
	if b.Puzzle == nil {
		return fmt.Errorf("Validation error: Puzzle must be built before validation")
	}

	if len(b.raw) < 2 {
		return fmt.Errorf("Validation error: too short")
	}

	// validate cib
	cib := b.cib()
	cibCksum := checksumRegion(b.raw[0x2C:0x2C+8], 0)
	if cib != cibCksum {
		return fmt.Errorf("Validation error: cib")
	}

	//validate check sum
	target := binary.LittleEndian.Uint16(b.raw[0:2])
	cksum := b.checkSum()

	if target != cksum {
		fmt.Println("final", target, cksum)
		return fmt.Errorf("Validation error: checksum")
	}

	return nil
}

func checksumRegion(region []byte, cksum uint16) uint16 {
	for _, val := range region {
		// Check the least significant bit
		if cksum&0x01 == 0x01 {
			cksum = (cksum >> 1) + 0x8000 // Rotate with carry
		} else {
			cksum = cksum >> 1
		}

		// Add the byte value, ensuring 16-bit result
		cksum = (cksum + uint16(val)) & 0xFFFF
	}
	return cksum
}

func (b *PuzBuilder) checkSum() uint16 {
	cksum := b.cib()

	//validate check sum
	cksum = checksumRegion(b.Puzzle.solution, cksum)
	cksum = checksumRegion(b.Puzzle.input, cksum)

	for _, metaField := range []string{b.Puzzle.Title, b.Puzzle.Author, b.Puzzle.Copyright} {
		if len(metaField) > 0 {
			cksum = checksumRegion([]byte(metaField+"\x00"), cksum) // Include null terminator
		}
	}

	for _, clue := range b.Puzzle.Clues {
		cksum = checksumRegion([]byte(clue.Text), cksum)
	}

	if len(b.Puzzle.Notes) > 0 {
		cksum = checksumRegion([]byte(b.Puzzle.Notes+"\x00"), cksum) // Include null terminator
	}

	return cksum
}

func (b *PuzBuilder) cib() uint16 {
	return binary.LittleEndian.Uint16(b.raw[0x0E : 0x0E+2])
}
