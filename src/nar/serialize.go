package nar

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MaxInt64 = 1<<(64-1) - 1
)

func readString(r io.Reader) (string, error) {
	size, err := readLongLong(r)
	if err != nil {
		return "", err
	}
	bs := make([]byte, size)
	n, err := r.Read(bs)
	if err != nil {
		return "", err
	}
	if int64(n) != size {
		return "", fmt.Errorf("expected %d bytes, not %d", size, n)
	}

	for _, char := range bs {
		if char == 0 {
			return "", fmt.Errorf("expected no zeros, got %d %v", size, bs)
		}
	}

	err = readPadding(r, size)
	if err != nil {
		return "", err
	}

	//fmt.Println("STR", string(bs))

	return string(bs), nil
}

func readPadding(r io.Reader, l int64) error {
	pad := 8 - (l % 8)
	if pad == 8 {
		// lucky! no need for padding here
		return nil
	}

	bs := make([]byte, pad)
	n, err := r.Read(bs)
	if err != nil {
		return err
	}
	if int64(n) != pad {
		return fmt.Errorf("expected to read %d, got %d", pad, n)
	}
	for _, char := range bs {
		if char != 0 {
			return fmt.Errorf("expected zero padding, got %v", bs)
		}
	}
	return nil
}

const maxInt64 = 1<<63 - 1

func readLongLong(r io.Reader) (int64, error) {
	bs := make([]byte, 8, 8)
	if _, err := io.ReadFull(r, bs); err != nil {
		return 0, err
	}

	// this is uint64, little endian
	// we use int64 later, so error out if it's larger than what we're able to represent.
	num := binary.LittleEndian.Uint64(bs)
	if num > maxInt64 {
		return 0, fmt.Errorf("number is too big: %d > %d", num, maxInt64)
	}

	return int64(num), nil
}

func expectString(r io.Reader, expected string) error {
	s, err := readString(r)
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return err
	}
	if s != expected {
		return fmt.Errorf("expected '%s' got '%s'", expected, s)
	}
	return nil
}
