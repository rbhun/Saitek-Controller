package main

import (
	"flag"
	"fmt"
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

	fmt.Println("Debug: Radio Panel Data")
	fmt.Println("=======================")

	// Encode the display data
	displayData := encodeRadioDisplay(*com1Active, *com1Standby, *com2Active, *com2Standby)

	fmt.Printf("Input frequencies:\n")
	fmt.Printf("  COM1 Active:   %s\n", *com1Active)
	fmt.Printf("  COM1 Standby:  %s\n", *com1Standby)
	fmt.Printf("  COM2 Active:   %s\n", *com2Active)
	fmt.Printf("  COM2 Standby:  %s\n", *com2Standby)

	fmt.Printf("\nEncoded data (22 bytes):\n")
	for i, b := range displayData {
		fmt.Printf("  [%2d]: 0x%02x (%3d) '%c'\n", i, b, b, b)
	}

	// Show each display's 5 bytes
	fmt.Printf("\nDisplay breakdown:\n")
	fmt.Printf("COM1 Active (bytes 0-4):   %v\n", displayData[0:5])
	fmt.Printf("COM1 Standby (bytes 5-9):  %v\n", displayData[5:10])
	fmt.Printf("COM2 Active (bytes 10-14): %v\n", displayData[10:15])
	fmt.Printf("COM2 Standby (bytes 15-19):%v\n", displayData[15:20])
	fmt.Printf("Padding (bytes 20-21):     %v\n", displayData[20:22])

	// Test individual encoding
	fmt.Printf("\nIndividual encoding test:\n")
	fmt.Printf("'118.00' -> %v\n", encodeDisplay("118.00"))
	fmt.Printf("'118.50' -> %v\n", encodeDisplay("118.50"))
	fmt.Printf("'121.30' -> %v\n", encodeDisplay("121.30"))
	fmt.Printf("'121.90' -> %v\n", encodeDisplay("121.90"))
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
