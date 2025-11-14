#!/bin/bash
set -euo pipefail

# Creates individual GitHub issues for each failing test
# Usage: ./create-test-issues.sh <owner> <repo> <run-id> <failures-json>

OWNER="${1:-}"
REPO="${2:-}"
RUN_ID="${3:-}"
FAILURES_JSON="${4:-test-failures.json}"

if [[ -z "$OWNER" || -z "$REPO" || -z "$RUN_ID" ]]; then
    echo "Usage: $0 <owner> <repo> <run-id> <failures-json>"
    echo "Example: $0 deploymenttheory terraform-provider-microsoft365 123456 test-failures.json"
    exit 1
fi

if [[ ! -f "$FAILURES_JSON" ]]; then
    echo "Error: Failures JSON file not found: $FAILURES_JSON"
    exit 1
fi

DATE=$(date -u +"%Y-%m-%d")
WORKFLOW_URL="https://github.com/${OWNER}/${REPO}/actions/runs/${RUN_ID}"

# Check if JSON is empty or has no failures
FAILURE_COUNT=$(jq 'length' "$FAILURES_JSON")

if [[ "$FAILURE_COUNT" -eq 0 ]]; then
    echo "✅ No test failures found"
    exit 0
fi

echo "Found ${FAILURE_COUNT} failing test(s)"

# Process each test failure
for i in $(seq 0 $((FAILURE_COUNT - 1))); do
    TEST_NAME=$(jq -r ".[$i].test_name" "$FAILURES_JSON")
    CATEGORY=$(jq -r ".[$i].category" "$FAILURES_JSON")
    SERVICE=$(jq -r ".[$i].service" "$FAILURES_JSON")
    CONTEXT=$(jq -r ".[$i].context" "$FAILURES_JSON")
    
    # Build service path for better context
    SERVICE_PATH="${CATEGORY}"
    if [[ -n "$SERVICE" && "$SERVICE" != "null" ]]; then
        SERVICE_PATH="${CATEGORY}/${SERVICE}"
    fi
    
    echo ""
    echo "Processing failure: ${TEST_NAME}"
    
    # Create issue title - standardized format for de-duplication
    ISSUE_TITLE="🔴 Nightly Test Failure: ${TEST_NAME}"
    
    # Check if issue already exists
    echo "Checking for existing issue..."
    EXISTING_ISSUE=$(gh issue list \
        --repo "${OWNER}/${REPO}" \
        --state open \
        --label "nightly-test-failure" \
        --search "in:title \"${TEST_NAME}\"" \
        --json number,title \
        --jq ".[0].number" || echo "")
    
    if [[ -n "$EXISTING_ISSUE" && "$EXISTING_ISSUE" != "null" ]]; then
        echo "⚠️  Issue already exists: #${EXISTING_ISSUE}"
        echo "   Updating issue with latest failure info..."
        
        # Add comment to existing issue with latest failure
        COMMENT_BODY="## Latest Failure: ${DATE}

**Workflow Run:** [${RUN_ID}](${WORKFLOW_URL})
**Service Path:** \`${SERVICE_PATH}\`

### Failure Context

\`\`\`
${CONTEXT}
\`\`\`

---
🤖 Updated by nightly test failure handler"

        gh issue comment "$EXISTING_ISSUE" \
            --repo "${OWNER}/${REPO}" \
            --body "${COMMENT_BODY}"
        
        echo "✅ Updated issue #${EXISTING_ISSUE}"
        continue
    fi
    
    # Create new issue
    echo "Creating new issue..."
    
    ISSUE_BODY="## Test Failure Details

**Test Name:** \`${TEST_NAME}\`
**Service Path:** \`${SERVICE_PATH}\`
**First Detected:** ${DATE}
**Workflow Run:** [${RUN_ID}](${WORKFLOW_URL})

### Failure Context

\`\`\`
${CONTEXT}
\`\`\`

### Action Items

- [ ] Review the [workflow logs](${WORKFLOW_URL})
- [ ] Identify root cause of test failure
- [ ] Fix the failing test
- [ ] Verify fix in next nightly run
- [ ] Close this issue when resolved

### Links

- [View Full Workflow Run](${WORKFLOW_URL})
- [View Test Coverage](https://codecov.io/gh/${OWNER}/${REPO})

---
🤖 Automatically created by nightly test failure handler

**Note:** This issue will be automatically updated if the test continues to fail in subsequent nightly runs."

    ISSUE_URL=$(gh issue create \
        --repo "${OWNER}/${REPO}" \
        --title "${ISSUE_TITLE}" \
        --body "${ISSUE_BODY}" \
        --label "bug,nightly-test-failure,automated")
    
    echo "✅ Created issue: ${ISSUE_URL}"
done

echo ""
echo "✅ Issue creation complete: ${FAILURE_COUNT} test failure(s) processed"

