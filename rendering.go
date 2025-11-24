package barcode

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/boombuler/barcode"
)

// createBlankLabel initializes a white RGBA image for the label.
func createBlankLabel(width, height int) *image.RGBA {
	bounds := image.Rect(0, 0, width, height)
	img := image.NewRGBA(bounds)

	// Fill with white background
	draw.Draw(img, bounds, &image.Uniform{color.White}, image.Point{}, draw.Src)

	return img
}

// drawBarcodeOnLabel composites a barcode image onto the label at the specified position.
func drawBarcodeOnLabel(label *image.RGBA, barcode barcode.Barcode, position image.Rectangle) {
	draw.Draw(label, position, barcode, barcode.Bounds().Min, draw.Over)
}
