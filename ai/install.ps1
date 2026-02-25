# Install AI dependencies (--user to avoid Scripts permission errors on Windows)
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptDir
Write-Host "Installing into user site-packages (pip --user)..." -ForegroundColor Cyan
pip install --user -r requirements.txt
if ($LASTEXITCODE -eq 0) {
    Write-Host "Done. Run start-ai.bat or: python -m uvicorn main:app --reload --port 8000" -ForegroundColor Green
} else {
    Write-Host "Install failed. Try: Run PowerShell as Administrator, then pip install -r requirements.txt" -ForegroundColor Yellow
    exit 1
}
