/*
Package barcode provides barcode and label generation for warehouse operations.

Supports Code128 and QR code formats with dual output:
  - PNG images (base64-encoded) for web display
  - ZPL (Zebra Programming Language) for thermal printer output

Key features:
  - DPI-aware scaling for standard thermal printers (203, 300, 600 DPI)
  - Automatic text positioning and font sizing
  - Recursive font scaling to fit text on labels
*/
package barcode

import (
	"fmt"
	"image"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
)

// Standard DPI values supported by most thermal printers
var standardDPIValues = []int{203, 300, 600}

// BarcodeType represents supported barcode formats
type BarcodeType string

const (
	BarcodeTypeCode128 BarcodeType = "CODE128"
	BarcodeTypeQR      BarcodeType = "QR"
)

// TextPosition defines where text appears relative to the barcode
type TextPosition string

const (
	TextPositionAbove TextPosition = "ABOVE"
	TextPositionBelow TextPosition = "BELOW"
)

// TextSize defines predefined text sizes
type TextSize string

const (
	TextSizeSmall  TextSize = "SMALL"
	TextSizeMedium TextSize = "MEDIUM"
	TextSizeLarge  TextSize = "LARGE"
)

// TextLine represents a line of text to render on the label
type TextLine struct {
	Text     string
	Position TextPosition
	Size     TextSize
}

// BarcodeInput contains all parameters needed to generate a barcode label
type BarcodeInput struct {
	BarcodeData string      // The data to encode in the barcode
	BarcodeType BarcodeType // Type of barcode (CODE128 or QR)
	Width       float64     // Label width in millimeters
	Height      float64     // Label height in millimeters
	Dpi         int         // Printer DPI (203, 300, or 600)
	TextLines   []TextLine  // Optional text lines to render
}

// BarcodeOutput contains the generated barcode in multiple formats
type BarcodeOutput struct {
	ImageBase64 string // Base64-encoded PNG image
	ZPL         string // ZPL (Zebra Programming Language) commands
}

// GenerateBarcode creates a barcode label with optional text lines.
// It returns both a PNG image (as base64) and ZPL commands for thermal printers.
//
// The function coordinates the barcode generation pipeline:
//  1. Validates input parameters
//  2. Encodes the barcode data
//  3. Calculates appropriate barcode dimensions
//  4. Renders barcode and text onto a label image
//  5. Exports to PNG and ZPL formats
func GenerateBarcode(input BarcodeInput) (*BarcodeOutput, error) {
	if err := validateInput(input); err != nil {
		return nil, err
	}

	bc, err := encodeBarcode(input)
	if err != nil {
		return nil, err
	}

	labelImg, barcodeRect, err := renderLabel(input, bc)
	if err != nil {
		return nil, err
	}

	if err := renderTextLines(labelImg, input, barcodeRect); err != nil {
		return nil, err
	}

	return generateOutputFormats(labelImg)
}

// validateInput checks that all input parameters are valid
func validateInput(input BarcodeInput) error {
	if err := validateDPI(input.Dpi); err != nil {
		return err
	}

	if err := validateBarcodeType(input.BarcodeType); err != nil {
		return err
	}

	return nil
}

// validateDPI ensures the DPI is a supported thermal printer value
func validateDPI(dpi int) error {
	for _, validDpi := range standardDPIValues {
		if dpi == validDpi {
			return nil
		}
	}
	return fmt.Errorf("invalid dpi value: %d. Supported dpi values are: %v", dpi, standardDPIValues)
}

// validateBarcodeType ensures the barcode type is supported
func validateBarcodeType(barcodeType BarcodeType) error {
	switch barcodeType {
	case BarcodeTypeCode128, BarcodeTypeQR:
		return nil
	default:
		return fmt.Errorf("invalid barcode type: %s. Supported types: CODE128, QR", barcodeType)
	}
}

// encodeBarcode creates the actual barcode from the input data
func encodeBarcode(input BarcodeInput) (barcode.Barcode, error) {
	switch input.BarcodeType {
	case BarcodeTypeCode128:
		return encodeCode128(input.BarcodeData)
	case BarcodeTypeQR:
		return encodeQRCode(input.BarcodeData)
	default:
		// This should never happen due to validation, but included for safety
		return nil, fmt.Errorf("unsupported barcode type: %s", input.BarcodeType)
	}
}

// encodeCode128 creates a Code128 barcode
func encodeCode128(data string) (barcode.Barcode, error) {
	bc, err := code128.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Code128 barcode: %w", err)
	}
	return bc, nil
}

// encodeQRCode creates a QR code
func encodeQRCode(data string) (barcode.Barcode, error) {
	bc, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code: %w", err)
	}
	return bc, nil
}

// renderLabel creates the label image and places the barcode on it
func renderLabel(input BarcodeInput, bc barcode.Barcode) (*image.RGBA, image.Rectangle, error) {
	labelWidth := mmToPixels(input.Width, input.Dpi)
	labelHeight := mmToPixels(input.Height, input.Dpi)

	barcodeSize := calculateBarcodeSize(input, labelWidth, labelHeight)
	scaledBc, err := scaleBarcodeToFit(bc, barcodeSize)
	if err != nil {
		return nil, image.Rectangle{}, err
	}

	img := createBlankLabel(labelWidth, labelHeight)
	barcodeRect := centerBarcodeOnLabel(img, scaledBc)

	drawBarcodeOnLabel(img, scaledBc, barcodeRect)

	return img, barcodeRect, nil
}

// renderTextLines adds all text lines to the label image
func renderTextLines(img *image.RGBA, input BarcodeInput, barcodeRect image.Rectangle) error {
	for _, textLine := range input.TextLines {
		textY := calculateTextYPosition(barcodeRect, textLine.Position)
		addTextLine(img, textLine.Text, img.Bounds().Dx()/2, textY, textLine.Size, input.Dpi, textLine.Position)
	}
	return nil
}

// generateOutputFormats converts the label image to PNG and ZPL formats
func generateOutputFormats(img *image.RGBA) (*BarcodeOutput, error) {
	base64Image, err := imageToBase64(img)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image to base64: %w", err)
	}

	zplCode := imageToZPL(img)

	return &BarcodeOutput{
		ImageBase64: base64Image,
		ZPL:         zplCode,
	}, nil
}
