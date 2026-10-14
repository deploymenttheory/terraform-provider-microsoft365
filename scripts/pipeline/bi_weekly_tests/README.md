# Bi-weekly test diagnostics

Every test job streams resource samples to the test step's stdout about every
15 seconds, including while Go compiles without producing test output. Python
output is unbuffered so package discovery, the current package, and command PID
appear immediately. No additional debug setting is required.

Samples include CPU usage and I/O wait (`vmstat`), load, available memory and
swap, disk space and inodes, the largest processes by RSS, Linux pressure
metrics, and cgroup memory limits, peaks, OOM counters, and CPU throttling.
Cgroup ancestors are included because a parent may impose the effective limit.
Process samples use executable names, without arguments or environment values
that could contain credentials.

The `test-results-*` artifact for each job includes a `diagnostics/` directory:

- `environment.log`: initial runner capacity and tool versions.
- `resources.log`: timestamped samples and a final sample, with kernel OOM
  evidence when the runner permits access to its kernel log.
- `command.log`: the complete streamed Python/test runner output.
- `time.log`: GNU time's elapsed time, CPU time, and peak process RSS.
- `exit-status.log`: the command exit status.
- `termination.log`: signal and timestamp if the monitor receives SIGTERM/SIGINT.

Start with the last package/PID in `command.log` and the resource samples around
that timestamp. Compare `memory.events` against the initial snapshot: an increase
in `oom_kill` provides evidence of a cgroup OOM kill. Low available memory, swap
activity, sustained memory pressure, disk exhaustion, and CPU throttling help
distinguish resource constraints. Peak process RSS is not the sum of all
concurrent processes; use the host/cgroup samples for total memory consumption.

Exit 143 means termination by SIGTERM. It does not, by itself, identify OOM.
Check the job log for runner shutdown, cancellation, or timeout messages. A Go
command terminated by a signal stops further package execution and preserves
the signal exit status and partial test reports.

The upload steps run with `if: always()`, but abrupt runner shutdown or SIGKILL
can prevent cleanup and artifact uploads. Periodic stdout samples provide the
evidence already streamed before the runner disappears.

Run diagnostics checks without cloud credentials or provider builds:

```sh
python3 -m unittest discover -s scripts/pipeline/pr/tests -p test_monitor_resources.py -v
python3 -m unittest discover -s scripts/pipeline/bi_weekly_tests/tests -v
```

The monitor checks require Linux with Bash, GNU time, util-linux, and procps,
as provided by the hosted Ubuntu runners. `RESOURCE_MONITOR_INTERVAL_SECONDS`
can set a different positive whole-number sampling interval (default: 15).
