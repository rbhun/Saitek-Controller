# Saitek Controller Makefile

.PHONY: build run test clean examples fip-example standalone package release

# Build the main application
build:
	go build -o bin/saitek-controller cmd/main.go

# Run the main application
run: build
	./bin/saitek-controller

# Run with specific parameters
run-fip:
	go run cmd/main.go -title "FIP Test" -width 320 -height 240 -instrument artificial_horizon

run-image:
	go run cmd/main.go -title "FIP Image" -width 320 -height 240 -image assets/test.png

# Run examples
examples:
	go run examples/fip_example.go

fip-example:
	go run examples/fip_example.go

switch-example:
	go run examples/switch/main.go

test-switch:
	go run cmd/test_switch/main.go

# Test the application
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download
	go mod tidy

# Create assets directory and sample images
setup:
	mkdir -p assets
	mkdir -p bin

# Generate test images
generate-images:
	go run cmd/generate_images.go -output assets -width 320 -height 240

# Build all examples
build-examples:
	go build -o bin/fip-example examples/fip_example.go
	go build -o bin/switch-example examples/switch/main.go
	go build -o bin/test-switch cmd/test_switch/main.go

# Run with different instruments
run-artificial-horizon:
	go run cmd/main.go -instrument artificial_horizon

run-airspeed:
	go run cmd/main.go -instrument airspeed

run-altimeter:
	go run cmd/main.go -instrument altimeter

run-compass:
	go run cmd/main.go -instrument compass

run-vsi:
	go run cmd/main.go -instrument vsi

run-turn-coordinator:
	go run cmd/main.go -instrument turn_coordinator

# Build the GUI application
build-gui:
	go build -o bin/saitek-controller-gui cmd/saitek-controller-gui/main.go

# Run the GUI application
run-gui: build-gui
	./bin/saitek-controller-gui

# Run GUI with custom port
run-gui-port:
	go run cmd/saitek-controller-gui/main.go -port 8080

# Build standalone radio panel program
build-radio:
	go build -o bin/set-radio cmd/set_radio_dir/main.go

# Build all standalone programs
build-standalone: build-gui build-radio
	go build -o bin/saitek-controller cmd/main.go

# Create standalone application package
standalone: build-standalone
	@echo "Creating standalone application package..."
	@mkdir -p dist/saitek-controller
	@cp -r bin/* dist/saitek-controller/
	@cp -r assets dist/saitek-controller/
	@cp README.md dist/saitek-controller/
	@cp go.mod dist/saitek-controller/
	@cp go.sum dist/saitek-controller/
	@cp -r internal dist/saitek-controller/
	@cp -r cmd dist/saitek-controller/
	@cp -r examples dist/saitek-controller/
	@cp -r docs dist/saitek-controller/
	@echo "Standalone package created in dist/saitek-controller/"

# Create release package (minimal)
package: build-standalone
	@echo "Creating release package..."
	@mkdir -p dist/saitek-controller-release
	@cp bin/saitek-controller-gui dist/saitek-controller-release/
	@cp bin/set-radio dist/saitek-controller-release/
	@cp bin/saitek-controller dist/saitek-controller-release/
	@cp -r assets dist/saitek-controller-release/
	@cp README.md dist/saitek-controller-release/
	@echo "Release package created in dist/saitek-controller-release/"

# Create complete release with documentation
release: package
	@echo "Creating complete release package..."
	@mkdir -p dist/saitek-controller-complete
	@cp -r dist/saitek-controller-release/* dist/saitek-controller-complete/
	@cp -r docs dist/saitek-controller-complete/
	@cp -r examples dist/saitek-controller-complete/
	@cp Makefile dist/saitek-controller-complete/
	@cp go.mod dist/saitek-controller-complete/
	@cp go.sum dist/saitek-controller-complete/
	@echo "Complete release package created in dist/saitek-controller-complete/"

# Create macOS app bundle
macos-app: build-standalone
	@echo "Creating macOS app bundle..."
	@mkdir -p dist/SaitekController.app/Contents/MacOS
	@mkdir -p dist/SaitekController.app/Contents/Resources
	@cp bin/saitek-controller-gui dist/SaitekController.app/Contents/MacOS/
	@cp -r assets dist/SaitekController.app/Contents/Resources/
	@echo '<?xml version="1.0" encoding="UTF-8"?>' > dist/SaitekController.app/Contents/Info.plist
	@echo '<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' >> dist/SaitekController.app/Contents/Info.plist
	@echo '<plist version="1.0">' >> dist/SaitekController.app/Contents/Info.plist
	@echo '<dict>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <key>CFBundleExecutable</key>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <string>saitek-controller-gui</string>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <key>CFBundleIdentifier</key>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <string>com.saitek.controller</string>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <key>CFBundleName</key>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <string>Saitek Controller</string>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <key>CFBundleVersion</key>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <string>1.0</string>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <key>CFBundleShortVersionString</key>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <string>1.0</string>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <key>LSMinimumSystemVersion</key>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '    <string>10.15</string>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '</dict>' >> dist/SaitekController.app/Contents/Info.plist
	@echo '</plist>' >> dist/SaitekController.app/Contents/Info.plist
	@echo "macOS app bundle created in dist/SaitekController.app/"

# Clean all build artifacts
clean-all: clean
	rm -rf dist/

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build the main application"
	@echo "  run                - Run the main application"
	@echo "  build-gui          - Build the GUI application"
	@echo "  run-gui            - Run the GUI application"
	@echo "  run-gui-port       - Run GUI with custom port"
	@echo "  build-radio        - Build standalone radio panel program"
	@echo "  build-standalone   - Build all standalone programs"
	@echo "  standalone         - Create complete standalone package"
	@echo "  package            - Create minimal release package"
	@echo "  release            - Create complete release package"
	@echo "  macos-app          - Create macOS app bundle"
	@echo "  run-fip            - Run with artificial horizon"
	@echo "  run-image          - Run with image file"
	@echo "  examples           - Run FIP example"
	@echo "  switch-example     - Run switch panel example"
	@echo "  test-switch        - Run switch panel test"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  clean-all          - Clean all build artifacts"
	@echo "  deps               - Install dependencies"
	@echo "  setup              - Create directories"
	@echo "  build-examples     - Build example programs"
	@echo "  run-artificial-horizon - Run artificial horizon"
	@echo "  run-airspeed       - Run airspeed indicator"
	@echo "  run-altimeter      - Run altimeter"
	@echo "  run-compass        - Run compass"
	@echo "  run-vsi            - Run vertical speed indicator"
	@echo "  run-turn-coordinator - Run turn coordinator" 