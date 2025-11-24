package barcode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateDPI_ValidValues ensures standard DPI values pass validation
func TestValidateDPI_ValidValues(t *testing.T) {
	validDPIs := []int{203, 300, 600}

	for _, dpi := range validDPIs {
		t.Run(fmt.Sprintf("DPI_%d", dpi), func(t *testing.T) {
			err := validateDPI(dpi)
			assert.NoError(t, err, "DPI %d should be valid", dpi)
		})
	}
}

// TestValidateDPI_InvalidValue ensures invalid DPI values are rejected
func TestValidateDPI_InvalidValue(t *testing.T) {
	err := validateDPI(150)
	assert.Error(t, err, "Invalid DPI should return error")
	assert.Contains(t, err.Error(), "invalid dpi value")
}

// TestValidateBarcodeType_Valid ensures supported types pass validation
func TestValidateBarcodeType_Valid(t *testing.T) {
	validTypes := []BarcodeType{BarcodeTypeCode128, BarcodeTypeQR}

	for _, barcodeType := range validTypes {
		t.Run(string(barcodeType), func(t *testing.T) {
			err := validateBarcodeType(barcodeType)
			assert.NoError(t, err, "Barcode type %s should be valid", barcodeType)
		})
	}
}

// TestValidateBarcodeType_Invalid ensures invalid types are rejected
func TestValidateBarcodeType_Invalid(t *testing.T) {
	err := validateBarcodeType("INVALID_TYPE")
	assert.Error(t, err, "Invalid barcode type should return error")
	assert.Contains(t, err.Error(), "invalid barcode type")
}

// TestGenerateBarcode_Code128_Success verifies successful Code128 generation
func TestGenerateBarcode_Code128_Success(t *testing.T) {
	input := BarcodeInput{
		BarcodeData: "1234567890",
		BarcodeType: BarcodeTypeCode128,
		Width:       50.0,
		Height:      30.0,
		Dpi:         300,
		TextLines: []TextLine{
			{
				Text:     "Sample Text",
				Position: TextPositionAbove,
				Size:     TextSizeMedium,
			},
		},
	}

	output, err := GenerateBarcode(input)

	require.NoError(t, err, "Should successfully generate Code128 barcode")
	assert.NotNil(t, output, "Output should not be nil")
	assert.NotEmpty(t, output.ImageBase64, "Image base64 should not be empty")
	assert.NotEmpty(t, output.ZPL, "ZPL should not be empty")
	assert.Contains(t, output.ImageBase64, "iVBORw0KGgo", "Image should be valid PNG base64")
	assert.Contains(t, output.ZPL, "^XA", "ZPL should contain valid ZPL commands")
}

// TestGenerateBarcode_QR_Success verifies successful QR code generation
func TestGenerateBarcode_QR_Success(t *testing.T) {
	input := BarcodeInput{
		BarcodeData: "https://example.com/product/12345",
		BarcodeType: BarcodeTypeQR,
		Width:       50.0,
		Height:      50.0,
		Dpi:         203,
		TextLines: []TextLine{
			{
				Text:     "Product: 12345",
				Position: TextPositionBelow,
				Size:     TextSizeSmall,
			},
		},
	}

	output, err := GenerateBarcode(input)

	require.NoError(t, err, "Should successfully generate QR code")
	assert.NotNil(t, output, "Output should not be nil")
	assert.NotEmpty(t, output.ImageBase64, "Image base64 should not be empty")
	assert.NotEmpty(t, output.ZPL, "ZPL should not be empty")
}

// TestGenerateBarcode_InvalidInput verifies validation is performed
func TestGenerateBarcode_InvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		input       BarcodeInput
		expectedErr string
	}{
		{
			name: "Invalid DPI",
			input: BarcodeInput{
				BarcodeData: "test",
				BarcodeType: BarcodeTypeCode128,
				Width:       50.0,
				Height:      30.0,
				Dpi:         999,
			},
			expectedErr: "invalid dpi value",
		},
		{
			name: "Invalid Barcode Type",
			input: BarcodeInput{
				BarcodeData: "test",
				BarcodeType: "INVALID",
				Width:       50.0,
				Height:      30.0,
				Dpi:         300,
			},
			expectedErr: "invalid barcode type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := GenerateBarcode(tt.input)
			assert.Error(t, err, "Should return error")
			assert.Nil(t, output, "Output should be nil on error")
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

// TestGenerateBarcode_MultipleTextLines verifies multiple text lines are rendered
func TestGenerateBarcode_MultipleTextLines(t *testing.T) {
	input := BarcodeInput{
		BarcodeData: "LOC-A1-B2-C3",
		BarcodeType: BarcodeTypeCode128,
		Width:       75.0,
		Height:      40.0,
		Dpi:         600,
		TextLines: []TextLine{
			{
				Text:     "Warehouse A",
				Position: TextPositionAbove,
				Size:     TextSizeLarge,
			},
			{
				Text:     "LOC-A1-B2-C3",
				Position: TextPositionBelow,
				Size:     TextSizeMedium,
			},
		},
	}

	output, err := GenerateBarcode(input)

	require.NoError(t, err, "Should successfully generate barcode with multiple text lines")
	assert.NotNil(t, output, "Output should not be nil")
	assert.NotEmpty(t, output.ImageBase64, "Image should not be empty")
	assert.NotEmpty(t, output.ZPL, "ZPL should not be empty")
}

// TestMmToPixels verifies the conversion formula
func TestMmToPixels(t *testing.T) {
	tests := []struct {
		mm       float64
		dpi      int
		expected int
	}{
		{25.4, 203, 203},  // 1 inch at 203 DPI
		{25.4, 300, 300},  // 1 inch at 300 DPI
		{50.8, 203, 406},  // 2 inches at 203 DPI
		{10.0, 100, 39},   // 10mm at 100 DPI (approximately)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.1fmm_%dDPI", tt.mm, tt.dpi), func(t *testing.T) {
			result := mmToPixels(tt.mm, tt.dpi)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCalculateBarcodeSize_Code128 verifies Code128 sizing logic
func TestCalculateBarcodeSize_Code128(t *testing.T) {
	input := BarcodeInput{
		BarcodeType: BarcodeTypeCode128,
		Width:       100.0,
		Height:      60.0,
		Dpi:         300,
	}

	labelWidth := mmToPixels(input.Width, input.Dpi)
	labelHeight := mmToPixels(input.Height, input.Dpi)

	size := calculateBarcodeSize(input, labelWidth, labelHeight)

	assert.Greater(t, size.X, 0, "Width should be positive")
	assert.Greater(t, size.Y, 0, "Height should be positive")
	assert.Equal(t, labelWidth-(labelMarginPixels*2), size.X, "Code128 should use full width")
}

// TestCalculateBarcodeSize_QR verifies QR sizing logic (square)
func TestCalculateBarcodeSize_QR(t *testing.T) {
	input := BarcodeInput{
		BarcodeType: BarcodeTypeQR,
		Width:       50.0,
		Height:      50.0,
		Dpi:         300,
	}

	labelWidth := mmToPixels(input.Width, input.Dpi)
	labelHeight := mmToPixels(input.Height, input.Dpi)

	size := calculateBarcodeSize(input, labelWidth, labelHeight)

	assert.Equal(t, size.X, size.Y, "QR code must be square")
	assert.Greater(t, size.X, 0, "Size should be positive")
}

// TestGetFontSize verifies font sizing and scaling
func TestGetFontSize(t *testing.T) {
	tests := []struct {
		name     string
		size     TextSize
		dpi      int
		expected float64
	}{
		{name: "Small", size: TextSizeSmall, dpi: 300, expected: 8.0},
		{name: "Medium", size: TextSizeMedium, dpi: 300, expected: 10.0},
		{name: "Large", size: TextSizeLarge, dpi: 300, expected: 12.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fontSize, height := getFontSize(tt.size, tt.dpi, 200)
			assert.Greater(t, fontSize, 0.0, "Font size should be positive")
			assert.Greater(t, height, 0.0, "Font height should be positive")
		})
	}
}
