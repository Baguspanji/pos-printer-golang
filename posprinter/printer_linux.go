//go:build linux

package posprinter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PrintRaw sends raw bytes to the specified printer using CUPS
func PrintRaw(printerName string, data []byte) error {
	// Create a temporary file to hold the raw data
	tmpFile, err := os.CreateTemp("", "pos-print-*.raw")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("failed to write to temp file: %v", err)
	}

	// Close the file before printing
	tmpFile.Close()

	// Use lp command to print
	// -d specifies the printer name
	// -o raw tells CUPS to send the data without processing
	cmd := exec.Command("lp", "-d", printerName, "-o", "raw", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to print: %v, output: %s", err, string(output))
	}

	return nil
}

// ListPrinters returns a list of available printers using CUPS
func ListPrinters() ([]string, error) {
	// Use lpstat to list printers
	// -p lists all printers
	// -d shows the default printer
	cmd := exec.Command("lpstat", "-p")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list printers: %v, output: %s", err, string(output))
	}

	// Parse output
	// Format: "printer PrinterName is idle. enabled since ..."
	printers := []string{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract printer name from "printer NAME ..."
		if strings.HasPrefix(line, "printer ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				printerName := parts[1]
				printers = append(printers, printerName)
			}
		}
	}

	return printers, nil
}
