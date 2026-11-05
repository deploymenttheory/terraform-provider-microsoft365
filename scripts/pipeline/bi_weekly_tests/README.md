# Bi-weekly test diagnostics

## Test sharding

The matrix generator inventories top-level tests, fuzz seeds, and examples
from Go source without compiling providers. Packages above the test budget are
split into sorted, named batches. Small packages and large-package batches
share the same capacity-aware Rendezvous distribution across every service and
configuration block type. There are no resource-specific exceptions.

`TESTS_PER_SHARD` sets the maximum discovered test count per job (default: 25).
`ITEMS_PER_SHARD` caps package/batch count per job (default: 5). Both must be
positive. The generator starts with the capacity lower bound and adds shards
when indivisible packages cannot fit. Empty jobs are omitted, and a matrix
above GitHub Actions' 256-job limit fails with a configuration error.

Each job carries a JSON `test_filters` map from package paths to exact,
anchored Go test expressions. The runner applies each filter only to its own
package. Small packages stay unfiltered, as do packages with no discoverable
names, preserving coverage. Filtered packages missing from a shard fail
before execution. No test is deliberately omitted or run in multiple batches.

Assignment is deterministic for the same source inventory, limits, and
`RENDEZVOUS_SEED`. Capacity constraints and batch-boundary changes can move
additional work when the inventory changes. Job counts and compilation work
may increase; test-count budgets do not guarantee runtime because individual
test durations differ. The existing package and job timeouts remain in place.

## Diagnostics

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

## Memory budget

The workflow sets `GOMEMLIMIT=10GiB`, inherited by Go compilation, vet, and test
processes. Default vet remains enabled; package/build concurrency stays at one,
and provider-core retains race detection. This is a [soft per-process Go memory
budget](https://pkg.go.dev/runtime#hdr-Environment_Variables), not a limit on the
whole runner.

[Discovery results](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/1705c7bccf0b3916f39b4d32a2c7d4b67bcda434/scripts/pipeline/oom_discovery/RESULTS.md)
showed provider vet exhausting the approximately 16 GiB ARM64 runner. Both
unbudgeted cold builds failed; six budgeted cold builds passed, including two
provider-core race builds, with at least 4,116 MiB available throughout each.
These were compile-only experiments, not live acceptance executions.

## Result and incident integrity

Each `test-results-*-attempt-N` artifact includes `test-failures.json`,
`test-successes.json`, and `test-run.json`. The runner records package boundaries
in its output so identical test names in different packages remain distinct.
The manifest records whether all selected package commands finished. Passes from
an interrupted or failed package cannot resolve an incident. Skipped tests are
neither failures nor passes.

The summary job reads the current run attempt's job inventory and artifacts.
It always writes result arrays and `merged-status.json`, even when every artifact
is missing. Missing or corrupt result files, an interrupted manifest, failed
artifact downloads, failed jobs, and an artifact count that differs from the
number of executed test jobs disable automatic issue closure. Successful jobs
whose test step was skipped for missing credentials require no artifact and do
not create incidents. A failed jobs API query fails visibly and disables closure.

Only top-level Go tests (`--- FAIL: TestX` at column 0) are incident
identities. Subtest lines and CheckDestroy helper lines, which embed random
resource names, belong to the parent test's failure (and its truncated error
context), so a recurring failure always maps to the same identity.

New test incidents carry a hidden package-and-test identity in their body. An
incident closes only when that exact identity explicitly passes, its package
exits successfully, aggregation is complete, and no matching failure exists.
Legacy incidents without package identity remain open for manual reconciliation;
a title alone cannot establish which package recovered. A subsequent failure
creates a package-specific incident if none exists. Valid partial failures can
still create or update incidents while closure is disabled.

Job incidents preserve cancelled or timed-out step details. OOM classification
requires a positive `memory.events` `oom_kill` delta. Runner-loss annotations
identify lost communication without guessing its underlying cause. Unavailable
logs are recorded as missing evidence. The summary job needs `actions: read` for
[job inventory and logs](https://docs.github.com/en/rest/actions/workflow-jobs)
and `checks: read` for annotations. Its `acceptance-report-attempt-N` artifact
contains the merged arrays, completeness report, job inventory, and classified
job failures.

Run the aggregation and issue processor together (the latter writes to GitHub):

```sh
./scripts/pipeline/bi_weekly_tests/merge_test_results.py \
  test-artifacts merged-failures.json merged-successes.json \
  workflow-jobs.json merged-status.json
./scripts/pipeline/bi_weekly_tests/manage_test_issues.py \
  OWNER REPO RUN_ID merged-failures.json merged-successes.json merged-status.json
```

The regression suite mocks all GitHub writes. It covers package identity,
interruptions, explicit passes, expected skips, missing/corrupt artifacts,
stale run attempts, API pagination/failures, OOM deltas, and issue lifecycle.
