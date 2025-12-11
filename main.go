package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/baguspanji/pos-printer/posprinter"
)

func main() {
	http.HandleFunc("/print", posprinter.PrintHandler)
	http.HandleFunc("/printers", posprinter.PrintersHandler)
	http.HandleFunc("/test", posprinter.TestHandler)
	http.HandleFunc("/health", posprinter.HealthHandler)

	fmt.Println("Printer middleware running on 0.0.0.0:3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
