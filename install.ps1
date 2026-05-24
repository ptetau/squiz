#!/usr/bin/env pwsh
# Install squiz on Windows: binary on PATH + SKILL.md into ~\.claude\skills\squiz\.
#
# Usage:
#   irm https://raw.githubusercontent.com/ptetau/squiz/main/install.ps1 | iex
#
# Env overrides:
#   $env:VERSION          pin a specific version (default: latest GitHub release)
#   $env:SQUIZ_BIN_DIR    where to install squiz.exe (default: %LOCALAPPDATA%\Programs\squiz)
#   $env:SQUIZ_SKILL_DIR  where to install SKILL.md (default: %USERPROFILE%\.claude\skills\squiz)

$ErrorActionPreference = 'Stop'

$Owner    = 'ptetau'
$Repo     = 'squiz'
$BinDir   = if ($env:SQUIZ_BIN_DIR)   { $env:SQUIZ_BIN_DIR }   else { "$env:LOCALAPPDATA\Programs\squiz" }
$SkillDir = if ($env:SQUIZ_SKILL_DIR) { $env:SQUIZ_SKILL_DIR } else { "$env:USERPROFILE\.claude\skills\squiz" }

$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
  'AMD64' { 'x86_64' }
  'ARM64' { 'arm64' }
  default { throw "unsupported arch: $($env:PROCESSOR_ARCHITECTURE)" }
}

$version = $env:VERSION
if (-not $version) {
  $latest = Invoke-RestMethod "https://api.github.com/repos/$Owner/$Repo/releases/latest"
  $version = $latest.tag_name -replace '^v',''
  if (-not $version) { throw "could not resolve latest version (set `$env:VERSION to pin)" }
}

$archive = "squiz_${version}_Windows_${arch}.zip"
$baseUrl = "https://github.com/$Owner/$Repo/releases/download/v${version}"

$tmp = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "squiz-install-$(Get-Random)") -Force
try {
  Write-Host "-> downloading $archive"
  Invoke-WebRequest -Uri "$baseUrl/$archive"      -OutFile (Join-Path $tmp $archive)       -UseBasicParsing
  Invoke-WebRequest -Uri "$baseUrl/checksums.txt" -OutFile (Join-Path $tmp 'checksums.txt') -UseBasicParsing

  Write-Host "-> verifying checksum"
  $line     = Get-Content (Join-Path $tmp 'checksums.txt') | Where-Object { $_ -match " $archive$" } | Select-Object -First 1
  if (-not $line) { throw "no checksum entry for $archive" }
  $expected = ($line -split '\s+')[0].ToLower()
  $actual   = (Get-FileHash (Join-Path $tmp $archive) -Algorithm SHA256).Hash.ToLower()
  if ($expected -ne $actual) {
    throw "checksum mismatch for $archive`n  expected: $expected`n  actual:   $actual"
  }

  Write-Host "-> extracting"
  Expand-Archive -Path (Join-Path $tmp $archive) -DestinationPath $tmp -Force

  New-Item -ItemType Directory -Path $BinDir, $SkillDir -Force | Out-Null
  Move-Item -Path (Join-Path $tmp 'squiz.exe') -Destination (Join-Path $BinDir 'squiz.exe') -Force
  Move-Item -Path (Join-Path $tmp 'SKILL.md')  -Destination (Join-Path $SkillDir 'SKILL.md') -Force

  # Add BinDir to user PATH if missing.
  $userPath = [Environment]::GetEnvironmentVariable('Path','User')
  if ($userPath -notlike "*$BinDir*") {
    $newPath = if ($userPath) { "$userPath;$BinDir" } else { $BinDir }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host "-> added $BinDir to user PATH (open a new shell to pick it up)"
  }

  Write-Host ""
  Write-Host "OK installed squiz $version"
  Write-Host "  binary:  $BinDir\squiz.exe"
  Write-Host "  skill:   $SkillDir\SKILL.md"
  try { & (Join-Path $BinDir 'squiz.exe') version } catch { }
} finally {
  Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
