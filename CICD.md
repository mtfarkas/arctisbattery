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

Push a tag like `v1.0.0` and you'll get builds for:
- Windows (amd64 + arm64)
- Linux (amd64 + arm64)
- macOS (amd64 + arm64)

Each binary gets zipped with the README and attached to the release.

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

## Cross-compiling locally

CI does this automatically, but if you need to test it:

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build
CGO_ENABLED=1 GOOS=windows GOARCH=arm64 go build
```

Cross-compiling with CGO is annoying though. Let the CI handle it.
