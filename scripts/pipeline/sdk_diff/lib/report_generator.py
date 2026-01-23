#!/usr/bin/env python3
"""Multi-file report generation for SDK diff analysis.

Generates structured reports optimized for actionability and future PR automation.
"""

import json
from typing import Dict, List, Any
from datetime import datetime
from pathlib import Path


def generate_all_reports(
    sdk_name: str,
    current_version: str,
    latest_version: str,
    comparison_url: str,
    categorized_changes: Dict[str, List[Dict]],
    stats: Dict[str, int],
    breaking_changes: List[str],
    usage_data: Dict[str, Any],
    output_dir: Path
) -> List[str]:
    """Generate all report files.
    
    Args:
        sdk_name: SDK package name
        current_version: Current version in provider
        latest_version: Latest available version
        comparison_url: GitHub comparison URL
        categorized_changes: Changes categorized by impact
        stats: Summary statistics
        breaking_changes: Breaking changes from release notes
        usage_data: Provider's SDK usage baseline
        output_dir: Directory to write reports
        
    Returns:
        List of generated file paths
    """
    output_dir = Path(output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)
    
    generated_files = []
    metadata = {
        "sdk": sdk_name,
        "current_version": current_version,
        "latest_version": latest_version,
        "comparison_url": comparison_url,
        "generated_at": datetime.utcnow().isoformat() + "Z"
    }
    
    # 01 - Summary
    exec_summary = _generate_executive_summary(metadata, stats, breaking_changes)
    path = output_dir / "01_summary.md"
    path.write_text(exec_summary)
    generated_files.append(str(path))
    
    # 02 - Breaking Changes
    critical_report = _generate_critical_breaking_changes(
        metadata, categorized_changes, breaking_changes
    )
    path = output_dir / "02_breaking_changes.md"
    path.write_text(critical_report)
    generated_files.append(str(path))
    
    # 03 - Provider Updates Required
    updates_report = _generate_provider_updates_required(
        metadata, categorized_changes
    )
    path = output_dir / "03_provider_updates_required.md"
    path.write_text(updates_report)
    generated_files.append(str(path))
    
    # 04 - Type Structure Changes
    type_changes = _generate_type_structure_changes(
        metadata, categorized_changes
    )
    path = output_dir / "04_type_structure_changes.md"
    path.write_text(type_changes)
    generated_files.append(str(path))
    
    # 05 - New Opportunities
    opportunities = _generate_new_opportunities(
        metadata, categorized_changes
    )
    path = output_dir / "05_new_opportunities.md"
    path.write_text(opportunities)
    generated_files.append(str(path))
    
    # 06 - Repo Metadata Changes
    metadata_report = _generate_metadata_changes(
        metadata, categorized_changes
    )
    path = output_dir / "06_repo_metadata_changes.md"
    path.write_text(metadata_report)
    generated_files.append(str(path))
    
    # 07 - Noise Filtered
    noise_data = {
        "metadata": metadata,
        "filtered_count": stats.get("noise", 0),
        "filtered_files": categorized_changes.get("noise", []),
        "filter_reasons": _generate_filter_reasons()
    }
    path = output_dir / "07_noise_filtered.json"
    path.write_text(json.dumps(noise_data, indent=2))
    generated_files.append(str(path))
    
    # 08 - Actionable Changes (Resource-Grouped)
    actionable = _generate_actionable_changes_json(
        metadata, categorized_changes, usage_data
    )
    path = output_dir / "08_actionable_changes.json"
    path.write_text(json.dumps(actionable, indent=2))
    generated_files.append(str(path))
    
    # 09 - Usage Baseline
    path = output_dir / "09_usage_baseline.json"
    path.write_text(json.dumps(usage_data, indent=2))
    generated_files.append(str(path))
    
    return generated_files


def _generate_executive_summary(
    metadata: Dict[str, str],
    stats: Dict[str, int],
    breaking_changes: List[str]
) -> str:
    """Generate executive summary."""
    lines = [
        "# 📊 SDK Version Diff - Executive Summary",
        "",
        f"**Generated:** {metadata['generated_at']}",
        f"**SDK:** `{metadata['sdk']}`",
        f"**Version Change:** `{metadata['current_version']}` → `{metadata['latest_version']}`",
        "",
        "## 🎯 Quick Decision Guide",
        ""
    ]
    
    # Decision based on stats
    critical_count = stats.get("critical", 0)
    enum_removed = stats.get("enum_removed", 0)
    warning_count = stats.get("warning", 0)
    
    if critical_count > 0 or enum_removed > 0:
        lines.extend([
            "### 🔴 **ACTION REQUIRED - Breaking Changes Detected**",
            "",
            f"- **{critical_count}** critical breaking changes",
            f"- **{enum_removed}** enum values removed (breaking)",
            "",
            "**Decision:** Review `02_breaking_changes.md` before upgrading.",
            ""
        ])
    elif warning_count > 0:
        lines.extend([
            "### 🟡 **REVIEW RECOMMENDED - Changes May Affect Provider**",
            "",
            f"- **{warning_count}** changes in packages you use",
            "",
            "**Decision:** Review `03_provider_updates_required.md` for potential improvements.",
            ""
        ])
    else:
        lines.extend([
            "### 🟢 **SAFE TO UPGRADE - No Breaking Changes**",
            "",
            "No breaking changes detected in packages used by the provider.",
            "",
            "**Decision:** Safe to upgrade. Consider new features in `05_new_opportunities.md`.",
            ""
        ])
    
    # Summary statistics
    lines.extend([
        "## 📈 Change Statistics",
        "",
        f"| Category | Count |",
        f"|----------|-------|",
        f"| Total SDK Changes | {stats.get('total_changes', 0)} |",
        f"| Actionable Changes | {stats.get('relevant_changes', 0)} |",
        f"| 🔴 Critical | {critical_count} |",
        f"| 🟡 Warnings | {warning_count} |",
        f"| ✅ Safe Additions | {stats.get('safe', 0)} |",
        f"| ✨ New Opportunities | {stats.get('opportunity', 0)} |",
        f"| 🎯 Enum Values Added | {stats.get('enum_added', 0)} |",
        f"| 🔥 Enum Values Removed | {enum_removed} |",
        f"| ⚪ Metadata Changes | {stats.get('metadata', 0)} |",
        f"| 🗑️ Filtered as Noise | {stats.get('noise', 0)} |",
        "",
        "## 📂 Detailed Reports",
        "",
        "1. `02_breaking_changes.md` - Must address before upgrade",
        "2. `03_provider_updates_required.md` - Fields/methods needing review",
        "3. `04_type_structure_changes.md` - Type modifications in use",
        "4. `05_new_opportunities.md` - New features to consider",
        "5. `06_repo_metadata_changes.md` - Build/config file changes (informational)",
        "6. `08_actionable_changes.json` - Structured data for PR automation",
        "",
        "---",
        "_For questions about this analysis, see the SDK diff pipeline documentation._"
    ])
    
    return '\n'.join(lines)


def _generate_critical_breaking_changes(
    metadata: Dict[str, str],
    categorized: Dict[str, List[Dict]],
    breaking_from_notes: List[str]
) -> str:
    """Generate critical breaking changes report."""
    lines = [
        "# 🔴 Critical Breaking Changes",
        "",
        f"**SDK:** `{metadata['sdk']}`",
        f"**Version:** `{metadata['current_version']}` → `{metadata['latest_version']}`",
        f"**Generated:** {metadata['generated_at']}",
        "",
        "## ⚠️ IMMEDIATE ACTION REQUIRED",
        "",
        "These changes will break existing provider functionality if not addressed.",
        ""
    ]
    
    critical_files = categorized.get("critical", [])
    enum_removed = categorized.get("enum_removed", [])
    
    if not critical_files and not enum_removed and not breaking_from_notes:
        lines.extend([
            "## ✅ No Critical Breaking Changes",
            "",
            "No critical breaking changes detected. Safe to proceed with upgrade.",
            ""
        ])
        return '\n'.join(lines)
    
    # Breaking changes from release notes
    if breaking_from_notes:
        lines.extend([
            "## 📋 Breaking Changes from Release Notes",
            ""
        ])
        for change in breaking_from_notes:
            lines.append(f"- {change}")
        lines.append("")
    
    # Enum values removed
    if enum_removed:
        lines.extend([
            "## 🎯 Enum Values Removed (Breaking)",
            "",
            "These enum values have been removed from the SDK. Any provider code using these values must be updated.",
            ""
        ])
        for change in enum_removed:
            enum_name = change.get("enum", "Unknown")
            removed_values = change.get("removed_values", [])
            used_in = change.get("used_in_files", [])
            
            lines.extend([
                f"### `{enum_name.split('.')[-1]}`",
                "",
                f"**Removed values:** `{', '.join(removed_values)}`",
                "",
                "**Used in:**"
            ])
            for file_path in used_in[:5]:  # Show first 5
                rel_path = file_path.split("internal/")[-1] if "internal/" in file_path else file_path
                lines.append(f"- `{rel_path}`")
            if len(used_in) > 5:
                lines.append(f"- _...and {len(used_in) - 5} more files_")
            
            lines.extend([
                "",
                "**Action Required:**",
                "1. Remove removed values from resource schema validation",
                "2. Update any default values using removed values",
                "3. Add migration notes to documentation",
                ""
            ])
    
    # Files removed
    if critical_files:
        lines.extend([
            "## 🗑️ Files Removed from SDK",
            "",
            "These files in packages you use have been removed:",
            ""
        ])
        for file_change in critical_files:
            lines.append(f"- `{file_change['file']}`")
            if file_change.get('reason'):
                lines.append(f"  - {file_change['reason']}")
        lines.append("")
    
    lines.extend([
        "## 🔗 References",
        "",
        f"- [Full Comparison]({metadata['comparison_url']})",
        f"- [Release Notes](https://github.com/microsoftgraph/msgraph-beta-sdk-go/releases/tag/{metadata['latest_version']})",
        ""
    ])
    
    return '\n'.join(lines)


def _generate_provider_updates_required(
    metadata: Dict[str, str],
    categorized: Dict[str, List[Dict]]
) -> str:
    """Generate provider updates required report."""
    lines = [
        "# 🟡 Provider Updates Required",
        "",
        f"**SDK:** `{metadata['sdk']}`",
        f"**Version:** `{metadata['current_version']}` → `{metadata['latest_version']}`",
        "",
        "## 📝 Changes Requiring Review",
        "",
        "These changes may require updates to the Terraform provider.",
        ""
    ]
    
    warnings = categorized.get("warning", [])
    field_changes = categorized.get("opportunity", [])  # Fields added/removed
    enum_added = categorized.get("enum_added", [])
    
    if not warnings and not field_changes and not enum_added:
        lines.extend([
            "## ✅ No Updates Required",
            "",
            "No changes requiring provider updates detected.",
            ""
        ])
        return '\n'.join(lines)
    
    # Field additions (opportunities)
    if field_changes:
        lines.extend([
            "## ✨ New Fields in Used Types",
            "",
            "These fields have been added to types used by the provider. Consider adding them to resource schemas.",
            ""
        ])
        for change in field_changes:
            type_name = change.get("type", "Unknown").split(".")[-1]
            field_name = change.get("field", "unknown")
            
            lines.extend([
                f"### `{type_name}.{field_name}`",
                "",
                f"**Type:** `{type_name}`",
                f"**New Field:** `{field_name}`",
                "",
                "**Suggested Action:** Consider adding to resource schema if relevant to use case.",
                ""
            ])
    
    # Enum values added
    if enum_added:
        lines.extend([
            "## 🎯 Enum Values Added",
            "",
            "New values have been added to enums used by the provider.",
            ""
        ])
        for change in enum_added:
            enum_name = change.get("enum", "Unknown").split(".")[-1]
            added_values = change.get("added_values", [])
            
            lines.extend([
                f"### `{enum_name}`",
                "",
                f"**Added values:** `{', '.join(added_values)}`",
                "",
                "**Suggested Action:** Update schema validation to accept new values if applicable.",
                ""
            ])
    
    # Modified files in used packages
    if warnings:
        lines.extend([
            "## ⚠️ Modified Files in Used Packages",
            "",
            f"**Total:** {len(warnings)} files modified in packages used by the provider.",
            "",
            "These changes may affect provider functionality. Review recommended:",
            ""
        ])
        
        # Group by package
        by_package = {}
        for file_change in warnings:
            filename = file_change['file']
            # Extract package
            if '/models/' in filename:
                pkg = 'models'
            elif '/' in filename:
                pkg = filename.split('/')[0]
            else:
                pkg = 'root'
            
            if pkg not in by_package:
                by_package[pkg] = []
            by_package[pkg].append(file_change)
        
        for pkg, changes in sorted(by_package.items()):
            lines.append(f"### Package: `{pkg}` ({len(changes)} files)")
            lines.append("")
            for change in changes[:10]:  # Show first 10
                lines.append(f"- `{change['file']}`")
                if change.get('additions') or change.get('deletions'):
                    lines.append(f"  - +{change.get('additions', 0)} / -{change.get('deletions', 0)} lines")
            if len(changes) > 10:
                lines.append(f"- _...and {len(changes) - 10} more files_")
            lines.append("")
    
    lines.extend([
        "## 🔗 Review Changes",
        "",
        f"[View Full Diff]({metadata['comparison_url']})",
        ""
    ])
    
    return '\n'.join(lines)


def _generate_type_structure_changes(
    metadata: Dict[str, str],
    categorized: Dict[str, List[Dict]]
) -> str:
    """Generate type structure changes report."""
    lines = [
        "# 🟡 Type Structure Changes",
        "",
        f"**SDK:** `{metadata['sdk']}`",
        f"**Version:** `{metadata['current_version']}` → `{metadata['latest_version']}`",
        "",
        "## 📦 Modified Types in Used Packages",
        "",
        "Changes to type structures in packages used by the provider.",
        ""
    ]
    
    # This would need more detailed analysis
    # For now, show modified model files
    warnings = [w for w in categorized.get("warning", []) if '/models/' in w['file'] and w['file'].endswith('.go')]
    
    if not warnings:
        lines.extend([
            "## ✅ No Type Structure Changes",
            "",
            "No changes to type structures detected.",
            ""
        ])
        return '\n'.join(lines)
    
    lines.extend([
        f"**Total:** {len(warnings)} model files modified",
        ""
    ])
    
    for change in warnings[:20]:  # Show first 20
        type_name = change['file'].replace('models/', '').replace('.go', '')
        lines.extend([
            f"### `{type_name}`",
            f"- File: `{change['file']}`",
            f"- Changes: +{change.get('additions', 0)} / -{change.get('deletions', 0)} lines",
            ""
        ])
    
    if len(warnings) > 20:
        lines.append(f"_...and {len(warnings) - 20} more types_")
    
    lines.extend([
        "",
        "## 🔍 Recommended Action",
        "",
        "1. Review modified types used directly in provider code",
        "2. Check for breaking field changes or renames",
        "3. Update resource schemas if needed",
        "",
        f"[View Full Diff]({metadata['comparison_url']})",
        ""
    ])
    
    return '\n'.join(lines)


def _generate_new_opportunities(
    metadata: Dict[str, str],
    categorized: Dict[str, List[Dict]]
) -> str:
    """Generate new opportunities report."""
    lines = [
        "# 🟢 New Opportunities",
        "",
        f"**SDK:** `{metadata['sdk']}`",
        f"**Version:** `{metadata['current_version']}` → `{metadata['latest_version']}`",
        "",
        "## ✨ New Features Available",
        "",
        "New features and types added to packages used by the provider.",
        ""
    ]
    
    safe_changes = categorized.get("safe", [])
    opportunities = categorized.get("opportunity", [])
    enum_added = categorized.get("enum_added", [])
    
    if not safe_changes and not opportunities and not enum_added:
        lines.extend([
            "## ℹ️ No New Features",
            "",
            "No new features or types detected in this release.",
            ""
        ])
        return '\n'.join(lines)
    
    # New types/files
    if safe_changes:
        lines.extend([
            f"## 🆕 New Types ({len(safe_changes)})",
            "",
            "New types and features added to the SDK:",
            ""
        ])
        
        # Group by package
        by_package = {}
        for change in safe_changes:
            filename = change['file']
            if '/models/' in filename:
                pkg = 'models'
            elif '/' in filename:
                pkg = filename.split('/')[0]
            else:
                pkg = 'root'
            
            if pkg not in by_package:
                by_package[pkg] = []
            by_package[pkg].append(change)
        
        for pkg, changes in sorted(by_package.items()):
            lines.append(f"### Package: `{pkg}`")
            lines.append("")
            for change in changes[:15]:
                type_name = change['file'].replace('models/', '').replace('.go', '')
                lines.append(f"- `{type_name}`")
            if len(changes) > 15:
                lines.append(f"- _...and {len(changes) - 15} more_")
            lines.append("")
    
    # New fields
    if opportunities:
        lines.extend([
            f"## ➕ New Fields ({len(opportunities)})",
            "",
            "New fields added to existing types:",
            ""
        ])
        for opp in opportunities[:10]:
            lines.append(f"- `{opp.get('type', 'Unknown')}.{opp.get('field', 'unknown')}`")
        if len(opportunities) > 10:
            lines.append(f"- _...and {len(opportunities) - 10} more_")
        lines.append("")
    
    # New enum values
    if enum_added:
        lines.extend([
            f"## 🎯 New Enum Values ({len(enum_added)})",
            "",
            "New values added to existing enums:",
            ""
        ])
        for change in enum_added:
            enum_name = change.get("enum", "Unknown").split(".")[-1]
            added = change.get("added_values", [])
            lines.append(f"- `{enum_name}`: {', '.join(added)}")
        lines.append("")
    
    lines.extend([
        "## 💡 Suggested Actions",
        "",
        "1. Review new types for potential provider resources",
        "2. Consider adding new fields to existing resources",
        "3. Update validation to accept new enum values",
        "",
        f"[View Full Changes]({metadata['comparison_url']})",
        ""
    ])
    
    return '\n'.join(lines)


def _generate_metadata_changes(
    metadata: Dict[str, str],
    categorized: Dict[str, List[Dict]]
) -> str:
    """Generate metadata changes report."""
    lines = [
        "# ⚪ Metadata Changes",
        "",
        f"**SDK:** `{metadata['sdk']}`",
        f"**Version:** `{metadata['current_version']}` → `{metadata['latest_version']}`",
        "",
        "## ℹ️ Build and Configuration Changes",
        "",
        "These are metadata/build file changes. No code action required.",
        ""
    ]
    
    metadata_changes = categorized.get("metadata", [])
    
    if not metadata_changes:
        lines.extend([
            "No metadata file changes detected.",
            ""
        ])
        return '\n'.join(lines)
    
    lines.extend([
        f"**Total:** {len(metadata_changes)} metadata files changed",
        ""
    ])
    
    for change in metadata_changes:
        lines.append(f"- `{change['file']}` ({change['status']})")
        if change.get('additions') or change.get('deletions'):
            lines.append(f"  - +{change.get('additions', 0)} / -{change.get('deletions', 0)} lines")
    
    lines.extend([
        "",
        "## 📝 Note",
        "",
        "These changes are informational only. They typically include:",
        "- Dependency updates (`go.mod`, `go.sum`)",
        "- Build tool configuration (`kiota-lock.json`)",
        "- Release automation files",
        "",
        "No provider code changes are needed for these updates.",
        ""
    ])
    
    return '\n'.join(lines)


def _generate_filter_reasons() -> Dict[str, str]:
    """Generate filter reason explanations."""
    return {
        "test_files": "Test files (_test.go) do not affect runtime behavior",
        "test_data": "Test data directories contain fixtures, not production code",
        "examples": "Example code is documentation, not part of the SDK API",
        "documentation": "Documentation changes (README, CHANGELOG) are informational",
        "github_workflows": "CI/CD configuration does not affect the SDK API",
        "internal_packages": "Internal packages are not part of the public API",
        "unused_packages": "Changes in packages not imported by the provider"
    }


def _generate_actionable_changes_json(
    metadata: Dict[str, str],
    categorized: Dict[str, List[Dict]],
    usage_data: Dict[str, Any]
) -> Dict[str, Any]:
    """Generate resource-grouped actionable changes for PR automation.
    
    This is the key file for future PR generation.
    """
    # TODO: Implement resource mapping logic
    # For now, create a basic structure
    
    actionable = {
        "metadata": metadata,
        "summary": {
            "breaking_changes": len(categorized.get("critical", [])) + len(categorized.get("enum_removed", [])),
            "schema_updates_available": len(categorized.get("opportunity", [])) + len(categorized.get("enum_added", [])),
            "type_modifications": len([w for w in categorized.get("warning", []) if '/models/' in w.get('file', '')]),
            "total_actionable": len(categorized.get("critical", [])) + len(categorized.get("warning", [])) + len(categorized.get("opportunity", []))
        },
        "resource_impacts": {},  # Will be populated by resource mapping
        "critical_actions": [],
        "enhancement_opportunities": [],
        "pr_templates_suggested": []
    }
    
    # Add critical actions
    for enum_change in categorized.get("enum_removed", []):
        actionable["critical_actions"].append({
            "type": "enum_value_removed",
            "priority": "critical",
            "enum": enum_change.get("enum", ""),
            "removed_values": enum_change.get("removed_values", []),
            "affected_files": enum_change.get("used_in_files", []),
            "action_required": "Update resource schema validation to remove deprecated values"
        })
    
    # Add enhancement opportunities
    for field_add in categorized.get("opportunity", []):
        actionable["enhancement_opportunities"].append({
            "type": "field_added",
            "priority": "optional",
            "entity": field_add.get("type", ""),
            "field": field_add.get("field", ""),
            "action_suggested": "Consider adding to resource schema if relevant"
        })
    
    for enum_add in categorized.get("enum_added", []):
        actionable["enhancement_opportunities"].append({
            "type": "enum_value_added",
            "priority": "optional",
            "enum": enum_add.get("enum", ""),
            "added_values": enum_add.get("added_values", []),
            "action_suggested": "Update schema validation to accept new values"
        })
    
    return actionable
