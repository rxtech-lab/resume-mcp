#!/bin/bash

# Exit on any error
set -e

# Check if required variables are set
if [ -z "${INSTALLER_SIGNING_CERTIFICATE_NAME}" ]; then
  echo "Error: INSTALLER_SIGNING_CERTIFICATE_NAME is not set"
  exit 1
fi

if [ -z "${APPLE_ID}" ] || [ -z "${APPLE_ID_PWD}" ] || [ -z "${APPLE_TEAM_ID}" ]; then
  echo "Error: Apple ID credentials not set (APPLE_ID, APPLE_ID_PWD, APPLE_TEAM_ID)"
  exit 1
fi

# Import binary configuration
source "$(dirname "$0")/binaries.sh"

PKG_FILE="resume-mcp_macOS_arm64.pkg"
TMP_DIR="tmp_pkg_build"

# Create a temporary directory structure for pkgbuild
echo "Creating package structure"
mkdir -p "${TMP_DIR}/usr/local/bin"
mkdir -p "${TMP_DIR}_scripts"

# Verify and copy each binary
for binary in "${BINARIES[@]}"; do
  BINARY_PATH="bin/${binary}"
  
  # Verify the binary is signed
  echo "Verifying binary signature: ${BINARY_PATH}"
  codesign --verify --verbose "${BINARY_PATH}" || {
    echo "Error: Binary ${binary} is not properly signed. Run sign.sh first."
    exit 1
  }
  
  # Copy binary to temporary directory
  cp "${BINARY_PATH}" "${TMP_DIR}/usr/local/bin/"
done

# Copy post-install script
echo "Adding post-install script"
cp "scripts/post-install.sh" "${TMP_DIR}_scripts/postinstall"
chmod +x "${TMP_DIR}_scripts/postinstall"

# Create the pkg file
echo "Building pkg installer"
pkgbuild --root "${TMP_DIR}" \
  --scripts "${TMP_DIR}_scripts" \
  --identifier "com.rxtech-lab.resume-mcp" \
  --version "1.0" \
  --sign "${INSTALLER_SIGNING_CERTIFICATE_NAME}" \
  --install-location "/" \
  "${PKG_FILE}"

# Clean up temporary directories
rm -rf "${TMP_DIR}"
rm -rf "${TMP_DIR}_scripts"

# Notarize the pkg file
echo "Submitting for notarization"
xcrun notarytool submit "${PKG_FILE}" --verbose --apple-id "${APPLE_ID}" --team-id "${APPLE_TEAM_ID}" --password "${APPLE_ID_PWD}" --wait

# Staple the notarization ticket to the pkg
echo "Stapling notarization ticket"
xcrun stapler staple -v "${PKG_FILE}"

echo "Package created, signed, notarized and stapled successfully: ${PKG_FILE}" 