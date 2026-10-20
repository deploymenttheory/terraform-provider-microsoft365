#!/usr/bin/env python3
"""Summarize evidence without turning absent results into successful tests."""

import argparse
import json
from pathlib import Path
import re


def summarize_evidence(metadata, text, tools):
    available = [int(value) for value in re.findall(r"Mem:\s+(?:\d+\s+){5}(\d+)", text)]
    counters = {}
    for path, value in re.findall(r"(\S+/memory.events):[^\n]*?\boom_kill (\d+)", text):
        counters.setdefault(path, []).append(int(value))
    oom_delta = max((max(values) - values[0] for values in counters.values()), default=None)
    completed = [tool for tool in tools if tool.get("event") == "end"]
    ended = {tool["invocation"] for tool in completed}
    unfinished = [tool for tool in tools if tool.get("event") == "start" and tool["invocation"] not in ended]
    status = metadata.get("exit_code")
    headroom = min(available) if available else None
    return {
        **metadata,
        "oom_kill_delta": oom_delta,
        "min_available_mib": headroom,
        "vet_invocations": sum(tool.get("tool") == "vet" and tool.get("event") == "start" for tool in tools),
        "largest_completed_tools": sorted(completed, key=lambda tool: tool["peak_rss_kib"], reverse=True)[:10],
        "unfinished_tools": unfinished,
        "candidate_pass": status == 0 and oom_delta == 0 and headroom is not None and headroom >= 2048,
        "evidence_complete": status is not None and oom_delta is not None and headroom is not None,
    }


def summarize(report):
    metadata_path = report / "experiment.json"
    metadata = json.loads(metadata_path.read_text()) if metadata_path.exists() else {}
    text = "\n".join(
        path.read_text(errors="replace")
        for path in (report / "environment.log", report / "resources.log")
        if path.exists()
    )
    tools = []
    for path in sorted((report / "tools").glob("*.jsonl")):
        for line in path.read_text().splitlines():
            try:
                tools.append(json.loads(line))
            except json.JSONDecodeError:
                continue  # A killed runner may leave a partial final record.
    return summarize_evidence(metadata, text, tools)


def summarize_log(text):
    """Recover streamed evidence even if runner termination prevented upload."""
    # Bounded replay is a separate experiment; its deliberate OOM control must
    # not contaminate the preceding build's counters or tool attribution.
    text = text.split("[VET REPLAY BEGIN]", 1)[0]
    metadata = {}
    tools = []
    for line in text.splitlines():
        for marker in ("[EXPERIMENT] ", "[EXPERIMENT RESULT] ", "[GO TOOL] "):
            if marker not in line:
                continue
            try:
                record = json.loads(line.split(marker, 1)[1])
            except json.JSONDecodeError:
                continue
            if marker == "[GO TOOL] ":
                tools.append(record)
            else:
                metadata.update(record)
    if "exit_code" not in metadata:
        statuses = re.findall(r"Command exit status: (\d+)", text)
        if statuses:
            metadata["exit_code"] = int(statuses[-1])
            metadata["exit_status_source"] = "resource monitor"
    return summarize_evidence(metadata, text, tools)


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
