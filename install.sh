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
BIN_PATH="/usr/local/bin"
CONFIG_DIR="/etc/$APP_NAME"
CONFIG_FILE="atpajah.yaml"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"

# --- Build the application ---
echo "Building the application..."
go build -o $APP_NAME -buildvcs=false ./cmd/$APP_NAME
echo "Build successful."

# --- Install files ---
echo "Installing files..."

# Create configuration directory
mkdir -p "$CONFIG_DIR"

# Copy configuration file if it doesn't exist
if [ ! -f "$CONFIG_DIR/$CONFIG_FILE" ]; then
  echo "Copying default configuration..."
  cp "$CONFIG_FILE" "$CONFIG_DIR/"
else
  echo "Configuration file already exists. Skipping copy."
fi

# Move the binary to the bin path
mv "$APP_NAME" "$BIN_PATH/"

# --- Create systemd service file ---
echo "Creating systemd service file..."

cat > "$SERVICE_FILE" << EOL
[Unit]
Description=Multi-IPSIX Service to manage IPv6 addresses
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
ExecStart=$BIN_PATH/$APP_NAME
User=root
Group=root
RemainAfterExit=true

[Install]
WantedBy=multi-user.target
EOL

# --- Enable and start the service ---
echo "Reloading systemd daemon and enabling the service..."

systemctl daemon-reload
systemctl enable $APP_NAME.service

echo "Installation complete!"
echo "You can start the service with: sudo systemctl start $APP_NAME"
echo "Check its status with: sudo systemctl status $APP_NAME"
