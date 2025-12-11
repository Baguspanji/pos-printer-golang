package posprinter

// Printer interface defines the contract for platform-specific printer implementations
type Printer interface {
	PrintRaw(printerName string, data []byte) error
	ListPrinters() ([]string, error)
}
