package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/superkooks/polishedqr"
	"github.com/urfave/cli/v2"
	"gocv.io/x/gocv"

	"image/draw"
	_ "image/jpeg"
	"image/png"
)

func writeOut(outPath string, data io.Reader) {
	if outPath == "-" || outPath == "" {
		io.Copy(os.Stdout, data)
		fmt.Println()
	} else {
		f, err := os.Create(outPath)
		if err != nil {
			panic(err)
		}

		io.Copy(f, data)

		f.Close()
	}
}

func main() {
	app := &cli.App{
		Name:  "polishedqr",
		Usage: "create and read qr codes",
		Flags: []cli.Flag{cli.BashCompletionFlag},
		Commands: []*cli.Command{
			{
				Name:      "create",
				Aliases:   []string{"c"},
				Usage:     "create a qr code",
				ArgsUsage: "infile",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:  "out",
						Usage: "the path to save the decoded data",
					},
					&cli.IntFlag{
						Name:        "version",
						Usage:       "the version (size) of the qr code to generate",
						DefaultText: "0 (auto)",
					},
					&cli.StringFlag{
						Name:        "ec",
						Usage:       "the error correction level to use (one of L, M, Q, H)",
						DefaultText: "M",
						Value:       "M",
					},
					&cli.Float64Flag{
						Name:        "scale",
						Usage:       "the factor to scale the output image by",
						DefaultText: "1",
						Value:       1,
					},
				},

				Action: func(ctx *cli.Context) error {
					inPath := ctx.Args().First()
					if inPath == "" {
						panic("no input file")
					}

					var inReader io.Reader
					if inPath == "-" {
						inReader = os.Stdin
					} else {
						var err error
						inReader, err = os.Open(inPath)
						if err != nil {
							panic(err)
						}
					}

					b, err := io.ReadAll(inReader)
					if err != nil {
						panic(err)
					}

					opts := &polishedqr.CreateOptions{
						ErrorCorrectionLevel: ctx.String("ec"),
						Version:              ctx.Int("version"),
					}
					img := polishedqr.CreateQRCode(b, opts)

					if ctx.Float64("scale") != 1 {
						mat, err := gocv.ImageToMatRGBA(img)
						if err != nil {
							panic(err)
						}
						gocv.Resize(mat, &mat, image.Pt(0, 0), ctx.Float64("scale"), ctx.Float64("scale"), gocv.InterpolationNearestNeighbor)

						i, err := mat.ToImage()
						if err != nil {
							panic(err)
						}

						r := image.Rect(
							0, 0,
							int(float64(img.Rect.Dx())*ctx.Float64("scale")),
							int(float64(img.Rect.Dy())*ctx.Float64("scale")),
						)

						img = image.NewRGBA(r)
						draw.Draw(img, r, i, image.Pt(0, 0), draw.Over)
					}

					if ctx.Path("out") == "" {
						PrintQRCodeASCII(img)
					} else {
						pipeR, pipeW := io.Pipe()

						go writeOut(ctx.Path("out"), pipeR)

						err = png.Encode(pipeW, img)
						if err != nil {
							panic(err)
						}
						pipeW.Close()
					}

					return nil
				},
			},
			{
				Name:      "read",
				Aliases:   []string{"r"},
				Usage:     "read a qr code from an image",
				ArgsUsage: "image",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:  "out",
						Usage: "the path to save the decoded data",
					},
				},

				Action: func(ctx *cli.Context) error {
					inPath := ctx.Args().First()
					if inPath == "" {
						panic("no file input")
					}

					var inReader io.Reader
					if inPath == "-" {
						inReader = os.Stdin
					} else {
						var err error
						inReader, err = os.Open(inPath)
						if err != nil {
							panic(err)
						}
					}

					i, _, err := image.Decode(inReader)
					if err != nil {
						panic(err)
					}

					var rgba *image.RGBA
					var ok bool
					if rgba, ok = i.(*image.RGBA); !ok {
						rgba = image.NewRGBA(i.Bounds())
						draw.Draw(rgba, i.Bounds(), i, i.Bounds().Min, draw.Src)
					}

					result, err := polishedqr.ReadFromImage(rgba)
					if err != nil {
						panic(fmt.Errorf("error decoding qr code: %v", err))
					}

					fmt.Printf(
						"detected version %v code with error correction %v\n\n",
						result.Version,
						result.ErrorCorrectionLevel,
					)

					writeOut(ctx.Path("out"), bytes.NewBuffer(result.Data))

					return nil
				},
			},
			{
				Name:      "webcam",
				Aliases:   []string{"w"},
				Usage:     "read a qr code from the webcam",
				ArgsUsage: " ",
				Action: func(ctx *cli.Context) error {
					result, err := polishedqr.ReadFromWebcam(true)
					if err != nil {
						panic(err)
					}

					fmt.Printf(
						"detected version %v code with error correction %v\n\n",
						result.Version,
						result.ErrorCorrectionLevel,
					)

					writeOut(ctx.Path("out"), bytes.NewBuffer(result.Data))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
