#!/bin/bash

# Generate receipt data
RECEIPT_TEXT="Toko Wijaya
Jl. Widorowati, Surakarta
================================
TRX-00001 Kasir : Kasir-1
--------------------------------
Beng-beng 2x3000      6000
--------------------------------
Total Item            2
Jumlah                6000
================================
Terima Kasih sudah belanja!"

# Simple ESC/POS commands (base64 encoded)
# ESC @ (initialize) + Text + Line feeds
ESCPOS_BASE64="G0BUZW5nLWJlbmcgMngzMDAwICAgICAgNjAwMAoKVG90YWwgSXRlbSAgICAgICAgICAgIDIKSnVtbGFoICAgICAgICAgICAgIDYwMDAK"

# For a more complete example with center alignment
# Using printf to create ESC/POS commands properly
ESCPOS_RAW=$(printf '\x1B\x40')  # Initialize printer
ESCPOS_RAW+=$(printf '\x1B\x61\x01')  # Center alignment
ESCPOS_RAW+="Toko Wijaya"
ESCPOS_RAW+=$(printf '\n')
ESCPOS_RAW+="Jl. Widorowati, Surakarta"
ESCPOS_RAW+=$(printf '\n\x1B\x61\x00')  # Left alignment
ESCPOS_RAW+="================================"
ESCPOS_RAW+=$(printf '\n')
ESCPOS_RAW+="TRX-00001 Kasir : Kasir-1"
ESCPOS_RAW+=$(printf '\n')
ESCPOS_RAW+="--------------------------------"
ESCPOS_RAW+=$(printf '\n')
ESCPOS_RAW+="Beng-beng 2x3000      6000"
ESCPOS_RAW+=$(printf '\n')
ESCPOS_RAW+="--------------------------------"
ESCPOS_RAW+=$(printf '\n')
ESCPOS_RAW+="Total Item            2"
ESCPOS_RAW+=$(printf '\n')
ESCPOS_RAW+="Jumlah                6000"
ESCPOS_RAW+=$(printf '\n')
ESCPOS_RAW+="================================"
ESCPOS_RAW+=$(printf '\n\x1B\x61\x01')  # Center alignment
ESCPOS_RAW+="Terima Kasih sudah belanja!"
ESCPOS_RAW+=$(printf '\n\x1B\x69')  # Cut paper

# Encode to base64
ENCODED=$(echo -n "$ESCPOS_RAW" | base64)

echo "===== Curl Command for Testing ====="
echo ""
echo "curl -X POST http://localhost:3000/print \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"printerName\": \"Generic / Text Only\", \"escpos\": \"$ENCODED\"}'"
echo ""
echo "===== Or use this for your actual printer ====="
echo ""
echo "curl -X POST http://localhost:3000/print \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"printerName\": \"POS Printer\", \"escpos\": \"$ENCODED\"}'"
