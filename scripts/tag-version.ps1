# PowerShell script to create and push a new version tag

# Get the latest tag
$LatestTag = git describe --tags --abbrev=0 2>$null
if (-not $LatestTag) {
    $LatestTag = "v0.0.0"
}
Write-Host "Latest tag: $LatestTag"

# Remove 'v' prefix and split version
$Version = $LatestTag -replace '^v', ''
$VersionParts = $Version -split '\.'
$Major = if ($VersionParts[0]) { [int]$VersionParts[0] } else { 0 }
$Minor = if ($VersionParts[1]) { [int]$VersionParts[1] } else { 0 }
$Patch = if ($VersionParts[2]) { [int]$VersionParts[2] } else { 0 }

# Display menu
Write-Host ""
Write-Host "Select version increment:"
Write-Host "1) Major ($Major.$Minor.$Patch -> $($Major+1).0.0)"
Write-Host "2) Minor ($Major.$Minor.$Patch -> $Major.$($Minor+1).0)"
Write-Host "3) Patch ($Major.$Minor.$Patch -> $Major.$Minor.$($Patch+1))"
Write-Host "4) Custom version"
Write-Host "5) Cancel"
Write-Host ""

$choice = Read-Host "Enter choice [1-5]"

switch ($choice) {
    "1" {
        $NewVersion = "$($Major+1).0.0"
    }
    "2" {
        $NewVersion = "$Major.$($Minor+1).0"
    }
    "3" {
        $NewVersion = "$Major.$Minor.$($Patch+1)"
    }
    "4" {
        $CustomVersion = Read-Host "Enter custom version (e.g., 1.2.3)"
        $NewVersion = $CustomVersion
    }
    "5" {
        Write-Host "Cancelled."
        exit 0
    }
    default {
        Write-Host "Invalid choice."
        exit 1
    }
}

$NewTag = "v$NewVersion"

# Ask for confirmation
Write-Host ""
$confirm = Read-Host "Create and push tag $NewTag? (y/n)"
if ($confirm -ne "y") {
    Write-Host "Cancelled."
    exit 0
}

# Create annotated tag
$TagMessage = Read-Host "Enter tag message (optional, press Enter to skip)"
if ([string]::IsNullOrWhiteSpace($TagMessage)) {
    $TagMessage = "Release $NewTag"
}

git tag -a $NewTag -m $TagMessage
Write-Host "Created tag: $NewTag"

# Push tag
git push origin $NewTag
Write-Host "Pushed tag: $NewTag"
Write-Host ""
Write-Host "GitHub Actions will now build and create a release."
Write-Host "Check: https://github.com/mtfarkas/arctisbattery/actions"
