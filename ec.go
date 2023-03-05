package main

import (
	"fmt"

	"github.com/templexxx/reedsolomon"
)

// Returns the combined data words and error words
// code is messy because standard is messy (that is my excuse)
func generateErrorWords(codewords []uint8, version int, ecLevel string) []uint8 {
	type block struct {
		ecCount   int
		dataWords []uint8
		ecWords   []uint8
	}
	var blocks []*block

	// Split codewords into blocks
	var c int
	var maxDataBlockLength int
	for _, blockType := range codeWordTable[version][ecLevel].blocks {
		if maxDataBlockLength < blockType.dataWords {
			// Will be used later
			maxDataBlockLength = blockType.dataWords
		}

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
		fmt.Println(b)
		encoder, err := reedsolomon.New(len(b.dataWords), b.ecCount)
		if err != nil {
			panic(err)
		}

		shards := make([][]byte, len(b.dataWords)+b.ecCount)
		for k := range shards {
			if k < len(b.dataWords) {
				shards[k] = []byte{b.dataWords[k]}
			} else {
				shards[k] = []byte{0}
			}
		}

		err = encoder.Encode(shards)
		if err != nil {
			panic(err)
		}

		fmt.Println(shards)

		for _, v := range shards[len(b.dataWords):] {
			b.ecWords = append(b.ecWords, v...)
		}
	}

	// Mix data words then error words
	var out []uint8
	for i := 0; i < maxDataBlockLength; i++ {
		for _, v := range blocks {
			if i < len(v.dataWords) {
				out = append(out, v.dataWords[i])
			}
		}
	}

	for i := 0; i < blocks[0].ecCount; i++ {
		for _, v := range blocks {
			out = append(out, v.ecWords[i])
		}
	}

	return out
}
