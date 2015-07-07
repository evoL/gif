#!/bin/sh
set -ev

if [ "$GIF_OS" == "darwin" ]; then
  # Download osxcross to prepare the C toolchain
  git clone https://github.com/tpoechtrager/osxcross /tmp/osxcross
  cd /tmp/osxcross

  # Download prepared SDKs
  wget --quiet -O tarballs/MacOSX10.10.sdk.tar.xz https://dl.dropboxusercontent.com/u/2078673/osx-sdks/MacOSX10.10.sdk.tar.xz
  wget --quiet -O tarballs/MacOSX10.9.sdk.tar.xz https://dl.dropboxusercontent.com/u/2078673/osx-sdks/MacOSX10.9.sdk.tar.xz

  # Build the toolchain
  UNATTENDED=1 SDK_VERSION=10.10 OSX_VERSION_MIN=10.9 ./build.sh
else
  # Make sure the PATH doesn't break
  mkdir -p /tmp/osxcross/target/bin
fi
