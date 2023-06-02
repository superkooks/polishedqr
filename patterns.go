package main

import (
	"image"
	"image/color"
	"image/draw"
)

func drawFinderPattern(i *image.RGBA, x0, y0 int) {
	iterateRect(9, 9, func(x, y int) {
		i.SetRGBA(x0+x-1, y0+y-1, WHITE)
	})

	iterateRect(7, 7, func(x, y int) {
		i.SetRGBA(x0+x, y0+y, BLACK)
	})

	iterateRect(5, 5, func(x, y int) {
		i.SetRGBA(x0+x+1, y0+y+1, WHITE)
	})

	iterateRect(3, 3, func(x, y int) {
		i.SetRGBA(x0+x+2, y0+y+2, BLACK)
	})
}

func drawTimingPatterns(i *image.RGBA) {
	black := true
	iterateRect(1, i.Rect.Dx()-16, func(x, y int) {
		if black {
			i.SetRGBA(x+6, y+8, BLACK)
		} else {
			i.SetRGBA(x+6, y+8, WHITE)
		}
		black = !black
	})

	black = true
	iterateRect(i.Rect.Dy()-16, 1, func(x, y int) {
		if black {
			i.SetRGBA(x+8, y+6, BLACK)
		} else {
			i.SetRGBA(x+8, y+6, WHITE)
		}
		black = !black
	})
}

func drawAlignmentPattern(i *image.RGBA, x0, y0 int) {
	iterateRect(5, 5, func(x, y int) {
		i.SetRGBA(x0+x, y0+y, BLACK)
	})

	iterateRect(3, 3, func(x, y int) {
		i.SetRGBA(x0+x+1, y0+y+1, WHITE)
	})

	i.SetRGBA(x0+2, y0+2, BLACK)
}

func drawTempFormatBits(i *image.RGBA) {
	iterateRect(9, 1, func(x, y int) {
		i.SetRGBA(x, 8, WHITE)
	})

	iterateRect(1, 8, func(x, y int) {
		i.SetRGBA(8, y, WHITE)
	})

	iterateRect(9, 1, func(x, y int) {
		i.SetRGBA(i.Rect.Dx()-x, 8, WHITE)
	})

	iterateRect(1, 8, func(x, y int) {
		i.SetRGBA(8, i.Rect.Dy()-y, WHITE)
	})

	i.SetRGBA(8, i.Rect.Dy()-8, BLACK)
}

func drawTempVersionBits(i *image.RGBA) {
	iterateRect(3, 6, func(x, y int) {
		i.SetRGBA(i.Rect.Dx()-9-x, y, WHITE)
	})

	iterateRect(6, 3, func(x, y int) {
		i.SetRGBA(x, i.Rect.Dy()-9-y, WHITE)
	})
}

func quietZone(i *image.RGBA) *image.RGBA {
	n := image.NewRGBA(image.Rect(0, 0, i.Rect.Dx()+8, i.Rect.Dy()+8))
	iterateRect(n.Rect.Dx(), n.Rect.Dy(), func(x, y int) {
		n.SetRGBA(x, y, WHITE)
	})

	draw.Draw(n, image.Rect(4, 4, i.Rect.Dx()+4, i.Rect.Dy()+4), i, image.Pt(0, 0), draw.Over)

	return n
}

func addFormatAndVersionInfo(img *image.RGBA, ecLevel string, maskPattern int, version int) {
	// Generate 15 bits of format info
	var ecNum int
	switch ecLevel {
	case "L":
		ecNum = 1
	case "M":
		ecNum = 0
	case "Q":
		ecNum = 3
	case "H":
		ecNum = 2
	default:
		panic("invalid error correction level")
	}

	if maskPattern > 7 || maskPattern < 0 {
		panic("invalid mask pattern")
	}

	code := ecNum<<3 | maskPattern
	encodedFormat := ((code << 10) | checkFormat(code<<10))
	maskedFormat := 0b101010000010010 ^ encodedFormat

	// Convert into colors
	var colors []color.RGBA
	for i := 0; i < 15; i++ {
		if maskedFormat&(1<<i) > 0 {
			colors = append(colors, BLACK)
		} else {
			colors = append(colors, WHITE)
		}
	}

	// Write format info (top left)
	img.SetRGBA(8, 0, colors[0])
	img.SetRGBA(8, 1, colors[1])
	img.SetRGBA(8, 2, colors[2])
	img.SetRGBA(8, 3, colors[3])
	img.SetRGBA(8, 4, colors[4])
	img.SetRGBA(8, 5, colors[5])
	img.SetRGBA(8, 7, colors[6])
	img.SetRGBA(8, 8, colors[7])
	img.SetRGBA(7, 8, colors[8])
	img.SetRGBA(5, 8, colors[9])
	img.SetRGBA(4, 8, colors[10])
	img.SetRGBA(3, 8, colors[11])
	img.SetRGBA(2, 8, colors[12])
	img.SetRGBA(1, 8, colors[13])
	img.SetRGBA(0, 8, colors[14])

	// Write format info (top right)
	img.SetRGBA(img.Rect.Dx()-1, 8, colors[0])
	img.SetRGBA(img.Rect.Dx()-2, 8, colors[1])
	img.SetRGBA(img.Rect.Dx()-3, 8, colors[2])
	img.SetRGBA(img.Rect.Dx()-4, 8, colors[3])
	img.SetRGBA(img.Rect.Dx()-5, 8, colors[4])
	img.SetRGBA(img.Rect.Dx()-6, 8, colors[5])
	img.SetRGBA(img.Rect.Dx()-7, 8, colors[6])
	img.SetRGBA(img.Rect.Dx()-8, 8, colors[7])

	// Write format info (bottom left)
	img.SetRGBA(8, img.Rect.Dy()-7, colors[8])
	img.SetRGBA(8, img.Rect.Dy()-6, colors[9])
	img.SetRGBA(8, img.Rect.Dy()-5, colors[10])
	img.SetRGBA(8, img.Rect.Dy()-4, colors[11])
	img.SetRGBA(8, img.Rect.Dy()-3, colors[12])
	img.SetRGBA(8, img.Rect.Dy()-2, colors[13])
	img.SetRGBA(8, img.Rect.Dy()-1, colors[14])

	// Add version info
	if version >= 7 {
		encodedVersion := ((version << 12) | checkVersion(version<<12))

		// Convert to colors
		var colors []color.RGBA
		for i := 0; i < 18; i++ {
			if encodedVersion&(1<<i) > 0 {
				colors = append(colors, BLACK)
			} else {
				colors = append(colors, WHITE)
			}
		}

		// Write version info (top right)
		var i int
		for y := 0; y < 6; y++ {
			for x := img.Rect.Dx() - 11; x < img.Rect.Dx()-8; x++ {
				img.SetRGBA(x, y, colors[i])
				i++
			}
		}

		// Write version info (bottom left)
		i = 0
		for x := 0; x < 6; x++ {
			for y := img.Rect.Dy() - 11; y < img.Rect.Dy()-8; y++ {
				img.SetRGBA(x, y, colors[i])
				i++
			}
		}
	}
}

func writeData(img *image.RGBA, data []uint8) {
	// Convert data into a series of colors
	var colors []color.RGBA
	for _, v := range data {
		for i := 7; i >= 0; i-- {
			if v&(1<<i) > 0 {
				colors = append(colors, GREEN)
			} else {
				colors = append(colors, RED)
			}
		}
	}

	// Draw in a zig zag pat
	direction := 1
	currentBit := 0
	for x := img.Rect.Dx() - 1; x >= 0; x -= 2 {
		if x == 6 {
			// Skip the vertical timing pattern
			x--
		}

		if direction == 1 {
			// Upwards
			for y := img.Rect.Dy() - 1; y >= 0; y-- {
				// Right module
				if img.RGBAAt(x, y) == BLUE {
					img.SetRGBA(x, y, colors[currentBit])
					currentBit++
				}
				if currentBit >= len(colors) {
					colors = append(colors, RED)
				}

				// Left module
				if img.RGBAAt(x-1, y) == BLUE {
					img.SetRGBA(x-1, y, colors[currentBit])
					currentBit++
				}
				if currentBit >= len(colors) {
					colors = append(colors, RED)
				}
			}

			direction = 0
		} else {
			// Downwards
			for y := 0; y < img.Rect.Dy(); y++ {
				// Right module
				if img.RGBAAt(x, y) == BLUE {
					img.SetRGBA(x, y, colors[currentBit])
					currentBit++
				}
				if currentBit >= len(colors) {
					colors = append(colors, RED)
				}

				// Left module
				if img.RGBAAt(x-1, y) == BLUE {
					img.SetRGBA(x-1, y, colors[currentBit])
					currentBit++
				}
				if currentBit >= len(colors) {
					colors = append(colors, RED)
				}
			}

			direction = 1
		}
	}
}
