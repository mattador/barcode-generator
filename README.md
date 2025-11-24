# Barcode Generation Module

Excerpt of barcode and label generation service for warehouse operations taken from a larger project I developed in 2023 called OptiWMS.

## Overview

This module generates barcodes in multiple formats for warehouse management systems. It supports both 1D (Code128) and 2D (QR) barcodes with dual output formats:

- **PNG Images**: Base64-encoded for display in web interfaces
- **ZPL Commands**: Zebra Programming Language for direct thermal printer output

## Architecture

The module is organized into focused, single-responsibility files:

### Core Files

- **`barcode.go`** - Main API and orchestration
  - `GenerateBarcode()` - Primary entry point
  - Input validation functions
  - Barcode encoding coordination

- **`dimensions.go`** - Size and layout calculations
  - `mmToPixels()` - Unit conversion
  - `calculateBarcodeSize()` - Determine barcode dimensions by type
  - `centerBarcodeOnLabel()` - Position calculation
  - `calculateTextHeight()` - Text space requirements

- **`rendering.go`** - Image manipulation
  - `createBlankLabel()` - Initialize label image
  - `drawBarcodeOnLabel()` - Composite barcode onto label

- **`fonts.go`** - Text rendering and font management
  - `getFontSize()` - Calculate appropriate font size
  - `scaleFontByLabelWidth()` - Scale fonts for label size
  - `addTextLine()` - Render text with automatic sizing
  - `addTextLineRecursive()` - Recursive font reduction algorithm

- **`formatting.go`** - Output format conversion
  - `imageToBase64()` - PNG to base64 encoding
  - `imageToZPL()` - PNG to Zebra printer language

- **`barcode_test.go`** - Comprehensive test suite
  - Validation tests
  - Format-specific tests
  - Integration tests

## Key Features

### 1. Multi-Format Support
- **Code128 Barcodes**: Rectangular, optimal for location/product labels
- **QR Codes**: Square, optimal for URLs/complex data

### 2. DPI-Aware Scaling
Supports standard thermal printer DPI values:
- 203 DPI (entry-level printers)
- 300 DPI (standard printers)
- 600 DPI (high-resolution printers)

### 3. Automatic Text Sizing
The `addTextLineRecursive()` function intelligently sizes text:
- Calculates optimal font size for label width
- Recursively reduces font size if text overflows
- Ensures text always fits within label boundaries

### 4. Flexible Text Positioning
Position text relative to barcode:
- `TextPositionAbove` - Above barcode
- `TextPositionBelow` - Below barcode

Text size options:
- `TextSizeSmall` - 8pt base
- `TextSizeMedium` - 10pt base (default)
- `TextSizeLarge` - 12pt base

## Usage

```go
import "github.com/mattador/barcode-generator"

// Generate a Code128 barcode with text labels
input := barcode.BarcodeInput{
	BarcodeData: "LOC-A1-B2-C3",
	BarcodeType: barcode.BarcodeTypeCode128,
	Width:       75.0,  // millimeters
	Height:      40.0,  // millimeters
	Dpi:         300,
	TextLines: []barcode.TextLine{
		{
			Text:     "Warehouse A",
			Position: barcode.TextPositionAbove,
			Size:     barcode.TextSizeLarge,
		},
		{
			Text:     "LOC-A1-B2-C3",
			Position: barcode.TextPositionBelow,
			Size:     barcode.TextSizeMedium,
		},
	},
}

output, err := barcode.GenerateBarcode(input)
if err != nil {
	log.Fatal(err)
}

// Use output.ImageBase64 for web display
// Use output.ZPL for thermal printer
```

## Testing

Run tests with:

```bash
go test -v
```

Test coverage includes:
- Input validation
- DPI support verification
- Barcode type support
- Multiple text lines
- Size calculations
- Font rendering

## Design Patterns

### 1. Single Responsibility Principle
Each file handles one concern:
- Validation, encoding, sizing, rendering, formatting

### 2. Input Validation
Validation happens early in the pipeline with clear error messages

### 3. Pipeline Design
Input data moves through validation, encoding, sizing calculations, rendering, and finally formatting into output formats. Clean separation between stages makes each step easy to test independently.

### 4. Composable Functions
Small functions can be tested and reused independently

## Performance Considerations

- No external service calls - all operations are local
- Efficient image encoding with base64
- Minimal memory allocation
- Single-pass rendering

## Dependencies

External packages:
- `github.com/boombuler/barcode` - Barcode encoding
- `github.com/golang/freetype` - Font rendering
- `golang.org/x/image` - Image utilities
- `simonwaldherr.de/go/zplgfa` - ZPL conversion

## Error Handling

Clear, actionable error messages:
- Invalid DPI: Lists supported values
- Invalid barcode type: Lists supported types
- Encoding failures: Wraps underlying errors with context

## Future Improvements

Possible enhancements:
- Additional barcode types (Code39, EAN-13, etc.)
- Custom fonts support
- Color support for QR codes
- Barcode rotation
- Multi-barcode per label

## Author

Matthew Cooper (2023)

## License

Internal use - part of OptiWMS project
