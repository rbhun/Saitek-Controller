package main

import (
	"fmt"
	"log"

	"github.com/google/gousb"
)

func main() {
	fmt.Println("Testing USB connection to Saitek Radio Panel...")

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

	// Try to send a control message
	testData := []byte{0x0F, 0x01, 0x01, 0x08, 0x80, 0x0F, 0x01, 0x01, 0x08, 0x85, 0x0F, 0x01, 0x02, 0x01, 0x83, 0x0F, 0x01, 0x02, 0x01, 0x87, 0x00, 0x00}

	fmt.Printf("Sending control message...\n")
	_, err = dev.Control(0x21, 0x09, 0x0300, 0, testData)
	if err != nil {
		log.Fatalf("Failed to send control message: %v", err)
	}

	fmt.Printf("Successfully sent control message!\n")
	fmt.Printf("Radio panel should now display the test frequencies.\n")
}
