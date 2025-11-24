package barcode

import (
	"image"
	"image/color"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// getFontSize calculates the appropriate font size in points and pixel height.
// It scales the font proportionally for larger labels to maintain readability.
func getFontSize(size TextSize, dpi int, labelWidth int) (float64, float64) {
	baseFontSize := getBaseFontSize(size)
	scaledFontSize := scaleFontByLabelWidth(baseFontSize, labelWidth)

	fontHeight := calculateFontHeight(scaledFontSize, dpi)

	return scaledFontSize, fontHeight
}

// getBaseFontSize returns the base font size in points for the given text size enum.
func getBaseFontSize(size TextSize) float64 {
	switch size {
	case TextSizeSmall:
		return 8.0
	case TextSizeMedium:
		return 10.0
	case TextSizeLarge:
		return 12.0
	default:
		return 10.0
	}
}

// scaleFontByLabelWidth adjusts font size proportionally to label width.
// Larger labels get larger fonts to maintain visual balance.
func scaleFontByLabelWidth(fontSize float64, labelWidth int) float64 {
	// Scale factor based on deviation from 200px baseline
	scaleFactor := 1.0 + float64(labelWidth-200)/1000

	// Clamp scale factor to reasonable bounds
	if scaleFactor > 2.0 {
		scaleFactor = 2.0
	} else if scaleFactor < 1.1 {
		scaleFactor = 1.1
	}

	return fontSize * scaleFactor
}

// calculateFontHeight returns the pixel height of text at the given font size and DPI.
func calculateFontHeight(fontSize float64, dpi int) float64 {
	fontData, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return 0
	}

	face := truetype.NewFace(fontData, &truetype.Options{
		Size: fontSize,
		DPI:  float64(dpi),
	})

	return float64(face.Metrics().Height.Ceil())
}

// addTextLine renders a text string on the label image at the specified position.
// It uses a recursive approach: if the text is too wide for the label, it reduces
// the font size by 0.1 points and tries again. This ensures text always fits.
func addTextLine(img *image.RGBA, text string, centerX, baseY int, size TextSize, dpi float64, position TextPosition) {
	fontSize, fontHeight := getFontSize(size, int(dpi), img.Bounds().Dx())
	addTextLineRecursive(img, text, centerX, baseY, fontSize, fontHeight, dpi, position)
}

// addTextLineRecursive is the internal recursive function that handles text rendering
// with automatic font size reduction if text doesn't fit.
func addTextLineRecursive(img *image.RGBA, text string, centerX, baseY int, fontSize, fontHeight, dpi float64, position TextPosition) {
	fontData, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return
	}

	// Measure text width at current font size
	face := truetype.NewFace(fontData, &truetype.Options{
		Size: fontSize,
		DPI:  dpi,
	})

	textWidth := font.MeasureString(face, text).Ceil()

	// If text is too wide, reduce font size and retry
	maxWidth := img.Bounds().Dx() - labelMarginPixels*2
	if textWidth > maxWidth {
		newFontHeight := calculateFontHeight(fontSize-0.1, int(dpi))
		addTextLineRecursive(img, text, centerX, baseY, fontSize-0.1, newFontHeight, dpi, position)
		return
	}

	// Draw the text
	drawText(img, text, centerX, baseY, fontSize, fontHeight, dpi, position, color.Black)
}

// drawText renders the actual text on the image.
func drawText(img *image.RGBA, text string, centerX, baseY int, fontSize, fontHeight, dpi float64, position TextPosition, col color.Color) {
	fontData, _ := truetype.Parse(goregular.TTF)

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(fontData)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(col))

	// Calculate text position
	face := truetype.NewFace(fontData, &truetype.Options{
		Size: fontSize,
		DPI:  dpi,
	})

	textWidth := font.MeasureString(face, text).Ceil()
	adjustedX := centerX - (textWidth / 2)

	// Adjust Y position based on text position (above/below barcode)
	adjustedY := baseY
	margin := int(fontHeight) / 2

	if position == TextPositionAbove {
		adjustedY = baseY - margin
	} else if position == TextPositionBelow {
		adjustedY = baseY + margin*2 + 5
	}

	pt := freetype.Pt(adjustedX, adjustedY)
	c.DrawString(text, pt)
}
