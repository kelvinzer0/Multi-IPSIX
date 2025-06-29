#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Check if running as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Please use sudo." >&2
  exit 1
fi

# --- Configuration ---
APP_NAME="multi-ipsix"
BIN_PATH="/usr/local/bin/$APP_NAME"
CONFIG_DIR="/etc/$APP_NAME"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"

# --- Stop and disable the service ---
echo "Stopping and disabling systemd service..."

if systemctl is-active --quiet "$APP_NAME.service"; then
  systemctl stop "$APP_NAME.service"
fi

if systemctl is-enabled --quiet "$APP_NAME.service"; then
  systemctl disable "$APP_NAME.service"
fi

# --- Remove files ---
echo "Removing installed files..."

rm -f "$SERVICE_FILE"
rm -f "$BIN_PATH"

# Ask for confirmation before deleting the configuration directory
if [ -d "$CONFIG_DIR" ]; then
  read -p "Do you want to remove the configuration directory ($CONFIG_DIR)? [y/N] " -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -r "$CONFIG_DIR"
    echo "Configuration directory removed."
  else
    echo "Skipping removal of configuration directory."
  fi
fi

# --- Reload systemd ---
echo "Reloading systemd daemon..."
systemctl daemon-reload

echo "Uninstallation complete."
