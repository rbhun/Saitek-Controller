package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gousb"
)

// Digit encoding map based on fpanels library
var digitMap = map[rune]byte{
	'0': 0x00, '1': 0x01, '2': 0x02, '3': 0x03,
	'4': 0x04, '5': 0x05, '6': 0x06, '7': 0x07,
	'8': 0x08, '9': 0x09, ' ': 0x0F, '-': 0x0E,
}

func main() {
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

	// Test 1: Confirm code 16 is decimal point with 118.00
	fmt.Printf("\n=== Test 1: 118.00 with code 16 as decimal point ===\n")
	packet1 := make([]byte, 22)
	result1 := []byte{1, 1, 8, 16, 0} // 118.0 with decimal point code 16
	copy(packet1[0:5], result1)
	copy(packet1[5:10], result1)
	copy(packet1[10:15], result1)
	copy(packet1[15:20], result1)
	packet1[20] = 0x00
	packet1[21] = 0x00

	_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet1)
	if err != nil {
		log.Printf("Test 1 failed: %v", err)
	} else {
		fmt.Printf("Test 1 sent successfully - Should show 118.0\n")
	}

	time.Sleep(3 * time.Second)

	// Test 2: Test with 121.30
	fmt.Printf("\n=== Test 2: 121.30 with code 16 as decimal point ===\n")
	packet2 := make([]byte, 22)
	result2 := []byte{1, 2, 1, 16, 3} // 121.3 with decimal point code 16
	copy(packet2[0:5], result2)
	copy(packet2[5:10], result2)
	copy(packet2[10:15], result2)
	copy(packet2[15:20], result2)
	packet2[20] = 0x00
	packet2[21] = 0x00

	_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet2)
	if err != nil {
		log.Printf("Test 2 failed: %v", err)
	} else {
		fmt.Printf("Test 2 sent successfully - Should show 121.3\n")
	}

	time.Sleep(3 * time.Second)

	// Test 3: Test with 128.30 (like in the image)
	fmt.Printf("\n=== Test 3: 128.30 with code 16 as decimal point ===\n")
	packet3 := make([]byte, 22)
	result3 := []byte{1, 2, 8, 16, 3} // 128.3 with decimal point code 16
	copy(packet3[0:5], result3)
	copy(packet3[5:10], result3)
	copy(packet3[10:15], result3)
	copy(packet3[15:20], result3)
	packet3[20] = 0x00
	packet3[21] = 0x00

	_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet3)
	if err != nil {
		log.Printf("Test 3 failed: %v", err)
	} else {
		fmt.Printf("Test 3 sent successfully - Should show 128.3\n")
	}

	time.Sleep(3 * time.Second)

	// Test 4: Test with 113.70 (like in the image)
	fmt.Printf("\n=== Test 4: 113.70 with code 16 as decimal point ===\n")
	packet4 := make([]byte, 22)
	result4 := []byte{1, 1, 3, 16, 7} // 113.7 with decimal point code 16
	copy(packet4[0:5], result4)
	copy(packet4[5:10], result4)
	copy(packet4[10:15], result4)
	copy(packet4[15:20], result4)
	packet4[20] = 0x00
	packet4[21] = 0x00

	_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet4)
	if err != nil {
		log.Printf("Test 4 failed: %v", err)
	} else {
		fmt.Printf("Test 4 sent successfully - Should show 113.7\n")
	}

	time.Sleep(3 * time.Second)

	fmt.Printf("\nTest completed. Did any of these show decimal points?\n")
	fmt.Printf("If code 16 works, we've found the correct decimal point encoding!\n")
}
