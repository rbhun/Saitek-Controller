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

	// Test 1: Just numbers without decimal point
	fmt.Printf("\n=== Test 1: Just numbers (11800) ===\n")
	packet1 := make([]byte, 22)
	result1 := []byte{1, 1, 8, 0, 0} // 11800
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
		fmt.Printf("Test 1 sent successfully - Should show 11800\n")
	}

	time.Sleep(3 * time.Second)

	// Test 2: Try different bit positions for decimal point
	fmt.Printf("\n=== Test 2: Bit 5 (0x20) ===\n")
	packet2 := make([]byte, 22)
	result2 := []byte{1, 1, 40, 0, 0} // 118.0 with bit 5 on 8
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
		fmt.Printf("Test 2 sent successfully - Check for decimal point\n")
	}

	time.Sleep(3 * time.Second)

	// Test 3: Try bit 4 (0x10)
	fmt.Printf("\n=== Test 3: Bit 4 (0x10) ===\n")
	packet3 := make([]byte, 22)
	result3 := []byte{1, 1, 24, 0, 0} // 118.0 with bit 4 on 8
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
		fmt.Printf("Test 3 sent successfully - Check for decimal point\n")
	}

	time.Sleep(3 * time.Second)

	// Test 4: Try different decimal point codes
	fmt.Printf("\n=== Test 4: Decimal point code 0x0B ===\n")
	packet4 := make([]byte, 22)
	result4 := []byte{1, 1, 8, 11, 0} // 118. with code 0x0B
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
		fmt.Printf("Test 4 sent successfully - Check for decimal point\n")
	}

	time.Sleep(3 * time.Second)

	// Test 5: Try decimal point code 0x0C
	fmt.Printf("\n=== Test 5: Decimal point code 0x0C ===\n")
	packet5 := make([]byte, 22)
	result5 := []byte{1, 1, 8, 12, 0} // 118. with code 0x0C
	copy(packet5[0:5], result5)
	copy(packet5[5:10], result5)
	copy(packet5[10:15], result5)
	copy(packet5[15:20], result5)
	packet5[20] = 0x00
	packet5[21] = 0x00

	_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet5)
	if err != nil {
		log.Printf("Test 5 failed: %v", err)
	} else {
		fmt.Printf("Test 5 sent successfully - Check for decimal point\n")
	}

	time.Sleep(3 * time.Second)

	fmt.Printf("\nTest completed. Which test showed a decimal point?\n")
	fmt.Printf("If none worked, the hardware might not support decimal points.\n")
}
