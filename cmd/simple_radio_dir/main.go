package main

import (
	"flag"
	"fmt"
	"log"
	"time"

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

	fmt.Println("Saitek Radio Panel Controller")
	fmt.Println("=============================")

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

	fmt.Printf("Successfully sent display data to radio panel!\n")
	fmt.Printf("Radio panel should now display the specified frequencies.\n")

	// Wait a moment to observe
	time.Sleep(2 * time.Second)

	// Try a different frequency to test updates
	fmt.Printf("Updating to test frequencies...\n")
	testData := encodeRadioDisplay("122.80", "122.90", "123.40", "123.50")
	_, err = dev.Control(0x21, 0x09, 0x0300, 0, testData)
	if err != nil {
		log.Printf("Failed to update display: %v", err)
	} else {
		fmt.Printf("Successfully updated display!\n")
	}

	time.Sleep(2 * time.Second)
	fmt.Printf("Test completed!\n")
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
	if textLen > 5 {
		text = text[:5] // Take first 5 characters
		textLen = 5
	}

	// Fill from left to right
	pos := 0
	for i := 0; i < textLen; i++ {
		ch := rune(text[i])

		if ch == '.' {
			// Look ahead for the next character to add dot to
			if i+1 < textLen {
				nextCh := rune(text[i+1])
				if val, ok := digitMap[nextCh]; ok {
					result[pos] = val | 0x80 // Set high bit for dot
					pos++
					i++ // Skip the next character since we've handled it
				}
			}
			continue
		}

		if val, ok := digitMap[ch]; ok {
			result[pos] = val
			pos++
		} else {
			result[pos] = digitMap[' '] // Default to space for unknown characters
			pos++
		}
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
