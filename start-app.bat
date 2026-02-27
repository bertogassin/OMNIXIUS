@echo off
title OMNIXIUS Launcher
cd /d "%~dp0"

(
echo App:     http://localhost:5173/app/
echo Backend: http://localhost:3000
echo.
echo Test login: test@test.com / Test123!
) > "%~dp0urls.txt"
echo Urls saved to urls.txt - open and copy from there.
echo.

echo Starting backend in new window...
start "OMNIXIUS Backend" cmd /k "cd /d %~dp0backend-go && echo Backend running. Close this window to stop. && echo. && go run ."

echo Waiting 10 sec for backend...
timeout /t 10 /nobreak >nul

echo Starting frontend in new window...
start "OMNIXIUS Web" cmd /k "cd /d %~dp0web && echo Frontend running. Close this window to stop. && echo. && npm run dev"

echo Waiting 8 sec for Vite...
timeout /t 8 /nobreak >nul

echo Opening browser...
start "" "http://localhost:5173/app/"

echo.
echo Done. If browser did not open, go to: http://localhost:5173/app/
echo Login: test@test.com / Test123! (click Create test user first if DB is empty)
echo.
pause
