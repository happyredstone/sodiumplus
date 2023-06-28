package helpers

import "encoding/binary"

type LengthHasher struct {
	length uint64
}

func (h *LengthHasher) Write(p []byte) (n int, err error) {
	h.length += uint64(len(p))

	return len(p), nil
}

func (h *LengthHasher) Sum(b []byte) []byte {
	ext := append(b, make([]byte, 8)...)

	binary.BigEndian.PutUint64(ext, h.length)

	return ext
}

func (h *LengthHasher) Size() int {
	return 8
}

func (h *LengthHasher) BlockSize() int {
	return 1
}

func (h *LengthHasher) Reset() {
	h.length = 0
}
