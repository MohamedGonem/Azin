#!/usr/bin/env sh
# Usage: ./scripts/docs/doc.sh [-o DIR] [--no-internal]

set -eu

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/../.." && pwd)
cd "$REPO_ROOT"

OUTPUT_DIR="docs/api"
INTERNAL_FLAG="-internal"

while [ $# -gt 0 ]; do
  case "$1" in
  -o | --output-dir)
    OUTPUT_DIR="$2"
    shift 2
    ;;
  --no-internal)
    INTERNAL_FLAG=""
    shift
    ;;
  -h | --help)
    echo "Usage: $0 [options]"
    echo "  -o, --output-dir DIR  Target directory (default: docs/api)"
    echo "  --no-internal         Exclude internal packages from documentation"
    echo "  -h, --help            Show this help message"
    exit 0
    ;;
  *)
    echo "Error: Unknown option: $1" >&2
    exit 1
    ;;
  esac
done

if [ -t 1 ]; then
  CYAN='\033[0;36m'
  GREEN='\033[0;32m'
  YELLOW='\033[0;33m'
  NC='\033[0m'
else
  CYAN='' GREEN='' YELLOW='' NC=''
fi

FULL_OUTPUT_DIR="$REPO_ROOT/$OUTPUT_DIR"
mkdir -p "$FULL_OUTPUT_DIR"

printf "%b\n" "${CYAN}Generating API documentation...${NC}"

if command -v doc2go >/dev/null 2>&1; then
  if [ -n "$INTERNAL_FLAG" ]; then
    doc2go "$INTERNAL_FLAG" -out "$FULL_OUTPUT_DIR" ./...
  else
    doc2go -out "$FULL_OUTPUT_DIR" ./...
  fi
elif command -v go >/dev/null 2>&1; then
  printf "%b\n" "${YELLOW}'doc2go' binary not found in PATH. Falling back to 'go run'...${NC}"
  if [ -n "$INTERNAL_FLAG" ]; then
    go run github.com/willabides/doc2go@latest "$INTERNAL_FLAG" -out "$FULL_OUTPUT_DIR" ./...
  else
    go run github.com/willabides/doc2go@latest -out "$FULL_OUTPUT_DIR" ./...
  fi
else
  echo "Error: Neither 'doc2go' nor 'go' was found in PATH." >&2
  exit 1
fi

printf "%b\n" "${GREEN}Successfully generated documentation at: ${OUTPUT_DIR}${NC}"
