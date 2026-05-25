#!/usr/bin/env pwsh
# Update squiz + squiz-plan on Windows. Mirrors update.sh.
#
# Usage:
#   irm https://raw.githubusercontent.com/ptetau/squiz/main/update.ps1 | iex
#
# With args, fetch the script then dot-source it:
#   $u = irm https://raw.githubusercontent.com/ptetau/squiz/main/update.ps1
#   Invoke-Expression "$u; Update-Squiz -Yes"
#   Invoke-Expression "$u; Update-Squiz -Version 0.5.0 -Yes"
#   Invoke-Expression "$u; Update-Squiz -DryRun"
#
# Env overrides:
#   $env:VERSION   pin a version
#   $env:YES = 1   skip prompt
#   $env:DRY_RUN=1 show plan only

function Update-Squiz {
  [CmdletBinding()]
  param(
    [string]$Version = $env:VERSION,
    [switch]$Yes,
    [switch]$DryRun
  )

  $ErrorActionPreference = 'Stop'
  $Owner = 'ptetau'
  $Repo  = 'squiz'
  if ($env:YES -eq '1')     { $Yes    = $true }
  if ($env:DRY_RUN -eq '1') { $DryRun = $true }

  # 1. Detect installed binaries on PATH.
  $squizCmd = Get-Command squiz.exe -ErrorAction SilentlyContinue
  $sqpCmd   = Get-Command squiz-plan.exe -ErrorAction SilentlyContinue
  if (-not $squizCmd -and -not $sqpCmd) {
    throw "no squiz binaries on PATH — run install.ps1 first:`n  irm https://raw.githubusercontent.com/$Owner/$Repo/main/install.ps1 | iex"
  }
  $squizBin = if ($squizCmd) { $squizCmd.Source } else { $null }
  $sqpBin   = if ($sqpCmd)   { $sqpCmd.Source }   else { $null }

  # Current versions ("squiz 0.6.0" → "0.6.0").
  function Get-CurrentVersion($exe) {
    if (-not $exe) { return $null }
    try {
      $out = & $exe version 2>$null
      if ($out -match '\s(\S+)\s*$') { return $matches[1] }
    } catch { }
    return '?'
  }
  $squizVer = Get-CurrentVersion $squizBin
  $sqpVer   = Get-CurrentVersion $sqpBin

  # 2. Resolve target version.
  if (-not $Version) {
    $latest  = Invoke-RestMethod "https://api.github.com/repos/$Owner/$Repo/releases/latest"
    $Version = $latest.tag_name -replace '^v',''
    if (-not $Version) { throw "could not resolve latest version (use -Version to pin)" }
  }

  # 3. Detect skill directories.
  $skillRoots = @()
  $globalRoot = Join-Path $env:USERPROFILE '.claude\skills'
  if (Test-Path $globalRoot)  { $skillRoots += $globalRoot }
  $localRoot = Join-Path (Get-Location) '.claude\skills'
  if ((Test-Path $localRoot) -and ($localRoot -ne $globalRoot)) { $skillRoots += $localRoot }

  $existingSkills = @()
  foreach ($root in $skillRoots) {
    foreach ($name in @('squiz','squiz-plan','squiz-update')) {
      $path = Join-Path $root "$name\SKILL.md"
      if (Test-Path $path) { $existingSkills += $path }
    }
  }

  # 4. Short-circuit if nothing to do.
  $needBinUpdate = $false
  if ($squizBin -and $squizVer -ne $Version) { $needBinUpdate = $true }
  if ($sqpBin   -and $sqpVer   -ne $Version) { $needBinUpdate = $true }
  if (-not $needBinUpdate -and $existingSkills.Count -eq 0) {
    Write-Host "already at v$Version — nothing to do"
    return
  }

  # 5. Show plan.
  Write-Host "update plan:"
  if ($squizBin) { Write-Host ("  squiz       {0} -> {1}   ({2})" -f $squizVer, $Version, $squizBin) }
  if ($sqpBin)   { Write-Host ("  squiz-plan  {0} -> {1}   ({2})" -f $sqpVer,   $Version, $sqpBin) }
  foreach ($s in $existingSkills) { Write-Host ("  skill       (refresh)         ({0})" -f $s) }

  if ($DryRun) { Write-Host "[dry-run] no changes made"; return }

  # 6. Confirm.
  if (-not $Yes) {
    $answer = Read-Host "proceed? [y/N]"
    if ($answer -notmatch '^(y|Y|yes|YES|Yes)$') { Write-Host "cancelled"; return }
  }

  # 7. Detect arch.
  $arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { 'x86_64' }
    'ARM64' { 'arm64' }
    default { throw "unsupported arch: $($env:PROCESSOR_ARCHITECTURE)" }
  }
  $archive = "squiz_${Version}_Windows_${arch}.zip"
  $baseUrl = "https://github.com/$Owner/$Repo/releases/download/v$Version"

  $tmp = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "squiz-update-$(Get-Random)") -Force
  try {
    # 8. Download + verify.
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

    # 9. Replace each binary at its current path. Windows can't overwrite
    # a running .exe; rename old to .old then write new. The .old file is
    # locked until the running process exits; OS cleans it up next reboot
    # if anything was left.
    function Replace-Binary($currentPath, $name) {
      if (-not $currentPath) { return }
      $newPath = Join-Path $tmp ($name + '.exe')
      if (-not (Test-Path $newPath)) {
        Write-Host "WARN archive missing $name.exe; skipping"
        return
      }
      $oldPath = "$currentPath.old"
      if (Test-Path $oldPath) {
        try { Remove-Item $oldPath -Force -ErrorAction SilentlyContinue } catch { }
      }
      try {
        # If the current binary is in use, rename first. If not, just overwrite.
        if (Test-Path $currentPath) {
          Rename-Item -Path $currentPath -NewName ([System.IO.Path]::GetFileName($oldPath)) -Force
        }
      } catch {
        throw "could not rename in-use binary $currentPath ($_)"
      }
      Copy-Item -Path $newPath -Destination $currentPath -Force
      Write-Host "  binary: $currentPath"
    }
    Replace-Binary $squizBin 'squiz'
    Replace-Binary $sqpBin   'squiz-plan'

    # 10. Replace SKILL.mds only at locations where they already exist.
    foreach ($path in $existingSkills) {
      $name = Split-Path -Path (Split-Path -Path $path -Parent) -Leaf
      $src = Join-Path $tmp "skills\$name\SKILL.md"
      if (Test-Path $src) {
        Copy-Item -Path $src -Destination $path -Force
        Write-Host "  skill:  $path"
      } else {
        Write-Host "  WARN archive missing skills/$name/SKILL.md; left $path untouched"
      }
    }

    Write-Host ""
    Write-Host "OK updated to v$Version"
    if ($squizBin) { try { & $squizBin version } catch { } }
    if ($sqpBin)   { try { & $sqpBin   version } catch { } }
  } finally {
    Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
  }
}

# When piped through `irm | iex` with no params, run with defaults
# (prompt before applying). Comment this out if you want to always
# dot-source first.
if ($MyInvocation.InvocationName -eq '&' -or $MyInvocation.InvocationName -eq '') {
  Update-Squiz
}
