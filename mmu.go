package gomu8080

import (
	"errors"
)

// MMU - Memory Management Unit
type MMU struct {
	Memory [65536]byte
}

func NewMMU() *MMU {
	return &MMU{}
}

func (m *MMU) Load(length int, data []byte, pos int) error {

	if pos < 0 || pos >= 65536 {
		return errors.New("MMU: Load: Error: invalid memory position")
	}
	if length < 0 || length > 65535 {
		return errors.New("MMU: Load: Error: invalid memory length")
	}

	for i := 0; i < length; i++ {
		m.Memory[pos+i] = data[i]
	}
	return nil
}
