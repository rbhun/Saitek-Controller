package fip

import (
	"testing"
)

func TestNewFIPPanel(t *testing.T) {
	panel, err := NewFIPPanel("Test Panel", 320, 240)
	if err != nil {
		t.Fatalf("Failed to create FIP panel: %v", err)
	}
	defer panel.Close()

	if panel.width != 320 {
		t.Errorf("Expected width 320, got %d", panel.width)
	}
	if panel.height != 240 {
		t.Errorf("Expected height 240, got %d", panel.height)
	}
	if panel.title != "Test Panel" {
		t.Errorf("Expected title 'Test Panel', got '%s'", panel.title)
	}
}

func TestInstrumentData(t *testing.T) {
	data := InstrumentData{
		Pitch:        5.0,
		Roll:         10.0,
		Airspeed:     120.0,
		Altitude:     5000.0,
		Pressure:     29.92,
		Heading:      180.0,
		VerticalSpeed: 500.0,
		TurnRate:     3.0,
		Slip:         0.0,
	}

	if data.Pitch != 5.0 {
		t.Errorf("Expected pitch 5.0, got %f", data.Pitch)
	}
	if data.Airspeed != 120.0 {
		t.Errorf("Expected airspeed 120.0, got %f", data.Airspeed)
	}
}

func TestImageGenerator(t *testing.T) {
	generator := NewImageGenerator(320, 240)

	// Test test pattern creation
	testPattern := generator.CreateTestPattern()
	if testPattern.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", testPattern.Bounds().Dx())
	}
	if testPattern.Bounds().Dy() != 240 {
		t.Errorf("Expected height 240, got %d", testPattern.Bounds().Dy())
	}

	// Test color bars creation
	colorBars := generator.CreateColorBars()
	if colorBars.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", colorBars.Bounds().Dx())
	}

	// Test gradient creation
	gradient := generator.CreateGradient()
	if gradient.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", gradient.Bounds().Dx())
	}
}

func TestInstrumentImageCreation(t *testing.T) {
	generator := NewImageGenerator(320, 240)
	data := InstrumentData{
		Pitch:        5.0,
		Roll:         10.0,
		Airspeed:     120.0,
		Altitude:     5000.0,
		Pressure:     29.92,
		Heading:      180.0,
		VerticalSpeed: 500.0,
		TurnRate:     3.0,
		Slip:         0.0,
	}

	// Test artificial horizon
	ahImg := generator.CreateInstrumentImage(InstrumentArtificialHorizon, data)
	if ahImg.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", ahImg.Bounds().Dx())
	}

	// Test airspeed indicator
	asImg := generator.CreateInstrumentImage(InstrumentAirspeed, data)
	if asImg.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", asImg.Bounds().Dx())
	}

	// Test altimeter
	altImg := generator.CreateInstrumentImage(InstrumentAltimeter, data)
	if altImg.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", altImg.Bounds().Dx())
	}

	// Test compass
	compImg := generator.CreateInstrumentImage(InstrumentCompass, data)
	if compImg.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", compImg.Bounds().Dx())
	}

	// Test vertical speed indicator
	vsiImg := generator.CreateInstrumentImage(InstrumentVerticalSpeed, data)
	if vsiImg.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", vsiImg.Bounds().Dx())
	}

	// Test turn coordinator
	tcImg := generator.CreateInstrumentImage(InstrumentTurnCoordinator, data)
	if tcImg.Bounds().Dx() != 320 {
		t.Errorf("Expected width 320, got %d", tcImg.Bounds().Dx())
	}
}

func TestPanelInterface(t *testing.T) {
	panel, err := NewFIPPanel("Test Panel", 320, 240)
	if err != nil {
		t.Fatalf("Failed to create FIP panel: %v", err)
	}
	defer panel.Close()

	// Test panel interface methods
	if panel.GetType() != PanelTypeFIP {
		t.Errorf("Expected panel type FIP, got %v", panel.GetType())
	}

	if panel.GetName() != "Test Panel" {
		t.Errorf("Expected panel name 'Test Panel', got '%s'", panel.GetName())
	}

	// Test connection (should work in virtual mode)
	if err := panel.Connect(); err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	if err := panel.Disconnect(); err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}
} 