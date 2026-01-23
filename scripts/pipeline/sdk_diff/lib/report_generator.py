#!/usr/bin/env python3
"""Report generation for SDK diff analysis.

Generates human-readable reports in Markdown and JSON formats.
"""

import json
from typing import Dict, List, Any
from datetime import datetime
from pathlib import Path


def generate_markdown_report(
    current_versions: Dict[str, str],
    latest_versions: Dict[str, str],
    categorized_changes: Dict[str, List[Dict]],
    stats: Dict[str, int],
    breaking_changes: List[str]
) -> str:
    """Generate markdown report.
    
    Args:
        current_versions: Current SDK versions in provider
        latest_versions: Latest available SDK versions
        categorized_changes: Categorized file changes
        stats: Summary statistics
        breaking_changes: Extracted breaking changes from release notes
        
    Returns:
        Markdown formatted report
    """
    sections = []
    
    sections.append(_generate_header())
    sections.append(_generate_version_table(current_versions, latest_versions))
    sections.append(_generate_summary_stats(stats))
    sections.append(_generate_breaking_changes_section(breaking_changes))
    sections.append(_generate_critical_section(categorized_changes))
    sections.append(_generate_warnings_section(categorized_changes))
    sections.append(_generate_safe_changes_section(categorized_changes))
    sections.append(_generate_opportunities_section(categorized_changes))
    sections.append(_generate_actions_section(stats))
    sections.append(_generate_footer())
    
    return '\n'.join(sections)


def _generate_header() -> str:
    """Generate report header."""
    lines = []
    lines.append("# SDK Version Diff Analysis Report")
    lines.append(f"\n**Generated:** {datetime.now().strftime('%Y-%m-%d %H:%M:%S UTC')}")
    lines.append("")
    return '\n'.join(lines)


def _generate_version_table(
    current_versions: Dict[str, str],
    latest_versions: Dict[str, str]
) -> str:
    """Generate version comparison table."""
    lines = []
    lines.append("## 📦 Version Comparison")
    lines.append("")
    lines.append("| SDK Package | Current Version | Latest Version | Status |")
    lines.append("|-------------|-----------------|----------------|--------|")
    
    for sdk in sorted(current_versions.keys()):
        current = current_versions[sdk]
        latest = latest_versions.get(sdk, "Unknown")
        
        if latest == current:
            status = "✅ Up to date"
        elif latest == "Unknown":
            status = "❓ Unknown"
        else:
            status = "⚠️ Update available"
        
        lines.append(f"| `{sdk}` | {current} | {latest} | {status} |")
    
    lines.append("")
    return '\n'.join(lines)


def _generate_summary_stats(stats: Dict[str, int]) -> str:
    """Generate summary statistics section."""
    lines = []
    lines.append("## 📊 Change Summary")
    lines.append("")
    lines.append(f"- **Total SDK changes:** {stats['total_changes']:,} files")
    lines.append(f"- **Relevant to provider:** {stats['relevant_changes']} files ({_percentage(stats['relevant_changes'], stats['total_changes'])})")
    lines.append(f"  - 🚨 **Critical:** {stats['critical']} changes")
    lines.append(f"  - ⚠️  **Warnings:** {stats['warning']} changes")
    lines.append(f"  - ✅ **Safe:** {stats['safe']} changes")
    lines.append(f"  - ✨ **Opportunities:** {stats.get('opportunity', 0)} new fields in used types")
    lines.append(f"- **Filtered as noise:** {stats['noise']:,} files")
    lines.append("")
    return '\n'.join(lines)


def _generate_breaking_changes_section(breaking_changes: List[str]) -> str:
    """Generate breaking changes from release notes section."""
    if not breaking_changes:
        return ""
    
    lines = []
    lines.append("## 🚨 Breaking Changes (from Release Notes)")
    lines.append("")
    for change in breaking_changes:
        lines.append(f"- {change}")
    lines.append("")
    return '\n'.join(lines)


def _generate_critical_section(categorized_changes: Dict[str, List[Dict]]) -> str:
    """Generate critical changes section."""
    critical = categorized_changes.get("critical", [])
    if not critical:
        return ""
    
    lines = []
    lines.append(f"## 🚨 Critical Changes ({len(critical)})")
    lines.append("")
    lines.append("These changes directly affect APIs used by the provider and **require immediate attention**.")
    lines.append("")
    
    for i, change in enumerate(critical, 1):
        lines.append(f"### {i}. `{change['file']}`")
        lines.append(f"- **Status:** {change['status'].upper()}")
        lines.append(f"- **Reason:** {change['reason']}")
        lines.append(f"- **Changes:** +{change['additions']} / -{change['deletions']}")
        lines.append("")
    
    return '\n'.join(lines)


def _generate_warnings_section(categorized_changes: Dict[str, List[Dict]]) -> str:
    """Generate warnings section."""
    warnings = categorized_changes.get("warning", [])
    if not warnings:
        return ""
    
    lines = []
    lines.append(f"## ⚠️  Warnings ({len(warnings)})")
    lines.append("")
    lines.append("These changes may affect the provider. Review recommended.")
    lines.append("")
    
    display_warnings = warnings[:10]
    for change in display_warnings:
        lines.append(f"- `{change['file']}` - {change['status']} - {change['reason']}")
    
    if len(warnings) > 10:
        lines.append(f"\n_... and {len(warnings) - 10} more warnings_")
    lines.append("")
    
    return '\n'.join(lines)


def _generate_safe_changes_section(categorized_changes: Dict[str, List[Dict]]) -> str:
    """Generate safe changes section."""
    safe = categorized_changes.get("safe", [])
    if not safe:
        return ""
    
    lines = []
    lines.append(f"## ✅ Safe Changes ({len(safe)})")
    lines.append("")
    lines.append("New features and additions in packages used by the provider.")
    lines.append("")
    
    by_package = {}
    for change in safe:
        pkg = _extract_package(change['file'])
        if pkg not in by_package:
            by_package[pkg] = []
        by_package[pkg].append(change['file'])
    
    for pkg, files in sorted(by_package.items()):
        lines.append(f"### Package: `{pkg}`")
        for file in files[:5]:
            lines.append(f"- {file}")
        if len(files) > 5:
            lines.append(f"- _... and {len(files) - 5} more files_")
        lines.append("")
    
    return '\n'.join(lines)


def _generate_opportunities_section(categorized_changes: Dict[str, List[Dict]]) -> str:
    """Generate opportunities section for new fields."""
    opportunities = categorized_changes.get("opportunity", [])
    if not opportunities:
        return ""
    
    lines = []
    lines.append(f"## ✨ Opportunities ({len(opportunities)})")
    lines.append("")
    lines.append("New fields added to SDK types you're already using. Consider adopting these:")
    lines.append("")
    
    by_type = {}
    for opp in opportunities:
        type_name = opp.get('type', 'Unknown')
        if type_name not in by_type:
            by_type[type_name] = []
        by_type[type_name].append(opp)
    
    for type_name, fields in sorted(by_type.items()):
        lines.append(f"### `{type_name}`")
        lines.append("")
        for field in fields:
            field_name = field.get('field', '')
            field_type = field.get('field_type', '')
            description = field.get('description', '')
            
            lines.append(f"- **`{field_name}`** (`{field_type}`)")
            if description:
                lines.append(f"  - {description}")
        lines.append("")
    
    return '\n'.join(lines)


def _generate_actions_section(stats: Dict[str, int]) -> str:
    """Generate recommended actions section."""
    lines = []
    lines.append("## 🎯 Recommended Actions")
    lines.append("")
    
    if stats['critical'] > 0:
        lines.append(f"1. ❗ **Address {stats['critical']} critical changes** before upgrading")
    
    if stats['warning'] > 0:
        lines.append(f"2. ⚠️  Review {stats['warning']} warnings for potential issues")
    
    if stats['safe'] > 0:
        lines.append(f"3. ✅ Consider using {stats['safe']} new features")
    
    if stats.get('opportunity', 0) > 0:
        lines.append(f"4. ✨ Review {stats['opportunity']} new field(s) in types you use")
    
    if stats['critical'] == 0 and stats['warning'] == 0:
        lines.append("✅ **No breaking changes detected!** Safe to upgrade.")
    
    lines.append("")
    return '\n'.join(lines)


def _generate_footer() -> str:
    """Generate report footer."""
    lines = []
    lines.append("---")
    lines.append("_This report was automatically generated by the SDK diff analysis tool._")
    return '\n'.join(lines)


def _percentage(part: int, total: int) -> str:
    """Calculate percentage as formatted string.
    
    Args:
        part: Numerator
        total: Denominator
        
    Returns:
        Formatted percentage like "1.2%"
    """
    if total == 0:
        return "0.0%"
    return f"{(part / total * 100):.1f}%"


def _extract_package(filepath: str) -> str:
    """Extract package name from file path.
    
    Args:
        filepath: Path like "models/user.go"
        
    Returns:
        Package name like "models"
    """
    parts = filepath.split('/')
    if len(parts) > 1:
        return parts[0]
    return "root"


def save_json_report(
    data: Dict[str, Any],
    output_path: Path
) -> None:
    """Save analysis data as JSON for programmatic access.
    
    Args:
        data: Analysis data to save
        output_path: Path to output JSON file
    """
    with open(output_path, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2)
    
    print(f"💾 JSON report saved: {output_path}")
