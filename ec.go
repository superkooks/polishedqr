package main

import (
	"fmt"
)

// Retunrs the combined data words and error words
func generateErrorWords(codewords []uint8, version int, ecLevel string) []uint8 {
	type block struct {
		ecCount    int
		dataWords  []uint8
		errorWords []uint8
	}
	var blocks []*block

	// Split codewords into blocks
	var c int
	for _, blockType := range codeWordTable[version][ecLevel].blocks {
		for i := 0; i < blockType.count; i++ {
			blocks = append(blocks, &block{
				dataWords: codewords[c : c+blockType.dataWords],
				ecCount:   codeWordTable[version][ecLevel].ecWordsPerBlock,
			})
			c += blockType.dataWords
		}
	}

	// Generate error codes for each block
	for _, b := range blocks {
		b.errorWords = rsEncode(b.dataWords, b.ecCount)
		fmt.Println("ec codewords", rsEncode(b.dataWords, b.ecCount))
	}

	fmt.Println(blocks[0])

	// Assemble the final sequence, taking each block in turn
	var out []uint8
	for i := 0; i < len(blocks[len(blocks)-1].dataWords); i++ {
		for _, v := range blocks {
			if i < len(v.dataWords) {
				out = append(out, v.dataWords[i])
			}
		}
	}
	for i := 0; i < len(blocks[len(blocks)-1].errorWords); i++ {
		for _, v := range blocks {
			if i < len(v.errorWords) {
				out = append(out, v.errorWords[i])
			}
		}
	}

	return out
}

func rsGenerator(symbols int) []int {
	g := []int{1}
	for i := 0; i < symbols; i++ {
		g = gfPolyMul(g, []int{1, gfPow(2, i)})
	}
	return g
}

// Returns the ec codewords
func rsEncode(msg []uint8, symbols int) []uint8 {
	gen := rsGenerator(symbols)

	m := make([]int, len(msg))
	for k, v := range msg {
		m[k] = int(v)
	}

	_, remainder := gfPolyDiv(append(m, make([]int, len(gen)-1)...), gen)

	r := make([]uint8, len(remainder))
	for k, v := range remainder {
		r[k] = uint8(v)
	}
	return r
}

// Positive modulo, returns non negative solution to x `%` d
func pmod(x, b int) int {
	x = x % b
	if x >= 0 {
		return x
	}
	if b < 0 {
		return x - b
	}
	return x + b
}

var gfExp = make([]int, 512)
var gfLog = make([]int, 256)

func gfAdd(x, y int) int {
	return x ^ y
}

func gfSub(x, y int) int {
	return x ^ y
}

// Uses Russian Peasant Multiplication. Buggered if I know.
func gfMul(x int, y int, prim int, fieldCharacFull int, carryless bool) int {
	r := 0
	for y > 0 {
		if y&1 > 0 {
			if carryless {
				r = r ^ x
			} else {
				r += x
			}
		}

		y = y >> 1
		x = x << 1
		if prim > 0 && x&fieldCharacFull > 0 {
			x = x ^ prim
		}
	}

	return r
}

func gfQuickMul(x, y int) int {
	if x == 0 || y == 0 {
		return 0
	}

	return gfExp[gfLog[x]+gfLog[y]]
}

func gfQuickDiv(x, y int) int {
	if y == 0 {
		panic("cannot divide by 0")
	}

	if x == 0 {
		return 0
	}

	return gfExp[pmod(gfLog[x]+255-gfLog[y], 255)]
}

func gfPow(x, power int) int {
	t := pmod(gfLog[x]*power, 255)
	if t < 0 {
		return gfExp[len(gfExp)+t]
	}
	return gfExp[t]
}

func gfInverse(x int) int {
	return gfExp[255-gfLog[x]]
}

func gfPolyScale(p []int, x int) []int {
	r := make([]int, len(p))
	for i := 0; i < len(p); i++ {
		r[i] = gfQuickMul(p[i], x)
	}
	return r
}

func gfPolyAdd(p, q []int) []int {
	var r []int
	if len(p) > len(q) {
		r = make([]int, len(p))
	} else {
		r = make([]int, len(q))
	}

	for i := 0; i < len(p); i++ {
		r[i+len(r)-len(p)] = p[i]
	}

	for i := 0; i < len(q); i++ {
		r[i+len(r)-len(q)] ^= q[i]
	}

	return r
}

func gfPolyMul(p, q []int) []int {
	r := make([]int, len(p)+len(q)-1)

	for j := 0; j < len(q); j++ {
		for i := 0; i < len(p); i++ {
			r[i+j] ^= gfQuickMul(p[i], q[j])
		}
	}

	return r
}

func gfPolyDiv(dividend, divisor []int) ([]int, []int) {
	msgOut := make([]int, len(dividend))
	copy(msgOut, dividend)

	for i := 0; i < len(dividend)-len(divisor)+1; i++ {
		coef := msgOut[i]
		if coef != 0 {
			for j := 1; j < len(divisor); j++ {
				if divisor[j] != 0 {
					msgOut[i+j] ^= gfQuickMul(divisor[j], coef)
				}
			}
		}
	}

	separator := len(divisor) - 1
	return msgOut[:len(msgOut)-separator], msgOut[len(msgOut)-separator:]
}

func gfPolyEval(poly []int, x int) int {
	y := poly[0]
	for i := 1; i < len(poly); i++ {
		y = gfQuickMul(y, x) ^ poly[i]
	}

	return y
}

var _ = initGFTables()

func initGFTables() bool {
	prim := 0x11d

	x := 1
	for i := 0; i < 255; i++ {
		gfExp[i] = x
		gfLog[x] = i
		x = gfMul(x, 2, prim, 256, true)
	}

	for i := 255; i < 512; i++ {
		gfExp[i] = gfExp[i-255]
	}

	return true
}

// I love plagiarism
