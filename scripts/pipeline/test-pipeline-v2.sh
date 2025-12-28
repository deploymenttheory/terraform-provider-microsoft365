#!/bin/bash
# Test script for the URL-based API changes monitoring pipeline
# This script runs the complete pipeline with the new URL extraction method

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "============================================================"
echo "API Changes Monitoring Pipeline - v2 (URL-Based) Test"
echo "============================================================"
echo ""
echo "Script directory: $SCRIPT_DIR"
echo "Workspace directory: $WORKSPACE_DIR"
echo ""

# Check dependencies
echo "Checking dependencies..."
if ! command -v python3 &> /dev/null; then
    echo "Error: python3 is required but not installed"
    exit 1
fi

# Check if required Python packages are installed
python3 -c "import feedparser, bs4, lxml" 2>/dev/null || {
    echo "Warning: Some Python packages are missing"
    echo "Installing required packages..."
    python3 -m pip install --quiet --user -r "$SCRIPT_DIR/requirements.txt" || {
        echo "Error: Failed to install dependencies"
        echo "Please run: pip install -r $SCRIPT_DIR/requirements.txt"
        exit 1
    }
}

echo "✓ Dependencies OK"
echo ""

# Step 1: Scan provider
echo "Step 1: Scanning provider endpoints..."
echo "------------------------------------------------------------"
python3 "$SCRIPT_DIR/scan-provider-endpoints.py" \
    --output "$SCRIPT_DIR/test-provider-endpoints.json" \
    --base-path "$WORKSPACE_DIR"

if [ $? -eq 0 ]; then
    echo "✓ Provider endpoints scanned successfully"
    echo ""
else
    echo "✗ Failed to scan provider endpoints"
    exit 1
fi

# Step 2: Parse changelog with URL-based extraction
echo "Step 2: Parsing Microsoft Graph API changelog (URL-based)..."
echo "------------------------------------------------------------"
python3 "$SCRIPT_DIR/parse-graph-changelog.py" \
    --output "$SCRIPT_DIR/test-changelog-data.json" \
    --lookback-days 120 \
    --provider-resources "$SCRIPT_DIR/test-provider-endpoints.json" \
    --debug \
    --verbose

if [ $? -eq 0 ]; then
    echo "✓ Changelog parsed successfully with URL extraction"
    echo ""
else
    echo "✗ Failed to parse changelog"
    exit 1
fi

# Step 3: Compare and identify gaps (dry run - no issues created)
echo "Step 3: Comparing and identifying gaps..."
echo "------------------------------------------------------------"
python3 "$SCRIPT_DIR/compare-and-create-issues.py" \
    --changelog "$SCRIPT_DIR/test-changelog-data.json" \
    --provider "$SCRIPT_DIR/test-provider-endpoints.json" \
    --output "$SCRIPT_DIR/test-gaps-report.json" \
    --create-issues false

if [ $? -eq 0 ]; then
    echo "✓ Gaps analysis completed successfully"
    echo ""
else
    echo "✗ Failed to analyze gaps"
    exit 1
fi

# Display results
echo "============================================================"
echo "PIPELINE V2 TEST COMPLETED SUCCESSFULLY"
echo "============================================================"
echo ""
echo "Generated files:"
echo "  - test-provider-endpoints.json"
echo "  - test-changelog-data.json (URL-based extraction)"
echo "  - test-gaps-report.json"
echo ""

if [ -f "$SCRIPT_DIR/test-changelog-data.json" ]; then
    echo "Summary from changelog (URL-based method):"
    python3 -c "
import json
with open('$SCRIPT_DIR/test-changelog-data.json', 'r') as f:
    data = json.load(f)
    print(f\"  Total changes: {data['total_changes']}\")
    print(f\"  Extraction method: {data.get('extraction_method', 'unknown')}\")
    
    # Count unique resources
    all_resources = set()
    all_methods = set()
    for change in data.get('changes', []):
        all_resources.update(change.get('resources', []))
        all_methods.update(change.get('methods', []))
    
    print(f\"  Unique resources extracted: {len(all_resources)}\")
    print(f\"  Unique methods extracted: {len(all_methods)}\")
    
    # Show sample resources
    if all_resources:
        print(f\"\")
        print(f\"  Sample extracted resources:\")
        for i, resource in enumerate(sorted(list(all_resources))[:10], 1):
            print(f\"    {i}. {resource}\")
"
fi

echo ""

if [ -f "$SCRIPT_DIR/test-gaps-report.json" ]; then
    echo "Summary from gaps report:"
    python3 -c "
import json
with open('$SCRIPT_DIR/test-gaps-report.json', 'r') as f:
    report = json.load(f)
    print(f\"  Total API changes: {report.get('total_changes', 0)}\")
    print(f\"  Relevant changes: {report.get('relevant_changes', 0)}\")
    print(f\"  Already implemented: {report.get('already_implemented', 0)}\")
    print(f\"  Gaps identified: {report.get('gaps_identified', 0)}\")
    
    if report.get('gaps_identified', 0) > 0:
        print(f\"\")
        print(f\"  Top gaps:\")
        for i, gap in enumerate(report.get('gaps', [])[:5], 1):
            priority_emoji = '🔴' if gap.get('priority') == 'high' else '🟡' if gap.get('priority') == 'medium' else '🟢'
            print(f\"    {i}. {priority_emoji} [{gap.get('priority', 'unknown').upper()}] {gap.get('title', 'Unknown')[:80]}\")
"
fi

echo ""
echo "✨ New Features in v2:"
echo "  ✓ URL-based extraction (not keyword matching)"
echo "  ✓ Exact resource names from documentation URLs"
echo "  ✓ Automatic API version detection"
echo "  ✓ Documentation links captured"
echo "  ✓ Method extraction from URLs"
echo ""
echo "To create GitHub issues, run:"
echo "  export GITHUB_REPOSITORY='deploymenttheory/terraform-provider-microsoft365'"
echo "  export GITHUB_TOKEN='your-token'"
echo "  python3 compare-and-create-issues.py \\"
echo "    --changelog test-changelog-data.json \\"
echo "    --provider test-provider-endpoints.json \\"
echo "    --output test-gaps-report.json \\"
echo "    --create-issues true \\"
echo "    --check-existing"
echo ""
echo "To clean up test files, run:"
echo "  rm $SCRIPT_DIR/test-*.json"

