package puzzle

import (
	"bytes"
)

type ByteStream struct {
	Raw     []byte
	Pointer int
	Size    int
}

func NewByteStream(raw []byte) *ByteStream {
	return &ByteStream{Raw: raw, Size: len(raw)}
}

func (s *ByteStream) IncPointer(n int) {
	s.Pointer += n
}

func (s *ByteStream) ChompN(n int) ([]byte, int) {
	if s.Pointer == s.Size || n < 0 {
		return []byte{}, 0
	}

	if n+s.Pointer > s.Size {
		n = s.Size - s.Pointer
	}

	out := s.Raw[s.Pointer : s.Pointer+n]
	s.IncPointer(n)

	return out, int(n)
}

func (s *ByteStream) Chomp() ([]byte, int) {
	return s.ChompN(1)
}

func (s *ByteStream) ReadString() (string, int) {
	nullIndex := bytes.IndexByte(s.Raw[s.Pointer:], 0)
	if nullIndex == -1 {
		return "", -1
	}

	str := string(s.Raw[s.Pointer : s.Pointer+nullIndex])
	offset := nullIndex + 1 // Move past the null terminator

	s.IncPointer(offset)

	return str, offset
}
