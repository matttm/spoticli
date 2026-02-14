#!/usr/bin/env bash
set -euo pipefail
# Helper to run the integration test locally.
# Usage: ./scripts/run-integration.sh

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "Bringing up docker-compose services..."
docker compose -f docker-compose.yml up -d --build
# Only teardown automatically when not debugging
if [ "${INTEGRATION_DEBUG:-}" != "1" ]; then
  trap 'echo "Tearing down docker-compose..."; docker compose -f docker-compose.yml down -v' EXIT
else
  echo "INTEGRATION_DEBUG=1; leaving docker-compose stack running for inspection"
fi

# Ensure test tone fixture is available for integration tests
echo "Ensuring test tone fixture exists..."
# First check if the test tone MP3 already exists to avoid unnecessary work
if [ ! -f "./assets/test_tone.mp3" ]; then
  # Test tone not found, check if ffmpeg is available to generate it
  if command -v ffmpeg >/dev/null 2>&1; then
    echo "Generating test tone into assets/"
    ./scripts/generate_tone.sh ./assets || true
  else
    echo "ffmpeg not found; skipping test tone generation"
  fi
else
  echo "Test tone already exists, skipping generation"
fi

echo "Waiting for backend health endpoint..."
for i in $(seq 1 90); do
  if curl -fsS http://localhost:4200/ >/dev/null 2>&1; then
    echo "backend ready"
    break
  fi
  sleep 2
done

echo "Running integration test..."
INTEGRATION=1 go test . -run TestIntegration_UploadDownloadStream -v
echo "Integration test finished"

# Ensure docker-compose is torn down after the test completes
echo "Tearing down docker-compose..."
docker compose -f docker-compose.yml down -v || true
