#!/bin/bash
set -e

VERSION="2024.12.2"
BASE_URL="https://github.com/cloudflare/cloudflared/releases/download/${VERSION}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BIN_DIR="${SCRIPT_DIR}/../pkg/cloudflare/bin"

mkdir -p "${BIN_DIR}"

download() {
    local os=$1
    local arch=$2
    local ext=$3
    local is_tgz=$4
    local filename="cloudflared-${os}-${arch}${ext}"
    local url="${BASE_URL}/${filename}"
    local output="${BIN_DIR}/cloudflared-${os}-${arch}"

    if [ "${os}" = "windows" ]; then
        output="${output}.exe"
    fi

    if [ -f "${output}" ]; then
        echo "✓ cloudflared-${os}-${arch}${ext} already exists"
        return 0
    fi

    echo "↓ Downloading ${filename}..."
    
    if [ "${is_tgz}" = "true" ]; then
        local tmpdir=$(mktemp -d)
        curl -L -o "${tmpdir}/${filename}" "${url}"
        tar -xzf "${tmpdir}/${filename}" -C "${tmpdir}"
        mv "${tmpdir}/cloudflared" "${output}"
        rm -rf "${tmpdir}"
    else
        curl -L -o "${output}" "${url}"
    fi
    
    chmod +x "${output}"
    echo "✓ cloudflared-${os}-${arch} downloaded"
}

download "darwin" "arm64" ".tgz" "true"
download "darwin" "amd64" ".tgz" "true"
download "linux" "amd64" "" "false"
download "linux" "arm64" "" "false"
download "windows" "amd64" ".exe" "false"

echo ""
echo "All cloudflared binaries downloaded to ${BIN_DIR}"
ls -lh "${BIN_DIR}"
