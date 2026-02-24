# OMNIXIUS Backend — запуск из PowerShell
$root = $PSScriptRoot
Set-Location "$root\backend-go"
Write-Host "Starting API on http://localhost:3000"
Write-Host "Press Ctrl+C to stop."
Write-Host ""
go run .
