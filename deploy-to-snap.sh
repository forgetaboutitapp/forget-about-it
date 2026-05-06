#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: ./deploy-to-snap.sh [channel]

Build the Linux server and provision executables, package them into
`forget-about-it` snaps, and upload those snaps to Snapcraft.

Arguments:
  channel                 Snap release channel. Default: edge

Environment variables:
  SNAP_RELEASE_CHANNEL    Release channel override. Default: edge
  SNAP_VERSION            Optional snap version override
  SNAP_ARCHES             Space-separated Debian architectures. Default: amd64 arm64
  SNAPCRAFT_STORE_CREDENTIALS
                          Exported Snapcraft credentials for non-interactive auth

Requirements:
  - Linux host with buf, dart, flutter, sqlc, go, and snapcraft installed
  - Snapcraft authentication already configured, or SNAPCRAFT_STORE_CREDENTIALS set
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "$1 is required" >&2
    exit 1
  fi
}

if [[ "$(uname -s)" != "Linux" ]]; then
  echo "deploy-to-snap.sh must run on Linux" >&2
  exit 1
fi

require_cmd buf
require_cmd dart
require_cmd flutter
require_cmd sqlc
require_cmd go
require_cmd snapcraft

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
repo_root=$(cd "$script_dir/.." && pwd)
protobuf_dir="$repo_root/protobufs"
frontend_dir="$script_dir/frontend"
server_dir="$script_dir/server"
snap_local_dir="$script_dir/snap/local"
snapcraft_yaml="$script_dir/snap/snapcraft.yaml"
channel="${SNAP_RELEASE_CHANNEL:-${1:-edge}}"
snap_arches="${SNAP_ARCHES:-amd64 arm64}"

pubspec_version=$(grep '^version:' "$frontend_dir/pubspec.yaml" | head -n1 | awk '{print $2}')
snap_version="${SNAP_VERSION:-${pubspec_version/+/-}}"
snap_output_dir="$script_dir/.snap-output"

mkdir -p "$snap_local_dir/bin" "$snap_output_dir"
rm -f "$snap_output_dir"/forget-about-it_*.snap

pushd "$protobuf_dir" >/dev/null
buf generate --template buf.gen.go.yaml
buf generate --template buf.gen.dart.yaml
popd >/dev/null

pushd "$frontend_dir" >/dev/null
dart run build_runner build -d
flutter build web --release --wasm
popd >/dev/null

rm -rf "$server_dir/web"
mkdir -p "$server_dir/web"
cp -R "$frontend_dir/build/web/." "$server_dir/web/"

pushd "$server_dir" >/dev/null
sqlc generate
popd >/dev/null

original_snapcraft_yaml=$(mktemp)
cp "$snapcraft_yaml" "$original_snapcraft_yaml"
trap 'cp "$original_snapcraft_yaml" "$snapcraft_yaml"; rm -f "$original_snapcraft_yaml"' EXIT
sed -i "s/^version: .*/version: \"$snap_version\"/" "$snapcraft_yaml"

pushd "$script_dir" >/dev/null
if [[ -n "${SNAPCRAFT_STORE_CREDENTIALS:-}" ]]; then
  export SNAPCRAFT_STORE_CREDENTIALS
fi

built_snaps=()
for snap_arch in $snap_arches; do
  case "$snap_arch" in
    amd64) go_arch=amd64 ;;
    arm64) go_arch=arm64 ;;
    *)
      echo "Unsupported snap architecture: $snap_arch" >&2
      exit 1
      ;;
  esac

  rm -f "$snap_local_dir/bin/server" "$snap_local_dir/bin/provision"
  pushd "$server_dir" >/dev/null
  GOOS=linux GOARCH="$go_arch" go build -o "$snap_local_dir/bin/server" ./cmd/server
  GOOS=linux GOARCH="$go_arch" go build -o "$snap_local_dir/bin/provision" ./cmd/provision
  popd >/dev/null

  snapcraft clean
  snap_file="$snap_output_dir/forget-about-it_${snap_version}_${snap_arch}.snap"
  snapcraft --destructive-mode --build-for="$snap_arch" --output "$snap_file"
  if [[ ! -f "$snap_file" ]]; then
    echo "snapcraft did not produce $snap_file" >&2
    exit 1
  fi
  built_snaps+=("$snap_file")
done

for snap_file in "${built_snaps[@]}"; do
  snapcraft upload "$snap_file" --release "$channel"
done
popd >/dev/null

echo "Published ${built_snaps[*]} to channel $channel"
