#!/usr/bin/env bash
# Runs the backend tests against a local `kyverno serve`, the Playground's evaluation backend.
# Usage: scripts/test-backend-serve.sh <kyverno-cli-binary> [port]
set -euo pipefail

CLI="${1:-kyverno}"
PORT="${2:-8081}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

KYVERNO_EXPERIMENTAL=true "$CLI" serve --allow-request-policies --address "127.0.0.1:$PORT" &
SERVE_PID=$!
trap 'kill "$SERVE_PID" 2>/dev/null || true' EXIT

for _ in $(seq 1 50); do
  if curl -sf "http://127.0.0.1:$PORT/readyz" >/dev/null; then
    break
  fi
  sleep 0.2
done

cd "$ROOT/backend"
KYVERNO_SERVE_URL="http://127.0.0.1:$PORT" go test ./... -race -coverprofile=coverage.out -covermode=atomic
