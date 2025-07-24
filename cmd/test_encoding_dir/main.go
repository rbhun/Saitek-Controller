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

// Method 1: Current method - add dot to next digit
func encodeDisplay1(text string) []byte {
	result := make([]byte, 5)

	// Initialize with spaces
	for i := 0; i < 5; i++ {
		result[i] = digitMap[' ']
	}

	// Process the text from left to right
	textLen := len(text)
	if textLen > 5 {
		text = text[:5]
		textLen = 5
	}

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
			result[pos] = digitMap[' ']
			pos++
		}
	}

	// Ensure we use all 5 positions
	for pos < 5 {
		result[pos] = digitMap[' ']
		pos++
	}

	return result
}

// Method 2: Try adding dot to previous digit
func encodeDisplay2(text string) []byte {
	result := make([]byte, 5)

	// Initialize with spaces
	for i := 0; i < 5; i++ {
		result[i] = digitMap[' ']
	}

	// Process the text from left to right
	textLen := len(text)
	if textLen > 5 {
		text = text[:5]
		textLen = 5
	}

	pos := 0
	for i := 0; i < textLen; i++ {
		ch := rune(text[i])

		if ch == '.' {
			// Add dot to the previous digit if possible
			if pos > 0 {
				result[pos-1] |= 0x80
			}
			continue
		}

		if val, ok := digitMap[ch]; ok {
			result[pos] = val
			pos++
		} else {
			result[pos] = digitMap[' ']
			pos++
		}
	}

	// Ensure we use all 5 positions
	for pos < 5 {
		result[pos] = digitMap[' ']
		pos++
	}

	return result
}

// Method 3: Try using a separate position for decimal point
func encodeDisplay3(text string) []byte {
	result := make([]byte, 5)

	// Initialize with spaces
	for i := 0; i < 5; i++ {
		result[i] = digitMap[' ']
	}

	// Process the text from left to right
	textLen := len(text)
	if textLen > 5 {
		text = text[:5]
		textLen = 5
	}

	pos := 0
	for i := 0; i < textLen; i++ {
		ch := rune(text[i])

		if ch == '.' {
			// Use a special code for decimal point
			result[pos] = 0x0A // Try a different code for decimal point
			pos++
			continue
		}

		if val, ok := digitMap[ch]; ok {
			result[pos] = val
			pos++
		} else {
			result[pos] = digitMap[' ']
			pos++
		}
	}

	// Ensure we use all 5 positions
	for pos < 5 {
		result[pos] = digitMap[' ']
		pos++
	}

	return result
}

func main() {
	testFreqs := []string{"118.00", "118.50", "121.30", "121.90"}

	fmt.Println("Testing Different Encoding Methods")
	fmt.Println("=================================")

	for _, freq := range testFreqs {
		fmt.Printf("\nFrequency: %s\n", freq)
		fmt.Printf("Method 1 (dot on next): %v\n", encodeDisplay1(freq))
		fmt.Printf("Method 2 (dot on prev): %v\n", encodeDisplay2(freq))
		fmt.Printf("Method 3 (separate dot): %v\n", encodeDisplay3(freq))
	}
}
