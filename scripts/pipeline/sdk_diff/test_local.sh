#!/bin/bash
# Local test script for SDK diff analysis
set -e

echo "========================================"
echo "SDK Diff Analysis - Local Test"
echo "========================================"
echo ""

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$REPO_ROOT"

echo "📂 Repository root: $REPO_ROOT"
echo ""

# Step 1: Detect versions
echo "Step 1: Detecting SDK versions..."
python3 scripts/pipeline/sdk_diff/steps/detect_versions.py --repo-path .
echo ""

# Step 2: Map usage
echo "Step 2: Mapping SDK usage (this may take a minute)..."
python3 scripts/pipeline/sdk_diff/steps/map_usage.py \
  --repo-path . \
  --usage-output /tmp/sdk_usage_test.json
echo ""

echo "✅ Test complete!"
echo ""
echo "To run full analysis (requires GITHUB_TOKEN):"
echo "  1. Set token: export GITHUB_TOKEN='your_token'"
echo "  2. Run: python3 scripts/pipeline/sdk_diff/steps/analyze_changes.py \\"
echo "          --usage-file /tmp/sdk_usage_test.json \\"
echo "          --sdk msgraph-beta-sdk-go \\"
echo "          --current-version v0.157.0 \\"
echo "          --latest-version v0.158.0 \\"
echo "          --changes-output /tmp/sdk_changes_test.json"
