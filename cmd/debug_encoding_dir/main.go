package main

import (
	"fmt"
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
	
	// Ensure we use all 5 positions by padding with spaces if needed
	for pos < 5 {
		result[pos] = digitMap[' ']
		pos++
	}

	return result
}

func main() {
	testFreqs := []string{"118.00", "118.50", "121.30", "121.90"}

	fmt.Println("Debug: Display Encoding")
	fmt.Println("======================")

	for _, freq := range testFreqs {
		encoded := encodeDisplay(freq)
		fmt.Printf("Frequency: %s\n", freq)
		fmt.Printf("  Encoded bytes: [%02X %02X %02X %02X %02X]\n",
			encoded[0], encoded[1], encoded[2], encoded[3], encoded[4])

		// Decode back to see what it would display
		var decoded string
		for _, b := range encoded {
			if b == 0x0F {
				decoded += " "
			} else if b&0x80 != 0 {
				// Has decimal point
				digit := b & 0x7F
				for k, v := range digitMap {
					if v == digit {
						decoded += string(k) + "."
						break
					}
				}
			} else {
				// Regular digit
				for k, v := range digitMap {
					if v == b {
						decoded += string(k)
						break
					}
				}
			}
		}
		fmt.Printf("  Decoded back:  %s\n", decoded)
		fmt.Println()
	}
}
