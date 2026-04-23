#!/bin/bash

# ═══════════════════════════════════════════════════════════════════════════
# AccessVille-BE - Docker Build & Deployment Management Script
# ═══════════════════════════════════════════════════════════════════════════

# Get Git SHA for image tagging
GIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DEFAULT_TAG="${GIT_SHA}"

# Environment variables with defaults
PLATFORM="${PLATFORM:-linux/amd64}" 
REGISTRY="${DOCKER_REGISTRY:-docker.io/aprilpollo}"
IMAGE_TAG="${IMAGE_TAG:-$DEFAULT_TAG}"

set -euo pipefail

# ═══════════════════════════════════════════════════════════════════════════
# Colors & Logging
# ═══════════════════════════════════════════════════════════════════════════
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m' # No Color

print_status()   { echo -e "${GREEN}✓ [SUCCESS]${NC} $1"; }
print_warning()  { echo -e "${YELLOW}⚠ [WARNING]${NC} $1"; }
print_error()    { echo -e "${RED}✗ [ERROR]${NC} $1"; }
print_info()     { echo -e "${BLUE}ℹ [INFO]${NC} $1"; }
print_step()     { echo -e "${CYAN}▶ [STEP]${NC} $1"; }
print_header()   { echo -e "\n${BOLD}${MAGENTA}━━━ $1 ━━━${NC}\n"; }

# ═══════════════════════════════════════════════════════════════════════════
# Prerequisites Check
# ═══════════════════════════════════════════════════════════════════════════
if ! command -v docker &>/dev/null; then
  print_error "Docker is not installed. Please install Docker first."
  exit 1
fi

# Detect Compose command (v2 plugin preferred, then v1 binary)
if docker compose version &>/dev/null; then
  COMPOSE="docker compose"
elif command -v docker-compose &>/dev/null; then
  COMPOSE="docker-compose"
else
  print_error "Docker Compose is not installed. Install either 'docker compose' (v2) or 'docker-compose' (v1)."
  exit 1
fi

# Check Git availability for versioning
if ! command -v git &>/dev/null; then
  print_warning "Git not found. Using default tag without SHA."
fi

# ═══════════════════════════════════════════════════════════════════════════
# Path Configuration
# ═══════════════════════════════════════════════════════════════════════════
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "$SCRIPT_DIR"

COMPOSE_FILE="docker/docker-compose.yml"
ENV_FILE=".env"

if [ ! -f "$COMPOSE_FILE" ]; then
  print_error "Compose file not found at: $COMPOSE_FILE"
  exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
  print_warning ".env not found at: $ENV_FILE"
  print_warning "Services may not start properly without environment variables."
fi

# ═══════════════════════════════════════════════════════════════════════════
# Build Functions
# ═══════════════════════════════════════════════════════════════════════════
build() {
  print_header "Building All Docker Images"
  print_info "Tag: ${IMAGE_TAG}"
  $COMPOSE -f "$COMPOSE_FILE" build
  print_status "Build completed successfully!"
}

build_service() {
  local service=$1
  if [ -z "$service" ]; then
    print_error "Service name is required"
    print_info "Available services: backend, backend-lpr, backend-rtsp"
    exit 1
  fi
  
  print_status "Building $service..."
  $COMPOSE -f "$COMPOSE_FILE" build "$service"
  print_status "Build completed for $service!"
}

build_and_push() {
  print_header "Building & Pushing"
  
  local image_name="${REGISTRY}/task-manager:${IMAGE_TAG}"
  
  print_step "Building image: $image_name"
  print_info "Platform: $PLATFORM"
  
  DOCKER_BUILDKIT=1 docker build \
    -f docker/Dockerfile \
    -t "$image_name" \
    --platform "$PLATFORM" \
    . \
    --push
  
  print_status "Image pushed: $image_name"
}

# ═══════════════════════════════════════════════════════════════════════════
# Help & Documentation
# ═══════════════════════════════════════════════════════════════════════════
help() {
  echo ""
  echo -e "${BOLD}${MAGENTA}╔═══════════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${BOLD}${MAGENTA}║${NC}  ${BOLD}Task Management - Docker Management Script${NC}                       ${BOLD}${MAGENTA}║${NC}"
  echo -e "${BOLD}${MAGENTA}╚═══════════════════════════════════════════════════════════════════╝${NC}"
  echo ""
  echo -e "${YELLOW}Usage:${NC} $0 [COMMAND] [OPTIONS]"
  echo ""
  echo -e "${CYAN}Current Configuration:${NC}"
  echo "  Registry: ${REGISTRY}"
  echo "  Tag:      ${IMAGE_TAG}"
  echo "  Platform: ${PLATFORM}"
  echo "  Git SHA:  ${GIT_SHA}"
  echo ""

  echo -e "${YELLOW}Build Commands:${NC}"
  echo "  build                     - Build Docker images"
  echo "  push                      - Build and Push Docker images"
  echo ""
}

# ═══════════════════════════════════════════════════════════════════════════
# Main Execution
# ═══════════════════════════════════════════════════════════════════════════

# Parse command-line arguments
command="${1:-}"
shift || true

# Handle commands
case "$command" in
  build) build;;
  push) build_and_push;;
  tag) 
    if [ -n "$1" ]; then IMAGE_TAG="$1"; shift; fi
    build_and_push;;
  help) help;;
  *) 
    # print_error "Unknown command: $command"
    help;;
esac