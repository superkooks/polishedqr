package main

func main() {
	// f, err := os.Create("out.png")
	// if err != nil {
	// 	panic(err)
	// }

	// err = png.Encode(f, CreateQRCode([]byte("AC-42"), &CreateOptions{
	// 	ErrorCorrectionLevel: "L",
	// }))
	// if err != nil {
	// 	panic(err)
	// }

	// f, err := os.Open("image7.png")
	// if err != nil {
	// 	panic(err)
	// }

	// img, err := png.Decode(f)
	// if err != nil {
	// 	panic(err)
	// }

	// ReadFromImage(img.(*image.RGBA))

	ReadFromWebcam(true)
}
