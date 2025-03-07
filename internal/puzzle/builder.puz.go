package puzzle

import (
	"encoding/binary"
	"fmt"
	"os"
)

type PuzBuilder struct {
	raw      []byte
	filepath string
	Puzzle   *Puzzle

	// write locations
	cib         []byte
	puzzleInput []byte
	gext        []byte
}

func NewPuzBuilder(raw []byte, path string) *PuzBuilder {
	return &PuzBuilder{raw: raw, filepath: path, cib: raw[0x0E : 0x0E+2]}
}

func (b *PuzBuilder) Build() (*Puzzle, error) {
	puz := NewPuzzle()

	// extract basic data
	// checksum := binary.LittleEndian.Uint16(b.raw[0:2])
	puz.Width = int(b.raw[0x2c])
	puz.Height = int(b.raw[0x2d])

	// extract grid information
	gridSize := puz.Width * puz.Height
	stream := NewByteStream(b.raw)
	stream.ChompN(0x34)

	puz.solution = make([]byte, gridSize)
	if data, n := stream.ChompN(gridSize); n == gridSize {
		copy(puz.solution, data)
	} else {
		return nil, fmt.Errorf("Malformed .puz file: missing grid")
	}

	puz.input = make([]byte, gridSize)
	if data, n := stream.ChompN(gridSize); n == gridSize {
		copy(puz.input, data)
		b.puzzleInput = data
	} else {
		return nil, fmt.Errorf("Malformed .puz file: missing grid")
	}

	// extract title
	if title, n := stream.readString(); n == -1 {
		return nil, fmt.Errorf("Malformed .puz file: missing null terminator in title")
	} else {
		puz.Title = title
	}

	// extract author
	if author, n := stream.readString(); n == -1 {
		return nil, fmt.Errorf("Malformed .puz file: missing null terminator in author")
	} else {
		puz.Author = author
	}

	// extract copyright
	if copyright, n := stream.readString(); n == -1 {
		return nil, fmt.Errorf("Malformed .puz file: missing null terminator in copyright")
	} else {
		puz.Copyright = copyright
	}

	// extract clues
	clueCount := binary.LittleEndian.Uint16(b.raw[0x2e : 0x2e+2])
	puz.Clues = make([]*Clue, clueCount) // Prepare slice to hold clues

	for i := range clueCount {
		// Extract the clue string
		if clueText, n := stream.readString(); n == -1 {
			return nil, fmt.Errorf("Malformed .puz file: missing null terminator in clues")
		} else {
			puz.Clues[i] = NewClue(clueText)
		}
	}

	if notes, n := stream.readString(); n != -1 {
		puz.Notes = notes
	}

	// initialize cells
	InitPuzzle(puz)

	// extra info
	for i := 0; i < 4; i++ {
		if section, n := stream.ChompN(0x04); n == 0x04 {
			var data []byte
			if lenData, n := stream.ChompN(0x02); n == 0x02 {
				l := int(binary.LittleEndian.Uint16(lenData))
				// todo: chksum
				stream.ChompN(0x02)
				data, _ = stream.ChompN(l)
			} else {
				break
			}

			switch string(section) {
			case "GEXT":
				applyGEXT(puz, data)
				b.gext = data
			default:
				fmt.Println(string(section))
			}
		} else {
			break
		}
		stream.Chomp()
	}

	// cross reference
	b.Puzzle = puz
	puz.builder = b

	return puz, nil
}

func (b *PuzBuilder) Write() {
	b.updateRaw()
	os.WriteFile(b.filepath, b.raw, 0644)
}

func (b *PuzBuilder) updateRaw() {
	copy(b.Puzzle.input, b.puzzleInput)
	cib := b.getCIB()
	binary.LittleEndian.PutUint16(b.cib, cib)
}

func applyGEXT(puz *Puzzle, data []byte) {
	for i, cell := range puz.Grid {
		datum := data[i]
		if datum&0x80 != 0 {
			cell.IsCircled = true
		}
	}
}

func (b *PuzBuilder) Validate() error {
	if b.Puzzle == nil {
		return fmt.Errorf("Validation error: Puzzle must be built before validation")
	}

	if len(b.raw) < 2 {
		return fmt.Errorf("Validation error: too short")
	}

	// validate cib
	cib := b.getCIB()
	cibCksum := checksumRegion(b.raw[0x2C:0x2C+8], 0)
	if cib != cibCksum {
		return fmt.Errorf("Validation error: cib")
	}

	//validate check sum
	target := binary.LittleEndian.Uint16(b.raw[0:2])
	cksum := b.getCheckSum()

	if target != cksum {
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

func (b *PuzBuilder) getCheckSum() uint16 {
	cksum := b.getCIB()

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

func (b *PuzBuilder) getCIB() uint16 {
	return binary.LittleEndian.Uint16(b.raw[0x0E : 0x0E+2])
}
