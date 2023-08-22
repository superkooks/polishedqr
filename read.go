package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"gocv.io/x/gocv"
)

var windowSegmented *gocv.Window

type QRCodeResult struct {
	ErrorCorrectionLevel string
	CharacterSet         CharacterSet
	Version              int
	Data                 []byte
}

func ReadFromWebcam(displayIntermediates bool) {
	webcam, _ := gocv.VideoCaptureDevice(0)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	window := gocv.NewWindow("Original")
	img := gocv.NewMat()

	if displayIntermediates {
		windowSegmented = gocv.NewWindow("Segmented")
	}

	for {
		webcam.Read(&img)
		window.IMShow(img)
		readQRCode(img, displayIntermediates)
		window.WaitKey(1)
	}
}

func ReadFromImage(img *image.RGBA) {
	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		panic(err)
	}
	windowSegmented = gocv.NewWindow("Segmented")
	_, err = readQRCode(mat, true)
	if err != nil {
		fmt.Println(err)
	}
	for {
		windowSegmented.WaitKey(1)
	}
}

func readQRCode(img gocv.Mat, useWindows bool) (decoded QRCodeResult, err error) {
	// Convert into grayscale
	grayscale := gocv.NewMat()
	gocv.CvtColor(img, &grayscale, gocv.ColorRGBToGray)

	// Find the minimum and maximum reflectance, then threshold the image
	min, max, _, _ := gocv.MinMaxLoc(grayscale)
	thresheld := gocv.NewMat()
	gocv.Threshold(grayscale, &thresheld, (min+max)/2, 255, gocv.ThresholdBinary)

	// Detect edges
	hierarchy := gocv.NewMat()
	contours := gocv.FindContoursWithParams(thresheld, &hierarchy, gocv.RetrievalTree, gocv.ChainApproxSimple)

	// Find nested contours that roughly have areas in the ratio of a finder pattern
	var finderPatterns []gocv.RotatedRect
	var alignmentPatterns []gocv.RotatedRect
	for i := 0; i < contours.Size(); i++ {
		// Get first child contour
		child := int(hierarchy.GetVeciAt(0, i)[2])
		if child < 0 {
			continue
		}

		for {
			// Detect alignment pattern
			area1 := gocv.ContourArea(contours.At(i))
			area2 := gocv.ContourArea(contours.At(child))
			if area2/area1 < 1.0/3.0/3.0+0.05 && area2/area1 > 1.0/3.0/3.0-0.05 {
				// Alignment pattern has good ratio of areas
				pattern := gocv.MinAreaRect(contours.At(i))
				alignmentPatterns = append(alignmentPatterns, pattern)
				gocv.Circle(&img, pattern.Center, 1, color.RGBA{0, 255, 0, 255}, 3)
			}

			// Get first child contour of child countour
			child2 := int(hierarchy.GetVeciAt(0, child)[2])
			if child2 > 0 {
				for {
					area1 := gocv.ContourArea(contours.At(i))
					area2 := gocv.ContourArea(contours.At(child))
					area3 := gocv.ContourArea(contours.At(child2))

					if area3/area1 < 3.0*3.0/7.0/7.0+0.05 && area3/area1 > 3.0*3.0/7.0/7.0-0.05 &&
						area2/area1 < 5.0*5.0/7.0/7.0+0.15 && area2/area1 > 5.0*5.0/7.0/7.0-0.15 {
						// Finder pattern has good ratio of areas
						gocv.DrawContours(&img, contours, i, color.RGBA{255, 0, 0, 255}, 1)
						gocv.DrawContours(&img, contours, child, color.RGBA{255, 0, 0, 255}, 1)
						gocv.DrawContours(&img, contours, child2, color.RGBA{255, 0, 0, 255}, 1)

						pattern := gocv.MinAreaRect(contours.At(i))
						finderPatterns = append(finderPatterns, pattern)
						gocv.Circle(&img, pattern.Center, 1, color.RGBA{0, 255, 0, 255}, 3)
					}

					// Get sibling at previous heirarchy level
					child2 = int(hierarchy.GetVeciAt(0, child2)[0])
					if child2 < 0 {
						break
					}
				}
			}

			// Get sibling at previous heirarchy level
			child = int(hierarchy.GetVeciAt(0, child)[0])
			if child < 0 {
				break
			}
		}
	}

	if len(finderPatterns) < 3 {
		for i := 0; i < contours.Size(); i++ {
			gocv.DrawContours(&img, contours, i, color.RGBA{255, 0, 0, 255}, 2)
		}

		windowSegmented.IMShow(img)
		return QRCodeResult{}, fmt.Errorf("could not find qr code (only %v finder patterns)", len(finderPatterns))
	}

	// Take the first 3 finder patterns, determine the top-left one
	minCosTheta := 1.0
	var topLeft gocv.RotatedRect
	var pat2, pat3 gocv.RotatedRect
	for i := 0; i < 3; i++ {
		fp1 := finderPatterns[i]
		fp2 := finderPatterns[(i+1)%3]
		fp3 := finderPatterns[(i+2)%3]

		// Assuming fp1 is top left, get the line to the two other patterns
		l2 := fp2.Center.Sub(fp1.Center)
		l3 := fp3.Center.Sub(fp1.Center)

		// Find the angle between l2 and l3 using the dot product
		cosTheta := float64(l2.X*l3.X+l2.Y*l3.Y) / vecLen(l2) / vecLen(l3)
		if cosTheta < minCosTheta {
			topLeft = fp1
			minCosTheta = cosTheta

			pat2 = fp2
			pat3 = fp3
		}
	}

	gocv.Circle(&img, topLeft.Center, 1, color.RGBA{0, 0, 255, 255}, 3)

	// Determine the top right finder pattern
	var topRight gocv.RotatedRect
	var bottomLeft gocv.RotatedRect
	{
		l2 := pat2.Center.Sub(topLeft.Center)
		l3 := pat3.Center.Sub(topLeft.Center)

		// Find the z coordinate of the cross product
		z := l2.X*l3.Y - l2.Y*l3.X
		if z > 0 {
			topRight = pat2
			bottomLeft = pat3
		} else {
			topRight = pat3
			bottomLeft = pat2
		}
	}

	gocv.Circle(&img, topRight.Center, 1, color.RGBA{255, 0, 0, 255}, 3)

	// Find the version of the symbol
	xDim := float64(topLeft.Width+topRight.Width) / 14
	version := int(math.Round((vecLen(topRight.Center.Sub(topLeft.Center))/xDim - 10) / 4))

	if version < 1 {
		return QRCodeResult{}, errors.New("unable to determine provisional version")
	}

	if version > 6 {
		// We have to check the version information itself
		modSize := float64(topRight.Width) / 7

		// Find the unit module vector in the direction of +X
		l1 := topRight.Center.Sub(topLeft.Center)
		x1 := float64(l1.X) / vecLen(l1) * modSize
		y1 := float64(l1.Y) / vecLen(l1) * modSize

		// Find the unit module vector in the direction of +Y
		l2 := bottomLeft.Center.Sub(topLeft.Center)
		x2 := float64(l2.X) / vecLen(l2) * modSize
		y2 := float64(l2.Y) / vecLen(l2) * modSize

		// Read the version information
		var versionBits int
		var i int
		for y := -3.0; y < 3; y++ {
			for x := -7.0; x < -4; x++ {
				pt := topRight.Center.Add(image.Pt(int(x1*x+x2*y), int(y1*x+y2*y)))
				gocv.Circle(&img, pt, 0, color.RGBA{255, 0, 0, 255}, 1)

				if thresheld.GetUCharAt(pt.Y, pt.X) == 0 {
					versionBits |= (1 << i)
				}
				i++
			}
		}

		version = decodeVersion(versionBits)
		if version < 0 {
			// Technically we could go read the second version information block.
			// But, I am lazy.
			return QRCodeResult{}, errors.New("unable to decode version information")
		}
	}

	// Sample into image
	i := image.NewRGBA(image.Rect(0, 0, version*4+17, version*4+17))
	if version > 1 {
		// Find the bottom-rightmost alignment pattern for versions > 1
		modSizeX := vecLen(topRight.Center.Sub(topLeft.Center)) / float64(version*4+10)
		modSizeY := vecLen(bottomLeft.Center.Sub(topLeft.Center)) / float64(version*4+10)

		// Find the unit module vector in the direction of +X
		l1 := topRight.Center.Sub(topLeft.Center)
		x1 := float64(l1.X) / vecLen(l1) * modSizeX
		y1 := float64(l1.Y) / vecLen(l1) * modSizeY

		// Find the unit module vector in the direction of +Y
		l2 := bottomLeft.Center.Sub(topLeft.Center)
		x2 := float64(l2.X) / vecLen(l2) * modSizeX
		y2 := float64(l2.Y) / vecLen(l2) * modSizeY

		// Get the top-leftmost module
		pointA := topLeft.Center.Sub(image.Pt(int(x1*3+x2*3), int(y1*3+y2*3)))
		gocv.Circle(&img, pointA, 0, color.RGBA{255, 0, 0, 255}, 1)

		// Get the provisional position of the alignment pattern
		standardAligns := getAlignmentPositions(version)
		bottomStandardAlign := standardAligns[len(standardAligns)-1]

		provisional := pointA.Add(image.Pt(int(x1*float64(bottomStandardAlign[0])+x2*float64(bottomStandardAlign[1])), int(y1*float64(bottomStandardAlign[0])+y2*float64(bottomStandardAlign[1]))))

		// Find closest alignment pattern to provisional centre
		var minPattern gocv.RotatedRect
		minDist := math.MaxFloat64
		for _, v := range alignmentPatterns {
			diff := provisional.Sub(v.Center)
			dist := math.Sqrt(math.Pow(float64(diff.X)/(x1+x2)/2, 2) + math.Pow(float64(diff.Y)/(y1+y2)/2, 2))
			if dist < minDist {
				minDist = dist
				minPattern = v
			}
		}
		gocv.Circle(&img, minPattern.Center, 1, color.RGBA{255, 0, 0, 255}, 3)

		// Generate a perspective transform from the 4 points
		src := gocv.NewPointVectorFromPoints([]image.Point{
			topLeft.Center, topRight.Center, bottomLeft.Center, minPattern.Center,
		})
		dst := gocv.NewPointVectorFromPoints([]image.Point{
			{35, 35},
			{version*40 + 135, 35},
			{35, version*40 + 135},
			{bottomStandardAlign[0]*10 + 5, bottomStandardAlign[1]*10 + 5},
		})
		transform := gocv.GetPerspectiveTransform(src, dst)

		// Warp the image with matrix
		warped := gocv.NewMat()
		gocv.WarpPerspective(thresheld, &warped, transform, image.Pt(version*40+170, version*40+170))

		warpedColor := gocv.NewMat()
		gocv.CvtColor(warped, &warpedColor, gocv.ColorGrayToBGR)

		for x := 0.0; x < float64(i.Rect.Dx()); x++ {
			for y := 0.0; y < float64(i.Rect.Dy()); y++ {
				pt := image.Pt(int(10*x+5), int(10*y+5))
				gocv.Circle(&warpedColor, pt, 0, color.RGBA{255, 0, 0, 255}, 1)

				if warped.GetUCharAt(pt.Y, pt.X) == 0 {
					i.SetRGBA(int(x), int(y), BLACK)
				} else {
					i.SetRGBA(int(x), int(y), WHITE)
				}

			}
		}

		windowSegmented.IMShow(warpedColor)

	} else {
		modSizeX := vecLen(topRight.Center.Sub(topLeft.Center)) / float64(version*4+10)
		modSizeY := vecLen(bottomLeft.Center.Sub(topLeft.Center)) / float64(version*4+10)

		// Find the unit module vector in the direction of +X
		l1 := topRight.Center.Sub(topLeft.Center)
		x1 := float64(l1.X) / vecLen(l1) * modSizeX
		y1 := float64(l1.Y) / vecLen(l1) * modSizeY

		// Find the unit module vector in the direction of +Y
		l2 := bottomLeft.Center.Sub(topLeft.Center)
		x2 := float64(l2.X) / vecLen(l2) * modSizeX
		y2 := float64(l2.Y) / vecLen(l2) * modSizeY

		// The most top left module
		pointA := topLeft.Center.Sub(image.Pt(int(x1*3+x2*3), int(y1*3+y2*3)))
		gocv.Circle(&img, pointA, 0, color.RGBA{255, 0, 0, 255}, 1)

		// Sample every module
		for x := 0.0; x < float64(i.Rect.Dx()); x++ {
			for y := 0.0; y < float64(i.Rect.Dy()); y++ {
				pt := pointA.Add(image.Pt(int(x1*x+x2*y), int(y1*x+y2*y)))
				gocv.Circle(&img, pt, 0, color.RGBA{255, 0, 0, 255}, 1)

				if thresheld.GetUCharAt(pt.Y, pt.X) == 0 {
					i.SetRGBA(int(x), int(y), BLACK)
				} else {
					i.SetRGBA(int(x), int(y), WHITE)
				}

			}
		}
	}

	// Get format info
	var ecLevel string
	var maskPattern int
	{
		// Convert into series of colours
		var colors []color.RGBA
		colors = append(colors, i.RGBAAt(8, 0))
		colors = append(colors, i.RGBAAt(8, 1))
		colors = append(colors, i.RGBAAt(8, 2))
		colors = append(colors, i.RGBAAt(8, 3))
		colors = append(colors, i.RGBAAt(8, 4))
		colors = append(colors, i.RGBAAt(8, 5))
		colors = append(colors, i.RGBAAt(8, 7))
		colors = append(colors, i.RGBAAt(8, 8))
		colors = append(colors, i.RGBAAt(7, 8))
		colors = append(colors, i.RGBAAt(5, 8))
		colors = append(colors, i.RGBAAt(4, 8))
		colors = append(colors, i.RGBAAt(3, 8))
		colors = append(colors, i.RGBAAt(2, 8))
		colors = append(colors, i.RGBAAt(1, 8))
		colors = append(colors, i.RGBAAt(0, 8))

		// Convert colors into bits
		var formatBits int
		for k, v := range colors {
			if v == BLACK {
				formatBits |= (1 << k)
			}
		}

		// Apply EC and convert
		formatBits = decodeFormat(formatBits ^ 0b101010000010010)

		switch formatBits >> 3 {
		case 0:
			ecLevel = "M"
		case 1:
			ecLevel = "L"
		case 2:
			ecLevel = "H"
		case 3:
			ecLevel = "Q"
		default:
			return QRCodeResult{}, errors.New("unable to determine ec level")
		}

		maskPattern = formatBits & 0b111
	}

	f, _ := os.Create("test.png")
	png.Encode(f, i)

	// Mask off fixed patterns
	{
		// Finder patterns + format information
		iterateRect(9, 9, func(x, y int) {
			i.SetRGBA(x, y, BLUE)
		})
		iterateRect(9, 8, func(x, y int) {
			i.SetRGBA(x, i.Rect.Dy()-1-y, BLUE)
		})
		iterateRect(8, 9, func(x, y int) {
			i.SetRGBA(i.Rect.Dx()-1-x, y, BLUE)
		})

		// Timing patterns
		iterateRect(1, i.Rect.Dy(), func(x, y int) {
			i.SetRGBA(6, y, BLUE)
		})
		iterateRect(i.Rect.Dx(), 1, func(x, y int) {
			i.SetRGBA(x, 6, BLUE)
		})

		// Alignment pattersn
		alignments := getAlignmentPositions(version)
		for _, v := range alignments {
			iterateRect(5, 5, func(x, y int) {
				i.SetRGBA(v[0]+x-2, v[1]+y-2, BLUE)
			})
		}

		if version > 6 {
			// Version information
			iterateRect(3, 6, func(x, y int) {
				i.SetRGBA(i.Rect.Dx()-11+x, y, BLUE)
			})
			iterateRect(6, 3, func(x, y int) {
				i.SetRGBA(x, i.Rect.Dy()-11+y, BLUE)
			})
		}
	}

	// Read data
	var data []uint8
	{
		// Read data in zig zag path
		var colors []uint8
		direction := 1
		for x := i.Rect.Dx() - 1; x >= 0; x -= 2 {
			if x == 6 {
				// Skip the vertical timing pattern
				x--
			}

			if direction == 1 {
				// Upwards
				for y := i.Rect.Dy() - 1; y >= 0; y-- {
					// Right module
					if i.RGBAAt(x, y) != BLUE {
						var h uint8
						if Masks[maskPattern](x, y) {
							h = 1
						}

						if i.RGBAAt(x, y) == BLACK {
							colors = append(colors, h^1)
						} else {
							colors = append(colors, h)
						}
					}

					// Left module
					if i.RGBAAt(x-1, y) != BLUE {
						var h uint8
						if Masks[maskPattern](x-1, y) {
							h = 1
						}

						if i.RGBAAt(x-1, y) == BLACK {
							colors = append(colors, h^1)
						} else {
							colors = append(colors, h)
						}
					}
				}

				direction = 0
			} else {
				// Downwards
				for y := 0; y < i.Rect.Dy(); y++ {
					// Right module
					if i.RGBAAt(x, y) != BLUE {
						var h uint8
						if Masks[maskPattern](x, y) {
							h = 1
						}

						if i.RGBAAt(x, y) == BLACK {
							colors = append(colors, h^1)
						} else {
							colors = append(colors, h)
						}
					}

					// Left module
					if i.RGBAAt(x-1, y) != BLUE {
						var h uint8
						if Masks[maskPattern](x-1, y) {
							h = 1
						}

						if i.RGBAAt(x-1, y) == BLACK {
							colors = append(colors, h^1)
						} else {
							colors = append(colors, h)
						}
					}
				}

				direction = 1
			}
		}

		// Convert colors into bits
		for i := 0; i < len(colors)-7; i += 8 {
			var byt uint8
			for j := 0; j < 8; j++ {
				if colors[i+j] == 1 {
					byt |= (1 << (7 - j))
				}
			}
			data = append(data, byt)
		}
	}

	// Split codewords into blocks & error correct
	datawords, err := correctDataWords(data, version, ecLevel)
	if err != nil {
		return QRCodeResult{}, err
	}

	// Convert back into bits
	var bits Bits
	for _, v := range datawords {
		for i := 7; i >= 0; i-- {
			if v&(1<<i) > 0 {
				bits = append(bits, 1)
			} else {
				bits = append(bits, 0)
			}
		}
	}

	// Decode the character set
	// NB: We only support one character set per symbol
	var decodedData []byte
	var charset CharacterSet
	if bytes.Equal(bits[:4], Bits{0, 0, 0, 1}) {
		// Numeric
		charset = Numeric
		charBits := CharacterCountBitCapacity(Numeric, version)

		// Get character count
		var charCount int
		for k, v := range bits[4 : 4+charBits] {
			if v > 0 {
				charCount |= 1 << (charBits - 1 - k)
			}
		}

		// Divide into groups of 3 digits (or less) and convert
		for i := 0; i < charCount; i += 3 {
			bc := 10
			if charCount-i == 2 {
				bc = 7
			} else if charCount-i == 1 {
				bc = 4
			}

			var group int
			for j := 0; j < bc; j++ {
				if bits[4+charBits+i/3*10+j] > 0 {
					group |= 1 << (bc - 1 - j)
				}
			}

			digits := fmt.Sprintf("%03v", group)
			if charCount-i == 2 {
				decodedData = append(decodedData, digits[1])
				decodedData = append(decodedData, digits[2])
			} else if charCount-i == 1 {
				decodedData = append(decodedData, digits[2])
			} else {
				decodedData = append(decodedData, digits[0])
				decodedData = append(decodedData, digits[1])
				decodedData = append(decodedData, digits[2])
			}
		}

	} else if bytes.Equal(bits[:4], Bits{0, 0, 1, 0}) {
		// Alphanumeric
		charset = Alphanumeric
		charBits := CharacterCountBitCapacity(Alphanumeric, version)

		// Get character count
		var charCount int
		for k, v := range bits[4 : 4+charBits] {
			if v > 0 {
				charCount |= 1 << (charBits - 1 - k)
			}
		}

		// Divide into groups of 2 characters (or less) and convert
		for i := 0; i < charCount; i += 2 {
			bc := 11
			if charCount-i == 1 {
				bc = 6
			}

			var group int
			for j := 0; j < bc; j++ {
				if bits[4+charBits+i/2*11+j] > 0 {
					group |= 1 << (bc - 1 - j)
				}
			}

			if charCount-i == 1 {
				decodedData = append(decodedData, byte(alphanumericTableReverse[group]))
			} else {
				decodedData = append(decodedData, byte(alphanumericTableReverse[group/45]))
				decodedData = append(decodedData, byte(alphanumericTableReverse[group%45]))
			}
		}

	} else if bytes.Equal(bits[:4], Bits{0, 1, 0, 0}) {
		// Bytes
		charset = Bytes
		charBits := CharacterCountBitCapacity(Bytes, version)

		// Get character count
		var charCount int
		for k, v := range bits[4 : 4+charBits] {
			if v > 0 {
				charCount |= 1 << (charBits - 1 - k)
			}
		}

		// Convert bits to bytes
		for i := 0; i < charCount; i++ {
			var byt byte
			for j := 0; j < 8; j++ {
				if bits[4+charBits+i*8+j] > 0 {
					byt |= 1 << (7 - j)
				}
			}

			decodedData = append(decodedData, byt)
		}

	} else {
		return QRCodeResult{}, errors.New("unknown character set")
	}

	fmt.Println(ecLevel, maskPattern)
	fmt.Println("version", version)

	if useWindows {
		// mat, err := gocv.ImageToMatRGB()
		// if err != nil {
		// 	panic(err)
		// }
		windowSegmented.IMShow(img)
	}

	fmt.Println("decoded data:", string(decodedData))

	return QRCodeResult{
		ErrorCorrectionLevel: ecLevel,
		CharacterSet:         charset,
		Version:              version,
		Data:                 decodedData,
	}, nil
}
