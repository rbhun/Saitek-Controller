package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/google/gousb"
)

func main() {
	// Parse command line flags
	var (
		com1Active  = flag.String("com1a", "118.00", "COM1 Active frequency")
		com1Standby = flag.String("com1s", "118.50", "COM1 Standby frequency")
		com2Active  = flag.String("com2a", "121.30", "COM2 Active frequency")
		com2Standby = flag.String("com2s", "121.90", "COM2 Standby frequency")
	)
	flag.Parse()

	fmt.Println("Saitek Radio Panel - Set Frequencies")
	fmt.Println("====================================")

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
	displayData := encodeRadioDisplay(*com1Active, *com1Standby, *com2Active, *com2Standby)

	fmt.Printf("Setting display:\n")
	fmt.Printf("  COM1 Active:   %s\n", *com1Active)
	fmt.Printf("  COM1 Standby:  %s\n", *com1Standby)
	fmt.Printf("  COM2 Active:   %s\n", *com2Active)
	fmt.Printf("  COM2 Standby:  %s\n", *com2Standby)

	// Send control message
	fmt.Printf("Sending control message...\n")
	_, err = dev.Control(0x21, 0x09, 0x0300, 0, displayData)
	if err != nil {
		log.Fatalf("Failed to send control message: %v", err)
	}

	fmt.Printf("Successfully set radio panel display!\n")
	fmt.Printf("Radio panel should now display: %s, %s, %s, %s\n",
		*com1Active, *com1Standby, *com2Active, *com2Standby)
}

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

	// Fill from left to right
	pos := 0
	for i := 0; i < textLen; i++ {
		ch := rune(text[i])

		if ch == '.' {
			// Decimal point: set high nibble to 0xD (1101) and keep digit value in low nibble
			// This encodes the previous digit with a decimal point
			if pos > 0 {
				// Set the high nibble to 0xD while keeping the digit value
				result[pos-1] = 0xD0 | (result[pos-1] & 0x0F)
			}
			// Don't increment pos for decimal point
		} else if val, ok := digitMap[ch]; ok {
			if pos < 5 {
				result[pos] = val
				pos++
			}
		} else {
			if pos < 5 {
				result[pos] = digitMap[' '] // Default to space for unknown characters
				pos++
			}
		}
	}

	// Ensure we use all 5 positions by padding with spaces if needed
	for pos < 5 {
		result[pos] = digitMap[' ']
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
