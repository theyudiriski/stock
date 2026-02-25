#!/usr/bin/env bash
# Deploy stock app to VPS: pull latest code and rebuild on the server.
# Usage: ./deploy/deploy.sh
# Edit VPS_USER, VPS_HOST, VPS_APP_DIR below. config.yaml is never touched.

set -e

# --- Edit these ---
VPS_USER="${VPS_USER:-root}"
VPS_HOST="${VPS_HOST:-}"
VPS_APP_DIR="${VPS_APP_DIR:-/opt/stock}"
SSH_KEY="${SSH_KEY:-}"   # e.g. -i ~/.ssh/stock_vps

# --- No need to edit below ---
if [[ -z "$VPS_HOST" ]]; then
  echo "Set VPS_HOST (and optionally VPS_USER, VPS_APP_DIR) in this script or env."
  echo "Example: VPS_HOST=1.2.3.4 ./deploy/deploy.sh"
  exit 1
fi

SSH_OPTS=()
[[ -n "$SSH_KEY" ]] && SSH_OPTS=(-i "$SSH_KEY")

echo "Deploying to $VPS_USER@$VPS_HOST ($VPS_APP_DIR) ..."
ssh "${SSH_OPTS[@]}" "$VPS_USER@$VPS_HOST" "cd $VPS_APP_DIR && git pull && go build -o stock ./cmd/"
echo "Done. Binary updated at $VPS_APP_DIR/stock"
echo "If you added migrations, run on VPS: ./stock -type=migrate -config config.yaml"
