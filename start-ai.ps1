# OMNIXIUS AI — run with python -m uvicorn (no need for uvicorn.exe in PATH)
$ErrorActionPreference = "Stop"
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location (Join-Path $scriptDir "ai")

Write-Host ""
Write-Host "  AI:  http://localhost:8000" -ForegroundColor Cyan
Write-Host "  Docs: http://localhost:8000/docs" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press Ctrl+C to stop." -ForegroundColor Gray
Write-Host ""

try {
    python -m uvicorn main:app --reload --port 8000
} catch {
    Write-Host ""
    Write-Host "If 'No module named uvicorn': install first:" -ForegroundColor Yellow
    Write-Host "  pip install --user -r requirements.txt" -ForegroundColor White
    Write-Host "Then run this script again." -ForegroundColor Yellow
    exit 1
}
