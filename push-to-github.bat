@echo off
title Push to GitHub
cd /d "%~dp0"

echo Adding all changes...
git add -A

echo Committing...
set MSG=Update %date% %time%
if not "%~1"=="" set MSG=%~1
git commit -m "%MSG%" 2>nul
if errorlevel 1 (
  echo Nothing to commit or commit failed.
) else (
  echo Pushing to GitHub...
  git push origin main
  if errorlevel 1 (
    echo Push failed. Check remote and branch.
  ) else (
    echo Done. Changes are on GitHub.
  )
)
echo.
pause
