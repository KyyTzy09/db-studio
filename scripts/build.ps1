# DBStudio Automated Build Script for Windows PowerShell
$ErrorActionPreference = "Stop"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " 🚀 Starting DBStudio Full Build Pipeline " -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

$RepoRoot = Split-Path -Path $PSScriptRoot -Parent

# Step 1: Build Frontend (SvelteKit)
Write-Host "`n[1/3] Building SvelteKit Frontend Static Assets..." -ForegroundColor Yellow
Set-Location "$RepoRoot\web"
bun run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Frontend build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Frontend build completed successfully." -ForegroundColor Green

# Step 2: Ensure bin/ directory exists
Set-Location $RepoRoot
if (!(Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Step 3: Cross-Compile Go Binaries
Write-Host "`n[2/3] Compiling Go Single Binaries (CGO_ENABLED=0)..." -ForegroundColor Yellow

$env:CGO_ENABLED = "0"

# Target 1: Windows x64
Write-Host " -> Building Windows x64 (bin/dbstudio-win-x64.exe)..." -ForegroundColor Gray
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o bin/dbstudio-win-x64.exe main.go
Copy-Item bin/dbstudio-win-x64.exe bin/dbstudio.exe -Force

# Target 2: Linux x64
Write-Host " -> Building Linux x64 (bin/dbstudio-linux-x64)..." -ForegroundColor Gray
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o bin/dbstudio-linux-x64 main.go

# Target 3: macOS ARM64 (Apple Silicon)
Write-Host " -> Building macOS ARM64 (bin/dbstudio-darwin-arm64)..." -ForegroundColor Gray
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -ldflags="-s -w" -o bin/dbstudio-darwin-arm64 main.go

# Target 4: macOS x64 (Intel)
Write-Host " -> Building macOS x64 (bin/dbstudio-darwin-x64)..." -ForegroundColor Gray
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o bin/dbstudio-darwin-x64 main.go

# Reset environment variables
Remove-Item env:GOOS -ErrorAction SilentlyContinue
Remove-Item env:GOARCH -ErrorAction SilentlyContinue

Write-Host "`n[3/3] Build Summary:" -ForegroundColor Yellow
Get-ChildItem -Path "bin" | Select-Object Name, @{Name="Size (MB)"; Expression={[math]::round($_.Length / 1MB, 2)}} | Format-Table -AutoSize

Write-Host "`n🎉 All DBStudio binaries compiled successfully in bin/!" -ForegroundColor Green
