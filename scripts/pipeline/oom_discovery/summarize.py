#!/usr/bin/env python3
"""Summarize evidence without turning absent results into successful tests."""

import argparse
import json
from pathlib import Path
import re


def summarize(report):
    metadata_path = report / "experiment.json"
    metadata = json.loads(metadata_path.read_text()) if metadata_path.exists() else {}
    text = "\n".join(
        path.read_text(errors="replace")
        for path in (report / "environment.log", report / "resources.log")
        if path.exists()
    )
    available = [int(value) for value in re.findall(r"Mem:\s+(?:\d+\s+){5}(\d+)", text)]
    counters = {}
    for path, value in re.findall(r"(\S+/memory.events):[^\n]*?\boom_kill (\d+)", text):
        counters.setdefault(path, []).append(int(value))
    oom_delta = max((max(values) - values[0] for values in counters.values()), default=None)
    tools = []
    for path in sorted((report / "tools").glob("*.jsonl")):
        for line in path.read_text().splitlines():
            try:
                tools.append(json.loads(line))
            except json.JSONDecodeError:
                continue  # A killed runner may leave a partial final record.
    completed = [tool for tool in tools if tool.get("event") == "end"]
    ended = {tool["invocation"] for tool in completed}
    unfinished = [tool for tool in tools if tool.get("event") == "start" and tool["invocation"] not in ended]
    status = metadata.get("exit_code")
    headroom = min(available) if available else None
    result = {
        **metadata,
        "oom_kill_delta": oom_delta,
        "min_available_mib": headroom,
        "vet_invocations": sum(tool.get("tool") == "vet" and tool.get("event") == "start" for tool in tools),
        "largest_completed_tools": sorted(completed, key=lambda tool: tool["peak_rss_kib"], reverse=True)[:10],
        "unfinished_tools": unfinished,
        "candidate_pass": status == 0 and oom_delta == 0 and headroom is not None and headroom >= 2048,
        "evidence_complete": status is not None and oom_delta is not None and headroom is not None,
    }
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("report", type=Path)
    args = parser.parse_args()
    result = summarize(args.report)
    args.report.mkdir(parents=True, exist_ok=True)
    (args.report / "result.json").write_text(json.dumps(result, indent=2) + "\n")
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
