#!/usr/bin/env python3
"""Replay the largest completed vet invocation inside a bounded cgroup."""

import argparse
import json
import os
from pathlib import Path
import subprocess
import sys
import time


def service_command(unit, report, replay, memory_limit):
    command = [
        "sudo", "-n", "systemd-run", "--wait", "--pipe",
        f"--unit={unit}", "--service-type=exec", "--expand-environment=no",
        f"--uid={os.getuid()}", f"--gid={os.getgid()}",
        f"--working-directory={replay['cwd']}",
        "--property=MemoryAccounting=yes", "--property=MemoryMax=13G",
        "--property=MemorySwapMax=0", "--property=RuntimeMaxSec=300",
        "--property=OOMPolicy=stop",
        "--setenv=GOMAXPROCS=1", "--setenv=PYTHONUNBUFFERED=1",
        f"--setenv=OOM_DISCOVERY_REPORT={report}",
        f"--setenv=TOOLEXEC_IMPORTPATH={replay['package']}",
    ]
    if memory_limit:
        command.append(f"--setenv=GOMEMLIMIT={memory_limit}")
    return [*command, sys.executable, str(Path(__file__).with_name("tool_exec.py")), *replay["argv"]]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("report", type=Path)
    args = parser.parse_args()
    report = args.report.resolve()
    completed = []
    for path in (report / "tools").glob("*.jsonl"):
        completed.extend(
            record for record in map(json.loads, path.read_text().splitlines())
            if record["tool"] == "vet" and record["event"] == "end" and record["exit_code"] == 0
        )
    if not completed:
        raise SystemExit("No completed vet invocation available to replay.")
    largest = max(completed, key=lambda record: record["peak_rss_kib"])
    replay = json.loads((report / "tools" / f"{largest['invocation']}.replay.json").read_text())
    print("[VET REPLAY BEGIN] " + json.dumps({
        "package": replay["package"], "original_peak_rss_kib": largest["peak_rss_kib"],
        "hard_cgroup_limit_bytes": 13 * 1024**3,
        "note": "Both replays have the same 13 GiB hard bound; neither executes tests.",
    }), flush=True)
    results = []
    for name, memory_limit in (("budget", "10GiB"), ("default", None)):
        target = report / "replays" / name
        target.mkdir(parents=True)
        unit = f"oom-discovery-vet-{os.getpid()}-{name}"
        started = time.monotonic()
        command = service_command(unit, target, replay, memory_limit)
        process = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
        (target / "command.log").write_text(process.stdout)
        print(process.stdout, flush=True)
        properties = subprocess.run(
            ["sudo", "-n", "systemctl", "show", unit, "--property=Result",
             "--property=ExecMainStatus", "--property=MemoryPeak"],
            capture_output=True, text=True,
        )
        result = {
            "name": name, "package": replay["package"],
            "memory_limit": memory_limit or "unset", "hard_cgroup_limit_bytes": 13 * 1024**3,
            "exit_code": process.returncode,
            "elapsed_seconds": round(time.monotonic() - started, 3),
            "systemd_properties": properties.stdout.strip(),
            "tool_records": [
                json.loads(line) for path in (target / "tools").glob("*.jsonl")
                for line in path.read_text().splitlines()
            ],
        }
        (target / "result.json").write_text(json.dumps(result, indent=2) + "\n")
        print("[VET REPLAY RESULT] " + json.dumps(result), flush=True)
        results.append(result)
        subprocess.run(["sudo", "-n", "systemctl", "stop", unit], capture_output=True)
        subprocess.run(["sudo", "-n", "systemctl", "reset-failed", unit], capture_output=True)
    (report / "replay-results.json").write_text(json.dumps(results, indent=2) + "\n")
    # A measured OOM in the bounded control is an experimental result, not a
    # reason to discard the successful main build's result or stop collection.
    if not all(result["tool_records"] for result in results):
        raise SystemExit("A replay failed before the Go tool started; inspect command.log.")


if __name__ == "__main__":
    main()
