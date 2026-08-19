#!/bin/sh
set -eu

release_tag=${1:?release tag is required}
if ! printf '%s\n' "$release_tag" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$'; then
  echo "Invalid release tag: $release_tag" >&2
  exit 2
fi

if [ -e dist ]; then
  echo "Refusing to replace existing dist directory." >&2
  exit 1
fi
mkdir -p dist

for target in \
  darwin/amd64 darwin/arm64 \
  linux/386 linux/amd64 linux/arm linux/arm64 \
  windows/386 windows/amd64 windows/arm64
do
  target_os=${target%/*}
  target_arch=${target#*/}
  extension=""
  if [ "$target_os" = "windows" ]; then
    extension=".exe"
  fi
  output="dist/${target_os}-${target_arch}${extension}"
  echo "Building $output"
  GOOS=$target_os GOARCH=$target_arch CGO_ENABLED=0 \
    go build -trimpath -ldflags="-s -w -X main.version=$release_tag" \
    -o "$output" ./cmd/gh-pr-quality-gate
done
