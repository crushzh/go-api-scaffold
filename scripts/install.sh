#!/bin/bash
#
# Installation script
#
# Usage: ./install.sh [install directory]
# Default install to /opt/myapp
#

set -e

# ==================== Configuration ====================
APP_NAME="myapp"
DEFAULT_INSTALL_DIR="/opt/${APP_NAME}"
INSTALL_DIR="${1:-$DEFAULT_INSTALL_DIR}"

# ==================== Color output ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# ==================== Detection ====================
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PACKAGE_DIR="$(dirname "$SCRIPT_DIR")"

# Check binary file
if [ ! -f "${PACKAGE_DIR}/${APP_NAME}" ]; then
    error "${APP_NAME} executable not found, please build first: make build"
fi

# ==================== Install ====================
info "Installing ${APP_NAME} to ${INSTALL_DIR} ..."

# If already installed, stop old service first
if [ -f "${INSTALL_DIR}/scripts/manage.sh" ]; then
    info "Existing installation detected, stopping old service..."
    "${INSTALL_DIR}/scripts/manage.sh" stop 2>/dev/null || true
    sleep 2
fi

# Create install directories
mkdir -p "${INSTALL_DIR}"
mkdir -p "${INSTALL_DIR}/data"
mkdir -p "${INSTALL_DIR}/logs"
mkdir -p "${INSTALL_DIR}/configs"
mkdir -p "${INSTALL_DIR}/scripts"

# Copy files
info "Copying program files..."
cp "${PACKAGE_DIR}/${APP_NAME}" "${INSTALL_DIR}/"
chmod +x "${INSTALL_DIR}/${APP_NAME}"

# Copy config (do not overwrite existing)
if [ ! -f "${INSTALL_DIR}/configs/config.yaml" ]; then
    cp "${PACKAGE_DIR}/configs/config.yaml" "${INSTALL_DIR}/configs/"
    info "Default config file installed"
else
    warn "Config file already exists, skipping overwrite"
fi

# Copy management script
cp "${PACKAGE_DIR}/scripts/manage.sh" "${INSTALL_DIR}/scripts/"
chmod +x "${INSTALL_DIR}/scripts/manage.sh"

# Copy frontend files (if present)
if [ -d "${PACKAGE_DIR}/dist" ]; then
    cp -r "${PACKAGE_DIR}/dist" "${INSTALL_DIR}/"
    info "Frontend files installed"
fi

# ==================== Complete ====================
info ""
info "=========================================="
info "  ${APP_NAME} installation complete!"
info "=========================================="
info ""
info "  Install directory:  ${INSTALL_DIR}"
info "  Config file:        ${INSTALL_DIR}/configs/config.yaml"
info "  Data directory:     ${INSTALL_DIR}/data/"
info "  Log directory:      ${INSTALL_DIR}/logs/"
info ""
info "  Management commands:"
info "    ${INSTALL_DIR}/scripts/manage.sh start    Start"
info "    ${INSTALL_DIR}/scripts/manage.sh stop     Stop"
info "    ${INSTALL_DIR}/scripts/manage.sh restart  Restart"
info "    ${INSTALL_DIR}/scripts/manage.sh status   Status"
info ""
