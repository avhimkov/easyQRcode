package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

//QRreader - Encode qrcode
func QRreader(path string) error {

	// open and decode image file JPG
	fileEncode, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	img, _, _ := image.Decode(fileEncode)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)

	fmt.Println(result)

	return err
}
