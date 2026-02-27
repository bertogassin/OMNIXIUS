@echo off
title OMNIXIUS Rust (stack 1)
cd /d "%~dp0services\rust"
echo.
echo   Rust service:  http://localhost:8081
echo   Stack order: Rust ^> Go ^> ...
echo.
echo Press Ctrl+C to stop.
echo.
cargo run
