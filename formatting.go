package barcode

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"

	"simonwaldherr.de/go/zplgfa"
)

// imageToBase64 converts an image to a base64-encoded PNG string.
// This allows the image to be easily transmitted in JSON or HTML data URLs.
func imageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// imageToZPL converts an image to ZPL (Zebra Programming Language) commands.
// ZPL is the standard language for Zebra thermal printers.
// The conversion uses image flattening and ASCII compression for efficiency.
func imageToZPL(img image.Image) string {
	// Convert to RGBA if needed
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		bounds := img.Bounds()
		rgbaImg = image.NewRGBA(bounds)
		// Copy image data (simplified - in production might need more robust conversion)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
	}

	flat := zplgfa.FlattenImage(rgbaImg)
	return zplgfa.ConvertToZPL(flat, zplgfa.CompressedASCII)
}
