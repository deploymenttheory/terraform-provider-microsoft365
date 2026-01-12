# Microsoft Graph OpenAPI Schema Detection

Detects schema changes in the Microsoft Graph Beta API by monitoring the OpenAPI specification from [microsoftgraph/msgraph-metadata](https://github.com/microsoftgraph/msgraph-metadata).

## Overview

This tool provides **early warning** of API changes by directly monitoring the OpenAPI spec, which is updated **before** SDK releases. It complements the SDK detection by providing an additional validation layer.

## Quick Start

```bash
cd scripts/pipeline/graph_openapi_schema_detection

# Test with two specific commits from msgraph-metadata repo
./detect_openapi_changes.py \
  --old-commit abc123 \
  --new-commit def456 \
  --dry-run

# Filter to only provider-used models
./detect_openapi_changes.py \
  --old-commit abc123 \
  --new-commit def456 \
  --filter-by-usage ../sdk_schema_detection/provider_model_usage.json \
  --dry-run

# Save results for analysis
./detect_openapi_changes.py \
  --old-commit abc123 \
  --new-commit def456 \
  --save-results openapi_changes.json \
  --dry-run
```

## Prerequisites

```bash
# Install PyYAML (for parsing OpenAPI spec)
pip3 install pyyaml --break-system-packages  # macOS
# or
pip3 install pyyaml  # Linux/CI
```

## Architecture

Mirrors the proven SDK detection architecture:

```
graph_openapi_schema_detection/
├── core/                           # Core components
│   ├── progress_reporter.py       # Output formatting (reused from SDK detection)
│   ├── spec_fetcher.py            # Download OpenAPI specs
│   ├── version_detector.py        # Extract version from spec
│   ├── schema_parser.py           # Parse YAML efficiently
│   ├── schema_comparer.py         # Detect property changes
│   └── issue_builder.py           # Generate issue content
│
├── models/                         # Data models
│   ├── results.py                 # DetectionResult, VersionResult
│   ├── changes.py                 # PropertyChange, SchemaChange
│   └── statistics.py              # Parsing statistics
│
└── detect_openapi_changes.py      # Main entry point
```

## What It Detects

### Property Changes
- ✅ **Added properties** → Add to Terraform schema
- ⚠️ **Removed properties** → Deprecate (BREAKING)
- ⚠️ **Type changes** → Update schema type (BREAKING)
- ⚠️ **Required changes** → Update validation (BREAKING if newly required)
- ℹ️  **Nullable changes** → Update schema configuration

### Breaking Changes
Automatically detected:
- Properties removed from schema
- Property type changed (e.g., string → integer)
- Property changed from optional → required

## Usage Examples

### Compare Two Commits

```bash
# Get commit SHAs from msgraph-metadata repo
# https://github.com/microsoftgraph/msgraph-metadata/commits/master/openapi/beta/openapi.yaml

./detect_openapi_changes.py \
  --old-commit <older_commit_sha> \
  --new-commit <newer_commit_sha> \
  --dry-run
```

### Filter by Provider Usage

Uses the same `provider_model_usage.json` from SDK detection:

```bash
# Step 1: Analyze provider usage (reuse from SDK detection)
cd ../sdk_schema_detection
./analyze_provider_model_usage.py \
  --provider-path ../../../ \
  --output provider_model_usage.json

# Step 2: Detect OpenAPI changes with filtering
cd ../graph_openapi_schema_detection
./detect_openapi_changes.py \
  --old-commit abc123 \
  --new-commit def456 \
  --filter-by-usage ../sdk_schema_detection/provider_model_usage.json \
  --save-results openapi_changes.json \
  --dry-run
```

## Output

### Console Output
```
🚀 OpenAPI Schema Change Detection

📥 Downloading latest OpenAPI spec from master...
   Downloaded 59.87 MB
📄 Extracting schemas section...
   Found schemas at line 1,201,913
   Schemas end at line 1,431,917
   Extracted 230,004 lines
🔍 Parsing schemas YAML...
   Parsed 9,520 schemas in 6.40s

🔬 Comparing schemas...
   Found 15 schema(s) with changes

  📊 Comparison Statistics:
     Schemas compared: 9,520
     Schemas with changes: 15
     Schemas added: 5
     Schemas removed: 2

  🔧 Property Changes:
     Properties added: 45
     Properties removed: 3
     Type changes: 2
     Required changes: 5

🔍 Filtering changes by provider usage...
   Loaded 329 model(s) from provider usage data
   ✓ Kept 7 relevant model(s)
   ✗ Filtered 8 unused model(s)

✅ Analysis complete!
```

### JSON Output (`openapi_changes.json`)
```json
{
  "spec_version": "beta",
  "previous_version": "beta",
  "detection_timestamp": "2026-01-12T15:30:00",
  "total_schemas_changed": 7,
  "schemas_with_changes": 15,
  "filtered_schemas": 8,
  "breaking_changes_count": 2,
  "schema_changes": [
    {
      "schema_name": "microsoft.graph.user",
      "model_name": "User",
      "file_path": "models/user.go",
      "has_breaking_changes": false,
      "change_summary": "+3 properties",
      "added_properties": [
        {"name": "newField", "type": "string", "required": false}
      ]
    }
  ]
}
```

## Integration

### With SDK Detection Pipeline

The OpenAPI detection can run **alongside** SDK detection:

1. **Weekly**: Monitor OpenAPI spec for changes
2. **On PR**: SDK detection validates changes made it to SDK
3. **Cross-validation**: Ensure API → SDK → Provider alignment

### GitHub Actions (Future)

```yaml
name: OpenAPI Schema Monitor

on:
  schedule:
    - cron: '0 2 * * 1'  # Monday 2 AM UTC
  workflow_dispatch:

jobs:
  monitor-openapi:
    runs-on: ubuntu-latest
    steps:
      - name: Detect OpenAPI changes
        run: |
          cd scripts/pipeline/graph_openapi_schema_detection
          ./detect_openapi_changes.py \
            --old-commit $OLD_COMMIT \
            --new-commit $NEW_COMMIT \
            --filter-by-usage ../sdk_schema_detection/provider_model_usage.json
```

## Performance

- **Download time**: ~3 seconds (compressed)
- **Parse time**: ~6 seconds (YAML parsing)
- **Compare time**: ~1 second
- **Total**: ~10-15 seconds

Efficient enough for CI/CD pipelines.

## Advantages Over SDK Detection

| Aspect | SDK Detection | OpenAPI Detection |
|--------|---------------|-------------------|
| **Timing** | After SDK release | Before SDK release ⚡ |
| **Coverage** | SDK-exposed types | Full API surface |
| **Parsing** | Complex (Go code) | Simpler (YAML) |
| **Source** | Generated code | API spec (source of truth) |
| **Breaking Changes** | Inferred | Explicit |

## Limitations

- Requires commit SHAs (caching not yet implemented)
- GitHub issue creation not yet wired up (TODO)
- Weekly monitoring script not yet created (TODO)

## Future Enhancements

- [ ] Automatic caching of previous spec version
- [ ] GitHub issue creation via CLI
- [ ] Weekly monitoring script
- [ ] Cross-validation with SDK detection
- [ ] Trend analysis dashboard

## Testing

See commits in [msgraph-metadata](https://github.com/microsoftgraph/msgraph-metadata/commits/master/openapi/beta/openapi.yaml) for real commit SHAs to test with.

```bash
# Example: Test with real commits
./detect_openapi_changes.py \
  --old-commit <sha_from_last_week> \
  --new-commit <sha_from_this_week> \
  --filter-by-usage ../sdk_schema_detection/provider_model_usage.json \
  --save-results test_results.json \
  --dry-run
```

## License

Part of the terraform-provider-microsoft365 project.
