"""Measure PR added lines using complete, exact-head Codecov file reports."""

import re
import subprocess

from lib.git_operations import get_changed_files, get_merge_base


def added_go_lines(base_ref, expected_commit):
    """Read added Go lines from the actual PR merge base, never the last commit."""
    head = subprocess.run(
        ["git", "rev-parse", "HEAD"], capture_output=True, text=True, check=True
    ).stdout.strip()
    if not expected_commit or head != expected_commit or not base_ref:
        raise ValueError("Patch coverage requires the checked-out exact head and Git base")
    base = get_merge_base(base_ref)
    added = {}
    for name in get_changed_files(base_ref):
        if name.endswith("_test.go"):
            continue
        diff = subprocess.run(
            [
                "git",
                "diff",
                "--no-ext-diff",
                "--no-textconv",
                "--no-renames",
                "--unified=0",
                base,
                head,
                "--",
                name,
            ],
            capture_output=True,
            text=True,
            check=True,
        ).stdout
        if "Binary files " in diff or "GIT binary patch" in diff:
            raise ValueError(f"Cannot verify added executable lines in binary file: {name}")
        lines = set()
        for match in re.finditer(r"^@@ .* \+(\d+)(?:,(\d+))? @@", diff, re.MULTILINE):
            start, count = int(match[1]), int(match[2]) if match[2] is not None else 1
            lines.update(range(start, start + count))
        if lines:
            added[name] = lines
    return added


def _integer(value):
    """Accept integer counters while rejecting booleans and coerced numeric strings."""
    return isinstance(value, int) and not isinstance(value, bool)


def _file_lines(file_report):
    """Validate every executable line against Codecov's complete per-file totals."""
    rows, totals = file_report.get("line_coverage"), file_report.get("totals")
    if not isinstance(rows, list) or not isinstance(totals, dict):
        raise ValueError("Missing complete file coverage")
    lines, counts = {}, [0, 0, 0]
    for row in rows:
        if not isinstance(row, list) or len(row) != 2:
            raise ValueError("Invalid executable line coverage shape")
        if (
            not all(_integer(value) for value in row)
            or row[0] < 1
            or row[1] not in (0, 1, 2)
            or row[0] in lines
        ):
            raise ValueError("Invalid or duplicate executable line coverage")
        lines[row[0]] = row[1]
        counts[row[1]] += 1
    expected = [totals.get(key) for key in ("hits", "misses", "partials")]
    if (
        any(not _integer(value) or value < 0 for value in expected)
        or counts != expected
        or not _integer(totals.get("lines"))
        or totals["lines"] != len(lines)
    ):
        raise ValueError("Incomplete file coverage totals")
    return lines


def measure_patch(report, added, expected_commit):
    """Count changed executable lines; missed and partially covered lines both fail coverage."""
    if (
        not isinstance(report, dict)
        or not expected_commit
        or report.get("commitid") != expected_commit
        or report.get("state") != "complete"
    ):
        raise ValueError("Codecov commit report is stale or incomplete")
    body = report.get("report")
    if not isinstance(body, dict):
        raise ValueError("Missing Codecov report")
    files = body.get("files")
    if not isinstance(files, list):
        raise ValueError("Missing Codecov file reports")
    by_name = {}
    for item in files:
        if not isinstance(item, dict) or not isinstance(item.get("name"), str):
            raise ValueError("Invalid Codecov file report")
        if item["name"] in by_name:
            raise ValueError("Duplicate Codecov file report")
        by_name[item["name"]] = item
    total = covered = 0
    for name, added_lines in added.items():
        if name not in by_name:
            raise ValueError(f"Missing coverage for changed file: {name}")
        lines = _file_lines(by_name[name])
        changed = added_lines.intersection(lines)
        total += len(changed)
        covered += sum(lines[line] == 0 for line in changed)
    return {
        "coverage_pct": covered * 100 / total if total else None,
        "total_lines": total,
        "covered_lines": covered,
        "no_executable_changes": total == 0,
    }
