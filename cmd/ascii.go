package main

import (
	"fmt"
	"image"
)

const lowHalf = "▄"
const upperHalf = "▀"
const empty = " "
const full = "█"

func PrintQRCodeASCII(img *image.RGBA) {
	// Set foreground and background
	fmt.Print("\033[37;40m")

	// Print the qr code using block characters
	for y := 0; y < img.Rect.Dy(); y += 2 {
		for x := 0; x < img.Rect.Dx(); x++ {
			if img.RGBAAt(x, y).R == 255 {
				if img.RGBAAt(x, y+1).R == 255 {
					fmt.Print(full)
				} else {
					fmt.Print(upperHalf)
				}
			} else {
				if img.RGBAAt(x, y+1).R == 255 {
					fmt.Print(lowHalf)
				} else {
					fmt.Print(empty)
				}
			}
		}

		// ANSI reset
		fmt.Println("\033[0m")
	}
}
