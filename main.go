package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/baguspanji/pos-printer/middleware"
	"github.com/baguspanji/pos-printer/posprinter"
)

func main() {
	http.HandleFunc("/print", middleware.CORS(posprinter.PrintHandler))
	http.HandleFunc("/printers", middleware.CORS(posprinter.PrintersHandler))
	http.HandleFunc("/test", middleware.CORS(posprinter.TestHandler))
	http.HandleFunc("/health", middleware.CORS(posprinter.HealthHandler))

	fmt.Println("Printer middleware running on 0.0.0.0:3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
