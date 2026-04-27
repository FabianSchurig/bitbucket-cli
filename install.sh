#!/bin/sh
# install.sh — Install bb-cli and/or bb-mcp from GitHub Releases.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh -s -- --binary bb-mcp
#   curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh -s -- --version v1.2.3
#
# Options:
#   --binary NAME    Binary to install: bb-cli (default), bb-mcp, or all
#   --version TAG    Version tag to install (default: latest)
#   --install-dir DIR  Installation directory (default: /usr/local/bin)

set -e

REPO="FabianSchurig/bitbucket-cli"
BINARY="bb-cli"
VERSION=""
INSTALL_DIR="/usr/local/bin"

usage() {
  cat >&2 <<EOF
Usage:
  sh install.sh [--binary NAME] [--version TAG] [--install-dir DIR]

Options:
  --binary NAME      Binary to install: bb-cli (default), bb-mcp, or all
  --version TAG      Version tag to install (default: latest)
  --install-dir DIR  Installation directory (default: /usr/local/bin)
EOF
}

require_value() {
  if [ $# -lt 2 ] || [ -z "$2" ]; then
    echo "Missing value for option: $1" >&2
    usage
    exit 1
  fi
}

# Parse arguments
while [ $# -gt 0 ]; do
  case "$1" in
    --binary)
      require_value "$1" "${2:-}"
      BINARY="$2"
      shift 2
      ;;
    --version)
      require_value "$1" "${2:-}"
      VERSION="$2"
      shift 2
      ;;
    --install-dir)
      require_value "$1" "${2:-}"
      INSTALL_DIR="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

detect_os() {
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$os" in
    linux)  echo "linux" ;;
    darwin) echo "darwin" ;;
    *)
      echo "Unsupported OS: $os" >&2
      echo "On Windows, download binaries directly from https://github.com/${REPO}/releases" >&2
      exit 1
      ;;
  esac
}

detect_arch() {
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *)
      echo "Unsupported architecture: $arch" >&2
      exit 1
      ;;
  esac
}

extract_tag_from_release_url() {
  printf '%s\n' "$1" | sed -n 's#^.*/tag/\([^/?[:space:]]*\).*$#\1#p'
}

get_latest_version_from_redirect() {
  latest_url="https://github.com/${REPO}/releases/latest"

  if command -v curl >/dev/null 2>&1; then
    final_url="$(curl -fsSIL -o /dev/null -w '%{url_effective}' "$latest_url" 2>/dev/null || true)"
    version="$(extract_tag_from_release_url "$final_url")"
    if [ -n "$version" ]; then
      printf '%s\n' "$version"
      return 0
    fi
  fi

  if command -v wget >/dev/null 2>&1; then
    redirect_url="$(wget -S --max-redirect=0 -O /dev/null "$latest_url" 2>&1 | sed -n 's/^[[:space:]]*Location: \(.*\)$/\1/p' | tail -n 1)"
    version="$(extract_tag_from_release_url "$redirect_url")"
    if [ -n "$version" ]; then
      printf '%s\n' "$version"
      return 0
    fi
  fi

  return 1
}

get_latest_version_from_api() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
  else
    echo "Error: curl or wget is required" >&2
    exit 1
  fi
}

get_latest_version() {
  version="$(get_latest_version_from_redirect || true)"
  if [ -n "$version" ]; then
    printf '%s\n' "$version"
    return 0
  fi

  get_latest_version_from_api
}

download() {
  url="$1"
  output="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$output" "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$output" "$url"
  else
    echo "Error: curl or wget is required" >&2
    exit 1
  fi
}

verify_checksum() {
  archive_path="$1"
  archive_name="$2"
  version="$3"

  checksums_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"
  checksums_file="$(dirname "$archive_path")/checksums.txt"

  echo "Verifying checksum..."
  download "$checksums_url" "$checksums_file"

  expected="$(grep "${archive_name}" "$checksums_file" | awk '{print $1}')"
  if [ -z "$expected" ]; then
    echo "Warning: checksum for ${archive_name} not found in checksums.txt, skipping verification" >&2
    return
  fi

  if command -v sha256sum >/dev/null 2>&1; then
    actual="$(sha256sum "$archive_path" | awk '{print $1}')"
  elif command -v shasum >/dev/null 2>&1; then
    actual="$(shasum -a 256 "$archive_path" | awk '{print $1}')"
  else
    echo "Warning: sha256sum or shasum not found, skipping checksum verification" >&2
    return
  fi

  if [ "$expected" != "$actual" ]; then
    echo "Error: checksum mismatch for ${archive_name}" >&2
    echo "  expected: ${expected}" >&2
    echo "  actual:   ${actual}" >&2
    exit 1
  fi
  echo "Checksum verified."
}

install_binary() {
  bin_name="$1"
  os="$2"
  arch="$3"
  version="$4"
  version_num="${version#v}"

  echo "Installing ${bin_name} ${version} for ${os}/${arch}..."

  archive_name="${bin_name}_${version_num}_${os}_${arch}.tar.gz"
  download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

  tmpdir="$(mktemp -d)"

  echo "Downloading ${download_url}..."
  download "$download_url" "${tmpdir}/${archive_name}"

  verify_checksum "${tmpdir}/${archive_name}" "$archive_name" "$version"

  echo "Extracting..."
  tar -xzf "${tmpdir}/${archive_name}" -C "$tmpdir"

  if [ ! -f "${tmpdir}/${bin_name}" ]; then
    echo "Error: ${bin_name} binary not found in archive" >&2
    rm -rf "$tmpdir"
    exit 1
  fi

  echo "Installing to ${INSTALL_DIR}/${bin_name}..."
  if mkdir -p "$INSTALL_DIR" 2>/dev/null && [ -w "$INSTALL_DIR" ]; then
    mv "${tmpdir}/${bin_name}" "${INSTALL_DIR}/${bin_name}"
    chmod +x "${INSTALL_DIR}/${bin_name}"
  else
    echo "Elevated permissions required to install to ${INSTALL_DIR}. Use --install-dir to choose a writable directory."
    sudo mkdir -p "$INSTALL_DIR"
    sudo mv "${tmpdir}/${bin_name}" "${INSTALL_DIR}/${bin_name}"
    sudo chmod +x "${INSTALL_DIR}/${bin_name}"
  fi

  rm -rf "$tmpdir"
  echo "${bin_name} ${version} installed successfully to ${INSTALL_DIR}/${bin_name}"
}

main() {
  os="$(detect_os)"
  arch="$(detect_arch)"

  if [ -z "$VERSION" ]; then
    echo "Fetching latest version..."
    VERSION="$(get_latest_version)"
    if [ -z "$VERSION" ]; then
      echo "Error: could not determine latest version" >&2
      exit 1
    fi
  fi

  echo "Version: ${VERSION}"

  case "$BINARY" in
    bb-cli)
      install_binary "bb-cli" "$os" "$arch" "$VERSION"
      ;;
    bb-mcp)
      install_binary "bb-mcp" "$os" "$arch" "$VERSION"
      ;;
    all)
      install_binary "bb-cli" "$os" "$arch" "$VERSION"
      install_binary "bb-mcp" "$os" "$arch" "$VERSION"
      ;;
    *)
      echo "Unknown binary: ${BINARY}. Use bb-cli, bb-mcp, or all." >&2
      exit 1
      ;;
  esac
}

main
