@echo off
title OMNIXIUS Backend
cd /d "%~dp0backend-go"
echo.
echo   Site + API:  http://localhost:3000
echo.
echo Press Ctrl+C to stop.
echo.
go run .
