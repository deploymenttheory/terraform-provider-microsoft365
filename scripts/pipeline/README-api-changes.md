# Microsoft Graph API Changes Monitoring Pipeline

This directory contains scripts and workflows for automatically monitoring Microsoft Graph API changelog and identifying gaps in the Terraform provider implementation.

## Overview

The API changes monitoring system consists of three main components:

1. **RSS Feed Parser** (`parse-graph-changelog.py`) - Fetches and parses the Microsoft Graph API changelog RSS feed
2. **Provider Scanner** (`scan-provider-endpoints.py`) - Scans the Terraform provider codebase to identify implemented endpoints
3. **Gap Analyzer** (`compare-and-create-issues.py`) - Compares API changes with implementation and creates GitHub issues for gaps

## Architecture

```
┌─────────────────────────────────┐
│  Microsoft Graph API Changelog  │
│         (RSS Feed)              │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│   parse-graph-changelog.py      │
│  - Fetch RSS feed               │
│  - Extract API changes          │
│  - Filter relevant changes      │
└────────────┬────────────────────┘
             │
             ▼
       changelog-data.json
             │
             ├──────────────────────┐
             │                      │
             ▼                      ▼
┌─────────────────────┐   ┌─────────────────────┐
│ scan-provider-      │   │ compare-and-        │
│ endpoints.py        │   │ create-issues.py    │
│ - Scan Go code      │   │ - Compare data      │
│ - Extract endpoints │   │ - Identify gaps     │
└──────┬──────────────┘   │ - Create issues     │
       │                  └──────┬──────────────┘
       ▼                         │
provider-endpoints.json          ▼
                          gaps-report.json
                                 │
                                 ▼
                          GitHub Issues
```

## GitHub Actions Workflow

The `api-changes.yml` workflow runs automatically:

- **Schedule**: Every Monday at 9 AM UTC
- **Manual Trigger**: Can be run manually with custom parameters

### Workflow Inputs

When manually triggering the workflow, you can specify:

- `lookback_days` - Number of days to look back for changes (default: 30)
- `create_issues` - Whether to create GitHub issues (default: true)

## Scripts

### 1. parse-graph-changelog.py

Parses the Microsoft Graph API changelog RSS feed and extracts relevant changes.

**Usage:**
```bash
python parse-graph-changelog.py \
  --output changelog-data.json \
  --lookback-days 30
```

**Options:**
- `--output` - Output JSON file path (default: changelog-data.json)
- `--lookback-days` - Number of days to look back (default: 30)
- `--url` - RSS feed URL (default: Microsoft Graph changelog)

**Features:**
- Filters changes relevant to device management and Intune
- Extracts resources, methods, endpoints, and properties
- Identifies changes that support CRUD or minimal operations
- Supports both v1.0 and beta API versions

**Output Format:**
```json
{
  "generated_at": "2025-12-28T...",
  "lookback_days": 30,
  "total_changes": 150,
  "changes": [
    {
      "guid": "...",
      "title": "Device and app management",
      "description": "...",
      "pub_date": "2025-12-16T...",
      "categories": ["Device and app management"],
      "api_version": "v1.0",
      "resources": ["deviceManagementScript"],
      "methods": ["create", "update"],
      "endpoints": ["/deviceManagement/scripts"],
      "properties": ["displayName"],
      "change_type": "added",
      "is_relevant": true,
      "supports_crud_or_minimal": true
    }
  ]
}
```

### 2. scan-provider-endpoints.py

Scans the Terraform provider codebase to identify implemented Graph API endpoints.

**Usage:**
```bash
python scan-provider-endpoints.py \
  --output provider-endpoints.json \
  --base-path .
```

**Options:**
- `--output` - Output JSON file path (default: provider-endpoints.json)
- `--base-path` - Base path of provider codebase (default: current directory)

**Features:**
- Scans Go resource files
- Extracts Graph API endpoint calls
- Identifies CRUD operations
- Creates lookup structures for efficient comparison

**Output Format:**
```json
{
  "generated_at": "2025-12-28T...",
  "total_resources": 50,
  "total_endpoints": 120,
  "total_graph_resources": 80,
  "resources": [
    {
      "resource_name": "device_management_windows_remediation_script",
      "file_path": "internal/services/.../resource.go",
      "graph_resources": ["deviceManagement", "deviceHealthScripts"],
      "graph_methods": ["Get", "Post", "Patch"],
      "endpoints": ["/deviceManagement/deviceHealthScripts"],
      "operations": ["create", "read", "update", "delete"]
    }
  ],
  "lookup": {
    "resources": {...},
    "endpoints": {...},
    "operations": {...}
  }
}
```

### 3. compare-and-create-issues.py

Compares API changes with provider implementation and creates GitHub issues for gaps.

**Usage:**
```bash
python compare-and-create-issues.py \
  --changelog changelog-data.json \
  --provider provider-endpoints.json \
  --output gaps-report.json \
  --create-issues true
```

**Options:**
- `--changelog` - Path to changelog JSON file (required)
- `--provider` - Path to provider endpoints JSON file (required)
- `--output` - Output JSON file path (default: gaps-report.json)
- `--create-issues` - Whether to create GitHub issues (default: false)
- `--check-existing` - Check for existing issues before creating new ones

**Environment Variables:**
- `GITHUB_REPOSITORY` - GitHub repository (format: owner/repo)
- `GITHUB_TOKEN` - GitHub token for API access

**Features:**
- Identifies three types of gaps:
  1. **New resources** - API resources not yet implemented
  2. **Updated resources** - Existing resources with new functionality
  3. **Missing operations** - Resources missing CRUD operations
- Prioritizes gaps (high, medium, low)
- Creates detailed GitHub issues with implementation checklists
- Avoids duplicate issues

**Output Format:**
```json
{
  "generated_at": "2025-12-28T...",
  "total_changes": 150,
  "relevant_changes": 80,
  "already_implemented": 60,
  "gaps_identified": 20,
  "issues_created": 15,
  "gaps": [
    {
      "title": "[v1.0] Implement new Graph API resource: ...",
      "body": "## Summary\n...",
      "labels": ["enhancement", "api-change", "priority-high"],
      "priority": "high",
      "gap_type": "new_resource",
      "api_version": "v1.0"
    }
  ],
  "created_issues": [
    {
      "url": "https://github.com/...",
      "title": "...",
      "priority": "high"
    }
  ]
}
```

## Gap Types and Priorities

### Gap Types

1. **new_resource**
   - Completely new Graph API resource
   - Not implemented in provider
   - Supports CRUD or minimal operations

2. **updated_resource**
   - Existing resource with new methods/properties
   - Partially implemented in provider
   - Needs updates to match API

3. **missing_operation**
   - Resource exists but missing CRUD operations
   - Needs additional implementation

### Priority Calculation

- **High Priority**:
  - New resources in v1.0 API with full CRUD support
  - Deprecated resources still in use
  
- **Medium Priority**:
  - Updates to existing resources
  - New resources in beta API with CRUD support
  
- **Low Priority**:
  - Minor changes
  - Properties without operational impact

## Issue Labels

Created issues are automatically labeled:

- `enhancement` - Feature request
- `api-change` - API change related
- `priority-high|medium|low` - Priority level
- `graph-v1.0|beta` - API version
- `device-management|identity|security` - Category
- `new-resource|updated-resource|missing-operation` - Gap type

## Manual Execution

You can run the pipeline locally:

### Prerequisites

```bash
# Install Python dependencies
pip install feedparser beautifulsoup4 lxml requests

# Install GitHub CLI (for issue creation)
brew install gh  # macOS
# or
apt-get install gh  # Linux

# Authenticate with GitHub
gh auth login
```

### Run Pipeline

```bash
# Step 1: Parse changelog
python parse-graph-changelog.py \
  --output changelog-data.json \
  --lookback-days 30

# Step 2: Scan provider
python scan-provider-endpoints.py \
  --output provider-endpoints.json

# Step 3: Compare and create issues (dry run)
python compare-and-create-issues.py \
  --changelog changelog-data.json \
  --provider provider-endpoints.json \
  --output gaps-report.json \
  --create-issues false

# Step 4: Create issues (if satisfied with gaps)
export GITHUB_REPOSITORY="deploymenttheory/terraform-provider-microsoft365"
export GITHUB_TOKEN="your-token"

python compare-and-create-issues.py \
  --changelog changelog-data.json \
  --provider provider-endpoints.json \
  --output gaps-report.json \
  --create-issues true \
  --check-existing
```

## Customization

### Adding New Categories

Edit `RELEVANT_CATEGORIES` in `parse-graph-changelog.py`:

```python
RELEVANT_CATEGORIES = [
    "Device and app management",
    "Identity and access",
    # Add your categories here
]
```

### Adjusting Filters

Modify the relevance check in `GraphAPIChange.is_relevant()`:

```python
relevant_keywords = [
    'device', 'intune', 'management',
    # Add your keywords here
]
```

### Customizing Issue Templates

Edit the `_generate_body()` method in the `Gap` class in `compare-and-create-issues.py`.

## Troubleshooting

### No Changes Detected

- Check the lookback period (`--lookback-days`)
- Verify RSS feed is accessible
- Review `RELEVANT_CATEGORIES` filter

### Issues Not Created

- Verify `GITHUB_TOKEN` has `issues: write` permission
- Check `GITHUB_REPOSITORY` format (owner/repo)
- Ensure GitHub CLI is installed and authenticated

### False Positives

- Review gap detection logic in `compare_changes_with_provider()`
- Adjust keyword matching in endpoint comparison
- Use `--check-existing` to avoid duplicates

## Artifacts

The workflow produces three JSON artifacts:

1. `changelog-data.json` - Parsed API changes
2. `provider-endpoints.json` - Provider implementation map
3. `gaps-report.json` - Identified gaps and created issues

These are available as workflow artifacts for 30 days.

## Future Enhancements

Potential improvements:

- [ ] Natural language processing for better change detection
- [ ] Integration with provider roadmap
- [ ] Automatic PR creation for simple changes
- [ ] Slack/Discord notifications for new gaps
- [ ] Historical gap tracking and metrics
- [ ] Priority scoring based on community requests

## Contributing

When modifying the pipeline:

1. Test locally before committing
2. Update this README with changes
3. Consider backward compatibility
4. Add logging for debugging

## Support

For issues or questions:

- Create a GitHub issue with label `pipeline`
- Tag `@maintainers` in discussions
- Review workflow logs in Actions tab

