#!/usr/bin/env bash
set -euo pipefail

# Ensure we're in a Git repo
if ! git rev-parse --git-dir > /dev/null 2>&1; then
  echo "❌ Not a git repository."
  exit 1
fi

# Ensure no uncommitted changes
if ! git diff-index --quiet HEAD --; then
  echo "❌ You have uncommitted changes. Commit or stash before bumping."
  exit 1
fi

# Ensure we're on main (or adjust branch name if needed)
branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$branch" != "main" ]; then
  echo "❌ You must be on 'main' branch (currently on '$branch')."
  exit 1
fi

# Get latest tag (default v0.0.0 if none exist)
latest=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Extract major, minor, patch
major=$(echo "$latest" | cut -d. -f1 | sed 's/v//')
minor=$(echo "$latest" | cut -d. -f2)
patch=$(echo "$latest" | cut -d. -f3)

# Determine bump type
case "${1:-patch}" in
  major)
    major=$((major+1))
    minor=0
    patch=0
    ;;
  minor)
    minor=$((minor+1))
    patch=0
    ;;
  patch)
    patch=$((patch+1))
    ;;
  *)
    echo "❌ Unknown bump type: $1 (use major|minor|patch)"
    exit 1
    ;;
esac

new="v$major.$minor.$patch"

# Build
mkdir -p dist

echo "➡️  Building binaries..."

GOOS=linux GOARCH=amd64 go build -o dist/myapp-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o dist/myapp-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o dist/myapp-windows-amd64.exe

echo "✅ Binaries built in dist/"

# Create GitHub release with gh CLI (optional)
if command -v gh >/dev/null 2>&1; then
    echo "➡️  Creating GitHub release..."
    gh release create "$new" dist/* --title "$new" --notes "Release $new"
    echo "✅ Release created on GitHub"
else
    echo "GH not installed"
fi

# Confirm & tag
echo "➡️  Bumping version: $latest → $new"
git tag "$new"
git push origin "$new"

echo "✅ Done! New version: $new"



