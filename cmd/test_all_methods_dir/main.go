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

	// Test Method 1: High bit (0x80)
	fmt.Printf("\n=== Testing Method 1: High bit (0x80) ===\n")
	packet1 := make([]byte, 22)
	result1 := []byte{1, 1, 136, 0, 0} // 118.0 with high bit on 8
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
		fmt.Printf("Method 1 sent successfully - Check if you see decimal point\n")
	}

	time.Sleep(3 * time.Second)

	// Test Method 2: Bit 6 (0x40)
	fmt.Printf("\n=== Testing Method 2: Bit 6 (0x40) ===\n")
	packet2 := make([]byte, 22)
	result2 := []byte{1, 1, 72, 0, 0} // 118.0 with bit 6 on 8
	copy(packet2[0:5], result2)
	copy(packet2[5:10], result2)
	copy(packet2[10:15], result2)
	copy(packet2[15:20], result2)
	packet2[20] = 0x00
	packet2[21] = 0x00

	_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet2)
	if err != nil {
		log.Printf("Method 2 failed: %v", err)
	} else {
		fmt.Printf("Method 2 sent successfully - Check if you see decimal point\n")
	}

	time.Sleep(3 * time.Second)

	// Test Method 3: Separate decimal point code
	fmt.Printf("\n=== Testing Method 3: Separate decimal point (0x0A) ===\n")
	packet3 := make([]byte, 22)
	result3 := []byte{1, 1, 8, 10, 0} // 118. with separate decimal point
	copy(packet3[0:5], result3)
	copy(packet3[5:10], result3)
	copy(packet3[10:15], result3)
	copy(packet3[15:20], result3)
	packet3[20] = 0x00
	packet3[21] = 0x00

	_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet3)
	if err != nil {
		log.Printf("Method 3 failed: %v", err)
	} else {
		fmt.Printf("Method 3 sent successfully - Check if you see decimal point\n")
	}

	time.Sleep(3 * time.Second)

	fmt.Printf("\nTest completed. Which method showed the decimal point?\n")
}
