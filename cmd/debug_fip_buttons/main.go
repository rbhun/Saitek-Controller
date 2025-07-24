package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"saitek-controller/internal/fip"
)

func main() {
	fmt.Println("FIP Button Debug Tool")
	fmt.Println("=====================")
	fmt.Println()

	// Create FIP panel
	panel, err := fip.NewIOKitFIPPanel()
	if err != nil {
		log.Fatalf("Failed to create FIP panel: %v", err)
	}
	defer panel.Close()

	// Connect to the device
	fmt.Println("Connecting to FIP device...")
	err = panel.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to FIP device: %v", err)
	}

	fmt.Println("✓ Successfully connected to FIP device")
	fmt.Println("Press buttons on the FIP device to analyze the data...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start detailed monitoring
	go monitorButtonData(panel)

	// Main loop
	for {
		select {
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return
		case <-time.After(100 * time.Millisecond):
			// Print detailed data analysis
			analyzeButtonData(panel)
		}
	}
}

func monitorButtonData(panel *fip.IOKitFIPPanel) {
	events := panel.GetButtonEvents()

	for event := range events {
		action := "pressed"
		if !event.Pressed {
			action = "released"
		}

		fmt.Printf("\n🎯 BUTTON EVENT: Button %d %s\n", event.ButtonID+1, action)
		fmt.Printf("   Raw Data: %v\n", panel.GetLastButtonData())
		fmt.Printf("   Binary: %08b %08b\n", panel.GetLastButtonData()[0], panel.GetLastButtonData()[1])
	}
}

func analyzeButtonData(panel *fip.IOKitFIPPanel) {
	rawData := panel.GetLastButtonData()
	if len(rawData) < 2 {
		return
	}

	// Check if data has changed from idle
	if rawData[0] != 0 || rawData[1] != 0 {
		fmt.Printf("\n🔍 DATA CHANGE DETECTED: %v\n", rawData)
		fmt.Printf("   Binary: %08b %08b\n", rawData[0], rawData[1])

		// Analyze each bit
		fmt.Println("   Bit Analysis:")
		for byteIndex, byteVal := range rawData {
			fmt.Printf("   Byte %d (%08b): ", byteIndex, byteVal)
			for bitIndex := 0; bitIndex < 8; bitIndex++ {
				bit := (byteVal >> bitIndex) & 1
				if bit == 1 {
					fmt.Printf("Bit%d=1 ", bitIndex)
				}
			}
			fmt.Println()
		}
	}

	// Check for any pressed buttons
	pressedButtons := []int{}
	for i := 0; i < 12; i++ {
		if panel.GetButtonState(i) {
			pressedButtons = append(pressedButtons, i+1)
		}
	}

	if len(pressedButtons) > 0 {
		fmt.Printf("   Pressed Buttons: %v\n", pressedButtons)
	}
}
