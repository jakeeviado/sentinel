# Sentinel - PowerShell Setup Script
# This script fixes the go.sum checksum issue and builds the project

Write-Host "===============================================" -ForegroundColor Cyan
Write-Host "  ⌀ Sentinel Setup - PowerShell Fix" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: REMOVE THE OLD go.sum
Write-Host "Step 1: Removing old go.sum file..." -ForegroundColor Yellow
if (Test-Path "go.sum")
{
    Remove-Item "go.sum" -Force
    Write-Host "  [OK] Removed old go.sum" -ForegroundColor Green
} else
{
    Write-Host "  [OK] No go.sum to remove" -ForegroundColor Green
}
Write-Host ""

# Step 2: GENERATE A FRESH go.sum
Write-Host "Step 2: Generating fresh go.sum..." -ForegroundColor Yellow
$output = go mod tidy 2>&1
if ($LASTEXITCODE -ne 0)
{
    Write-Host "  [ERROR] go mod tidy failed!" -ForegroundColor Red
    Write-Host $output
    Read-Host "Press Enter to exit"
    exit 1
}
Write-Host "  [OK] Dependencies downloaded" -ForegroundColor Green
Write-Host ""

# Step 3: BUILD
Write-Host "Step 3: Building Sentinel..." -ForegroundColor Yellow
$output = go build -o sentinel.exe 2>&1
if ($LASTEXITCODE -ne 0)
{
    Write-Host "  [ERROR] Build failed!" -ForegroundColor Red
    Write-Host $output
    Read-Host "Press Enter to exit"
    exit 1
}
Write-Host "  [OK] Build successful!" -ForegroundColor Green
Write-Host ""

# SUCCESS
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host "  Hell Yeah! Setup has been completed!" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Test it with:" -ForegroundColor Yellow
Write-Host "  .\sentinel.exe --version"
Write-Host "  .\sentinel.exe scan --path .\examples --verbose"
Write-Host ""
Read-Host "Press Enter to exit"
