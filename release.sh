#!/usr/bin/env bash
set -euo pipefail

# -----------------------
# CONFIG
# -----------------------
DIST_DIR="dist"
CHANGELOG="CHANGELOG.md"
DEFAULT_BRANCH="main"
APP_NAME="ChainGuardian"  # Change this to your binary name
DOCKER_USER="${DOCKER_USER:-mydockeruser}" # Docker Hub user/org
DOCKER_REPO="${DOCKER_REPO:-myapp}"        # Docker Hub repo

# -----------------------
# HELPERS
# -----------------------
function abort { echo "❌ $1"; exit 1; }

# Ensure we're in a Git repo
git rev-parse --git-dir > /dev/null 2>&1 || abort "Not a git repository."

# Ensure no uncommitted changes
git diff-index --quiet HEAD -- || abort "You have uncommitted changes."

# Ensure we're on main
branch=$(git rev-parse --abbrev-ref HEAD)
[[ "$branch" == "$DEFAULT_BRANCH" ]] || abort "Must be on '$DEFAULT_BRANCH' branch."

# Ensure CHANGELOG.md exists
[[ -f "$CHANGELOG" ]] || abort "$CHANGELOG not found."

# -----------------------
# GET LATEST VERSION
# -----------------------
latest=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
major=$(echo "$latest" | cut -d. -f1 | sed 's/v//')
minor=$(echo "$latest" | cut -d. -f2)
patch=$(echo "$latest" | cut -d. -f3)

# -----------------------
# BUMP VERSION
# -----------------------
bump_type="${1:-patch}"
case "$bump_type" in
    major) major=$((major+1)); minor=0; patch=0 ;;
    minor) minor=$((minor+1)); patch=0 ;;
    patch) patch=$((patch+1)) ;;
    *) abort "Unknown bump type: $bump_type (use major|minor|patch)" ;;
esac

new="v$major.$minor.$patch"

# -----------------------
# CHECK CHANGELOG HEADER
# -----------------------
if ! grep -q "^## $new" "$CHANGELOG"; then
    abort "$CHANGELOG does not contain a header '## [$new]'. Please add it before releasing."
fi

# -----------------------
# RUN TESTS
# -----------------------
echo "➡️  Running tests..."
go test ./... || abort "Tests failed."

# -----------------------
# TAG & PUSH
# -----------------------
echo "➡️  Tagging version: $latest → $new"
git tag "$new"
git push origin "$new"

# -----------------------
# BUILD BINARIES
# -----------------------
echo "➡️  Building binaries..."
mkdir -p "$DIST_DIR"

GOOS=linux   GOARCH=amd64 go build -o "$DIST_DIR/$APP_NAME-linux-amd64"
GOOS=darwin  GOARCH=amd64 go build -o "$DIST_DIR/$APP_NAME-darwin-amd64"
GOOS=windows GOARCH=amd64 go build -o "$DIST_DIR/$APP_NAME-windows-amd64.exe"

echo "✅ Binaries built in $DIST_DIR/"

# -----------------------
# CREATE GITHUB RELEASE
# -----------------------
if command -v gh >/dev/null 2>&1; then
    echo "➡️  Creating GitHub release..."
    gh release create "$new" "$DIST_DIR"/* --title "$new" --notes "Release $new"
    echo "✅ GitHub release created."
fi

# -----------------------
# DOCKER BUILD & PUSH
# -----------------------
echo "➡️  Building Docker image..."
docker build -t "$DOCKER_USER/$DOCKER_REPO:latest" -t "$DOCKER_USER/$DOCKER_REPO:$new" .

echo "➡️  Pushing Docker image..."
docker push "$DOCKER_USER/$DOCKER_REPO:latest"
docker push "$DOCKER_USER/$DOCKER_REPO:$new"

echo "✅ Docker images pushed: $DOCKER_USER/$DOCKER_REPO:$new and :latest"

echo "🎉 Release $new completed successfully!"
