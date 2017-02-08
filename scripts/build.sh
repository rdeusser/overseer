#!/usr/bin/env bash
#
# This script builds the application from source for multiple platforms.

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd "$DIR" || exit

# Host os and arch
OS="$(go env GOOS)"

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "$(git status --porcelain)" && echo "+CHANGES" || true)

# Delete the old dir
echo "==> Removing old directory..."
rm -rf build/*
mkdir -p build/

# Allow LD_FLAGS to be appended during development compilations
LD_FLAGS="-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.Version=${BIN_VERSION} $LD_FLAGS"
# In relase mode we don't want debug information in the binary
if [[ -n "${M_PROD}" ]]; then
    LD_FLAGS="-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.Version=${BIN_VERSION} -s -w"
fi

# Build!
echo "==> Building..."
mkdir -p build/linux  && GOOS=linux  go build -ldflags "${LD_FLAGS}" -o "build/linux/overseer"
mkdir -p build/darwin && GOOS=darwin go build -ldflags "${LD_FLAGS}" -o "build/darwin/overseer"

case "$OS" in
    "linux")
        cp "build/linux/overseer" "$GOPATH/bin"
        ;;
    "darwin")
        cp "build/darwin/overseer" "$GOPATH/bin"
        ;;
    *)
        echo "couldn't detect your os version using go env"
        exit 1
        ;;
esac

# Done!
echo
echo "==> Results:"
ls -hl build/
