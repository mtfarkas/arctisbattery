#!/bin/bash
# Script to create and push a new version tag

set -e

# Get the latest tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "Latest tag: $LATEST_TAG"

# Remove 'v' prefix and split version
VERSION=${LATEST_TAG#v}
IFS='.' read -r -a VERSION_PARTS <<< "$VERSION"

MAJOR=${VERSION_PARTS[0]:-0}
MINOR=${VERSION_PARTS[1]:-0}
PATCH=${VERSION_PARTS[2]:-0}

# Display menu
echo ""
echo "Select version increment:"
echo "1) Major (${MAJOR}.${MINOR}.${PATCH} -> $((MAJOR+1)).0.0)"
echo "2) Minor (${MAJOR}.${MINOR}.${PATCH} -> ${MAJOR}.$((MINOR+1)).0)"
echo "3) Patch (${MAJOR}.${MINOR}.${PATCH} -> ${MAJOR}.${MINOR}.$((PATCH+1)))"
echo "4) Custom version"
echo "5) Cancel"
echo ""
read -p "Enter choice [1-5]: " choice

case $choice in
  1)
    NEW_VERSION="$((MAJOR+1)).0.0"
    ;;
  2)
    NEW_VERSION="${MAJOR}.$((MINOR+1)).0"
    ;;
  3)
    NEW_VERSION="${MAJOR}.${MINOR}.$((PATCH+1))"
    ;;
  4)
    read -p "Enter custom version (e.g., 1.2.3): " CUSTOM_VERSION
    NEW_VERSION="$CUSTOM_VERSION"
    ;;
  5)
    echo "Cancelled."
    exit 0
    ;;
  *)
    echo "Invalid choice."
    exit 1
    ;;
esac

NEW_TAG="v${NEW_VERSION}"

# Ask for confirmation
echo ""
read -p "Create and push tag ${NEW_TAG}? (y/n): " confirm
if [ "$confirm" != "y" ]; then
  echo "Cancelled."
  exit 0
fi

# Create annotated tag
read -p "Enter tag message (optional, press Enter to skip): " TAG_MESSAGE
if [ -z "$TAG_MESSAGE" ]; then
  TAG_MESSAGE="Release ${NEW_TAG}"
fi

git tag -a "$NEW_TAG" -m "$TAG_MESSAGE"
echo "Created tag: $NEW_TAG"

# Push tag
git push origin "$NEW_TAG"
echo "Pushed tag: $NEW_TAG"
echo ""
echo "GitHub Actions will now build and create a release."
echo "Check: https://github.com/mtfarkas/arctisbattery/actions"
