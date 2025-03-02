#!/usr/bin/env bash
set -euo pipefail

# Usage:
#   ./scripts/calc_new_version.sh <major|minor|patch> [alpha|beta]
# Example:
#   ./scripts/calc_new_version.sh patch alpha
#
# The script reads the latest tag from the repository (defaulting to v0.0.0 if none exists),
# calculates the new version based on semver rules with our prerelease transition logic,
# and prints only the new version string.

if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <major|minor|patch> [alpha|beta]" >&2
  exit 1
fi

RELEASE_TYPE=$1
PRERELEASE=${2:-""}

if [[ "$RELEASE_TYPE" != "major" && "$RELEASE_TYPE" != "minor" && "$RELEASE_TYPE" != "patch" ]]; then
  echo "Error: release type must be one of: major, minor, patch" >&2
  exit 1
fi

if [ -n "$PRERELEASE" ]; then
  if [[ "$PRERELEASE" != "alpha" && "$PRERELEASE" != "beta" ]]; then
    echo "Error: prerelease, if provided, must be alpha or beta" >&2
    exit 1
  fi
fi

LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

VERSION_NO_V=${LATEST_TAG#v}
if [[ $VERSION_NO_V =~ ^([0-9]+)\.([0-9]+)\.([0-9]+)(-([a-z]+)\.([0-9]+))?$ ]]; then
  CURRENT_MAJOR=${BASH_REMATCH[1]}
  CURRENT_MINOR=${BASH_REMATCH[2]}
  CURRENT_PATCH=${BASH_REMATCH[3]}
  CURRENT_PRERELEASE_TYPE=${BASH_REMATCH[5]:-}
  CURRENT_PRERELEASE_NUM=${BASH_REMATCH[6]:-0}
else
  echo "Error: Latest tag '$LATEST_TAG' is not in a valid semantic version format." >&2
  exit 1
fi

NEW_MAJOR=$CURRENT_MAJOR
NEW_MINOR=$CURRENT_MINOR
NEW_PATCH=$CURRENT_PATCH
NEW_PRERELEASE=""

case "$RELEASE_TYPE" in
  major)
    NEW_MAJOR=$((CURRENT_MAJOR + 1))
    NEW_MINOR=0
    NEW_PATCH=0
    ;;
  minor)
    NEW_MINOR=$((CURRENT_MINOR + 1))
    NEW_PATCH=0
    ;;
  patch)
    if [ -n "$PRERELEASE" ]; then
      # For patch bump with prerelease:
      if [ -z "$CURRENT_PRERELEASE_TYPE" ]; then
        # If no prerelease exists yet, bump the patch.
        NEW_PATCH=$((CURRENT_PATCH + 1))
      else
        if [ "$CURRENT_PRERELEASE_TYPE" == "$PRERELEASE" ]; then
          # Same prerelease type: do not change numeric version.
          NEW_PATCH=$CURRENT_PATCH
        elif [ "$CURRENT_PRERELEASE_TYPE" == "alpha" ] && [ "$PRERELEASE" == "beta" ]; then
          # Allow switching from alpha to beta without bumping numeric part.
          NEW_PATCH=$CURRENT_PATCH
        elif [ "$CURRENT_PRERELEASE_TYPE" == "beta" ] && [ "$PRERELEASE" == "alpha" ]; then
          echo "Error: Switching from beta to alpha in a patch bump is not allowed without a numeric bump." >&2
          exit 1
        else
          # Fallback: bump numeric part.
          NEW_PATCH=$((CURRENT_PATCH + 1))
        fi
      fi
    else
      NEW_PATCH=$((CURRENT_PATCH + 1))
    fi
    ;;
esac

# Set prerelease part if a prerelease argument is provided.
if [ -n "$PRERELEASE" ]; then
  if [ -z "$CURRENT_PRERELEASE_TYPE" ]; then
    NEW_PRERELEASE="${PRERELEASE}.1"
  else
    if [ "$CURRENT_PRERELEASE_TYPE" == "$PRERELEASE" ]; then
      NEW_PRERELEASE="${PRERELEASE}.$((CURRENT_PRERELEASE_NUM + 1))"
    elif [ "$CURRENT_PRERELEASE_TYPE" == "alpha" ] && [ "$PRERELEASE" == "beta" ]; then
      NEW_PRERELEASE="${PRERELEASE}.1"
    fi
  fi
fi

if [ -n "$NEW_PRERELEASE" ]; then
  NEW_TAG="v${NEW_MAJOR}.${NEW_MINOR}.${NEW_PATCH}-${NEW_PRERELEASE}"
else
  NEW_TAG="v${NEW_MAJOR}.${NEW_MINOR}.${NEW_PATCH}"
fi

echo "$NEW_TAG"
