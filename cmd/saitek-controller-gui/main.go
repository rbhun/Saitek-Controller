package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"saitek-controller/internal/fip"
)

// PanelManager manages all connected panels
type PanelManager struct {
	radio  *fip.RadioPanel
	multi  *fip.MultiPanel
	switch_ *fip.SwitchPanel
	
	radioConnected  bool
	multiConnected  bool
	switchConnected bool
	
	mu sync.RWMutex
}

// PanelState represents the current state of all panels
type PanelState struct {
	Radio struct {
		Connected bool   `json:"connected"`
		COM1Active string `json:"com1Active"`
		COM1Standby string `json:"com1Standby"`
		COM2Active string `json:"com2Active"`
		COM2Standby string `json:"com2Standby"`
	} `json:"radio"`
	
	Multi struct {
		Connected bool   `json:"connected"`
		TopRow    string `json:"topRow"`
		BottomRow string `json:"bottomRow"`
		LEDs      uint8  `json:"leds"`
	} `json:"multi"`
	
	Switch struct {
		Connected bool `json:"connected"`
		Lights    struct {
			GreenN bool `json:"greenN"`
			GreenL bool `json:"greenL"`
			GreenR bool `json:"greenR"`
			RedN   bool `json:"redN"`
			RedL   bool `json:"redL"`
			RedR   bool `json:"redR"`
		} `json:"lights"`
	} `json:"switch"`
}

// NewPanelManager creates a new panel manager
func NewPanelManager() *PanelManager {
	return &PanelManager{
		radio:  fip.NewRadioPanel(),
		multi:  fip.NewMultiPanel(),
		switch_: fip.NewSwitchPanel(),
	}
}

// ConnectAll attempts to connect to all panels
func (pm *PanelManager) ConnectAll() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	// Connect to radio panel
	if err := pm.radio.Connect(); err != nil {
		log.Printf("Failed to connect to radio panel: %v", err)
		pm.radioConnected = false
	} else {
		log.Printf("Successfully connected to radio panel")
		pm.radioConnected = true
	}
	
	// Connect to multi panel
	if err := pm.multi.Connect(); err != nil {
		log.Printf("Failed to connect to multi panel: %v", err)
		pm.multiConnected = false
	} else {
		log.Printf("Successfully connected to multi panel")
		pm.multiConnected = true
	}
	
	// Connect to switch panel
	if err := pm.switch_.Connect(); err != nil {
		log.Printf("Failed to connect to switch panel: %v", err)
		pm.switchConnected = false
	} else {
		log.Printf("Successfully connected to switch panel")
		pm.switchConnected = true
	}
}

// GetState returns the current state of all panels
func (pm *PanelManager) GetState() PanelState {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	var state PanelState
	
	// Radio panel state
	state.Radio.Connected = pm.radioConnected
	if pm.radioConnected {
		// For now, we'll use default values - in a real app you'd cache the current display
		state.Radio.COM1Active = "118.00"
		state.Radio.COM1Standby = "118.50"
		state.Radio.COM2Active = "121.30"
		state.Radio.COM2Standby = "121.90"
	}
	
	// Multi panel state
	state.Multi.Connected = pm.multiConnected
	if pm.multiConnected {
		state.Multi.TopRow = "0000"
		state.Multi.BottomRow = "0000"
		state.Multi.LEDs = 0
	}
	
	// Switch panel state
	state.Switch.Connected = pm.switchConnected
	if pm.switchConnected {
		state.Switch.Lights.GreenN = false
		state.Switch.Lights.GreenL = false
		state.Switch.Lights.GreenR = false
		state.Switch.Lights.RedN = false
		state.Switch.Lights.RedL = false
		state.Switch.Lights.RedR = false
	}
	
	return state
}

// SetRadioDisplay sets the radio panel display
func (pm *PanelManager) SetRadioDisplay(com1Active, com1Standby, com2Active, com2Standby string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if !pm.radioConnected {
		return fmt.Errorf("radio panel not connected")
	}
	
	return pm.radio.SetDisplay(com1Active, com1Standby, com2Active, com2Standby)
}

// SetMultiDisplay sets the multi panel display
func (pm *PanelManager) SetMultiDisplay(topRow, bottomRow string, leds uint8) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if !pm.multiConnected {
		return fmt.Errorf("multi panel not connected")
	}
	
	return pm.multi.SetDisplay(topRow, bottomRow, leds)
}

// SetSwitchLights sets the switch panel landing gear lights
func (pm *PanelManager) SetSwitchLights(lights fip.LandingGearLights) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if !pm.switchConnected {
		return fmt.Errorf("switch panel not connected")
	}
	
	return pm.switch_.SetLandingGearLights(lights)
}

// Close closes all panels
func (pm *PanelManager) Close() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if pm.radio != nil {
		pm.radio.Close()
	}
	if pm.multi != nil {
		pm.multi.Close()
	}
	if pm.switch_ != nil {
		pm.switch_.Close()
	}
}

// HTML template for the web interface
const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Saitek Controller</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
            color: white;
            padding: 20px;
            text-align: center;
        }
        
        .header h1 {
            margin: 0;
            font-size: 2.5em;
            font-weight: 300;
        }
        
        .content {
            padding: 30px;
        }
        
        .panel-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 30px;
            margin-bottom: 30px;
        }
        
        .panel {
            background: #f8f9fa;
            border-radius: 10px;
            padding: 25px;
            border: 2px solid #e9ecef;
            transition: all 0.3s ease;
        }
        
        .panel:hover {
            border-color: #007bff;
            box-shadow: 0 5px 15px rgba(0,123,255,0.2);
        }
        
        .panel-header {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
            padding-bottom: 15px;
            border-bottom: 2px solid #dee2e6;
        }
        
        .panel-title {
            font-size: 1.5em;
            font-weight: 600;
            color: #2c3e50;
            margin: 0;
        }
        
        .status-indicator {
            width: 12px;
            height: 12px;
            border-radius: 50%;
            margin-left: auto;
        }
        
        .status-connected {
            background: #28a745;
            box-shadow: 0 0 10px rgba(40, 167, 69, 0.5);
        }
        
        .status-disconnected {
            background: #dc3545;
            box-shadow: 0 0 10px rgba(220, 53, 69, 0.5);
        }
        
        .form-group {
            margin-bottom: 20px;
        }
        
        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #495057;
        }
        
        .form-control {
            width: 100%;
            padding: 12px;
            border: 2px solid #ced4da;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s ease;
            box-sizing: border-box;
        }
        
        .form-control:focus {
            outline: none;
            border-color: #007bff;
            box-shadow: 0 0 0 3px rgba(0,123,255,0.1);
        }
        
        .btn {
            background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            margin-right: 10px;
            margin-bottom: 10px;
        }
        
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(0,123,255,0.3);
        }
        
        .btn-secondary {
            background: linear-gradient(135deg, #6c757d 0%, #545b62 100%);
        }
        
        .btn-success {
            background: linear-gradient(135deg, #28a745 0%, #1e7e34 100%);
        }
        
        .btn-danger {
            background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
        }
        
        .btn-warning {
            background: linear-gradient(135deg, #ffc107 0%, #e0a800 100%);
            color: #212529;
        }
        
        .led-grid {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 10px;
            margin-top: 15px;
        }
        
        .led-item {
            text-align: center;
            padding: 10px;
            background: #e9ecef;
            border-radius: 8px;
        }
        
        .led-checkbox {
            margin-right: 8px;
        }
        
        .frequency-display {
            background: #2c3e50;
            color: #00ff00;
            font-family: 'Courier New', monospace;
            font-size: 18px;
            padding: 15px;
            border-radius: 8px;
            text-align: center;
            margin: 10px 0;
            font-weight: bold;
        }
        
        .alert {
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
        }
        
        .alert-success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        
        .alert-danger {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        
        .alert-info {
            background: #d1ecf1;
            color: #0c5460;
            border: 1px solid #bee5eb;
        }
        
        .footer {
            text-align: center;
            padding: 20px;
            background: #f8f9fa;
            border-top: 1px solid #dee2e6;
            color: #6c757d;
        }
        
        @media (max-width: 768px) {
            .panel-grid {
                grid-template-columns: 1fr;
            }
            
            .header h1 {
                font-size: 2em;
            }
            
            .content {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Saitek Flight Controller</h1>
            <p>Manage your Saitek Flight Radio, Multi, and Switch Panels</p>
        </div>
        
        <div class="content">
            <div id="alerts"></div>
            
            <div class="panel-grid">
                <!-- Radio Panel -->
                <div class="panel">
                    <div class="panel-header">
                        <h2 class="panel-title">Radio Panel</h2>
                        <div id="radio-status" class="status-indicator status-disconnected"></div>
                    </div>
                    
                    <div class="form-group">
                        <label for="com1-active">COM1 Active Frequency:</label>
                        <input type="text" id="com1-active" class="form-control" value="118.00" placeholder="e.g., 118.00">
                    </div>
                    
                    <div class="form-group">
                        <label for="com1-standby">COM1 Standby Frequency:</label>
                        <input type="text" id="com1-standby" class="form-control" value="118.50" placeholder="e.g., 118.50">
                    </div>
                    
                    <div class="form-group">
                        <label for="com2-active">COM2 Active Frequency:</label>
                        <input type="text" id="com2-active" class="form-control" value="121.30" placeholder="e.g., 121.30">
                    </div>
                    
                    <div class="form-group">
                        <label for="com2-standby">COM2 Standby Frequency:</label>
                        <input type="text" id="com2-standby" class="form-control" value="121.90" placeholder="e.g., 121.90">
                    </div>
                    
                    <button class="btn" onclick="setRadioDisplay()">Set Radio Display</button>
                    <button class="btn btn-secondary" onclick="clearRadioDisplay()">Clear Display</button>
                </div>
                
                <!-- Multi Panel -->
                <div class="panel">
                    <div class="panel-header">
                        <h2 class="panel-title">Multi Panel</h2>
                        <div id="multi-status" class="status-indicator status-disconnected"></div>
                    </div>
                    
                    <div class="form-group">
                        <label for="multi-top">Top Row Display:</label>
                        <input type="text" id="multi-top" class="form-control" value="0000" placeholder="e.g., 0000">
                    </div>
                    
                    <div class="form-group">
                        <label for="multi-bottom">Bottom Row Display:</label>
                        <input type="text" id="multi-bottom" class="form-control" value="0000" placeholder="e.g., 0000">
                    </div>
                    
                    <div class="form-group">
                        <label>Button LEDs:</label>
                        <div class="led-grid">
                            <div class="led-item">
                                <input type="checkbox" id="led-ap" class="led-checkbox">
                                <label for="led-ap">AP</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="led-hdg" class="led-checkbox">
                                <label for="led-hdg">HDG</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="led-nav" class="led-checkbox">
                                <label for="led-nav">NAV</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="led-ias" class="led-checkbox">
                                <label for="led-ias">IAS</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="led-alt" class="led-checkbox">
                                <label for="led-alt">ALT</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="led-vs" class="led-checkbox">
                                <label for="led-vs">VS</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="led-apr" class="led-checkbox">
                                <label for="led-apr">APR</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="led-rev" class="led-checkbox">
                                <label for="led-rev">REV</label>
                            </div>
                        </div>
                    </div>
                    
                    <button class="btn" onclick="setMultiDisplay()">Set Multi Display</button>
                    <button class="btn btn-secondary" onclick="clearMultiDisplay()">Clear Display</button>
                </div>
                
                <!-- Switch Panel -->
                <div class="panel">
                    <div class="panel-header">
                        <h2 class="panel-title">Switch Panel</h2>
                        <div id="switch-status" class="status-indicator status-disconnected"></div>
                    </div>
                    
                    <div class="form-group">
                        <label>Landing Gear Lights:</label>
                        <div class="led-grid">
                            <div class="led-item">
                                <input type="checkbox" id="light-green-n" class="led-checkbox">
                                <label for="light-green-n">Green N</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="light-green-l" class="led-checkbox">
                                <label for="light-green-l">Green L</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="light-green-r" class="led-checkbox">
                                <label for="light-green-r">Green R</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="light-red-n" class="led-checkbox">
                                <label for="light-red-n">Red N</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="light-red-l" class="led-checkbox">
                                <label for="light-red-l">Red L</label>
                            </div>
                            <div class="led-item">
                                <input type="checkbox" id="light-red-r" class="led-checkbox">
                                <label for="light-red-r">Red R</label>
                            </div>
                        </div>
                    </div>
                    
                    <button class="btn btn-success" onclick="setGearDown()">Gear Down (Green)</button>
                    <button class="btn btn-danger" onclick="setGearUp()">Gear Up (Red)</button>
                    <button class="btn btn-warning" onclick="setGearTransition()">Gear Transition (Yellow)</button>
                    <button class="btn btn-secondary" onclick="setAllLightsOff()">All Lights Off</button>
                </div>
            </div>
            
            <div style="text-align: center; margin-top: 30px;">
                <button class="btn" onclick="refreshStatus()">Refresh Status</button>
                <button class="btn btn-secondary" onclick="connectAll()">Reconnect All</button>
            </div>
        </div>
        
        <div class="footer">
            <p>Saitek Flight Controller - Web Interface</p>
        </div>
    </div>
    
    <script>
        // LED bit values
        const LED_AP = 0x01;
        const LED_HDG = 0x02;
        const LED_NAV = 0x04;
        const LED_IAS = 0x08;
        const LED_ALT = 0x10;
        const LED_VS = 0x20;
        const LED_APR = 0x40;
        const LED_REV = 0x80;
        
        function showAlert(message, type = 'info') {
            const alertsDiv = document.getElementById('alerts');
            const alertDiv = document.createElement('div');
            alertDiv.className = 'alert alert-' + type;
            alertDiv.textContent = message;
            alertsDiv.appendChild(alertDiv);
            
            setTimeout(function() {
                alertDiv.remove();
            }, 5000);
        }
        
        function updateStatus() {
            fetch('/api/status')
                .then(response => response.json())
                .then(data => {
                    // Update radio panel status
                    const radioStatus = document.getElementById('radio-status');
                    radioStatus.className = 'status-indicator ' + (data.radio.connected ? 'status-connected' : 'status-disconnected');
                    
                    // Update multi panel status
                    const multiStatus = document.getElementById('multi-status');
                    multiStatus.className = 'status-indicator ' + (data.multi.connected ? 'status-connected' : 'status-disconnected');
                    
                    // Update switch panel status
                    const switchStatus = document.getElementById('switch-status');
                    switchStatus.className = 'status-indicator ' + (data.switch.connected ? 'status-connected' : 'status-disconnected');
                })
                .catch(error => {
                    console.error('Error fetching status:', error);
                });
        }
        
        function setRadioDisplay() {
            const com1Active = document.getElementById('com1-active').value;
            const com1Standby = document.getElementById('com1-standby').value;
            const com2Active = document.getElementById('com2-active').value;
            const com2Standby = document.getElementById('com2-standby').value;
            
            fetch('/api/radio/set', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    com1Active: com1Active,
                    com1Standby: com1Standby,
                    com2Active: com2Active,
                    com2Standby: com2Standby
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showAlert('Radio display updated successfully!', 'success');
                } else {
                    showAlert('Failed to update radio display: ' + data.error, 'danger');
                }
            })
            .catch(error => {
                showAlert('Error updating radio display: ' + error.message, 'danger');
            });
        }
        
        function clearRadioDisplay() {
            document.getElementById('com1-active').value = '';
            document.getElementById('com1-standby').value = '';
            document.getElementById('com2-active').value = '';
            document.getElementById('com2-standby').value = '';
            setRadioDisplay();
        }
        
        function setMultiDisplay() {
            const topRow = document.getElementById('multi-top').value;
            const bottomRow = document.getElementById('multi-bottom').value;
            
            // Calculate LED value
            let leds = 0;
            if (document.getElementById('led-ap').checked) leds |= LED_AP;
            if (document.getElementById('led-hdg').checked) leds |= LED_HDG;
            if (document.getElementById('led-nav').checked) leds |= LED_NAV;
            if (document.getElementById('led-ias').checked) leds |= LED_IAS;
            if (document.getElementById('led-alt').checked) leds |= LED_ALT;
            if (document.getElementById('led-vs').checked) leds |= LED_VS;
            if (document.getElementById('led-apr').checked) leds |= LED_APR;
            if (document.getElementById('led-rev').checked) leds |= LED_REV;
            
            fetch('/api/multi/set', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    topRow: topRow,
                    bottomRow: bottomRow,
                    leds: leds
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showAlert('Multi panel display updated successfully!', 'success');
                } else {
                    showAlert('Failed to update multi panel display: ' + data.error, 'danger');
                }
            })
            .catch(error => {
                showAlert('Error updating multi panel display: ' + error.message, 'danger');
            });
        }
        
        function clearMultiDisplay() {
            document.getElementById('multi-top').value = '';
            document.getElementById('multi-bottom').value = '';
            
            // Clear all LED checkboxes
            document.querySelectorAll('.led-checkbox').forEach(checkbox => {
                checkbox.checked = false;
            });
            
            setMultiDisplay();
        }
        
        function setSwitchLights() {
            const lights = {
                greenN: document.getElementById('light-green-n').checked,
                greenL: document.getElementById('light-green-l').checked,
                greenR: document.getElementById('light-green-r').checked,
                redN: document.getElementById('light-red-n').checked,
                redL: document.getElementById('light-red-l').checked,
                redR: document.getElementById('light-red-r').checked
            };
            
            fetch('/api/switch/set', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(lights)
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showAlert('Switch panel lights updated successfully!', 'success');
                } else {
                    showAlert('Failed to update switch panel lights: ' + data.error, 'danger');
                }
            })
            .catch(error => {
                showAlert('Error updating switch panel lights: ' + error.message, 'danger');
            });
        }
        
        function setGearDown() {
            // Set all green lights on, red lights off
            document.getElementById('light-green-n').checked = true;
            document.getElementById('light-green-l').checked = true;
            document.getElementById('light-green-r').checked = true;
            document.getElementById('light-red-n').checked = false;
            document.getElementById('light-red-l').checked = false;
            document.getElementById('light-red-r').checked = false;
            setSwitchLights();
        }
        
        function setGearUp() {
            // Set all red lights on, green lights off
            document.getElementById('light-green-n').checked = false;
            document.getElementById('light-green-l').checked = false;
            document.getElementById('light-green-r').checked = false;
            document.getElementById('light-red-n').checked = true;
            document.getElementById('light-red-l').checked = true;
            document.getElementById('light-red-r').checked = true;
            setSwitchLights();
        }
        
        function setGearTransition() {
            // Set all lights on (creates yellow)
            document.getElementById('light-green-n').checked = true;
            document.getElementById('light-green-l').checked = true;
            document.getElementById('light-green-r').checked = true;
            document.getElementById('light-red-n').checked = true;
            document.getElementById('light-red-l').checked = true;
            document.getElementById('light-red-r').checked = true;
            setSwitchLights();
        }
        
        function setAllLightsOff() {
            // Set all lights off
            document.getElementById('light-green-n').checked = false;
            document.getElementById('light-green-l').checked = false;
            document.getElementById('light-green-r').checked = false;
            document.getElementById('light-red-n').checked = false;
            document.getElementById('light-red-l').checked = false;
            document.getElementById('light-red-r').checked = false;
            setSwitchLights();
        }
        
        function refreshStatus() {
            updateStatus();
            showAlert('Status refreshed', 'info');
        }
        
        function connectAll() {
            fetch('/api/connect', {
                method: 'POST'
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showAlert('Reconnected to all panels', 'success');
                    updateStatus();
                } else {
                    showAlert('Failed to reconnect: ' + data.error, 'danger');
                }
            })
            .catch(error => {
                showAlert('Error reconnecting: ' + error.message, 'danger');
            });
        }
        
        // Update status every 5 seconds
        setInterval(updateStatus, 5000);
        
        // Initial status update
        updateStatus();
    </script>
</body>
</html>
`

// Server handles the web interface
type Server struct {
	panelManager *PanelManager
	template     *template.Template
}

// NewServer creates a new server instance
func NewServer(pm *PanelManager) *Server {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}
	
	return &Server{
		panelManager: pm,
		template:     tmpl,
	}
}

// handleIndex serves the main page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	s.template.Execute(w, nil)
}

// handleStatus returns the current status of all panels
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	state := s.panelManager.GetState()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// handleRadioSet handles setting the radio panel display
func (s *Server) handleRadioSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request struct {
		COM1Active  string `json:"com1Active"`
		COM1Standby string `json:"com1Standby"`
		COM2Active  string `json:"com2Active"`
		COM2Standby string `json:"com2Standby"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	err := s.panelManager.SetRadioDisplay(request.COM1Active, request.COM1Standby, request.COM2Active, request.COM2Standby)
	
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	}
}

// handleMultiSet handles setting the multi panel display
func (s *Server) handleMultiSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request struct {
		TopRow    string `json:"topRow"`
		BottomRow string `json:"bottomRow"`
		LEDs      uint8  `json:"leds"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	err := s.panelManager.SetMultiDisplay(request.TopRow, request.BottomRow, request.LEDs)
	
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	}
}

// handleSwitchSet handles setting the switch panel lights
func (s *Server) handleSwitchSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request struct {
		GreenN bool `json:"greenN"`
		GreenL bool `json:"greenL"`
		GreenR bool `json:"greenR"`
		RedN   bool `json:"redN"`
		RedL   bool `json:"redL"`
		RedR   bool `json:"redR"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	lights := fip.LandingGearLights{
		GreenN: request.GreenN,
		GreenL: request.GreenL,
		GreenR: request.GreenR,
		RedN:   request.RedN,
		RedL:   request.RedL,
		RedR:   request.RedR,
	}
	
	err := s.panelManager.SetSwitchLights(lights)
	
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	}
}

// handleConnect handles reconnecting to all panels
func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	s.panelManager.ConnectAll()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func main() {
	var (
		port = flag.String("port", "8080", "Port to listen on")
		host = flag.String("host", "localhost", "Host to bind to (use 0.0.0.0 for network access)")
	)
	flag.Parse()
	
	// Create panel manager
	panelManager := NewPanelManager()
	
	// Connect to all panels
	fmt.Println("Connecting to Saitek panels...")
	panelManager.ConnectAll()
	
	// Create server
	server := NewServer(panelManager)
	
	// Set up routes
	http.HandleFunc("/", server.handleIndex)
	http.HandleFunc("/api/status", server.handleStatus)
	http.HandleFunc("/api/radio/set", server.handleRadioSet)
	http.HandleFunc("/api/multi/set", server.handleMultiSet)
	http.HandleFunc("/api/switch/set", server.handleSwitchSet)
	http.HandleFunc("/api/connect", server.handleConnect)
	
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Start server in a goroutine
	go func() {
		addr := *host + ":" + *port
		fmt.Printf("Starting Saitek Controller GUI on http://%s\n", addr)
		if *host == "0.0.0.0" {
			fmt.Printf("Network access enabled - other computers can connect\n")
		} else {
			fmt.Printf("Local access only\n")
		}
		fmt.Printf("Open your web browser and navigate to the URL above\n")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal("Server error:", err)
		}
	}()
	
	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down...")
	panelManager.Close()
} 