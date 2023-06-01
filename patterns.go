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

func quietZone(i *image.RGBA) *image.RGBA {
	n := image.NewRGBA(image.Rect(0, 0, i.Rect.Dx()+8, i.Rect.Dy()+8))
	iterateRect(n.Rect.Dx(), n.Rect.Dy(), func(x, y int) {
		n.SetRGBA(x, y, WHITE)
	})

	draw.Draw(n, image.Rect(4, 4, i.Rect.Dx()+4, i.Rect.Dy()+4), i, image.Pt(0, 0), draw.Over)

	return n
}

func writeData(img *image.RGBA, data []uint8) {
	// Convert data into a series of colors
	var colors []color.RGBA
	for _, v := range data {
		for i := 7; i >= 0; i-- {
			if v&(1<<i) > 0 {
				colors = append(colors, BLACK)
			} else {
				colors = append(colors, WHITE)
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
					colors = append(colors, WHITE)
				}

				// Left module
				if img.RGBAAt(x-1, y) == BLUE {
					img.SetRGBA(x-1, y, colors[currentBit])
					currentBit++
				}
				if currentBit >= len(colors) {
					colors = append(colors, WHITE)
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
					colors = append(colors, WHITE)
				}

				// Left module
				if img.RGBAAt(x-1, y) == BLUE {
					img.SetRGBA(x-1, y, colors[currentBit])
					currentBit++
				}
				if currentBit >= len(colors) {
					colors = append(colors, WHITE)
				}
			}

			direction = 1
		}
	}
}
