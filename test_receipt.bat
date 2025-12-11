@echo off
setlocal

rem Use PowerShell to generate the Base64 string
for /f "delims=" %%i in ('powershell -NoProfile -Command "$e=[char]0x1B; $raw = \"$e@\" + \"$e`a`1\" + \"Toko Wijaya`nJl. Widorowati, Surakarta`n\" + \"$e`a`0\" + \"================================`n\" + \"TRX-00001 Kasir : Kasir-1`n\" + \"--------------------------------`n\" + \"Beng-beng 2x3000      6000`n\" + \"--------------------------------`n\" + \"Total Item            2`n\" + \"Jumlah                6000`n\" + \"================================`n\" + \"$e`a`1\" + \"Terima Kasih sudah belanja!`n\" + \"$e`i\"; [Convert]::ToBase64String([System.Text.Encoding]::ASCII.GetBytes($raw))"') do set ENCODED=%%i

echo ===== Curl Command for Testing =====
echo.
echo curl -X POST http://localhost:3000/print ^
 -H "Content-Type: application/json" ^
 -d "{\"printerName\": \"Generic / Text Only\", \"escpos\": \"%ENCODED%\"}"
echo.
echo ===== Or use this for your actual printer =====
echo.
echo curl -X POST http://localhost:3000/print ^
 -H "Content-Type: application/json" ^
 -d "{\"printerName\": \"POS Printer\", \"escpos\": \"%ENCODED%\"}"

pause
