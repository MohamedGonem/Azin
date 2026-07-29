#!/bin/sh
# Azin compiler installer script
# Usage: ./install.sh [--prefix DIR] [--branch BRANCH_OR_TAG]

set -eu

REPO="https://github.com/azin-lang/Azin.git"
BIN="azc"
PREFIX="${PREFIX:-/usr/local}"
BRANCH=""

if [ -t 1 ]; then
    CYAN='\033[0;36m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    RED='\033[0;31m'
    NC='\033[0m'
else
    CYAN='' GREEN='' YELLOW='' RED='' NC=''
fi

log_info() { printf "%b\n" "${CYAN}info:${NC} $1"; }
log_warn() { printf "%b\n" "${YELLOW}warning:${NC} $1" >&2; }
log_err() { printf "%b\n" "${RED}error:${NC} $1" >&2; }

# Parse command-line options
while [ $# -gt 0 ]; do
    case "$1" in
    -p | --prefix)
        PREFIX="$2"
        shift 2
        ;;
    -b | --branch | --tag)
        BRANCH="$2"
        shift 2
        ;;
    -h | --help)
        echo "Azin compiler installer"
        echo "Usage: $0 [options]"
        echo "  -p, --prefix DIR     Set installation prefix (default: /usr/local)"
        echo "  -b, --branch TARGET  Clone specific git branch or tag"
        echo "  -h, --help           Show this help message"
        exit 0
        ;;
    *)
        log_err "Unknown option '$1'"
        exit 1
        ;;
    esac
done

INSTALL_DIR="$PREFIX/bin"

require() {
    if ! command -v "$1" >/dev/null 2>&1; then
        log_err "Required command '$1' not found. Please install it first."
        exit 1
    fi
}

case "$(uname -s)" in
Linux | Darwin | FreeBSD | OpenBSD | NetBSD | DragonFly | SunOS) ;;
*)
    log_err "Unsupported operating system: $(uname -s)"
    exit 1
    ;;
esac

require git
require go

HAVE_CC=0
if [ -n "${CC:-}" ] && command -v "$CC" >/dev/null 2>&1; then
    HAVE_CC=1
else
    for cc in gcc clang cc; do
        if command -v "$cc" >/dev/null 2>&1; then
            HAVE_CC=1
            break
        fi
    done
fi

if [ "$HAVE_CC" -eq 0 ]; then
    log_err "A C compiler (gcc, clang, or cc) is required."
    exit 1
fi

CLONE_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'azin.XXXXXX')
trap 'rm -rf "$CLONE_DIR"' EXIT INT TERM HUP

GIT_CLONE_ARGS="--depth 1"
if [ -n "$BRANCH" ]; then
    GIT_CLONE_ARGS="$GIT_CLONE_ARGS --branch $BRANCH"
    log_info "Cloning $REPO (branch/tag: $BRANCH)..."
else
    log_info "Cloning $REPO..."
fi

git clone $GIT_CLONE_ARGS "$REPO" "$CLONE_DIR"

cd "$CLONE_DIR"

if [ -f "scripts/build/build.sh" ]; then
    sh scripts/build/build.sh
else
    log_info "Fallback: Building binary directly with Go..."
    mkdir -p build
    go build -trimpath -o "build/$BIN" ./cmd/azc
fi

if [ ! -f "build/$BIN" ]; then
    log_err "Build failed: build/$BIN not found."
    exit 1
fi

find_privilege_command() {
    if [ "$(id -u)" -eq 0 ]; then
        PRIVCMD=""
        return 0
    fi
    for cmd in doas sudo pkexec pfexec; do
        if command -v "$cmd" >/dev/null 2>&1; then
            PRIVCMD="$cmd"
            return 0
        fi
    done
    return 1
}

run_privileged() {
    if [ "$(id -u)" -eq 0 ]; then
        "$@"
    elif [ -n "${PRIVCMD:-}" ]; then
        "$PRIVCMD" "$@"
    else
        log_err "Insufficient permissions to install to $INSTALL_DIR."
        log_err "Install sudo, doas, or run with --prefix ~/.local"
        exit 1
    fi
}

is_writable() {
    target="$1"
    while [ ! -d "$target" ] && [ "$target" != "/" ] && [ "$target" != "." ]; do
        target=$(dirname "$target")
    done
    [ -w "$target" ]
}

copy_binary() {
    src="$1"
    dst="$2"
    if command -v install >/dev/null 2>&1; then
        install -m 755 "$src" "$dst"
    else
        cp "$src" "$dst"
        chmod 755 "$dst"
    fi
}

log_info "Installing $BIN to $INSTALL_DIR..."

if is_writable "$INSTALL_DIR"; then
    mkdir -p "$INSTALL_DIR"
    copy_binary "build/$BIN" "$INSTALL_DIR/$BIN"
else
    find_privilege_command || {
        log_err "Cannot write to $INSTALL_DIR and no privilege escalation tool (sudo/doas) found."
        exit 1
    }
    run_privileged mkdir -p "$INSTALL_DIR"

    # Copy file via temporary location to ensure safe root privilege execution
    TMP_BIN=$(mktemp 2>/dev/null || echo "/tmp/azc_install_tmp")
    copy_binary "build/$BIN" "$TMP_BIN"
    run_privileged copy_binary "$TMP_BIN" "$INSTALL_DIR/$BIN"
    rm -f "$TMP_BIN"
fi

printf "%b\n" "${GREEN}Azin installed successfully!${NC}"
printf "Binary location: %s/%s\n\n" "$INSTALL_DIR" "$BIN"

if command -v "$BIN" >/dev/null 2>&1; then
    "$BIN" -version 2>/dev/null || "$BIN" --help 2>/dev/null || true
else
    log_warn "$INSTALL_DIR is not in your PATH environment variable."
    printf "Add it to your PATH by adding this line to your shell profile (~/.bashrc, ~/.zshrc, etc.):\n"
    printf "  export PATH=\"%s:\$PATH\"\n" "$INSTALL_DIR"
fi
