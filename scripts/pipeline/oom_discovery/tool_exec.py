#!/usr/bin/env python3
"""Attribute Go tool memory to packages; never execute a test binary."""

import json
import os
from pathlib import Path
import shutil
import signal
import subprocess
import sys
import time


def emit(record, path):
    line = json.dumps(record, sort_keys=True)
    with path.open("a", encoding="utf-8") as stream:
        stream.write(line + "\n")
    print("[GO TOOL] " + line, file=sys.stderr, flush=True)


def main():
    argv = sys.argv[1:]
    if not argv:
        raise SystemExit("usage: tool_exec.py <go-tool> [arguments...]")
    # Go uses these responses to calculate tool/cache identities. Do not add
    # diagnostics to them, or claim that instrumentation changes Go's version.
    if "-V=full" in argv:
        os.execv(argv[0], argv)

    directory = Path(os.environ["OOM_DISCOVERY_REPORT"]) / "tools"
    directory.mkdir(parents=True, exist_ok=True)
    invocation = f"{os.getpid()}-{time.time_ns()}"
    path = directory / f"{invocation}.jsonl"
    tool = Path(argv[0]).name
    record = {
        "invocation": invocation,
        "tool": tool,
        "package": os.environ.get("TOOLEXEC_IMPORTPATH", ""),
        "wrapper_pid": os.getpid(),
        "memory_limit": os.environ.get("GOMEMLIMIT", "unset"),
    }
    if tool == "vet":
        # Config and replay arguments contain build paths, not credentials.
        # Retain the work directory with go test -work for same-run replay.
        configs = [Path(arg) for arg in argv[1:] if arg.endswith(".cfg")]
        if configs and configs[-1].is_file():
            config = json.loads(configs[-1].read_text())
            record["vet_id"] = config.get("ID")
            record["vet_import_path"] = config.get("ImportPath")
            shutil.copyfile(configs[-1], directory / f"{invocation}.vet-config.json")
            (directory / f"{invocation}.replay.json").write_text(
                json.dumps({"argv": argv, "cwd": os.getcwd(), "package": record["package"]})
            )

    started = time.monotonic()
    process = subprocess.Popen(argv)
    record["pid"] = process.pid
    emit({**record, "event": "start", "time": time.time()}, path)
    _, status, usage = os.wait4(process.pid, 0)
    process.returncode = os.waitstatus_to_exitcode(status)
    # Linux reports KiB; macOS reports bytes (useful for local harness tests).
    peak_kib = usage.ru_maxrss / 1024 if sys.platform == "darwin" else usage.ru_maxrss
    emit({
        **record,
        "event": "end",
        "time": time.time(),
        "elapsed_seconds": round(time.monotonic() - started, 3),
        "peak_rss_kib": int(peak_kib),
        "user_seconds": usage.ru_utime,
        "system_seconds": usage.ru_stime,
        "exit_code": process.returncode,
    }, path)
    if process.returncode < 0:
        signum = -process.returncode
        signal.signal(signum, signal.SIG_DFL)
        os.kill(os.getpid(), signum)
    raise SystemExit(process.returncode)


if __name__ == "__main__":
    main()
