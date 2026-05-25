#!/usr/bin/env pwsh
# Install squiz on Windows: binary on PATH + SKILL.md files into ~\.claude\skills\<skill>\.
#
# Usage:
#   irm https://raw.githubusercontent.com/ptetau/squiz/main/install.ps1 | iex
#
# Env overrides:
#   $env:VERSION            pin a specific version (default: latest GitHub release)
#   $env:SQUIZ_BIN_DIR      where to install squiz.exe (default: %LOCALAPPDATA%\Programs\squiz)
#   $env:SQUIZ_SKILLS_ROOT  root for skill dirs (default: %USERPROFILE%\.claude\skills)

$ErrorActionPreference = 'Stop'

$Owner       = 'ptetau'
$Repo        = 'squiz'
$BinDir      = if ($env:SQUIZ_BIN_DIR)     { $env:SQUIZ_BIN_DIR }     else { "$env:LOCALAPPDATA\Programs\squiz" }
$SkillsRoot  = if ($env:SQUIZ_SKILLS_ROOT) { $env:SQUIZ_SKILLS_ROOT } else { "$env:USERPROFILE\.claude\skills" }

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

  New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
  # Install every .exe the archive ships at its top level (squiz.exe, squiz-plan.exe, …).
  foreach ($exe in 'squiz.exe','squiz-plan.exe') {
    $src = Join-Path $tmp $exe
    if (Test-Path $src) {
      Move-Item -Path $src -Destination (Join-Path $BinDir $exe) -Force
      Write-Host "  binary:  $BinDir\$exe"
    }
  }

  # Install every SKILL.md the archive ships under skills\<name>\.
  $skillsDir = Join-Path $tmp 'skills'
  if (Test-Path $skillsDir) {
    Get-ChildItem $skillsDir -Directory | ForEach-Object {
      $src = Join-Path $_.FullName 'SKILL.md'
      if (Test-Path $src) {
        $dst = Join-Path $SkillsRoot $_.Name
        New-Item -ItemType Directory -Path $dst -Force | Out-Null
        Move-Item -Path $src -Destination (Join-Path $dst 'SKILL.md') -Force
        Write-Host "  skill:   $dst\SKILL.md"
      }
    }
  }

  # Add BinDir to user PATH if missing.
  $userPath = [Environment]::GetEnvironmentVariable('Path','User')
  if ($userPath -notlike "*$BinDir*") {
    $newPath = if ($userPath) { "$userPath;$BinDir" } else { $BinDir }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host "-> added $BinDir to user PATH (open a new shell to pick it up)"
  }

  Write-Host ""
  Write-Host "OK installed squiz $version"
  try { & (Join-Path $BinDir 'squiz.exe') version } catch { }
  if (Test-Path (Join-Path $BinDir 'squiz-plan.exe')) {
    try { & (Join-Path $BinDir 'squiz-plan.exe') version } catch { }
  }
} finally {
  Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
