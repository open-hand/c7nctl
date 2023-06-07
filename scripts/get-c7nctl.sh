#!/usr/bin/env bash

: ${BINARY_NAME:="c7nctl"}
: ${USE_SUDO:="true"}
: ${DEBUG:="false"}
: ${VERIFY_CHECKSUM:="true"}
: ${VERIFY_SIGNATURES:="false"}
: ${C7NCTL_INSTALL_DIR:="."}

HAS_CURL="$(type "curl" &> /dev/null && echo true || echo false)"
HAS_WGET="$(type "wget" &> /dev/null && echo true || echo false)"
HAS_OPENSSL="$(type "openssl" &> /dev/null && echo true || echo false)"
HAS_GPG="$(type "gpg" &> /dev/null && echo true || echo false)"

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    #armv5*) ARCH="armv5";;
    #armv6*) ARCH="armv6";;
    #armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    #x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    #i686) ARCH="386";;
    #i386) ARCH="386";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
  esac
}

# runs the given command as root (detects if we are root already)
runAsRoot() {
  local CMD="$*"

  if [ $EUID -ne 0 -a $USE_SUDO = "true" ]; then
    CMD="sudo $CMD"
  fi

  $CMD
}

# verifySupported checks that the os/arch combination is supported for
# binary builds, as well whether or not necessary tools are present.
verifySupported() {
  local supported="darwin-amd64\nlinux-amd64\nlinux-arm64\nwindows-amd64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuilt binary for ${OS}-${ARCH}."
    echo "To build from source, go to https://github.com/openhand/c7nctl"
    exit 1
  fi

  if [ "${HAS_CURL}" != "true" ] && [ "${HAS_WGET}" != "true" ]; then
    echo "Either curl or wget is required"
    exit 1
  fi

  if [ "${VERIFY_CHECKSUM}" == "true" ] && [ "${HAS_OPENSSL}" != "true" ]; then
    echo "In order to verify checksum, openssl must first be installed."
    echo "Please install openssl or set VERIFY_CHECKSUM=false in your environment."
    exit 1
  fi

  if [ "${VERIFY_SIGNATURES}" == "true" ]; then
    if [ "${HAS_GPG}" != "true" ]; then
      echo "In order to verify signatures, gpg must first be installed."
      echo "Please install gpg or set VERIFY_SIGNATURES=false in your environment."
      exit 1
    fi
    if [ "${OS}" != "linux" ]; then
      echo "Signature verification is currently only supported on Linux."
      echo "Please set VERIFY_SIGNATURES=false or verify the signatures manually."
      exit 1
    fi
  fi
}

# checkDesiredVersion checks if the desired version is available.
checkDesiredVersion() {
  if [ "x$DESIRED_VERSION" == "x" ]; then
    # Get tag from release URL
    local latest_release_url="https://file.choerodon.com.cn/choerodon-install/c7nctl/"
    if [ "${HAS_CURL}" == "true" ]; then
      TAG=$(curl -Ls $latest_release_url | grep 'href="1.[0-9]*.[0-9]*\/"' | grep -v no-underline | tail -n 1 | cut -d '"' -f 2 | sed s'/.$//')
    elif [ "${HAS_WGET}" == "true" ]; then
      TAG=$(wget $latest_release_url -O - 2>&1 | grep 'href="1.[0-9]*.[0-9]*/\"' | grep -v no-underline | tail -n 1 | cut -d '"' -f 2 | sed s'/.$//')
    fi
  else
    TAG=$DESIRED_VERSION
  fi
}


# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  C7NCTL_DIST="c7nctl-$TAG-$OS-$ARCH.tar.gz"
  DOWNLOAD_URL="https://file.choerodon.com.cn/choerodon-install/c7nctl/$TAG/$C7NCTL_DIST"
  CHECKSUM_URL="$DOWNLOAD_URL.sha256"
  C7NCTL_TMP_ROOT="$(mktemp -dt c7nctl-installer-XXXXXX)"
  C7NCTL_TMP_FILE="$C7NCTL_TMP_ROOT/$C7NCTL_DIST"
  C7NCTL_SUM_FILE="$C7NCTL_TMP_ROOT/$C7NCTL_DIST.sha256"
  echo "Downloading $DOWNLOAD_URL"
  if [ "${HAS_CURL}" == "true" ]; then
    curl -SsL "$CHECKSUM_URL" -o "$C7NCTL_SUM_FILE"
    curl -SsL "$DOWNLOAD_URL" -o "$C7NCTL_TMP_FILE"
  elif [ "${HAS_WGET}" == "true" ]; then
    wget -q -O "$C7NCTL_SUM_FILE" "$CHECKSUM_URL"
    wget -q -O "$C7NCTL_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# verifyFile verifies the SHA256 checksum of the binary package
# and the GPG signatures for both the package and checksum file
# (depending on settings in environment).
verifyFile() {
  if [ "${VERIFY_CHECKSUM}" == "true" ]; then
    verifyChecksum
  fi
}

# installFile installs the Helm binary.
installFile() {
  C7NCTL_TMP="$C7NCTL_TMP_ROOT/$BINARY_NAME"
  mkdir -p "$C7NCTL_TMP"
  tar xf "$C7NCTL_TMP_FILE" -C "$C7NCTL_TMP"
  C7NCTL_TMP_BIN="$C7NCTL_TMP/c7nctl-$TAG-$OS-$ARCH/c7nctl"
  echo "Preparing to install $BINARY_NAME into ${C7NCTL_INSTALL_DIR}"
  runAsRoot cp "$C7NCTL_TMP_BIN" "$C7NCTL_INSTALL_DIR/$BINARY_NAME"
  echo "$BINARY_NAME installed into $C7NCTL_INSTALL_DIR/$BINARY_NAME"
}

# verifyChecksum verifies the SHA256 checksum of the binary package.
verifyChecksum() {
  printf "Verifying checksum... "
  local sum=$(openssl sha1 -sha256 ${C7NCTL_TMP_FILE} | awk '{print $2}')
  local expected_sum=$(cat ${C7NCTL_SUM_FILE})
  if [ "$sum" != "$expected_sum" ]; then
    echo "SHA sum of ${C7NCTL_TMP_FILE} does not match. Aborting."
    exit 1
  fi
  echo "Done."
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    if [[ -n "$INPUT_ARGUMENTS" ]]; then
      echo "Failed to install $BINARY_NAME with the arguments provided: $INPUT_ARGUMENTS"
      help
    else
      echo "Failed to install $BINARY_NAME"
    fi
    echo -e "\tFor support, go to https://github.com/openhand/c7nctl"
  fi
  cleanup
  exit $result
}

# help provides possible cli installation arguments
help () {
  echo "Accepted cli arguments are:"
  echo -e "\t[--help|-h ] ->> prints this help"
  echo -e "\t[--version|-v <desired_version>] . When not defined it fetches the latest release from GitHub"
  echo -e "\te.g. --version v3.0.0 or -v canary"
  echo -e "\t[--no-sudo]  ->> install without sudo"
}

# cleanup temporary files to avoid https://github.com/helm/helm/issues/2977
cleanup() {
  if [[ -d "${C7NCTL_TMP_ROOT:-}" ]]; then
    rm -rf "$C7NCTL_TMP_ROOT"
  fi
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e

# Set debug if desired
if [ "${DEBUG}" == "true" ]; then
  set -x
fi

# Parsing input arguments (if any)
export INPUT_ARGUMENTS="${@}"
set -u
while [[ $# -gt 0 ]]; do
  case $1 in
    '--version'|-v)
       shift
       if [[ $# -ne 0 ]]; then
           export DESIRED_VERSION="${1}"
       else
           echo -e "Please provide the desired version. e.g. --version v3.0.0 or -v canary"
           exit 0
       fi
       ;;
    '--no-sudo')
       USE_SUDO="false"
       ;;
    '--help'|-h)
       help
       exit 0
       ;;
    *) exit 1
       ;;
  esac
  shift
done
set +u

initArch
initOS
verifySupported
checkDesiredVersion
downloadFile
verifyFile
installFile

cleanup