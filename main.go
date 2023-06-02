package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type CreateOptions struct {
	// Selects the level of error correction to use.
	// "L": recover 7% of codewords
	// "M": recover 15%
	// "Q": recover 25%
	// "H": recover 30%
	// If unset, defaults to "M"
	ErrorCorrectionLevel string

	// The character set to encode the data with.
	// If unset, the character set will be chosen automatically,
	// however, encoding formats will not be mixed on the same code.
	CharacterSet *CharacterSet

	// The version (size) of the qr code.
	// If unset, the version will be the smallest that can fit the data
	Version int
}

// Create a qr code from data with options, which may be nil.
// Data that is numeric or alphanumeric should be passed in their ascii form
func CreateQRCode(data []byte, opts *CreateOptions) *image.RGBA {
	if opts == nil {
		opts = &CreateOptions{}
	}

	if opts.ErrorCorrectionLevel == "" {
		opts.ErrorCorrectionLevel = "M"
	}

	// Encode data in correct mode
	var mode CharacterSet
	if opts.CharacterSet != nil {
		mode = *opts.CharacterSet
	} else {
		mode = AutodetectCharacterSet(data)
	}

	var version int
	if opts.Version == 0 {
		// We will iteratively increase version until data fits
		version = 1
	} else {
		version = opts.Version
	}

	var dataBits Bits
	var totalDatawords int
	for {
		// Encode data
		switch mode {
		case Numeric:
			dataBits = ConvertToNumeric(data, version)
		case Alphanumeric:
			dataBits = ConvertToAlphanumeric(data, version)
		case Bytes:
			dataBits = ConvertToBytes(data, version)
		default:
			panic("unsupported encoding mode")
		}

		// Get total data size of this symbol
		totalDatawords = 0
		blocks := codeWordTable[version][opts.ErrorCorrectionLevel]
		for _, v := range blocks.blocks {
			totalDatawords += v.dataWords * v.count
		}

		// Check whether the data fits
		fmt.Println(len(dataBits), totalDatawords, version)
		if len(dataBits) > totalDatawords*8 {
			if version != opts.Version {
				// Version is unset in options, try a larger symbol size
				if version == 40 {
					panic("data cannot fit in largest qr code")
				}

				version++
				continue
			} else {
				panic("data cannot fit in designated size qr code")
			}
		}

		break
	}

	// Add terminator (if required)
	for i := 0; i < 4 && len(dataBits) < totalDatawords*8; i++ {
		dataBits = append(dataBits, 0)
	}

	// Pad to nearest 8 bits
	for i := 0; i < len(dataBits)%8; i++ {
		dataBits = append(dataBits, 0)
	}

	// Convert to codewords
	var codewords []uint8
	for i := 0; i < len(dataBits); i += 8 {
		var acc uint8
		acc += dataBits[i] << 7
		acc += dataBits[i+1] << 6
		acc += dataBits[i+2] << 5
		acc += dataBits[i+3] << 4
		acc += dataBits[i+4] << 3
		acc += dataBits[i+5] << 2
		acc += dataBits[i+6] << 1
		acc += dataBits[i+7] << 0
		codewords = append(codewords, acc)
	}

	// Add padding codewords
	for i := 0; len(codewords) < totalDatawords; i++ {
		if i%2 == 0 {
			codewords = append(codewords, 0b11101100)
		} else {
			codewords = append(codewords, 0b00010001)
		}
	}

	// Generate error correction
	allwords := generateErrorWords(codewords, version, opts.ErrorCorrectionLevel)

	// Create image
	i := image.NewRGBA(image.Rect(0, 0, 17+version*4, 17+version*4))

	// Draw background (so we can see coding area)
	iterateRect(i.Rect.Dx(), i.Rect.Dy(), func(x, y int) {
		i.SetRGBA(x, y, BLUE)
	})

	// Draw finder patterns in three corners
	drawFinderPattern(i, 0, 0)
	drawFinderPattern(i, 0, i.Rect.Dy()-7)
	drawFinderPattern(i, i.Rect.Dx()-7, 0)

	// Place temporary format bits
	drawTempFormatBits(i)
	if version >= 7 {
		drawTempVersionBits(i)
	}

	// Draw both timing patterns
	drawTimingPatterns(i)

	// Draw alignment patterns
	drawAlignmentPatterns(i)

	// Draw the data onto the qr code with a zig-zag pattern
	writeData(i, allwords)

	// Apply the best mask
	pattern := applyBestMask(i, opts.ErrorCorrectionLevel, version)
	addFormatAndVersionInfo(i, opts.ErrorCorrectionLevel, pattern, version)

	// Place qr code in quiet zone
	i = quietZone(i)

	return i
}

func main() {
	f, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}

	err = png.Encode(f, CreateQRCode([]byte("00000.UFF7THUFF7000001F8F7THUFF7UF00000000UFF7UFF7F7UFF7UF00000000UFF7UEUFF7T*000005F7UFF7UEUFF7UFF500000001F7T*00000.UFF7UF7QF7SK000.QOM:UPUFF7UFEA0000001+F7UFF7THUFF7UFEA0000001+F7UEUFF7UE0000003ZUFF7UF7QF7UFF7SK000000F7UF"), &CreateOptions{
		ErrorCorrectionLevel: "L",
	}))
	if err != nil {
		panic(err)
	}
}
