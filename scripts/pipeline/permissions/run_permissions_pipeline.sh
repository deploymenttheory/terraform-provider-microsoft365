#!/usr/bin/env bash
# run_permissions_pipeline.sh
#
# Runs the full Microsoft Graph API permissions update pipeline in order:
#
#   1. update_graph_service_permissions.py
#      Derives required Graph API permissions from endpoint metadata and
#      updates ReadPermissions / WritePermissions in every Go service
#      constructor (New*() function) under internal/services/.
#
#   2. update_graph_service_permission_doc_templates.py
#      Reads the (now-updated) Go constructor permissions and synchronises
#      the **Required:** section in each matching Terraform documentation
#      template under templates/.
#
# Usage:
#   ./run_permissions_pipeline.sh [--apply] [--refresh-metadata]
#                                 [--domain DOMAIN] [--category CATEGORY]
#
# Options:
#   --apply             Write changes to disk (default: dry-run preview).
#   --refresh-metadata  Force re-download of Graph command metadata cache.
#   --domain DOMAIN     Limit both steps to one service domain.
#   --category CATEGORY Limit both steps to one service category.
#
# Examples:
#   ./run_permissions_pipeline.sh                      # dry-run all
#   ./run_permissions_pipeline.sh --apply              # apply all changes
#   ./run_permissions_pipeline.sh --apply --refresh-metadata
#   ./run_permissions_pipeline.sh --domain applications --apply
#   ./run_permissions_pipeline.sh --category resources --apply

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ---------------------------------------------------------------------------
# Argument parsing — forward recognised flags to both Python scripts
# ---------------------------------------------------------------------------
APPLY_FLAG=""
REFRESH_FLAG=""
DOMAIN_FLAG=""
CATEGORY_FLAG=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --apply)
            APPLY_FLAG="--apply"
            shift
            ;;
        --refresh-metadata)
            REFRESH_FLAG="--refresh-metadata"
            shift
            ;;
        --domain)
            DOMAIN_FLAG="--domain $2"
            shift 2
            ;;
        --category)
            CATEGORY_FLAG="--category $2"
            shift 2
            ;;
        -h|--help)
            sed -n '/^# Usage:/,/^[^#]/p' "$0" | sed 's/^# \?//'
            exit 0
            ;;
        *)
            echo "Unknown option: $1" >&2
            exit 1
            ;;
    esac
done

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
banner() {
    echo ""
    echo "============================================================"
    echo "  $1"
    echo "============================================================"
    echo ""
}

run_step() {
    local step_num="$1"
    local step_name="$2"
    shift 2
    banner "Step ${step_num}: ${step_name}"
    python3 "$@"
    echo ""
}

# ---------------------------------------------------------------------------
# Pipeline
# ---------------------------------------------------------------------------
banner "Graph API Permissions Pipeline"
echo "  Script dir : ${SCRIPT_DIR}"
echo "  Apply      : ${APPLY_FLAG:-'(dry-run)'}"
echo "  Domain     : ${DOMAIN_FLAG:-'(all)'}"
echo "  Category   : ${CATEGORY_FLAG:-'(all)'}"
echo "  Refresh    : ${REFRESH_FLAG:-'(cached metadata)'}"

# Step 1 — update Go constructor permissions
# shellcheck disable=SC2086
run_step 1 "Update Go service constructor permissions" \
    "${SCRIPT_DIR}/update_graph_service_permissions.py" \
    ${APPLY_FLAG} \
    ${REFRESH_FLAG} \
    ${DOMAIN_FLAG} \
    ${CATEGORY_FLAG}

# Step 2 — sync Terraform documentation templates
# shellcheck disable=SC2086
run_step 2 "Sync Terraform documentation template permissions" \
    "${SCRIPT_DIR}/update_graph_service_permission_doc_templates.py" \
    ${APPLY_FLAG} \
    ${DOMAIN_FLAG} \
    ${CATEGORY_FLAG}

banner "Pipeline complete"
if [[ -z "${APPLY_FLAG}" ]]; then
    echo "  This was a dry-run. Re-run with --apply to write changes."
fi
