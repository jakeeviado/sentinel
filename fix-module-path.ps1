# Fix Module Path Script
# This replaces the placeholder module name with your actual module name

param(
    [Parameter(Mandatory=$true)]
    [string]$ModuleName
)

Write-Host "Updating module paths to: $ModuleName" -ForegroundColor Cyan
Write-Host ""

Write-Host "Updating go.mod..." -ForegroundColor Yellow
(Get-Content go.mod) -replace 'module github.com/jakeeviado/sentinel', "module $ModuleName" | Set-Content go.mod

$goFiles = Get-ChildItem -Path . -Filter *.go -Recurse

Write-Host "Updating imports in .go files..." -ForegroundColor Yellow
foreach ($file in $goFiles)
{
    $content = Get-Content $file.FullName -Raw
    $newContent = $content -replace 'github.com/jakeeviado/sentinel', $ModuleName

    if ($content -ne $newContent)
    {
        Set-Content -Path $file.FullName -Value $newContent
        Write-Host "  Updated: $($file.Name)" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "Done! Now run:" -ForegroundColor Cyan
Write-Host "  go mod tidy" -ForegroundColor White
Write-Host "  go build -o sentinel.exe" -ForegroundColor White
