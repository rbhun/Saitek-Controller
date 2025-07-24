package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gousb"
)

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

	// Test 1: 118.0 using correct decimal point encoding
	fmt.Printf("\n=== Test 1: 118.0 with 0xD8 (8.) ===\n")
	packet1 := make([]byte, 22)
	result1 := []byte{1, 1, 0xD8, 0, 0} // 118.0 with 8. (0xD8 = 11011000)
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

	// Test 2: 121.3 using correct decimal point encoding
	fmt.Printf("\n=== Test 2: 121.3 with 0xD1 (1.) ===\n")
	packet2 := make([]byte, 22)
	result2 := []byte{1, 2, 0xD1, 3, 0} // 121.3 with 1. (0xD1 = 11010001)
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

	// Test 3: 128.3 using correct decimal point encoding
	fmt.Printf("\n=== Test 3: 128.3 with 0xD8 (8.) ===\n")
	packet3 := make([]byte, 22)
	result3 := []byte{1, 2, 0xD8, 3, 0} // 128.3 with 8. (0xD8 = 11011000)
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

	// Test 4: 113.7 using correct decimal point encoding
	fmt.Printf("\n=== Test 4: 113.7 with 0xD3 (3.) ===\n")
	packet4 := make([]byte, 22)
	result4 := []byte{1, 1, 0xD3, 7, 0} // 113.7 with 3. (0xD3 = 11010011)
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

	fmt.Printf("\nTest completed. These should show decimal points!\n")
	fmt.Printf("Using the correct 0xDx encoding from the documentation.\n")
}
