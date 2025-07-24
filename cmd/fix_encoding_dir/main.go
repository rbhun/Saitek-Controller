package main

import (
	"fmt"
	"log"

	"github.com/google/gousb"
)

// Digit encoding map based on fpanels library
var digitMap = map[rune]byte{
	'0': 0x00, '1': 0x01, '2': 0x02, '3': 0x03,
	'4': 0x04, '5': 0x05, '6': 0x06, '7': 0x07,
	'8': 0x08, '9': 0x09, ' ': 0x0F, '-': 0x0E,
}

// encodeDisplay encodes a string into 5 bytes for the display
func encodeDisplay(text string) []byte {
	result := make([]byte, 5)

	// Initialize with spaces
	for i := 0; i < 5; i++ {
		result[i] = digitMap[' ']
	}

	// Process the text from left to right
	textLen := len(text)
	if textLen > 5 {
		text = text[:5] // Take first 5 characters
		textLen = 5
	}

	// First pass: collect all digits and note decimal point positions
	var digits []rune
	var decimalPositions []int

			for i := 0; i < textLen; i++ {
			ch := rune(text[i])
			if ch == '.' {
				// Mark the position where decimal point should be added
				decimalPositions = append(decimalPositions, len(digits)-1)
			} else if _, ok := digitMap[ch]; ok {
				digits = append(digits, ch)
			}
		}

	// Second pass: fill the result array
	pos := 0
	for i, digit := range digits {
		if pos >= 5 {
			break
		}

		// Check if this digit should have a decimal point
		hasDecimal := false
		for _, dp := range decimalPositions {
			if dp == i {
				hasDecimal = true
				break
			}
		}

		if hasDecimal {
			result[pos] = digitMap[digit] | 0x80 // Set high bit for decimal point
		} else {
			result[pos] = digitMap[digit]
		}
		pos++
	}

	return result
}

// encodeRadioDisplay encodes all four displays into a 22-byte packet
func encodeRadioDisplay(com1Active, com1Standby, com2Active, com2Standby string) []byte {
	packet := make([]byte, 22)

	// Top Left (COM1 Active)
	copy(packet[0:5], encodeDisplay(com1Active))

	// Top Right (COM1 Standby)
	copy(packet[5:10], encodeDisplay(com1Standby))

	// Bottom Left (COM2 Active)
	copy(packet[10:15], encodeDisplay(com2Active))

	// Bottom Right (COM2 Standby)
	copy(packet[15:20], encodeDisplay(com2Standby))

	// Last two bytes are zero for Windows compatibility
	packet[20] = 0x00
	packet[21] = 0x00

	return packet
}

func main() {
	// Parse command line flags
	var (
		com1Active  = "118.00"
		com1Standby = "118.50"
		com2Active  = "121.30"
		com2Standby = "121.90"
	)

	fmt.Println("Fixed Radio Panel Encoding Test")
	fmt.Println("==============================")

	// Test individual encoding
	fmt.Printf("Testing individual encoding:\n")
	fmt.Printf("'118.00' -> %v\n", encodeDisplay("118.00"))
	fmt.Printf("'118.50' -> %v\n", encodeDisplay("118.50"))
	fmt.Printf("'121.30' -> %v\n", encodeDisplay("121.30"))
	fmt.Printf("'121.90' -> %v\n", encodeDisplay("121.90"))

	// Create USB context
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Find the radio panel device
	dev, err := ctx.OpenDeviceWithVIDPID(gousb.ID(0x06A3), gousb.ID(0x0D05))
	if err != nil {
		log.Fatalf("Failed to find device: %v", err)
	}

	if dev == nil {
		log.Fatalf("Device not found")
	}

	fmt.Printf("Found device: %s\n", dev.Desc.Product)

	// Set auto detach
	if err := dev.SetAutoDetach(true); err != nil {
		log.Printf("Warning: failed to set auto detach: %v", err)
	}

	// Encode the display data
	displayData := encodeRadioDisplay(com1Active, com1Standby, com2Active, com2Standby)

	fmt.Printf("Setting display:\n")
	fmt.Printf("  COM1 Active:   %s\n", com1Active)
	fmt.Printf("  COM1 Standby:  %s\n", com1Standby)
	fmt.Printf("  COM2 Active:   %s\n", com2Active)
	fmt.Printf("  COM2 Standby:  %s\n", com2Standby)

	// Send control message
	fmt.Printf("Sending control message...\n")
	_, err = dev.Control(0x21, 0x09, 0x0300, 0, displayData)
	if err != nil {
		log.Fatalf("Failed to send control message: %v", err)
	}

	fmt.Printf("Successfully set radio panel display!\n")
	fmt.Printf("Radio panel should now display: %s, %s, %s, %s\n",
		com1Active, com1Standby, com2Active, com2Standby)
}
