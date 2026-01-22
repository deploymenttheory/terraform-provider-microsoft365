#!/usr/bin/env python3
"""Coverage calculation and analysis.

Provides functions for calculating coverage statistics from Go coverage files.
"""

from pathlib import Path
from typing import Dict, Any


def calculate_coverage(coverage_file: Path) -> Dict[str, Any]:
    """Calculate coverage statistics from Go coverage file.
    
    Parses coverage file in Go's coverage format and calculates:
    - Coverage percentage
    - Total statements
    - Covered statements
    
    Args:
        coverage_file: Path to Go coverage file.
    
    Returns:
        Dict with coverage statistics:
        {
            "coverage_pct": float,      # Coverage percentage (0-100)
            "total_lines": int,         # Total statements
            "covered_lines": int        # Covered statements
        }
    """
    if not coverage_file.exists():
        print(f"⚠️  Coverage file not found: {coverage_file}")
        return {"coverage_pct": 0.0, "total_lines": 0, "covered_lines": 0}
    
    total_statements = 0
    covered_statements = 0
    
    with open(coverage_file, 'r', encoding='utf-8') as f:
        for line in f:
            if line.startswith('mode:'):
                continue
            
            # Parse coverage line format: file.go:startline.col,endline.col num_statements count
            parts = line.strip().split()
            if len(parts) < 3:
                continue
            
            try:
                num_statements = int(parts[1])
                hit_count = int(parts[2])
                
                total_statements += num_statements
                if hit_count > 0:
                    covered_statements += num_statements
            except (ValueError, IndexError):
                continue
    
    if total_statements == 0:
        return {"coverage_pct": 0.0, "total_lines": 0, "covered_lines": 0}
    
    coverage_pct = round((covered_statements / total_statements) * 100, 2)
    
    return {
        "coverage_pct": coverage_pct,
        "total_lines": total_statements,
        "covered_lines": covered_statements
    }


def check_coverage_threshold(coverage_pct: float, min_threshold: float) -> bool:
    """Check if coverage meets minimum threshold.
    
    Args:
        coverage_pct: Coverage percentage to check.
        min_threshold: Minimum required coverage percentage.
    
    Returns:
        True if coverage meets threshold, False otherwise.
    """
    return coverage_pct >= min_threshold


def format_coverage_summary(stats: Dict[str, Any]) -> str:
    """Format coverage statistics as human-readable string.
    
    Args:
        stats: Coverage statistics dict.
    
    Returns:
        Formatted summary string.
    """
    return (
        f"Coverage: {stats['coverage_pct']}%\n"
        f"Total Statements: {stats['total_lines']}\n"
        f"Covered Statements: {stats['covered_lines']}"
    )
