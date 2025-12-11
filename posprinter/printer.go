package posprinter

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	PRINTER_ENUM_LOCAL       = 0x00000002
	PRINTER_ENUM_CONNECTIONS = 0x00000004
)

var (
	winspool             = windows.NewLazySystemDLL("winspool.drv")
	procOpenPrinter      = winspool.NewProc("OpenPrinterW")
	procClosePrinter     = winspool.NewProc("ClosePrinter")
	procStartDocPrinter  = winspool.NewProc("StartDocPrinterW")
	procEndDocPrinter    = winspool.NewProc("EndDocPrinter")
	procStartPagePrinter = winspool.NewProc("StartPagePrinter")
	procEndPagePrinter   = winspool.NewProc("EndPagePrinter")
	procWritePrinter     = winspool.NewProc("WritePrinter")
	procEnumPrinters     = winspool.NewProc("EnumPrintersW")
)

type DOC_INFO_1 struct {
	DocName    *uint16
	OutputFile *uint16
	DataType   *uint16
}

type PRINTER_INFO_2 struct {
	pServerName         *uint16
	pPrinterName        *uint16
	pShareName          *uint16
	pPortName           *uint16
	pDriverName         *uint16
	pComment            *uint16
	pLocation           *uint16
	pDevMode            uintptr
	pSepFile            *uint16
	pPrintProcessor     *uint16
	pDatatype           *uint16
	pParameters         *uint16
	pSecurityDescriptor uintptr
	Attributes          uint32
	Priority            uint32
	DefaultPriority     uint32
	StartTime           uint32
	UntilTime           uint32
	Status              uint32
	cJobs               uint32
	AveragePPM          uint32
}

// PrintRaw sends raw bytes to the specified printer
func PrintRaw(printerName string, data []byte) error {
	// Convert printer name to UTF16 pointer
	printer, err := windows.UTF16PtrFromString(printerName)
	if err != nil {
		return err
	}

	var hPrinter windows.Handle

	// Open printer
	r1, _, err := procOpenPrinter.Call(
		uintptr(unsafe.Pointer(printer)),
		uintptr(unsafe.Pointer(&hPrinter)),
		0,
	)
	if r1 == 0 {
		return fmt.Errorf("failed to open printer: %v", err)
	}
	defer procClosePrinter.Call(uintptr(hPrinter))

	// Prepare document info
	docInfo := DOC_INFO_1{
		DocName:    windows.StringToUTF16Ptr("POS Print"),
		OutputFile: nil,
		DataType:   windows.StringToUTF16Ptr("RAW"),
	}

	// StartDocPrinter
	r2, _, err := procStartDocPrinter.Call(
		uintptr(hPrinter),
		1,
		uintptr(unsafe.Pointer(&docInfo)),
	)
	if r2 == 0 {
		return fmt.Errorf("failed StartDocPrinter: %v", err)
	}
	defer procEndDocPrinter.Call(uintptr(hPrinter))

	// StartPagePrinter
	r3, _, err := procStartPagePrinter.Call(uintptr(hPrinter))
	if r3 == 0 {
		return fmt.Errorf("failed StartPagePrinter: %v", err)
	}
	defer procEndPagePrinter.Call(uintptr(hPrinter))

	// WritePrinter
	var written uint32
	r4, _, err := procWritePrinter.Call(
		uintptr(hPrinter),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
	)
	if r4 == 0 {
		return fmt.Errorf("failed WritePrinter: %v", err)
	}

	return nil
}

// ListPrinters returns a list of available printers
func ListPrinters() ([]string, error) {
	var flags uint32 = PRINTER_ENUM_LOCAL | PRINTER_ENUM_CONNECTIONS

	var needed, returned uint32

	// 1. First call: query buffer size
	r1, _, err := procEnumPrinters.Call(
		uintptr(flags),
		0,
		2, // level PRINTER_INFO_2
		0,
		0,
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)),
	)
	if r1 == 0 && needed == 0 {
		return nil, fmt.Errorf("EnumPrinters query failed: %v", err)
	}

	// 2. Allocate buffer
	buf := make([]byte, needed)

	// 3. Second call: get printers
	r2, _, err := procEnumPrinters.Call(
		uintptr(flags),
		0,
		2,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(needed),
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)),
	)
	if r2 == 0 {
		return nil, fmt.Errorf("EnumPrinters failed: %v", err)
	}

	// Parse results
	printers := []string{}

	entrySize := unsafe.Sizeof(PRINTER_INFO_2{})
	for i := uint32(0); i < returned; i++ {
		entry := (*PRINTER_INFO_2)(unsafe.Pointer(&buf[int(i)*int(entrySize)]))
		name := windows.UTF16PtrToString(entry.pPrinterName)
		if name != "" {
			printers = append(printers, name)
		}
	}

	return printers, nil
}
