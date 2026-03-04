#!/usr/bin/env bash
set -euo pipefail

# Load .env if present
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
if [[ -f "$PROJECT_ROOT/.env" ]]; then
    set -a
    source "$PROJECT_ROOT/.env"
    set +a
fi

BINARY="$PROJECT_ROOT/ha-ctl"

if [[ ! -x "$BINARY" ]]; then
    echo "FAIL: binary not found at $BINARY (run 'make build' first)"
    exit 1
fi

PASS=0
FAIL=0

run_test() {
    local name="$1"
    shift
    if output=$("$@" 2>&1); then
        echo "PASS: $name"
        PASS=$((PASS + 1))
    else
        echo "FAIL: $name"
        echo "  Output: $output"
        FAIL=$((FAIL + 1))
    fi
}

run_test_json() {
    local name="$1"
    shift
    if output=$("$@" 2>/dev/null) && echo "$output" | python3 -m json.tool > /dev/null 2>&1; then
        echo "PASS: $name"
        PASS=$((PASS + 1))
    else
        echo "FAIL: $name (invalid JSON or command failed)"
        echo "  Output: $output"
        FAIL=$((FAIL + 1))
    fi
}

echo "=== ha-ctl smoke tests ==="
echo ""

# Version
run_test "version" "$BINARY" version

# Entities (JSON output)
run_test_json "entities" "$BINARY" entities

# Entities with domain filter
run_test_json "entities --domain input_boolean" "$BINARY" entities --domain input_boolean

# State
if [[ -n "${HA_TEST_ENTITY:-}" ]]; then
    run_test_json "state $HA_TEST_ENTITY" "$BINARY" state "$HA_TEST_ENTITY"

    # Call service
    run_test_json "call turn_on" "$BINARY" call input_boolean turn_on --entity "$HA_TEST_ENTITY"
    run_test_json "call turn_off" "$BINARY" call input_boolean turn_off --entity "$HA_TEST_ENTITY"
else
    echo "SKIP: state and call tests (HA_TEST_ENTITY not set)"
fi

# Context (compact default)
run_test_json "context" "$BINARY" context

# Context --full
run_test_json "context --full" "$BINARY" context --full

# Context --domain
run_test_json "context --domain sensor" "$BINARY" context --domain sensor

# Find
run_test_json "find" "$BINARY" find "light"

# Find with domain filter
run_test_json "find --domain light" "$BINARY" find "light" --domain light

# Entities with state filter
run_test_json "entities --state on" "$BINARY" entities --state on

# Entities with domain + state filter
run_test_json "entities --domain light --state on" "$BINARY" entities --domain light --state on

# Cache refresh
run_test_json "cache refresh" "$BINARY" cache refresh

echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="

if [[ $FAIL -gt 0 ]]; then
    exit 1
fi
