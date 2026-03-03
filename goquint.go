// Package goquint implements proquint encoding and decoding.
//
// Proquints are pronounceable representations of integers, as described
// in the specification: https://arxiv.org/html/0901.4016
//
// A proquint encodes a 16-bit value as a 5-character "quintuplet" of
// alternating consonants and vowels (CVCVC). Two quintuplets separated
// by a hyphen encode a 32-bit value (e.g. "lusab-babad").
package goquint

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var (
	consonants = "bdfghjklmnpqrstvz"
	vowels     = "aiou"
)

// Encode converts a 32-bit number to a proquint string (e.g. "lusab-babad").
func Encode(n uint32) string {
	q1 := encodeQuintuplet(byte((n>>24)&0xFF), byte((n>>16)&0xFF))
	q2 := encodeQuintuplet(byte((n>>8)&0xFF), byte(n&0xFF))
	return q1 + "-" + q2
}

// Decode converts a proquint string back to a 32-bit number.
// Hyphens in the input are ignored.
func Decode(proquint string) (uint32, error) {
	clean := ""
	for _, ch := range proquint {
		if ch != '-' {
			clean += string(ch)
		}
	}

	if len(clean) != 10 {
		return 0, fmt.Errorf("goquint: invalid proquint length: expected 10 characters, got %d", len(clean))
	}

	high, err := decodeQuintuplet(clean[:5])
	if err != nil {
		return 0, fmt.Errorf("goquint: invalid first quintuplet: %w", err)
	}

	low, err := decodeQuintuplet(clean[5:])
	if err != nil {
		return 0, fmt.Errorf("goquint: invalid second quintuplet: %w", err)
	}

	return uint32(high[0])<<24 | uint32(high[1])<<16 | uint32(low[0])<<8 | uint32(low[1]), nil
}

// Random generates a random proquint using crypto/rand.
func Random() string {
	num, _ := rand.Int(rand.Reader, big.NewInt(0xFFFFFFFF))
	return Encode(uint32(num.Int64()))
}

// EncodeHex converts the first 32 bits of a hex string to a proquint.
// If the hex string is shorter than 8 characters, a random proquint is returned.
func EncodeHex(hexStr string) string {
	if len(hexStr) < 8 {
		return Random()
	}

	var num uint32
	fmt.Sscanf(hexStr[:8], "%x", &num)
	return Encode(num)
}

func encodeQuintuplet(high, low byte) string {
	val := uint16(high)<<8 | uint16(low)

	c1 := (val >> 12) & 0x0F
	v1 := (val >> 10) & 0x03
	c2 := (val >> 6) & 0x0F
	v2 := (val >> 4) & 0x03
	c3 := val & 0x0F

	return string([]byte{
		consonants[c1],
		vowels[v1],
		consonants[c2],
		vowels[v2],
		consonants[c3],
	})
}

func decodeQuintuplet(q string) ([2]byte, error) {
	if len(q) != 5 {
		return [2]byte{}, fmt.Errorf("quintuplet must be 5 characters, got %d", len(q))
	}

	c1 := findIndex(consonants, q[0])
	v1 := findIndex(vowels, q[1])
	c2 := findIndex(consonants, q[2])
	v2 := findIndex(vowels, q[3])
	c3 := findIndex(consonants, q[4])

	if c1 < 0 || v1 < 0 || c2 < 0 || v2 < 0 || c3 < 0 {
		return [2]byte{}, fmt.Errorf("invalid characters in quintuplet: %s", q)
	}

	val := uint16(c1)<<12 | uint16(v1)<<10 | uint16(c2)<<6 | uint16(v2)<<4 | uint16(c3)
	return [2]byte{byte(val >> 8), byte(val & 0xFF)}, nil
}

func findIndex(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
