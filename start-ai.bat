@echo off
title OMNIXIUS AI
cd /d "%~dp0ai"
echo.
echo   AI:  http://localhost:8000
echo   Docs: http://localhost:8000/docs
echo.
echo Press Ctrl+C to stop.
echo.
python -m uvicorn main:app --reload --port 8000
if errorlevel 1 (
  echo.
  echo If "No module named uvicorn": run in this folder:
  echo   pip install --user -r requirements.txt
  echo Then run this file again.
  pause
)
