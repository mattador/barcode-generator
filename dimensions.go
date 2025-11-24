package barcode

import (
	"image"
	"math"

	"github.com/boombuler/barcode"
)

// Constants for label layout
const labelMarginPixels = 10

// mmToPixels converts millimeters to pixels based on the printer DPI.
// Formula: pixels = mm * dpi / 25.4 (25.4 mm per inch)
func mmToPixels(mm float64, dpi int) int {
	return int(mm * float64(dpi) / 25.4)
}

// calculateBarcodeSize determines the appropriate barcode dimensions based on type.
// Code128: Uses full width, constrained height
// QR: Must be square, sized to fit with text
func calculateBarcodeSize(input BarcodeInput, labelWidth, labelHeight int) image.Point {
	if input.BarcodeType == BarcodeTypeCode128 {
		return calculateCode128Size(labelWidth, labelHeight)
	}
	return calculateQRSize(input, labelWidth, labelHeight)
}

// calculateCode128Size determines dimensions for Code128 barcodes.
// Code128 can be rectangular, so we use full label width and constrain height.
func calculateCode128Size(labelWidth, labelHeight int) image.Point {
	barcodeWidth := labelWidth - (labelMarginPixels * 2)
	barcodeHeight := int(math.Min(float64(labelHeight/2), 200))
	return image.Pt(barcodeWidth, barcodeHeight)
}

// calculateQRSize determines dimensions for QR codes.
// QR codes must be square, so we calculate the largest square that fits.
func calculateQRSize(input BarcodeInput, labelWidth, labelHeight int) image.Point {
	// Start with the smaller of width or height
	maxSize := int(math.Min(float64(labelWidth), float64(labelHeight)))

	// Calculate space needed for text
	textHeight := calculateTextHeight(input)

	// Reduce available space for text
	availableHeight := float64(labelHeight) - textHeight
	finalSize := int(math.Min(float64(maxSize), availableHeight))

	return image.Pt(finalSize, finalSize)
}

// calculateTextHeight returns the total pixel height needed for all text lines.
func calculateTextHeight(input BarcodeInput) float64 {
	totalHeight := 0.0
	for _, textLine := range input.TextLines {
		_, height := getFontSize(textLine.Size, input.Dpi, 200)
		totalHeight += height * 2
	}
	return totalHeight
}

// scaleBarcodeToFit resizes a barcode to the specified dimensions.
func scaleBarcodeToFit(bc barcode.Barcode, size image.Point) (barcode.Barcode, error) {
	scaled, err := barcode.Scale(bc, size.X, size.Y)
	if err != nil {
		return nil, err
	}
	return scaled, nil
}

// centerBarcodeOnLabel calculates the position to center a barcode on the label.
// Returns the bounding rectangle where the barcode should be drawn.
func centerBarcodeOnLabel(img *image.RGBA, bc barcode.Barcode) image.Rectangle {
	imgBounds := img.Bounds()
	bcBounds := bc.Bounds()

	offsetX := (imgBounds.Dx() - bcBounds.Dx()) / 2
	offsetY := (imgBounds.Dy() - bcBounds.Dy()) / 2

	return bcBounds.Add(image.Pt(offsetX, offsetY))
}

// calculateTextYPosition determines the Y coordinate for text based on position relative to barcode.
func calculateTextYPosition(barcodeRect image.Rectangle, position TextPosition) int {
	if position == TextPositionAbove {
		return barcodeRect.Min.Y
	}
	return barcodeRect.Max.Y
}
