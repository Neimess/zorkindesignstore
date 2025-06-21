#!/usr/bin/env bash
set -euo pipefail

SERVICE_NAME="zorkin-backend"
UNIT_SRC="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/deploy/${SERVICE_NAME}.service"
UNIT_DST="/etc/systemd/system/${SERVICE_NAME}.service"
SERVICE_DIR="/opt/${SERVICE_NAME}"

echo "› Checking root privileges..."
if [[ $EUID -ne 0 ]]; then
  echo "Please run as root (sudo ./scripts/install_service.sh)"
  exit 1
fi

echo "› Creating service directory: $SERVICE_DIR"
mkdir -p "$SERVICE_DIR"
chown "$SUDO_USER":"$SUDO_USER" "$SERVICE_DIR"

echo "› Copying unit file to $UNIT_DST"
cp "$UNIT_SRC" "$UNIT_DST"

echo "› Reloading systemd daemon"
systemctl daemon-reload

echo "› Enabling $SERVICE_NAME service"
systemctl enable "$SERVICE_NAME"

echo "› Done. You can now start it with:"
echo "    sudo systemctl start $SERVICE_NAME"
