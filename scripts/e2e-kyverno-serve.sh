#!/usr/bin/env bash
# End-to-end parity check: the same fixtures through `kyverno apply`, `kyverno serve` and the
# Playground backend (which evaluates through kyverno serve) must produce the same rule results.
#
# Usage: scripts/e2e-kyverno-serve.sh <kyverno-cli-binary> <kyverno-repo>/test/cli/serve
# Requires: curl, jq, go.
set -euo pipefail

CLI="${1:?usage: $0 <kyverno-cli-binary> <fixtures-dir>}"
FIX="$(cd "${2:?usage: $0 <kyverno-cli-binary> <fixtures-dir>}" && pwd)"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVE_PORT="${SERVE_PORT:-18581}"
PG_PORT="${PG_PORT:-18580}"
WORK="$(mktemp -d)"
PIDS=()
FAILURES=0

cleanup() {
  for pid in "${PIDS[@]:-}"; do
    [ -n "$pid" ] && kill "$pid" 2>/dev/null || true
  done
}
trap cleanup EXIT

wait_for() {
  for _ in $(seq 1 100); do
    curl -sf "$1" >/dev/null 2>&1 && return 0
    sleep 0.2
  done
  echo "timed out waiting for $1" >&2
  exit 1
}

KYVERNO_EXPERIMENTAL=true "$CLI" serve --allow-request-policies --address "127.0.0.1:$SERVE_PORT" >"$WORK/serve.log" 2>&1 &
PIDS+=("$!")
wait_for "http://127.0.0.1:$SERVE_PORT/readyz"

(cd "$ROOT/backend" && go build -o "$WORK/playground" .)
"$WORK/playground" run --server-host 127.0.0.1 --server-port "$PG_PORT" --engine-serve-url "http://127.0.0.1:$SERVE_PORT" >"$WORK/playground.log" 2>&1 &
PIDS+=("$!")
wait_for "http://127.0.0.1:$PG_PORT/api/config"

# apply's summary line -> "pass fail warn error skip"
apply_summary() {
  # apply exits non-zero on policy violations; only its summary line matters here
  { "$CLI" apply "$@" --remove-color 2>/dev/null || true; } |
    sed -nE 's/.*pass: ([0-9]+), fail: ([0-9]+), warn: ([0-9]+), error: ([0-9]+), skip: ([0-9]+).*/\1 \2 \3 \4 \5/p' | tail -1
}

# count statuses across every result list of a Playground response -> "pass fail warn error skip"
playground_summary() {
  jq -r '[.mutation, .validation, .generation, .imageVerification, .deletion] | map(. // []) | add
         | [.[].policyResponse.rules[].status] as $s
         | [($s|map(select(.=="pass"))|length), ($s|map(select(.=="fail"))|length), ($s|map(select(.=="warn"))|length),
            ($s|map(select(.=="error"))|length), ($s|map(select(.=="skip"))|length)] | map(tostring) | join(" ")' "$1"
}

check() { # name policy-file payload-file apply-flag
  local name="$1" policy="$2" payload="$3" flag="$4"
  local want got serve_got
  want="$(apply_summary "$policy" "$flag" "$payload")"
  jq -n --rawfile p "$policy" --rawfile r "$payload" '{policies: $p, resources: $r, context: ""}' >"$WORK/req.json"
  curl -s -X POST -H 'Content-Type: application/json' --data-binary "@$WORK/req.json" "http://127.0.0.1:$PG_PORT/api/engine" >"$WORK/pg.json"
  got="$(playground_summary "$WORK/pg.json" 2>/dev/null || echo "invalid: $(head -c 300 "$WORK/pg.json")")"
  serve_got="$(jq -n --rawfile p "$policy" --rawfile r "$payload" --arg m "$5" '{mode: $m, policies: $p, payload: $r}' |
    curl -s -X POST -H 'Content-Type: application/json' --data-binary @- "http://127.0.0.1:$SERVE_PORT/v1alpha1/evaluate" |
    jq -r '.summary | [.pass, .fail, .warn, .error, .skip] | map(tostring) | join(" ")')"
  if [ -n "$want" ] && [ "$want" = "$got" ] && [ "$want" = "$serve_got" ]; then
    echo "PASS  $name  (pass fail warn error skip: $want)"
  else
    echo "FAIL  $name  apply=[$want] serve=[$serve_got] playground=[$got]"
    FAILURES=$((FAILURES + 1))
  fi
}

P="$FIX/policies"
D="$FIX/payloads"
for f in "$D"/kubernetes/*; do
  for pol in "$P"/kubernetes/*.yaml; do
    check "kubernetes $(basename "$pol") / $(basename "$f")" "$pol" "$f" -r Kubernetes
  done
done
for f in "$D"/json/*; do check "json / $(basename "$f")" "$P/json/check-dockerfile.yaml" "$f" --json JSON; done
for f in "$D"/http/*; do check "http / $(basename "$f")" "$P/http/allow-only-get.yaml" "$f" --http-payload HTTP; done
for f in "$D"/envoy/*; do check "envoy / $(basename "$f")" "$P/envoy/forbid-admin-path.yaml" "$f" --envoy-payload Envoy; done

if [ "$FAILURES" -gt 0 ]; then
  echo "$FAILURES parity check(s) failed; logs in $WORK" >&2
  exit 1
fi
rm -rf "$WORK"
echo "apply, kyverno serve and the Playground agree on every fixture"
