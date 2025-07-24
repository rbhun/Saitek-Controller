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
	fmt.Println("FIP Panel Test using IOKit")
	fmt.Println("===========================")
	fmt.Println()

	// Create FIP panel
	panel, err := fip.NewIOKitFIPPanel("FIP Test", 320, 240)
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
	fmt.Println("Press buttons on the FIP device to test...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start button event monitoring
	go monitorButtonEvents(panel)

	// Main loop
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return
		case <-ticker.C:
			// Print current button states every second
			printButtonStates(panel)
		}
	}
}

func monitorButtonEvents(panel *fip.IOKitFIPPanel) {
	events := panel.GetButtonEvents()

	for event := range events {
		action := "pressed"
		if !event.Pressed {
			action = "released"
		}

		fmt.Printf("🎯 FIP Button %d %s\n", event.ButtonID+1, action)

		// You can add specific actions for each button here
		switch event.ButtonID {
		case 0:
			fmt.Println("   → Button 1: Navigation")
		case 1:
			fmt.Println("   → Button 2: Menu")
		case 2:
			fmt.Println("   → Button 3: Enter")
		case 3:
			fmt.Println("   → Button 4: Escape")
		case 4:
			fmt.Println("   → Button 5: Up")
		case 5:
			fmt.Println("   → Button 6: Down")
		case 6:
			fmt.Println("   → Button 7: Left")
		case 7:
			fmt.Println("   → Button 8: Right")
		case 8:
			fmt.Println("   → Button 9: Function 1")
		case 9:
			fmt.Println("   → Button 10: Function 2")
		case 10:
			fmt.Println("   → Button 11: Function 3")
		case 11:
			fmt.Println("   → Button 12: Function 4")
		}
	}
}

func printButtonStates(panel *fip.IOKitFIPPanel) {
	fmt.Print("Button States: [")
	for i := 0; i < 12; i++ {
		if panel.GetButtonState(i) {
			fmt.Print("●")
		} else {
			fmt.Print("○")
		}
		if i < 11 {
			fmt.Print(" ")
		}
	}
	fmt.Println("]")

	// Print raw data for debugging
	rawData := panel.GetLastButtonData()
	if len(rawData) > 0 {
		fmt.Printf("Raw Data: %v\n", rawData)
	}
}
