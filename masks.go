package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

var Masks = []func(int, int) bool{mask1, mask2, mask3, mask4, mask5, mask6, mask7, mask8}

func mask1(x, y int) bool {
	return (x+y)%2 == 0
}
func mask2(x, y int) bool {
	return y%2 == 0
}
func mask3(x, y int) bool {
	return x%3 == 0
}
func mask4(x, y int) bool {
	return (x+y)%3 == 0
}
func mask5(x, y int) bool {
	return (y/2+x/3)%2 == 0
}
func mask6(x, y int) bool {
	return (x*y)%2+(x*y)%3 == 0
}
func mask7(x, y int) bool {
	return ((x*y)%2+(x*y)%3)%2 == 0
}
func mask8(x, y int) bool {
	return (((x+y)%2)+((x*y)%3))%2 == 0
}

func applyBestMask(img *image.RGBA, ecLevel string, version int) int {
	lowestPenalty := math.MaxInt
	bestMask := 0
	for k, v := range Masks {
		// Make a copy of the image
		masked := image.NewRGBA(img.Rect)
		draw.Draw(masked, img.Rect, img, image.Pt(0, 0), draw.Over)

		// Add format info
		addFormatAndVersionInfo(masked, ecLevel, k, version)

		// Apply the mask
		applyMask(masked, v)

		// Determine penalty
		p := determinePenalty(masked)
		fmt.Println("mask pen", k, p)
		if p < lowestPenalty {
			lowestPenalty = p
			bestMask = k
		}
	}

	applyMask(img, Masks[bestMask])

	return bestMask
}

func applyMask(img *image.RGBA, mask func(int, int) bool) {
	for x := 0; x < img.Rect.Dx(); x++ {
		for y := 0; y < img.Rect.Dy(); y++ {
			if mask(x, y) {
				// Flip it
				if img.RGBAAt(x, y) == GREEN {
					img.Set(x, y, WHITE)
				} else if img.RGBAAt(x, y) == RED {
					img.Set(x, y, BLACK)
				}
			} else {
				if img.RGBAAt(x, y) == RED {
					img.Set(x, y, WHITE)
				} else if img.RGBAAt(x, y) == GREEN {
					img.Set(x, y, BLACK)
				}
			}
		}
	}
}

func determinePenalty(masked *image.RGBA) int {
	var penalty int

	// Evaluation condition 1a (rows)
	for y := 0; y < masked.Rect.Dy(); y++ {
		consecutiveCount := 0
		consecutiveBit := true
		for x := 0; x < masked.Rect.Dx(); x++ {
			if masked.RGBAAt(x, y) == WHITE {
				if consecutiveBit {
					consecutiveBit = false
					consecutiveCount = 0
				}
				consecutiveCount++
			} else {
				if !consecutiveBit {
					consecutiveBit = true
					consecutiveCount = 0
				}
				consecutiveCount++
			}

			if consecutiveCount == 5 {
				penalty += 3
			} else if consecutiveCount > 5 {
				penalty++
			}
		}
	}

	// Evaluation condition 1b (columns)
	for x := 0; x < masked.Rect.Dx(); x++ {
		consecutiveCount := 0
		consecutiveBit := true
		for y := 0; y < masked.Rect.Dy(); y++ {
			if masked.RGBAAt(x, y) == WHITE {
				if consecutiveBit {
					consecutiveBit = false
					consecutiveCount = 0
				}
				consecutiveCount++
			} else {
				if !consecutiveBit {
					consecutiveBit = true
					consecutiveCount = 0
				}
				consecutiveCount++
			}

			if consecutiveCount == 5 {
				penalty += 3

			} else if consecutiveCount > 5 {
				penalty++
			}
		}
	}

	// Evaluation condition 2 (squares)
	for y := 0; y < masked.Rect.Dy()-1; y++ {
		for x := 0; x < masked.Rect.Dx()-1; x++ {
			// Get the 2x2 area with (x,y) as the top-left module
			modules := [][2]int{{x + 1, y}, {x, y + 1}, {x + 1, y + 1}}

			bit := true
			if masked.RGBAAt(x, y) == WHITE {
				bit = false
			}

			contiguous := true
			for _, v := range modules {
				if (masked.RGBAAt(v[0], v[1]) != BLACK && bit) ||
					(masked.RGBAAt(v[0], v[1]) != WHITE && !bit) {
					contiguous = false
					break
				}
			}

			if contiguous {
				penalty += 3
			}
		}
	}

	// Evaluation condition 3 (similar to finder pattern)
	for y := 0; y < masked.Rect.Dy(); y++ {
		for x := 0; x < masked.Rect.Dx()-6; x++ {
			// Check in the horizontal direction
			if masked.RGBAAt(x, y) == BLACK &&
				masked.RGBAAt(x+1, y) == WHITE &&
				masked.RGBAAt(x+2, y) == BLACK &&
				masked.RGBAAt(x+3, y) == BLACK &&
				masked.RGBAAt(x+4, y) == BLACK &&
				masked.RGBAAt(x+5, y) == WHITE &&
				masked.RGBAAt(x+6, y) == BLACK {

				// Check whether there is 4 white spaces on either side
				if checkFourPlusOneWhite(masked, x-4, y, x-1, y, x+7, y) {
					penalty += 40
				}
				if checkFourPlusOneWhite(masked, x+7, y, x+10, y, x-1, y) {
					penalty += 40
				}
			}
		}
	}

	for x := 0; x < masked.Rect.Dx(); x++ {
		for y := 0; y < masked.Rect.Dy()-6; y++ {
			// Check in the vertical direction
			if masked.RGBAAt(x, y) == BLACK &&
				masked.RGBAAt(x, y+1) == WHITE &&
				masked.RGBAAt(x, y+2) == BLACK &&
				masked.RGBAAt(x, y+3) == BLACK &&
				masked.RGBAAt(x, y+4) == BLACK &&
				masked.RGBAAt(x, y+5) == WHITE &&
				masked.RGBAAt(x, y+6) == BLACK {

				if checkFourPlusOneWhite(masked, x, y-4, x, y-1, x, y+7) {
					penalty += 40
				}
				if checkFourPlusOneWhite(masked, x, y+7, x, y+10, x, y-1) {
					penalty += 40
				}
			}
		}
	}

	// Evaluation condition 4 (white-dark module ratio)
	total := 0
	dark := 0
	for y := 0; y < masked.Rect.Dy(); y++ {
		for x := 0; x < masked.Rect.Dx(); x++ {
			if masked.RGBAAt(x, y) == BLACK {
				dark++
			}

			total++
		}
	}

	ratio := float64(dark) / float64(total)
	if ratio < 0.5 {
		dev := 0.5 - ratio
		penalty += int(dev*100) / 5 * 10
	} else if ratio > 0.55 {
		dev := ratio - 0.5
		penalty += (int(dev*100)/5 - 1) * 10
	}

	return penalty
}

func checkFourPlusOneWhite(img *image.RGBA, x1, y1, x2, y2, x3, y3 int) bool {
	var foundBlack bool
	oob := color.RGBA{}
	iterateRect(x2-x1+1, y2-y1+1, func(x, y int) {
		c := img.RGBAAt(x1+x, y1+y)
		if c != WHITE && c != oob {
			foundBlack = true
		}
	})

	return !foundBlack && (img.RGBAAt(x3, y3) == WHITE || img.RGBAAt(x3, y3) == oob)
}
