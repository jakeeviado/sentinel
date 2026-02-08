@echo off
REM Sentinel - Windows Setup Script
REM This script fixes the go.sum checksum issue and builds the project

echo ===============================================
echo        ⌀ Sentinel Setup - Windows Fix
echo ===============================================
echo.

echo Step 1: Removing old go.sum file...
if exist go.sum (
    del go.sum
    echo   [OK] Removed old go.sum
) else (
    echo   [OK] No go.sum to remove
)
echo.

echo Step 2: Generating fresh go.sum...
go mod tidy
if %errorlevel% neq 0 (
    echo   [ERROR] go mod tidy failed!
    pause
    exit /b 1
)
echo   [OK] Dependencies downloaded
echo.

echo Step 3: Building Sentinel...
go build -o sentinel.exe
if %errorlevel% neq 0 (
    echo   [ERROR] Build failed!
    pause
    exit /b 1
)
echo   [OK] Build successful!
echo.

echo ===============================================
echo               Setup Complete!
echo ===============================================
echo.
echo Test it with:
echo   sentinel.exe --version
echo   sentinel.exe scan --path ./examples --verbose
echo.
pause
