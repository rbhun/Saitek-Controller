package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gousb"
)

func main() {
	fmt.Println("Testing Multi Panel Button LEDs")
	fmt.Println("===============================")

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

	// Test different button LED combinations
	testCases := []struct {
		name string
		leds byte
		desc string
	}{
		{"Test 1: All off", 0x00, "All buttons should be off"},
		{"Test 2: AP only", 0x01, "Only AP button should be lit"},
		{"Test 3: HDG only", 0x02, "Only HDG button should be lit"},
		{"Test 4: NAV only", 0x04, "Only NAV button should be lit"},
		{"Test 5: IAS only", 0x08, "Only IAS button should be lit"},
		{"Test 6: ALT only", 0x10, "Only ALT button should be lit"},
		{"Test 7: VS only", 0x20, "Only VS button should be lit"},
		{"Test 8: APR only", 0x40, "Only APR button should be lit"},
		{"Test 9: REV only", 0x80, "Only REV button should be lit"},
		{"Test 10: AP + HDG", 0x03, "AP and HDG buttons should be lit"},
		{"Test 11: AP + HDG + NAV", 0x07, "AP, HDG, and NAV buttons should be lit"},
		{"Test 12: All on", 0xFF, "All buttons should be lit"},
		{"Test 13: Alternating", 0x55, "AP, NAV, ALT, REV buttons should be lit"},
		{"Test 14: Other alternating", 0xAA, "HDG, IAS, VS, APR buttons should be lit"},
	}

	for _, test := range testCases {
		fmt.Printf("\n=== %s ===\n", test.name)
		fmt.Printf("LEDs: 0x%02X (%08b) - %s\n", test.leds, test.leds, test.desc)

		// Create packet with some display data and the LED value
		packet := make([]byte, 12)
		// Set some display data (12345 on top, 67890 on bottom)
		copy(packet[0:5], []byte{0x01, 0x02, 0x03, 0x04, 0x05})
		copy(packet[5:10], []byte{0x06, 0x07, 0x08, 0x09, 0x00})
		packet[10] = test.leds
		packet[11] = 0xFF

		// Send control message
		_, err = dev.Control(0x21, 0x09, 0x0300, 0, packet)
		if err != nil {
			log.Printf("Failed to send control message: %v", err)
		} else {
			fmt.Printf("Sent successfully - Check button LEDs\n")
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Printf("\nButton LED test completed!\n")
}
