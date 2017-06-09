package core

import (
	"crypto/rand"
	"fmt"
	"io"
)

// MaxInt finds returns value between a and b
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MaxIntSlice returns max value in int slice
func MaxIntSlice(s []int) int {
	max := s[0]
	for _, n := range s {
		if n > max {
			max = n
		}
	}
	return max
}

// MinInt returns min value between a nd b
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MinIntSlice finds min value in int slice
func MinIntSlice(s []int) int {
	min := s[0]
	for _, n := range s {
		if n < min {
			min = n
		}
	}
	return min
}

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
