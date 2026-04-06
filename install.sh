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

# Parse arguments
while [ $# -gt 0 ]; do
  case "$1" in
    --binary)
      BINARY="$2"
      shift 2
      ;;
    --version)
      VERSION="$2"
      shift 2
      ;;
    --install-dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 1
      ;;
  esac
done

detect_os() {
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$os" in
    linux)  echo "linux" ;;
    darwin) echo "darwin" ;;
    mingw*|msys*|cygwin*) echo "windows" ;;
    *)
      echo "Unsupported OS: $os" >&2
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

get_latest_version() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
  else
    echo "Error: curl or wget is required" >&2
    exit 1
  fi
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

install_binary() {
  bin_name="$1"
  os="$2"
  arch="$3"
  version="$4"
  version_num="${version#v}"

  echo "Installing ${bin_name} ${version} for ${os}/${arch}..."

  ext="tar.gz"
  if [ "$os" = "windows" ]; then
    ext="zip"
  fi

  archive_name="bb-cli_${version_num}_${os}_${arch}.${ext}"
  download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT

  echo "Downloading ${download_url}..."
  download "$download_url" "${tmpdir}/${archive_name}"

  echo "Extracting..."
  if [ "$ext" = "zip" ]; then
    if command -v unzip >/dev/null 2>&1; then
      unzip -q "${tmpdir}/${archive_name}" -d "$tmpdir"
    else
      echo "Error: unzip is required for Windows archives" >&2
      exit 1
    fi
  else
    tar -xzf "${tmpdir}/${archive_name}" -C "$tmpdir"
  fi

  if [ ! -f "${tmpdir}/${bin_name}" ]; then
    echo "Error: ${bin_name} binary not found in archive" >&2
    exit 1
  fi

  echo "Installing to ${INSTALL_DIR}/${bin_name}..."
  if [ -w "$INSTALL_DIR" ]; then
    mv "${tmpdir}/${bin_name}" "${INSTALL_DIR}/${bin_name}"
    chmod +x "${INSTALL_DIR}/${bin_name}"
  else
    echo "Elevated permissions required to install to ${INSTALL_DIR}. Use --install-dir to choose a writable directory."
    sudo mv "${tmpdir}/${bin_name}" "${INSTALL_DIR}/${bin_name}"
    sudo chmod +x "${INSTALL_DIR}/${bin_name}"
  fi

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
