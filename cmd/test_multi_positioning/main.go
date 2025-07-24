package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gousb"
)

func main() {
	fmt.Println("Testing Multi Panel Positioning")
	fmt.Println("==============================")

	// Create USB context
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Find the multi panel device
	dev, err := ctx.OpenDeviceWithVIDPID(gousb.ID(0x06A3), gousb.ID(0x0D06))
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

	// Test positioning - try different starting positions
	testCases := []struct {
		name   string
		top    []byte
		bottom []byte
		leds   byte
	}{
		{"Test 1: Right-aligned 123", []byte{0x0F, 0x0F, 0x01, 0x02, 0x03}, []byte{0x0F, 0x0F, 0x0F, 0x0F, 0x0F}, 0x0F},
		{"Test 2: Center 123", []byte{0x0F, 0x01, 0x02, 0x03, 0x0F}, []byte{0x0F, 0x0F, 0x0F, 0x0F, 0x0F}, 0x0F},
		{"Test 3: Left-aligned 123", []byte{0x01, 0x02, 0x03, 0x0F, 0x0F}, []byte{0x0F, 0x0F, 0x0F, 0x0F, 0x0F}, 0x0F},
		{"Test 4: All positions", []byte{0x01, 0x02, 0x03, 0x04, 0x05}, []byte{0x0F, 0x0F, 0x0F, 0x0F, 0x0F}, 0x0F},
		{"Test 5: Bottom row test", []byte{0x0F, 0x0F, 0x0F, 0x0F, 0x0F}, []byte{0x01, 0x02, 0x03, 0x04, 0x05}, 0x0F},
		{"Test 6: Both rows", []byte{0x01, 0x02, 0x03, 0x04, 0x05}, []byte{0x06, 0x07, 0x08, 0x09, 0x00}, 0x0F},
	}

	for _, test := range testCases {
		fmt.Printf("\n=== %s ===\n", test.name)
		fmt.Printf("Top: %v, Bottom: %v\n", test.top, test.bottom)

		// Create packet
		packet := make([]byte, 12)
		copy(packet[0:5], test.top)
		copy(packet[5:10], test.bottom)
		packet[10] = test.leds
		packet[11] = 0xFF

		// Send control message
		_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet)
		if err != nil {
			log.Printf("Failed to send control message: %v", err)
		} else {
			fmt.Printf("Sent successfully - Check display\n")
		}

		time.Sleep(3 * time.Second)
	}

	fmt.Printf("\nPositioning test completed!\n")
}
