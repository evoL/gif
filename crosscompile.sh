#!/bin/bash
set -ev

# Setup the C toolchain
osarch="${GIF_OS}/${GIF_ARCH}"

case "$osarch" in
  'darwin/386')
    export CC=o32-clang
    ;;
  'darwin/amd64')
    export CC=o64-clang
    ;;
  'windows/386')
    export CC=i686-w64-mingw32-gcc
    ;;
  'windows/amd64')
    export CC=x86_64-w64-mingw32-gcc
    ;;
esac

# Build it!
GIT_VERSION="$(git describe --tags HEAD)"

gox -osarch="$osarch" -cgo -ldflags="-X github.com/evoL/gif/version.Version ${GIT_VERSION}"
