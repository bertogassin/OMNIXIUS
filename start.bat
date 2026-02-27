@echo off
title OMNIXIUS — Stack order: Rust ^> Go
set RUST_SERVICE_URL=http://localhost:8081
echo.
echo   Stack order: 1. Rust (8081)  2. Go (3000)
echo   Starting Rust first...
echo.
start "OMNIXIUS Rust" cmd /c "%~dp0start-rust.bat"
timeout /t 3 /nobreak >nul
echo   Starting Go API...
echo.
call "%~dp0start-backend.bat"
