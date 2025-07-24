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

// Test different decimal point encoding methods
func testDecimalEncoding() {
	testFreq := "118.00"

	fmt.Printf("Testing decimal point encoding for: %s\n", testFreq)

	// Method 1: High bit (0x80)
	fmt.Printf("\nMethod 1 - High bit (0x80):\n")
	result1 := make([]byte, 5)
	for i := 0; i < 5; i++ {
		result1[i] = digitMap[' ']
	}

	pos := 0
	for i := 0; i < len(testFreq); i++ {
		ch := rune(testFreq[i])

		if ch == '.' {
			if pos > 0 {
				result1[pos-1] |= 0x80
			}
			continue
		}

		if val, ok := digitMap[ch]; ok {
			if pos < 5 {
				result1[pos] = val
				pos++
			}
		}
	}
	fmt.Printf("Result: %v\n", result1)

	// Method 2: Different bit position (0x40)
	fmt.Printf("\nMethod 2 - Bit 6 (0x40):\n")
	result2 := make([]byte, 5)
	for i := 0; i < 5; i++ {
		result2[i] = digitMap[' ']
	}

	pos = 0
	for i := 0; i < len(testFreq); i++ {
		ch := rune(testFreq[i])

		if ch == '.' {
			if pos > 0 {
				result2[pos-1] |= 0x40
			}
			continue
		}

		if val, ok := digitMap[ch]; ok {
			if pos < 5 {
				result2[pos] = val
				pos++
			}
		}
	}
	fmt.Printf("Result: %v\n", result2)

	// Method 3: Separate decimal point code
	fmt.Printf("\nMethod 3 - Separate decimal point:\n")
	result3 := make([]byte, 5)
	for i := 0; i < 5; i++ {
		result3[i] = digitMap[' ']
	}

	pos = 0
	for i := 0; i < len(testFreq); i++ {
		ch := rune(testFreq[i])

		if ch == '.' {
			if pos < 5 {
				result3[pos] = 0x0A // Try different code for decimal point
				pos++
			}
			continue
		}

		if val, ok := digitMap[ch]; ok {
			if pos < 5 {
				result3[pos] = val
				pos++
			}
		}
	}
	fmt.Printf("Result: %v\n", result3)
}

func main() {
	testDecimalEncoding()

	// Test with hardware
	fmt.Printf("\nTesting with hardware...\n")

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

	// Test Method 1 (high bit)
	fmt.Printf("\nTesting Method 1 (high bit) on hardware...\n")
	packet1 := make([]byte, 22)

	// Fill with test data using Method 1
	result1 := []byte{1, 1, 136, 0, 0} // 118.0 with high bit
	copy(packet1[0:5], result1)
	copy(packet1[5:10], result1)
	copy(packet1[10:15], result1)
	copy(packet1[15:20], result1)
	packet1[20] = 0x00
	packet1[21] = 0x00

	_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet1)
	if err != nil {
		log.Printf("Method 1 failed: %v", err)
	} else {
		fmt.Printf("Method 1 sent successfully\n")
	}
}
