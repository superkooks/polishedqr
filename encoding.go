package main

import (
	"bytes"
	"regexp"
	"strconv"
)

type CharacterSet int

const (
	Numeric = iota
	Alphanumeric
	Bytes
	// Kanji not currently supported
)

var numericRE = regexp.MustCompile(`\d+`)
var alphanumericRE = regexp.MustCompile(`[\dA-Z\ $%*+\-.\/:]+`)

func AutodetectCharacterSet(data []byte) CharacterSet {
	// Feels hacky, but the most sane way to do this
	if bytes.Equal(numericRE.Find(data), data) {
		return Numeric
	} else if bytes.Equal(alphanumericRE.Find(data), data) {
		return Alphanumeric
	} else {
		return Bytes
	}
}

type Bits []uint8

func CharacterCount(count int, mode CharacterSet, version int) Bits {
	var bitCapacity int

	// Different versions and encodings have different numbers of bits
	// in the character count
	if version < 10 {
		switch mode {
		case Numeric:
			bitCapacity = 10
		case Alphanumeric:
			bitCapacity = 9
		case Bytes:
			bitCapacity = 8
		}
	} else if version < 27 {
		switch mode {
		case Numeric:
			bitCapacity = 12
		case Alphanumeric:
			bitCapacity = 11
		case Bytes:
			bitCapacity = 16
		}
	} else {
		switch mode {
		case Numeric:
			bitCapacity = 14
		case Alphanumeric:
			bitCapacity = 13
		case Bytes:
			bitCapacity = 16
		}
	}

	var out Bits
	for i := bitCapacity - 1; i >= 0; i-- {
		// cursed masking
		out = append(out, uint8((count&int(1<<i))>>i))
	}

	return out
}

func ConvertToNumeric(data []byte, version int) Bits {
	var b Bits

	// Add mode indicator
	b = append(b, Bits{0, 0, 0, 1}...)

	// Add character count
	b = append(b, CharacterCount(len(data), Numeric, version)...)

	// Divide into groups of three digits and convert to bits
	for i := 0; i < len(data); i += 3 {
		// Group digits
		var g []byte
		if len(data)-i < 3 {
			g = data[i:]
		} else {
			g = data[i : i+3]
		}

		// Convert group
		gi, err := strconv.Atoi(string(g))
		if err != nil {
			panic(err)
		}

		// Small groups are encoded in less bits
		bc := 10 - 1
		if len(g) == 2 {
			bc = 7 - 1
		} else if len(g) == 1 {
			bc = 4 - 1
		}

		// Add bits
		for i := bc; i >= 0; i-- {
			b = append(b, uint8((gi&int(1<<i))>>i))
		}
	}

	return b
}

func ConvertToAlphanumeric(data []byte, version int) Bits {
	var b Bits

	// Add mode indicator
	b = append(b, Bits{0, 0, 1, 0}...)

	// Add character count
	b = append(b, CharacterCount(len(data), Alphanumeric, version)...)

	// Divide into groups of three digits and convert to bits
	for i := 0; i < len(data); i += 2 {
		// Group digits
		var g []byte
		if len(data)-i < 2 {
			g = data[i:]
		} else {
			g = data[i : i+2]
		}

		// Convert group
		var gi int
		if len(g) > 1 {
			gi = 45 * alphanumericTable[g[0]]
			gi += alphanumericTable[g[1]]
		} else {
			gi = alphanumericTable[g[0]]
		}

		// Small groups are encoded in less bits
		bc := 11 - 1
		if len(g) == 1 {
			bc = 6 - 1
		}

		// Add bits
		for i := bc; i >= 0; i-- {
			b = append(b, uint8((gi&int(1<<i))>>i))
		}
	}

	return b
}

func ConvertToBytes(data []byte, version int) Bits {
	var b Bits

	// Add mode indicator
	b = append(b, Bits{0, 1, 0, 0}...)

	// Add character count
	b = append(b, CharacterCount(len(data), Bytes, version)...)

	// Convert each byte into bits
	for _, v := range data {
		for i := 7; i >= 0; i-- {
			b = append(b, uint8((v&(1<<i)))>>i)
		}
	}

	return b
}

// I love manually encoding tables into code by hand
var alphanumericTable = map[byte]int{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'A': 10,
	'B': 11,
	'C': 12,
	'D': 13,
	'E': 14,
	'F': 15,
	'G': 16,
	'H': 17,
	'I': 18,
	'J': 19,
	'K': 20,
	'L': 21,
	'M': 22,
	'N': 23,
	'O': 24,
	'P': 25,
	'Q': 26,
	'R': 27,
	'S': 28,
	'T': 29,
	'U': 30,
	'V': 31,
	'W': 32,
	'X': 33,
	'Y': 34,
	'Z': 35,
	' ': 36,
	'$': 37,
	'%': 38,
	'*': 39,
	'+': 40,
	'-': 41,
	'.': 42,
	'/': 43,
	':': 44,
}
