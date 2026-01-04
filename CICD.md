# Releases

## Creating a release

Just run the tag script and it'll walk you through it:

```bash
# Windows
.\scripts\tag-version.ps1

# Linux/Mac
bash scripts/tag-version.sh

# Or if you have Make
make release
```

Pick whether you want to bump major, minor, or patch version. Script creates the tag and pushes it, then GitHub Actions does the rest.

## What gets built

Push a tag like `v1.0.0` and you'll get Windows builds for amd64 and arm64. Each binary gets zipped with the README and attached to the release.

## Manual tagging

If you'd rather do it yourself:

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

## Version numbers

Using semver here:
- Major (1.0.0 → 2.0.0) - breaking stuff
- Minor (1.0.0 → 1.1.0) - new features
- Patch (1.0.0 → 1.0.1) - bug fixes

## If something breaks

Check https://github.com/mtfarkas/arctisbattery/actions to see what went wrong.

Run `make check` before tagging to catch issues locally.

If you need to redo a tag:
```bash
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
# then create it again
```

## Building ARM64 locally

If you want to test the ARM64 build on Windows:

```powershell
$Env:GOARCH = "arm64"; go build
```