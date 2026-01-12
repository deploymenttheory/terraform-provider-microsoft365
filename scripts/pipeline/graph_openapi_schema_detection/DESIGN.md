# Graph OpenAPI Schema Detection - Design

## Overview

Detect schema changes in the Microsoft Graph Beta API by monitoring the OpenAPI specification. This provides **early warning** of API changes before SDK updates are released.

## Current SDK Detection Flow (Reference)

```
1. analyze_provider_model_usage.py
   ↓ provider_model_usage.json
   
2. detect_schema_changes.py
   ├─ Parse go.mod or PR to get SDK version change (v0.156.0 → v0.157.0)
   ├─ Fetch diff between SDK versions (Go code)
   ├─ Parse Go structs/interfaces/embedded types
   ├─ Filter by provider usage
   ├─ Save results (schema_changes.json)
   └─ Create GitHub issues

Modules:
- core/: progress_reporter, version_parser, github_client, diff_fetcher, struct_parser, issue_builder
- models/: results, changes, statistics
```

## Proposed OpenAPI Detection Flow

```
Weekly Monitor:
1. Monitor msgraph-metadata repo for OpenAPI spec updates
2. Detect version changes in openapi.yaml
3. Compare schemas between versions
4. Filter to provider-used models
5. Generate change reports + GitHub issues

Architecture (mirrors SDK detection):

graph_openapi_schema_detection/
├── __init__.py
├── detect_openapi_changes.py          # Main entry point
├── monitor_spec_updates.py             # Weekly check for spec updates
│
├── core/                               # Core components
│   ├── __init__.py
│   ├── progress_reporter.py           # ✅ Reuse from sdk_schema_detection
│   ├── version_detector.py            # Detect OpenAPI version from spec
│   ├── github_client.py               # ✅ Reuse from sdk_schema_detection
│   ├── spec_fetcher.py                # Fetch OpenAPI specs (current/previous)
│   ├── schema_parser.py               # Parse YAML schemas efficiently
│   ├── schema_comparer.py             # Compare two versions of schemas
│   └── issue_builder.py               # Build GitHub issue content
│
├── models/                             # Data models
│   ├── __init__.py
│   ├── results.py                     # Result types (version, detection, etc)
│   ├── changes.py                     # PropertyChange, SchemaChange
│   └── statistics.py                  # Parsing/comparison statistics
│
└── regex_patterns.py                  # Regex for version detection
```

## Key Differences from SDK Detection

| Aspect | SDK Detection | OpenAPI Detection |
|--------|---------------|-------------------|
| **Trigger** | PR changes go.mod | Weekly schedule + manual |
| **Source** | SDK diff (Go code) | OpenAPI spec diff (YAML) |
| **Version Detection** | Parse go.mod / PR diff | Parse openapi.yaml version |
| **Fetching** | GitHub compare API (Go files) | Raw file download (YAML) |
| **Parsing** | Go structs/interfaces | YAML schemas |
| **Changes** | Fields, methods, embedded types | Properties, types, required, enums |
| **Complexity** | High (code parsing) | Lower (structured data) |

## Data Models

### PropertyChange
```python
@dataclass
class PropertyChange:
    property_name: str
    change_type: str  # 'added', 'removed', 'type_changed', 'required_changed'
    old_type: Optional[str] = None
    new_type: Optional[str] = None
    old_required: bool = False
    new_required: bool = False
    old_nullable: Optional[bool] = None
    new_nullable: Optional[bool] = None
```

### SchemaChange
```python
@dataclass
class SchemaChange:
    schema_name: str  # e.g., 'microsoft.graph.user'
    added_properties: List[PropertyChange]
    removed_properties: List[PropertyChange]
    type_changed_properties: List[PropertyChange]
    required_changed_properties: List[PropertyChange]
    
    @property
    def has_breaking_changes(self) -> bool:
        return bool(
            self.removed_properties or
            self.type_changed_properties or
            [p for p in self.required_changed_properties if p.new_required]
        )
```

### DetectionResult
```python
@dataclass
class DetectionResult:
    spec_version: str
    previous_version: str
    detection_timestamp: str
    total_schemas_changed: int
    schemas_with_changes: int
    filtered_schemas: int
    breaking_changes_count: int
    schema_changes: List[SchemaChange]
    statistics: ParseStatistics
```

## Core Components

### 1. spec_fetcher.py
```python
class SpecFetcher:
    """Fetches OpenAPI specs from msgraph-metadata repo."""
    
    def fetch_latest_spec(self) -> str:
        """Download current openapi.yaml from master branch."""
        
    def fetch_spec_at_commit(self, commit_sha: str) -> str:
        """Download openapi.yaml from specific commit."""
        
    def fetch_spec_history(self, days: int = 7) -> List[Tuple[str, str, datetime]]:
        """Get list of spec versions from commit history."""
        # Returns: [(commit_sha, commit_message, commit_date), ...]
```

### 2. version_detector.py
```python
class VersionDetector:
    """Detects version information from OpenAPI spec."""
    
    def extract_version_from_spec(self, spec_content: str) -> str:
        """Extract version from 'info.version' in OpenAPI spec."""
        
    def detect_version_change(self) -> VersionChangeResult:
        """Compare current spec version vs cached previous version."""
```

### 3. schema_parser.py
```python
class SchemaParser:
    """Parses OpenAPI schemas efficiently."""
    
    def extract_schemas_section(self, spec_content: str) -> str:
        """Extract just 'components.schemas' section (~10MB vs 60MB)."""
        
    def parse_schemas(self, schemas_content: str) -> Dict[str, Any]:
        """Parse YAML schemas section."""
        
    def extract_model_properties(self, schema_name: str, schema: Dict) -> Dict[str, Any]:
        """Extract properties with types, required status, nullable."""
```

### 4. schema_comparer.py
```python
class SchemaComparer:
    """Compares two versions of schemas."""
    
    def compare_schemas(
        self,
        old_schemas: Dict[str, Any],
        new_schemas: Dict[str, Any]
    ) -> List[SchemaChange]:
        """Compare all schemas and detect changes."""
        
    def compare_model(
        self,
        model_name: str,
        old_schema: Dict,
        new_schema: Dict
    ) -> SchemaChange:
        """Compare single model between versions."""
        
    def detect_property_changes(
        self,
        old_props: Dict,
        new_props: Dict,
        old_required: List[str],
        new_required: List[str]
    ) -> Tuple[List[PropertyChange], ...]:
        """Detect added, removed, type changes, required changes."""
```

### 5. issue_builder.py
```python
class IssueBuilder:
    """Builds GitHub issue content from OpenAPI changes."""
    
    def build_title(self, spec_version: str) -> str:
        """Build issue title."""
        # "API Schema Update Required: Microsoft Graph OpenAPI [version]"
        
    def build_body(
        self,
        spec_version: str,
        schema_changes: List[SchemaChange],
        breaking: bool
    ) -> str:
        """Build issue body with property changes."""
```

## Workflow Integration

### Option 1: Weekly Monitor (Recommended)
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
      - name: Check for OpenAPI updates
        run: |
          cd scripts/pipeline/graph_openapi_schema_detection
          ./monitor_spec_updates.py --check-weekly
          
      - name: Analyze if changed
        if: spec_changed
        run: |
          # Step 1: Use existing provider usage analysis
          cd ../sdk_schema_detection
          ./analyze_provider_model_usage.py --output usage.json
          
          # Step 2: Detect OpenAPI changes
          cd ../graph_openapi_schema_detection
          ./detect_openapi_changes.py \
            --filter-by-usage ../sdk_schema_detection/usage.json \
            --save-results openapi_changes.json
```

### Option 2: On Demand
```bash
# Manually check for changes
./detect_openapi_changes.py --current-spec [sha] --new-spec [sha]

# Compare with specific commits
./monitor_spec_updates.py --compare [old_sha] [new_sha]
```

## Advantages

1. **Early Detection**: Catch API changes before SDK release
2. **Simpler Parsing**: Structured YAML vs Go code
3. **Comprehensive**: Full API surface, not just SDK exposure
4. **Validation**: Cross-check with SDK detection
5. **Breaking Changes**: Explicit detection of breaking changes

## Implementation Phases

### Phase 1: Core Detection (MVP)
- ✅ spec_fetcher.py
- ✅ schema_parser.py
- ✅ schema_comparer.py
- ✅ Reuse progress_reporter from SDK detection
- ✅ Basic detect_openapi_changes.py

### Phase 2: Integration
- ✅ Reuse provider usage filter
- ✅ issue_builder.py (OpenAPI-specific)
- ✅ GitHub Actions weekly monitor

### Phase 3: Enhancement
- 📊 Dashboard showing OpenAPI vs SDK detection comparison
- 🔔 Slack/Teams notifications
- 📈 Trend analysis (API growth rate)

## Files to Create

1. ✅ `graph_openapi_schema_detection/__init__.py`
2. ✅ `graph_openapi_schema_detection/core/__init__.py`
3. ✅ `graph_openapi_schema_detection/models/__init__.py`
4. ✅ `core/spec_fetcher.py`
5. ✅ `core/version_detector.py`
6. ✅ `core/schema_parser.py`
7. ✅ `core/schema_comparer.py`
8. ✅ `core/issue_builder.py`
9. ✅ `models/results.py`
10. ✅ `models/changes.py`
11. ✅ `models/statistics.py`
12. ✅ `detect_openapi_changes.py`
13. ✅ `monitor_spec_updates.py`

## Code Reuse from SDK Detection

- ✅ `progress_reporter.py` - Direct reuse
- ✅ `github_client.py` - Reuse for issue creation
- ✅ `analyze_provider_model_usage.py` - Reuse for filtering
- ✅ Package structure and naming conventions

## Next Steps

1. Create package structure
2. Implement Phase 1 (Core Detection)
3. Test with real OpenAPI spec versions
4. Integrate with provider usage filter
5. Add GitHub Actions workflow
6. Documentation and examples
