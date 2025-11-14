#!/bin/bash
set -euo pipefail

# Test runner for nightly acceptance tests
# Usage: ./run-tests.sh <type> [service] [coverage-file]

TYPE="${1:-}"
SERVICE="${2:-}"
COVERAGE_FILE="${3:-coverage.txt}"
TEST_OUTPUT_FILE="${4:-test-output.log}"

if [[ -z "$TYPE" ]]; then
    echo "Usage: $0 <type> [service] [coverage-file] [test-output-file]"
    echo "Types: provider-core, resources, datasources"
    exit 1
fi

# Check if tests should be skipped (set by map-credentials.sh)
if [[ "${SKIP_TESTS:-false}" == "true" ]]; then
    echo "⏭️  Skipping tests - no credentials configured"
    echo "mode: atomic" > "$COVERAGE_FILE"
    exit 0
fi

run_provider_core_tests() {
    echo "Running provider core tests..."
    
    # Run tests and capture output, but allow failures
    set +e
    go test -v -race \
        -coverprofile="$COVERAGE_FILE" \
        -covermode=atomic \
        ./internal/client/... \
        ./internal/helpers/... \
        ./internal/provider/... \
        ./internal/utilities/... \
        2>&1 | tee "$TEST_OUTPUT_FILE"
    
    TEST_EXIT_CODE=$?
    set -e
    
    # Parse test failures and create JSON report
    parse_test_failures "$TEST_OUTPUT_FILE" "provider-core" ""
    
    return $TEST_EXIT_CODE
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

    # Check for test files recursively
    local test_count
    test_count=$(find "${test_dir}" -name "*_test.go" -type f | wc -l | tr -d ' ')

    if [[ "$test_count" -eq 0 ]]; then
        echo "No test files found in ${test_dir}, creating empty coverage file"
        echo "mode: atomic" > "$COVERAGE_FILE"
        return 0
    fi

    echo "Found ${test_count} test files"

    # Run tests and capture output, but allow failures
    set +e
    go test -v -race \
        -coverprofile="$COVERAGE_FILE" \
        -covermode=atomic \
        "${test_dir}/..." \
        2>&1 | tee "$TEST_OUTPUT_FILE"
    
    TEST_EXIT_CODE=$?
    set -e
    
    # Parse test failures and create JSON report
    parse_test_failures "$TEST_OUTPUT_FILE" "$category" "$service"
    
    return $TEST_EXIT_CODE
}

parse_test_failures() {
    local output_file="$1"
    local category="$2"
    local service="$3"
    
    local json_file="test-failures.json"
    
    # Initialize JSON array
    echo "[" > "$json_file"
    
    local first_entry=true
    
    # Parse test output for FAIL lines
    # Matches patterns like: "--- FAIL: TestName (0.01s)"
    while IFS= read -r line; do
        if echo "$line" | grep -q "^--- FAIL:"; then
            # Extract test name
            test_name=$(echo "$line" | sed 's/^--- FAIL: \([^ ]*\).*/\1/')
            
            # Get context (next few lines for error details)
            context=$(grep -A 10 "^--- FAIL: ${test_name}" "$output_file" | head -n 11 | tail -n 10 | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')
            
            # Add comma if not first entry
            if [ "$first_entry" = false ]; then
                echo "," >> "$json_file"
            fi
            first_entry=false
            
            # Add JSON object
            cat >> "$json_file" <<EOF
{
  "test_name": "${test_name}",
  "category": "${category}",
  "service": "${service}",
  "context": "${context}"
}
EOF
        fi
    done < "$output_file"
    
    echo "" >> "$json_file"
    echo "]" >> "$json_file"
    
    echo "✅ Test failure report created: $json_file"
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

echo "Tests completed"
