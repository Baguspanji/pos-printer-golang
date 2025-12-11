# POS Printer Middleware

A cross-platform HTTP middleware for printing to POS (Point of Sale) thermal printers. Supports Windows, macOS, and Linux.

## Features

- 🖨️ **Multi-Platform Support**: Works on Windows, macOS, and Linux
- 🔌 **Simple HTTP API**: Easy integration with any application
- 📝 **ESC/POS Compatible**: Supports standard ESC/POS commands
- 🚀 **Lightweight**: Single binary with no external dependencies (except CUPS on Unix-like systems)
- 💚 **Health Monitoring**: Built-in health check endpoint for monitoring and load balancers

## Platform Support

| Platform | Implementation | Requirements |
|----------|---------------|--------------|
| Windows  | Windows Spooler API | None (native) |
| macOS    | CUPS | Pre-installed on macOS |
| Linux    | CUPS | Usually pre-installed |

## Installation

### Download Pre-built Binaries

Download the appropriate binary for your platform:

- **Windows**: `pos-printer-windows.exe`
- **macOS (Intel)**: `pos-printer-darwin-amd64`
- **macOS (Apple Silicon)**: `pos-printer-darwin-arm64`
- **Linux (x64)**: `pos-printer-linux-amd64`
- **Linux (ARM64)**: `pos-printer-linux-arm64`

### Build from Source

Requirements:
- Go 1.24 or later

```bash
# Clone the repository
git clone <repository-url>
cd pos-printer

# Build for your current platform
go build

# Or use the Makefile for cross-compilation
make build-all
```

## Usage

### Starting the Server

```bash
# On Windows
./pos-printer-windows.exe

# On macOS/Linux
./pos-printer-darwin-amd64  # or appropriate binary
```

The server will start on `http://0.0.0.0:3000`

### API Endpoints

#### 1. List Available Printers

```bash
GET /printers
```

**Response:**
```json
["POS-80", "Thermal Printer", "Another Printer"]
```

#### 2. Print Raw ESC/POS Data

```bash
POST /print
Content-Type: application/json

{
  "printerName": "POS-80",
  "escpos": "<base64-encoded-escpos-data>"
}
```

**Response:**
```json
{"success": true}
```

#### 3. Test Print

```bash
GET /test
```

Prints a test receipt to the default printer (POS-80).

**Response:**
```json
{
  "status": "ok",
  "message": "Test print sent to printer"
}
```

#### 4. Health Check

```bash
GET /health
```

Check server health and status. Useful for monitoring and load balancers.

**Response:**
```json
{
  "status": "healthy",
  "uptime": "1h23m45s",
  "platform": "darwin/arm64",
  "printersCount": 1,
  "serverTime": "2025-12-11T20:56:38+07:00"
}
```

## Platform-Specific Notes

### Windows

- Uses the native Windows Spooler API
- No additional dependencies required
- Supports all printers installed in Windows

### macOS

- Uses CUPS (Common Unix Printing System)
- CUPS is pre-installed on macOS
- Printers must be configured in System Preferences

To list printers:
```bash
lpstat -p
```

To add a printer:
```bash
System Preferences → Printers & Scanners → Add Printer
```

### Linux

- Uses CUPS (Common Unix Printing System)
- CUPS is usually pre-installed on most distributions
- If not installed, install via package manager:

```bash
# Ubuntu/Debian
sudo apt-get install cups

# Fedora/RHEL
sudo dnf install cups

# Arch Linux
sudo pacman -S cups
```

To list printers:
```bash
lpstat -p
```

To add a printer:
```bash
sudo lpadmin -p PrinterName -E -v usb://path/to/printer
```

## Development

### Project Structure

```
pos-printer/
├── main.go                      # HTTP server entry point
├── posprinter/
│   ├── handler.go              # HTTP handlers (platform-agnostic)
│   ├── printer_interface.go   # Printer interface definition
│   ├── printer_windows.go     # Windows implementation
│   ├── printer_darwin.go      # macOS implementation
│   └── printer_linux.go       # Linux implementation
├── Makefile                    # Build automation
└── README.md
```

### Build Tags

The project uses Go build tags to compile platform-specific code:

- `//go:build windows` - Windows-only code
- `//go:build darwin` - macOS-only code
- `//go:build linux` - Linux-only code

### Building for Specific Platforms

```bash
# Build for Windows
make build-windows

# Build for macOS
make build-darwin

# Build for Linux
make build-linux

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

## Example Usage

### Python Example

```python
import requests
import base64

# List printers
response = requests.get('http://localhost:3000/printers')
printers = response.json()
print(f"Available printers: {printers}")

# Print receipt
escpos_data = b'\x1B\x40'  # ESC @ - Initialize
escpos_data += b'Hello, World!\n'
escpos_data += b'\x1B\x69'  # ESC i - Cut paper

encoded = base64.b64encode(escpos_data).decode('utf-8')

response = requests.post('http://localhost:3000/print', json={
    'printerName': 'POS-80',
    'escpos': encoded
})
print(response.json())
```

### JavaScript Example

```javascript
// List printers
fetch('http://localhost:3000/printers')
  .then(res => res.json())
  .then(printers => console.log('Available printers:', printers));

// Print receipt
const escposData = new Uint8Array([
  0x1B, 0x40,  // ESC @ - Initialize
  ...new TextEncoder().encode('Hello, World!\n'),
  0x1B, 0x69   // ESC i - Cut paper
]);

const base64Data = btoa(String.fromCharCode(...escposData));

fetch('http://localhost:3000/print', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    printerName: 'POS-80',
    escpos: base64Data
  })
})
  .then(res => res.json())
  .then(data => console.log(data));
```

## Troubleshooting

### Windows

**Problem**: Printer not found
- Ensure the printer is installed and visible in "Devices and Printers"
- Check the exact printer name (case-sensitive)

### macOS/Linux

**Problem**: "failed to list printers" or "failed to print"
- Ensure CUPS is running: `sudo systemctl status cups` (Linux) or check System Preferences (macOS)
- Verify printer is configured: `lpstat -p`
- Check printer permissions: `lpstat -v`

**Problem**: Permission denied
- On Linux, ensure your user is in the `lp` or `lpadmin` group:
  ```bash
  sudo usermod -a -G lp $USER
  ```

## License

[Your License Here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
