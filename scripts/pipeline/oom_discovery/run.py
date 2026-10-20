#!/usr/bin/env python3
"""Build-only comparison of Go vet and compiler memory on hosted runners."""

import argparse
import json
import os
from pathlib import Path
import subprocess
import sys
import time

from stream_tools import ToolEventStream

VARIANTS = ("baseline", "no-vet", "memory-budget", "combined")


def build_command(package, report, race, variant):
    wrapper = Path(__file__).with_name("tool_exec.py").resolve()
    command = [
        "go", "test", "-c", "-v", "-work", "-p=1", "-parallel=1",
        "-covermode=atomic", f"-coverprofile={report / 'coverage.txt'}",
        f"-o={report / 'discovery.test'}", f"-toolexec={wrapper}",
    ]
    if race:
        command.append("-race")
    if variant in ("no-vet", "combined"):
        command.append("-vet=off")
    return [*command, package]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--name", required=True)
    parser.add_argument("--variant", choices=VARIANTS, required=True)
    parser.add_argument("--package", required=True)
    parser.add_argument("--report", type=Path, required=True)
    parser.add_argument("--race", action="store_true")
    args = parser.parse_args()
    report = args.report.resolve()
    report.mkdir(parents=True, exist_ok=True)
    env = dict(os.environ)
    # These runs must not acquire tenant access or inherit acceptance settings.
    for key in list(env):
        if key.startswith("M365_"):
            del env[key]
    env.update(TF_ACC="0", GOMAXPROCS="1", GOFLAGS="-p=1",
               PYTHONUNBUFFERED="1", RESOURCE_MONITOR_INTERVAL_SECONDS="1",
               OOM_DISCOVERY_REPORT=str(report))
    event_log = report / "tool-events.jsonl"
    env["OOM_DISCOVERY_EVENT_LOG"] = str(event_log)
    # The existing monitor samples vmstat for one second, then sleeps one
    # second: approximately two seconds between resource snapshots.
    env.pop("GOGC", None)
    env.pop("GOMEMLIMIT", None)
    if args.variant in ("memory-budget", "combined"):
        env["GOMEMLIMIT"] = "10GiB"
    cache = Path(env["RUNNER_TEMP"]) / f"oom-build-cache-{args.name}"
    cache.mkdir()
    env["GOCACHE"] = str(cache)
    command = build_command(args.package, report, args.race, args.variant)
    metadata = {
        "name": args.name, "variant": args.variant, "package": args.package,
        "race": args.race, "command": command,
        "source_sha": subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip(),
        "go_version": subprocess.check_output(["go", "version"], text=True).strip(),
        "memory_limit": env.get("GOMEMLIMIT", "unset"),
        "cache": str(cache), "initial_cache_empty": not any(cache.iterdir()),
        "image_version": env.get("ImageVersion"),
        "run_id": env.get("GITHUB_RUN_ID"),
        "run_attempt": env.get("GITHUB_RUN_ATTEMPT"),
        "tool_events_streamed": True,
        "started_at": time.time(),
    }
    (report / "experiment.json").write_text(json.dumps(metadata, indent=2) + "\n")
    print("[EXPERIMENT] " + json.dumps(metadata), flush=True)
    monitor = Path(__file__).resolve().parents[1] / "pr/steps/monitor_resources.sh"
    started = time.monotonic()
    with ToolEventStream(event_log, sys.stdout):
        result = subprocess.run(["bash", str(monitor), str(report), *command], env=env)
    metadata.update(exit_code=result.returncode, elapsed_seconds=round(time.monotonic() - started, 3))
    (report / "experiment.json").write_text(json.dumps(metadata, indent=2) + "\n")
    print("[EXPERIMENT RESULT] " + json.dumps(metadata), flush=True)
    raise SystemExit(result.returncode if result.returncode >= 0 else 128 - result.returncode)


if __name__ == "__main__":
    main()
