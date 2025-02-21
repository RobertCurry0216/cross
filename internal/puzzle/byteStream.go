package puzzle

import (
	"bytes"
)

type ByteStream struct {
	raw     []byte
	pointer int
	size    int
}

func NewByteStream(raw []byte) *ByteStream {
	return &ByteStream{raw: raw, size: len(raw)}
}

func (s *ByteStream) incPointer(n int) {
	s.pointer += n
}

func (s *ByteStream) ChompN(n int) ([]byte, int) {
	if s.pointer == s.size || n < 0 {
		return []byte{}, 0
	}

	if n+s.pointer > s.size {
		n = s.size - s.pointer
	}

	out := s.raw[s.pointer : s.pointer+n]
	s.incPointer(n)

	return out, int(n)
}

func (s *ByteStream) Chomp() ([]byte, int) {
	return s.ChompN(1)
}

func (s *ByteStream) readString() (string, int) {
	nullIndex := bytes.IndexByte(s.raw[s.pointer:], 0)
	if nullIndex == -1 {
		return "", -1
	}

	str := string(s.raw[s.pointer : s.pointer+nullIndex])
	offset := nullIndex + 1 // Move past the null terminator

	s.incPointer(offset)

	return str, offset
}
