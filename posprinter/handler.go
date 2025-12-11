package posprinter

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

type PrintRequest struct {
	PrinterName  string `json:"printerName"`
	EscposBase64 string `json:"escpos"`
}

// PrintHandler handles HTTP requests for printing
func PrintHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var req PrintRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if req.PrinterName == "" {
		http.Error(w, "printerName required", 400)
		return
	}

	raw, err := base64.StdEncoding.DecodeString(req.EscposBase64)
	if err != nil {
		http.Error(w, "Invalid base64", 400)
		return
	}

	err = PrintRaw(req.PrinterName, raw)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte(`{"success":true}`))
}

// PrintersHandler handles HTTP requests to list available printers
func PrintersHandler(w http.ResponseWriter, r *http.Request) {
	printers, err := ListPrinters()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp, _ := json.Marshal(printers)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// TestHandler handles GET /test for testing printer printing
func TestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get printer name from query parameter
	printerName := "POS-80" // default printer name
	receiptText := `Toko Wijaya
Jl. Widorowati, Surakarta
================================
TRX-00001 Kasir : Kasir-1
--------------------------------
Beng-beng 2x3000      6000
--------------------------------
Total Item            2
Jumlah                6000
================================
Terima Kasih sudah belanja!`

	// Create test ESC/POS commands (initialization + text + line feed + cut)
	testContent := []byte{0x1B, 0x40} // ESC @ - Initialize printer
	testContent = append(testContent, []byte(receiptText)...)
	testContent = append(testContent, 0x0A, 0x0A, 0x0A) // Line feeds
	testContent = append(testContent, 0x1B, 0x69)       // ESC i - Partial cut

	err := PrintRaw(printerName, testContent)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","message":"Test print sent to printer"}`))
}
