#!/bin/bash
set -euo pipefail

# Test runner for nightly acceptance tests
# Usage: ./test-runner.sh <type> [service] [coverage-file]

TYPE="${1:-}"
SERVICE="${2:-}"
COVERAGE_FILE="${3:-coverage.txt}"

if [[ -z "$TYPE" ]]; then
    echo "Usage: $0 <type> [service] [coverage-file]"
    echo "Types: provider-core, resources, datasources"
    exit 1
fi

run_provider_core_tests() {
    echo "Running provider core tests..."
    go test -v -race \
        -coverprofile="$COVERAGE_FILE" \
        -covermode=atomic \
        ./internal/client/... \
        ./internal/helpers/... \
        ./internal/provider/... \
        ./internal/utilities/...
}

run_service_tests() {
    local category="$1"
    local service="$2"

    echo "Running tests for ${category}/${service}..."

    local test_dir="./internal/services/${category}/${service}"

    if [[ ! -d "$test_dir" ]]; then
        echo "Directory not found: ${test_dir}, creating empty coverage file"
        echo "mode: atomic" > "$COVERAGE_FILE"
        return 0
    fi

    if ! ls "${test_dir}"/*_test.go >/dev/null 2>&1; then
        echo "No test files found in ${test_dir}, creating empty coverage file"
        echo "mode: atomic" > "$COVERAGE_FILE"
        return 0
    fi

    local test_count=$(ls "${test_dir}"/*_test.go 2>/dev/null | wc -l | tr -d ' ')
    echo "Found ${test_count} test files"

    go test -v -race \
        -coverprofile="$COVERAGE_FILE" \
        -covermode=atomic \
        "${test_dir}/..."
}

case "$TYPE" in
    provider-core)
        run_provider_core_tests
        ;;
    resources)
        if [[ -z "$SERVICE" ]]; then
            echo "Error: service name required for resources tests"
            exit 1
        fi
        run_service_tests "resources" "$SERVICE"
        ;;
    datasources)
        if [[ -z "$SERVICE" ]]; then
            echo "Error: service name required for datasources tests"
            exit 1
        fi
        run_service_tests "datasources" "$SERVICE"
        ;;
    *)
        echo "Error: unknown test type: $TYPE"
        echo "Valid types: provider-core, resources, datasources"
        exit 1
        ;;
esac

echo "Tests completed successfully"
