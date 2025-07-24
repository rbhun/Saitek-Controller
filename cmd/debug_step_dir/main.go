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

func main() {
	testFreq := "118.00"

	fmt.Printf("Debugging encoding for: %s\n", testFreq)
	fmt.Printf("Length: %d\n", len(testFreq))

	result := make([]byte, 5)

	// Initialize with spaces
	for i := 0; i < 5; i++ {
		result[i] = digitMap[' ']
	}

	// Process the text from left to right
	textLen := len(testFreq)
	fmt.Printf("Original length: %d\n", textLen)
	if textLen > 5 {
		testFreq = testFreq[:5] // Take first 5 characters
		textLen = 5
		fmt.Printf("Truncated to 5 characters\n")
	}

		// Simple approach: process each character and handle decimal points
	pos := 0
	fmt.Printf("Will process %d characters (positions 0 to %d)\n", textLen, textLen-1)
	for i := 0; i < textLen; i++ {
		ch := rune(testFreq[i])
		fmt.Printf("Processing char '%c' at position %d, current pos=%d\n", ch, i, pos)
		
		if ch == '.' {
			// Add decimal point to the previous digit
			if pos > 0 {
				fmt.Printf("  Adding decimal point to previous digit at pos %d\n", pos-1)
				result[pos-1] |= 0x80
			}
			continue
		}
		
		if val, ok := digitMap[ch]; ok {
			if pos < 5 {
				fmt.Printf("  Adding digit %c (0x%02x) at pos %d\n", ch, val, pos)
				result[pos] = val
				pos++
			} else {
				fmt.Printf("  Skipping digit %c because pos >= 5\n", ch)
			}
		}
	}

	fmt.Printf("Final result: %v\n", result)

	// Decode back to see what it would display
	var decoded string
	for _, b := range result {
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
	fmt.Printf("Decoded back: '%s'\n", decoded)
}
