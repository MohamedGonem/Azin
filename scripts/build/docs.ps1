<#
.SYNOPSIS
    Generates static API documentation using doc2go.
.DESCRIPTION
    Scans Go packages in the repository and generates static HTML documentation.
    Works regardless of current working directory.
.PARAMETER OutputDir
    Target directory for generated documentation (default: "docs/api").
.PARAMETER IncludeInternal
    Includes internal packages in the generated output (default: $true).
.EXAMPLE
    .\scripts\docs\doc.ps1
.EXAMPLE
    .\scripts\docs\doc.ps1 -OutputDir "public/docs"
#>
[CmdletBinding()]
param (
    [string]$OutputDir = "docs/api",
    [bool]$IncludeInternal = $true
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
Set-Location $RepoRoot

$Doc2GoCmd = Get-Command "doc2go" -ErrorAction SilentlyContinue

if (-not $Doc2GoCmd) {
    if (Get-Command "go" -ErrorAction SilentlyContinue) {
        Write-Host "'doc2go' binary not found in PATH. Falling back to 'go run github.com/abhinav/doc2go@latest'..." -ForegroundColor Yellow
    } else {
        Write-Error "'doc2go' and 'go' were not found in PATH. Please install doc2go: go install github.com/abhinav/doc2go@latest"
        exit 1
    }
}

$FullOutputDir = Join-Path $RepoRoot $OutputDir
if (-not (Test-Path $FullOutputDir)) {
    New-Item -ItemType Directory -Path $FullOutputDir -Force | Out-Null
}

$DocArgs = @()
if ($IncludeInternal) {
    $DocArgs += "-internal"
}
$DocArgs += "-out", $FullOutputDir
$DocArgs += "./..."

Write-Host "Generating API documentation..." -ForegroundColor Cyan

$Timer = [System.Diagnostics.Stopwatch]::StartNew()

if ($Doc2GoCmd) {
    & doc2go @DocArgs
} else {
    & go run github.com/abhinav/doc2go@latest @DocArgs
}

$Timer.Stop()

Write-Host ("Successfully generated documentation in $OutputDir ({0:N2}s)" -f $Timer.Elapsed.TotalSeconds) -ForegroundColor Green