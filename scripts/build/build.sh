#!/usr/bin/env bash
# Usage: ./scripts/build/build.sh [--release] [--output-dir DIR]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
cd "${REPO_ROOT}"

OUTPUT_DIR="build"
RELEASE=false

while [[ $# -gt 0 ]]; do
  case "$1" in
  -r|--release)
    RELEASE=true
    shift
    ;;
  -o|--output-dir)
    OUTPUT_DIR="$2"
    shift 2
    ;;
  -h|--help)
    echo "Usage: $0 [options]"
    echo "  -r, --release     Build stripped release binary (-s -w)"
    echo "  -o, --output-dir  Directory for output binary (default: build)"
    echo "  -h, --help        Show this help message"
    exit 0
    ;;
  *)
    echo "Error: Unknown option: $1" >&2
    exit 1
    ;;
  esac
done

if ! command -v go >/dev/null 2>&1; then
  echo "Error: 'go' binary not found in PATH." >&2
  exit 1
fi

BINARY_NAME="azc"
if [[ "${OSTYPE:-}" == "msys" || "${OSTYPE:-}" == "cygwin" || "${OSTYPE:-}" == "win32" ]]; then
  BINARY_NAME="azc.exe"
fi

FULL_OUTPUT_DIR="${REPO_ROOT}/${OUTPUT_DIR}"
mkdir -p "${FULL_OUTPUT_DIR}"

OUTPUT_PATH="${FULL_OUTPUT_DIR}/${BINARY_NAME}"
SOURCE_PATH="./cmd/azc"

if [[ -t 1 ]]; then
  CYAN='\033[0;36m'
  YELLOW='\033[0;33m'
  GREEN='\033[0;32m'
  NC='\033[0m'
else
  CYAN=''
  YELLOW=''
  GREEN=''
  NC=''
fi

GO_FLAGS=("-trimpath")
if [[ "${RELEASE}" == true ]]; then
  GO_FLAGS+=("-ldflags" "-s -w")
fi

echo -e "${CYAN}Building Azin compiler (azc)...${NC}"
if [[ "${RELEASE}" == true ]]; then
  echo -e "${YELLOW}Mode: Release (symbols stripped)${NC}"
fi

START_TIME=$SECONDS

go build "${GO_FLAGS[@]}" -o "${OUTPUT_PATH}" "${SOURCE_PATH}"

ELAPSED=$(( SECONDS - START_TIME ))

echo -e "${GREEN}Successfully built ${OUTPUT_PATH} in ${ELAPSED}s${NC}"