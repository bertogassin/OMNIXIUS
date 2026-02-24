@echo off
title OMNIXIUS Backend
cd /d "%~dp0backend-go"
echo Starting API on http://localhost:3000
echo Press Ctrl+C to stop.
echo.
go run .
