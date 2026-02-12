#!/bin/bash
# Deploy resume-builder to VPS (resume.ivanpodgurskiy.name)
#
# Build strategy: build everything locally, then rsync to VPS
# - Go backend: built locally, binary copied to VPS
# - Frontend: built locally (npm run build), build/ copied to VPS
#
# STEP 1 (local): Run this script - builds all, shows transfer stats, copies to VPS
# STEP 2 (VPS):   SSH in and run docker compose (images use pre-built artifacts)
#
# Usage: ./deploy.sh

set -e

VPS_HOST="92.205.19.23"
VPS_USER="ivanpod"
REMOTE_PATH="~/docker/projects/resume-builder"

RSYNC_EXCLUDES=(
  --exclude 'node_modules'
  --exclude '.git'
  --exclude 'resume-frontend/node_modules'
  --exclude 'resume-frontend/.svelte-kit'
  --exclude 'resume-backend/.env'
)

echo "🔨 Building Go backend locally..."
cd resume-backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api
cd ..
echo "✅ Go binary built"
echo ""

echo "🔨 Building frontend locally..."
cd resume-frontend
npm run build
cd ..
echo "✅ Frontend built (resume-frontend/build)"
echo ""

echo "📊 Files to copy..."
echo "----------------------------------------"
STAGING_DIR=$(mktemp -d)
trap 'rm -rf "$STAGING_DIR"' EXIT
# Mirror what rsync would send (same excludes)
rsync -a --delete "${RSYNC_EXCLUDES[@]}" ./ "$STAGING_DIR/" >/dev/null 2>&1
FILE_COUNT=$(find "$STAGING_DIR" -type f | wc -l | tr -d ' ')
DIR_COUNT=$(find "$STAGING_DIR" -type d | wc -l | tr -d ' ')
TOTAL_KB=$(du -sk "$STAGING_DIR" 2>/dev/null | awk '{print $1}')
TOTAL_BYTES=$((TOTAL_KB * 1024))
echo "  Files:    $FILE_COUNT"
echo "  Dirs:     $DIR_COUNT"
if [ -n "$TOTAL_BYTES" ] && [ "$TOTAL_BYTES" -gt 0 ]; then
  if [ "$TOTAL_BYTES" -ge 1048576 ]; then
    SIZE_MB=$(echo "scale=2; $TOTAL_BYTES / 1048576" | bc 2>/dev/null || echo "$((TOTAL_BYTES / 1048576))")
    echo "  Size:     ${SIZE_MB} MB"
  elif [ "$TOTAL_BYTES" -ge 1024 ]; then
    SIZE_KB=$(echo "scale=2; $TOTAL_BYTES / 1024" | bc 2>/dev/null || echo "$((TOTAL_BYTES / 1024))")
    echo "  Size:     ${SIZE_KB} KB"
  else
    echo "  Size:     ${TOTAL_BYTES} B"
  fi
else
  echo "  Size:     (calculating...)"
fi
echo "----------------------------------------"
echo ""

echo "📤 Copying to VPS..."
rsync -avz --delete "${RSYNC_EXCLUDES[@]}" ./ "${VPS_USER}@${VPS_HOST}:${REMOTE_PATH}/"

echo ""
echo "✅ Code copied. Now SSH to VPS and run:"
echo ""
echo "   ssh ${VPS_USER}@${VPS_HOST}"
echo "   cd ~/docker/projects/resume-builder"
echo "   docker compose -f docker-compose.prod.yml build --no-cache"
echo "   docker compose -f docker-compose.prod.yml up -d"
echo "   docker compose -f docker-compose.prod.yml ps"
echo ""

# =============================================================================
# ALTERNATIVE: Git pull on VPS (builds backend on VPS - slow, 20+ min)
# =============================================================================
# For git pull flow, temporarily use full Dockerfile: in docker-compose.prod.yml
# change backend dockerfile to "Dockerfile" (not Dockerfile.runtime).
#
#   ssh ivanpod@92.205.19.23
#
#   # First time only:
#   mkdir -p ~/docker/projects && cd ~/docker/projects
#   git clone git@github.com:ivan-podgurskiy/resume-builder.git
#   cp resume-builder/resume-backend/.env.example resume-builder/resume-backend/.env
#   # Edit .env: nano resume-builder/resume-backend/.env
#
#   # Every deploy (backend build on VPS - slow):
#   cd ~/docker/projects/resume-builder
#   git pull
#   # Edit docker-compose.prod.yml: backend dockerfile -> "Dockerfile"
#   docker compose -f docker-compose.prod.yml build --no-cache
#   docker compose -f docker-compose.prod.yml up -d
#   docker compose -f docker-compose.prod.yml ps
#
# =============================================================================
# COPY COMMANDS TO RUN ON VPS (avoids SSH timeout during 5-10 min build):
# =============================================================================
#
#   ssh ivanpod@92.205.19.23
#   cd ~/docker/projects/resume-builder
#   docker compose -f docker-compose.prod.yml build --no-cache
#   docker compose -f docker-compose.prod.yml up -d
#   docker compose -f docker-compose.prod.yml ps
#
# Or in tmux: tmux new -s build, then run above, Ctrl+B D to detach
#
# =============================================================================
