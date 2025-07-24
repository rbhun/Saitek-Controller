package fip

import (
	"log"
	"time"

	"saitek-controller/internal/usb"
)

// SwitchPanel represents a Saitek Flight Switch Panel
type SwitchPanel struct {
	device    usb.USBDevice
	connected bool
	vendorID  uint16
	productID uint16
}

// LandingGearLights represents the landing gear indicator lights
type LandingGearLights struct {
	GreenN bool // Green N light
	GreenL bool // Green L light
	GreenR bool // Green R light
	RedN   bool // Red N light
	RedL   bool // Red L light
	RedR   bool // Red R light
}

// SwitchState represents the state of all switches on the panel
type SwitchState struct {
	// Byte 1 switches
	BAT      bool // Battery
	ALT      bool // Alternator
	AVIONICS bool // Avionics
	FUEL     bool // Fuel pump
	DEICE    bool // De-ice
	PITOT    bool // Pitot heat
	COWL     bool // Cowl flaps
	PANEL    bool // Panel lights

	// Byte 2 switches
	BEACON  bool // Beacon
	NAV     bool // Navigation lights
	STROBE  bool // Strobe lights
	TAXI    bool // Taxi lights
	LANDING bool // Landing lights
	OFF     bool // Off position
	R       bool // Right position
	L       bool // Left position

	// Byte 3 switches
	BOTH     bool // Both position
	START    bool // Start
	GEARUP   bool // Gear up
	GEARDOWN bool // Gear down
}

// NewSwitchPanel creates a new switch panel
func NewSwitchPanel() *SwitchPanel {
	return &SwitchPanel{
		vendorID:  0x06A3, // Logitech/Saitek vendor ID
		productID: 0x0D67, // Switch Panel product ID
	}
}

// NewSwitchPanelWithUSB creates a new switch panel with custom vendor/product IDs
func NewSwitchPanelWithUSB(vendorID, productID uint16) *SwitchPanel {
	return &SwitchPanel{
		vendorID:  vendorID,
		productID: productID,
	}
}

// Connect connects to the physical switch panel device
func (s *SwitchPanel) Connect() error {
	// Try the USB core approach first (like the Python code)
	log.Printf("Trying USB core approach...")
	if usbDev, err := usb.NewUSBCoreDevice(s.vendorID, s.productID); err != nil {
		log.Printf("USB core approach failed: %v", err)

		// Try the standard HID approach as fallback
		log.Printf("Trying HID approach...")
		device, err := usb.OpenDevice(s.vendorID, s.productID)
		if err != nil {
			log.Printf("HID approach also failed: %v", err)
			s.connected = false
			return err
		}
		s.device = device
	} else {
		log.Printf("USB core approach succeeded!")
		s.device = usbDev
	}

	s.connected = true
	return nil
}

// Disconnect disconnects from the switch panel device
func (s *SwitchPanel) Disconnect() error {
	if s.device != nil {
		s.device.Close()
		s.device = nil
	}
	s.connected = false
	return nil
}

// IsConnected returns whether the panel is connected
func (s *SwitchPanel) IsConnected() bool {
	return s.connected
}

// GetType returns the panel type
func (s *SwitchPanel) GetType() usb.PanelType {
	return usb.PanelTypeSwitch
}

// GetName returns the panel name
func (s *SwitchPanel) GetName() string {
	return "Saitek Flight Switch Panel"
}

// encodeLandingGearLights encodes the landing gear lights into a single byte
// Based on the fpanels documentation:
// 00000001 Green N
// 00000010 Green L
// 00000100 Green R
// 00001000 Red N
// 00010000 Red L
// 00100000 Red R
// xx000000 Not used
func encodeLandingGearLights(lights LandingGearLights) byte {
	var encoded byte

	if lights.GreenN {
		encoded |= 0x01
	}
	if lights.GreenL {
		encoded |= 0x02
	}
	if lights.GreenR {
		encoded |= 0x04
	}
	if lights.RedN {
		encoded |= 0x08
	}
	if lights.RedL {
		encoded |= 0x10
	}
	if lights.RedR {
		encoded |= 0x20
	}

	return encoded
}

// SetLandingGearLights sets the landing gear indicator lights
func (s *SwitchPanel) SetLandingGearLights(lights LandingGearLights) error {
	if s.device == nil {
		log.Printf("Mock: Setting landing gear lights - GreenN:%v GreenL:%v GreenR:%v RedN:%v RedL:%v RedR:%v",
			lights.GreenN, lights.GreenL, lights.GreenR, lights.RedN, lights.RedL, lights.RedR)
		return nil
	}

	encoded := encodeLandingGearLights(lights)

	// Send control message with the encoded byte
	// bmRequestType=0x21, bRequest=0x09, wValue=0x0300, wIndex=0
	packet := []byte{encoded}
	return s.device.SendControlMessage(0x21, 0x09, 0x0300, 0, packet)
}

// SetAllLightsOff turns off all landing gear lights
func (s *SwitchPanel) SetAllLightsOff() error {
	lights := LandingGearLights{
		GreenN: false,
		GreenL: false,
		GreenR: false,
		RedN:   false,
		RedL:   false,
		RedR:   false,
	}
	return s.SetLandingGearLights(lights)
}

// SetAllLightsGreen turns on all green lights and turns off all red lights
func (s *SwitchPanel) SetAllLightsGreen() error {
	lights := LandingGearLights{
		GreenN: true,
		GreenL: true,
		GreenR: true,
		RedN:   false,
		RedL:   false,
		RedR:   false,
	}
	return s.SetLandingGearLights(lights)
}

// SetAllLightsRed turns on all red lights and turns off all green lights
func (s *SwitchPanel) SetAllLightsRed() error {
	lights := LandingGearLights{
		GreenN: false,
		GreenL: false,
		GreenR: false,
		RedN:   true,
		RedL:   true,
		RedR:   true,
	}
	return s.SetLandingGearLights(lights)
}

// SetAllLightsYellow turns on both red and green lights (creates yellow)
func (s *SwitchPanel) SetAllLightsYellow() error {
	lights := LandingGearLights{
		GreenN: true,
		GreenL: true,
		GreenR: true,
		RedN:   true,
		RedL:   true,
		RedR:   true,
	}
	return s.SetLandingGearLights(lights)
}

// SetGearUp sets lights for gear up indication (typically red)
func (s *SwitchPanel) SetGearUp() error {
	return s.SetAllLightsRed()
}

// SetGearDown sets lights for gear down indication (typically green)
func (s *SwitchPanel) SetGearDown() error {
	return s.SetAllLightsGreen()
}

// SetGearTransition sets lights for gear in transition (typically yellow)
func (s *SwitchPanel) SetGearTransition() error {
	return s.SetAllLightsYellow()
}

// ReadSwitchState reads the current state of switches
func (s *SwitchPanel) ReadSwitchState() ([]byte, error) {
	if s.device == nil {
		log.Printf("Mock: Reading switch panel state")
		return make([]byte, 3), nil
	}

	// Read 3 bytes from endpoint 1
	return s.device.ReadBulkData(1, 3)
}

// ParseSwitchState parses the switch state bytes into readable format
func (s *SwitchPanel) ParseSwitchState(data []byte) *SwitchState {
	if len(data) < 3 {
		return nil
	}

	state := &SwitchState{}

	// Byte 1
	state.BAT = (data[0] & 0x01) != 0
	state.ALT = (data[0] & 0x02) != 0
	state.AVIONICS = (data[0] & 0x04) != 0
	state.FUEL = (data[0] & 0x08) != 0
	state.DEICE = (data[0] & 0x10) != 0
	state.PITOT = (data[0] & 0x20) != 0
	state.COWL = (data[0] & 0x40) != 0
	state.PANEL = (data[0] & 0x80) != 0

	// Byte 2
	state.BEACON = (data[1] & 0x01) != 0
	state.NAV = (data[1] & 0x02) != 0
	state.STROBE = (data[1] & 0x04) != 0
	state.TAXI = (data[1] & 0x08) != 0
	state.LANDING = (data[1] & 0x10) != 0
	state.OFF = (data[1] & 0x20) != 0
	state.R = (data[1] & 0x40) != 0
	state.L = (data[1] & 0x80) != 0

	// Byte 3
	state.BOTH = (data[2] & 0x01) != 0
	state.START = (data[2] & 0x02) != 0
	state.GEARUP = (data[2] & 0x04) != 0
	state.GEARDOWN = (data[2] & 0x08) != 0

	return state
}

// GetSwitchState returns the current switch state as a readable structure
func (s *SwitchPanel) GetSwitchState() (*SwitchState, error) {
	data, err := s.ReadSwitchState()
	if err != nil {
		return nil, err
	}
	return s.ParseSwitchState(data), nil
}

// Run starts a monitoring loop for the switch panel
func (s *SwitchPanel) Run() {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 Hz polling
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.IsConnected() {
				data, err := s.ReadSwitchState()
				if err != nil {
					log.Printf("Error reading switch panel state: %v", err)
					continue
				}

				state := s.ParseSwitchState(data)
				if state != nil {
					// Log any active switches
					if state.BAT {
						log.Printf("Switch Panel: BAT activated")
					}
					if state.ALT {
						log.Printf("Switch Panel: ALT activated")
					}
					if state.AVIONICS {
						log.Printf("Switch Panel: AVIONICS activated")
					}
					if state.FUEL {
						log.Printf("Switch Panel: FUEL activated")
					}
					if state.DEICE {
						log.Printf("Switch Panel: DE-ICE activated")
					}
					if state.PITOT {
						log.Printf("Switch Panel: PITOT activated")
					}
					if state.COWL {
						log.Printf("Switch Panel: COWL activated")
					}
					if state.PANEL {
						log.Printf("Switch Panel: PANEL activated")
					}
					if state.BEACON {
						log.Printf("Switch Panel: BEACON activated")
					}
					if state.NAV {
						log.Printf("Switch Panel: NAV activated")
					}
					if state.STROBE {
						log.Printf("Switch Panel: STROBE activated")
					}
					if state.TAXI {
						log.Printf("Switch Panel: TAXI activated")
					}
					if state.LANDING {
						log.Printf("Switch Panel: LANDING activated")
					}
					if state.OFF {
						log.Printf("Switch Panel: OFF activated")
					}
					if state.R {
						log.Printf("Switch Panel: R activated")
					}
					if state.L {
						log.Printf("Switch Panel: L activated")
					}
					if state.BOTH {
						log.Printf("Switch Panel: BOTH activated")
					}
					if state.START {
						log.Printf("Switch Panel: START activated")
					}
					if state.GEARUP {
						log.Printf("Switch Panel: GEAR UP activated")
					}
					if state.GEARDOWN {
						log.Printf("Switch Panel: GEAR DOWN activated")
					}
				}
			}
		}
	}
}

// Close closes the switch panel
func (s *SwitchPanel) Close() {
	s.Disconnect()
}
