# Acceptance runner investigation — 20 October 2026

This discovery branch tests preparation of acceptance-test binaries without
executing tests or accessing a tenant. Production workflows and issue state are
unchanged. Credential-based skips are expected and are not incidents.

The reproduced failure is a vet memory peak while analysing the full
`internal/provider` dependency in a service test build. Both unbudgeted controls
OOM-killed. All six builds with `GOMEMLIMIT=10GiB` retained vet and completed with
more than 4 GiB sampled available RAM, including two main service builds and two main
Provider Core builds with race instrumentation. Replaying that exact provider
vet invocation on one runner passed with the budget and OOM-killed without it.

## Incident evidence

The [18 October run](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/39089095871)
had 55 failed test jobs. Of these, 47 have downloadable stdout: every one records
`oom_kill` increasing from zero to one, and every last complete process table
before shutdown has `vet` as the largest resident process. Forty-six service jobs
never reached a test start. Provider Core completed earlier packages and failed
while preparing a later package. The other eight jobs report lost runner
communication but have no usable job log; their cause remains unconfirmed.

The runners had approximately 16 GiB RAM, no swap, and no Go memory budget.
Compilation and test parallelism were already limited to one. Representative
samples show vet near 14 GiB RSS, disk usage below 40%, and inode usage around 3%.
For example, [incident #4227](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/4227)
links to [identity/access shard 5/9](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/39089095871/job/117561805298),
whose last sampled vet RSS was 14,709,600 KiB.
The evidence identifies memory exhaustion during preparation, not disk exhaustion
or a failing acceptance assertion. The kernel victim cannot be identified from
an OOM counter alone. See the kernel's
[memory.events documentation](https://docs.kernel.org/admin-guide/cgroup-v2.html#memory-interface-files).

The [14 October run](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/38500096557)
has 53 similar shutdown logs but predates the OOM telemetry. Those failures are
consistent with the same mechanism, not individually proven OOM kills.
[Saved incident evidence](results/incident-runs.json) keeps these distinctions.

## Controlled experiments

All experiments use Go 1.26.0 on `ubuntu-24.04-arm`, a fresh private build cache,
atomic coverage, `GOMAXPROCS=1`, `-p=1`, and `-parallel=1`. The incident source is
`93ce6c1ef8e22c170148e0ed93ca3f23ab74feed`. Validation targets main at
`86b1728a077c49c6bdc3036956831369283100e2`. Source and harness are checked out
separately. At most four discovery runners run concurrently.

The representative service package is
`internal/services/resources/device_management/graph_beta/app_control_for_business_built_in_controls`.
Provider Core validation targets `internal/provider` with race instrumentation.
Every command uses `go test -c`; neither package initializers nor test bodies
execute. A local and CI integration check verifies this with a deliberately
panicking initializer and a deliberate vet error.

The [baseline/no-vet comparison](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/39430977234)
and [memory-budget comparison](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/39435724237)
completed as follows. Values within a cell are repetitions 1 / 2.

| Incident configuration | Outcome | OOM counter increase | Minimum available RAM (MiB) | Build duration (minutes) |
| --- | --- | --- | --- | --- |
| Default vet, no memory budget | Both runner shutdowns | 1 / 1 | 1,758 / 557 | No final timing record |
| `-vet=off`, no memory budget | Both pass | 0 / 0 | 4,304 / 4,618 | 30.76 / 29.63 |
| Default vet, `GOMEMLIMIT=10GiB` | Both pass | 0 / 0 | 4,146 / 4,116 | 45.85 / 49.23 |

The unbudgeted controls' last sampled vet RSS was 14,196,296 and 14,319,184 KiB
(13.5 and 13.7 GiB). These are samples, not a measured final peak. The budgeted
builds completed the provider vet invocation with `wait4` peak RSS of 10,491,944
and 10,492,052 KiB (approximately 10 GiB).

The two baseline and two budgeted runners all used image `20261018.154.1`;
the two vet-disabled runners used `20261012.149.1`. The no-vet comparison alone
therefore has an image-version confound. The budget comparison uses the same
image version, and the isolated replay additionally controls the runner and
retained input files.

The second budgeted job replayed the same provider vet command and config twice
inside a common 13 GiB cgroup, with no swap:

| Replay setting | Result | Measurement |
| --- | --- | --- |
| `GOMEMLIMIT=10GiB` | Exit 0 in 33.08 seconds | Vet peak RSS 10,490,720 KiB |
| Budget unset | `Result=oom-kill` in 25.35 seconds | Cgroup memory peak reached its 13 GiB limit |

The budget replay's systemd `MemoryPeak` was unavailable after successful service
cleanup; its RSS comes from the completed tool's `wait4` record. The failed
replay's memory peak is the cgroup total, not an individual process RSS. The
different measurements are not treated as interchangeable. The bounded OOM
left the runner alive and both replay reports were uploaded.

These same-input results establish that the default-memory provider vet
invocation can exceed the available budget and that the tested Go memory budget
avoids it. No test execution is needed to reproduce the incident mechanism.
The [saved control results](results/run-39430977234.json) and
[budget/replay results](results/run-39435724237.json) include command metadata,
tool attribution, counters, annotations, and image versions.

All four [main validation jobs](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/39444269849)
passed with default vet and `GOMEMLIMIT=10GiB`:

| Main target | Repetition | OOM counter increase | Minimum available RAM (MiB) | Build duration (minutes) |
| --- | --- | --- | --- | --- |
| Representative service | 1 | 0 | 4,158 | 47.86 |
| Representative service | 2 | 0 | 4,171 | 46.19 |
| Provider Core, race enabled | 1 | 0 | 4,204 | 45.65 |
| Provider Core, race enabled | 2 | 0 | 4,172 | 47.87 |

Each target passed on both runner-image versions, `20261012.149.1` and
`20261018.154.1`. Every experiment started with an empty build cache. The main
commit was still `86b1728a077c49c6bdc3036956831369283100e2` when validation
completed. [Saved validation results](results/run-39444269849.json) include the
per-tool measurements; all four meet the predefined no-OOM and 2 GiB headroom
criteria.

Across all six budgeted builds, the lowest sampled available RAM was 4,116 MiB.
The discovery tests, existing test-runner regressions, Linux resource-monitor
tests, and cgroup isolation check passed in CI. Local actionlint, Python syntax,
and whitespace checks also passed.

## Proposed execution fix

Set the following workflow-level environment in
`.github/workflows/bi-weekly-tests.yml`:

```yaml
env:
  GOMEMLIMIT: '10GiB'
```

The existing Python runner inherits this environment into `go test` and its
compiler/vet children. Retain default vet, serial build/test settings, coverage,
Provider Core race instrumentation, and credential-skip behaviour. Keep the
existing resource monitor and assess the next authorised live acceptance run
for OOM counters, available memory, completed tests, and complete results.
Include the explicit-pass incident-closure correction below in that rollout:
the current code can misinterpret a partially recovered run once some tests pass.

The combined `-vet=off` plus memory-budget variant is unnecessary for the
reproduced failure because the memory budget alone passes twice and retains
analysis. Disabling vet remains a tested preparation workaround, but is a larger
change in checking behaviour. Moving vet to lint would require separate memory
validation of that lint job; enabling `govet` in configuration alone is not proof
that it completes reliably.

The runners already had zero swap before cleanup, so merely disabling the
`swap-storage` cleanup option does not restore missing capacity. A longer timeout
also cannot prevent the observed OOM. A larger runner or deliberate swap
allocation could provide capacity, but neither was needed for the successful
incident comparisons and neither was tested here.

## Why a small shard still has a large build

The representative test imports `internal/mocks`, whose provider factories
import `internal/provider`. The provider registers the full service set, so
reducing the number of test packages in a shard does not isolate its dependency
graph. A local listing of the pinned main commit reaches 1,337 package entries,
including test variants: 355 under this provider and 236 under the Microsoft Graph
SDKs. That contextual listing used local Go 1.27.1 on macOS; the actual runner
measurements use Go 1.26.0 on Linux. The broad import graph explains the preparation
cost; the controlled experiments distinguish compiler, linker, and vet memory.

Sharding reduces how many test packages execute in a job but does not remove this
shared dependency graph. Even one package with serial compilation can exceed the
runner's RAM inside one vet process; reducing matrix concurrency does not change
that process's memory requirement. A longer-term improvement is to separate
lightweight HTTP mocks from provider factories and use narrower provider registrations where
test coverage permits. Acceptance coverage that depends on the complete provider
must remain covered; this architectural change needs its own validation.

The first vet-disabled build completed with no OOM and at least 4,304 MiB sampled
available memory. Its largest completed compiler invocation was the beta SDK
models package at about 10 GiB RSS. Compilation therefore remains expensive even
when the vet-related failure is removed.

## Reporting defects and proposed changes

1. **Missing artifacts cause a second failure.** In the 18 October summary,
   artifact download left no `test-artifacts` directory. The merge script exited
   with “Artifacts directory not found,” then issue processing failed because
   `merged-failures.json` did not exist. Create the result directory before
   download, always produce machine-readable merge outputs, and carry an explicit
   completeness status alongside them. Missing/corrupt results must be recorded
   as incomplete, not interpreted as a clean test run.
2. **Runner termination loses its classification.** A failed job with a cancelled
   test step becomes `Unknown` because `detect_job_failures.py` searches for
   `failure` steps. Preserve cancelled/timed-out step details and use available
   OOM counters, runner annotations, and log evidence to classify failures. Only
   classify OOM when supported; lost communication alone remains runner failure
   with unknown cause. Preserve explicit manual cancellation and credential skips.
3. **Missing results can close unrelated incidents.** Once any success exists,
   `manage_test_issues.py` closes open incidents for tests absent from failures,
   even if those tests never ran. Require an explicit pass for the same package
   and test identity, no contradictory failure, and trustworthy result
   completeness before automatic closure. Absence, skips, cancelled execution,
   and missing artifacts must preserve the open incident.

[Local probes](probe_reporting.py) reproduce all three behaviours with external
calls blocked. Their [saved output](results/reporting-probes.json) records one
closure attempt for a test that never ran, an `Unknown` classification for a
cancelled test step, and failure to create merged output for missing artifacts.
The probe also confirms that a successful job with a skipped test step is not
classified as an incident.

Run the probes without contacting GitHub:

```sh
python3 scripts/pipeline/oom_discovery/probe_reporting.py \
  --output /tmp/reporting-probes.json
```

These are separate reporting defects. The closure behaviour was reproduced
synthetically; this investigation does not claim it closed a particular historic
incident. Fixing them does not reduce compiler or vet memory.

For a production change, the reporting regression cases should include:

| Input | Required behaviour |
| --- | --- |
| No artifacts or a missing expected shard | Emit valid aggregate files, mark the run incomplete, preserve affected incidents |
| Corrupt result JSON | Preserve available evidence and explicitly report the corrupt input |
| Failed job, cancelled test step, positive OOM-counter delta | Identify the cancelled step and report confirmed OOM |
| Lost communication without an OOM counter | Report runner failure with unknown cause |
| Explicit cancellation or credential-based skip | Preserve that status; do not invent a test failure or pass |
| One passing test and an unrelated test with no result | Leave the unrelated incident open |
| Same test name in different packages | Match package plus test identity, not name alone |
| Explicit matching pass with complete results and no failure | Permit closure |

Stable package/test identity should be included in both test-result records and
structured incident metadata. The current result records carry test name,
configuration block type, and service, but not the precise Go package.
Older incidents without an unambiguous identity should require reconciliation
before automatic closure. An empty failures array is not sufficient evidence of
recovery. Explicit credential skips account for intentionally unexecuted work;
they are neither incomplete execution failures nor passing test evidence.

## Interpretation and rollout limits

GitHub supplied different runner-image versions across jobs; image versions are
recorded in results. Source, Go version, flags, runner class, and initial cache
state are controlled, but a historical hosted image cannot be restored by date.
The bounded replay compares memory settings on the same runner and build inputs.

Go [buffers tool stderr until an invocation exits](https://github.com/golang/go/blob/go1.26.0/src/cmd/go/internal/work/shell.go).
The first two phases therefore can lose an unfinished tool's package record
when the VM shuts down. Later
instrumentation streams tool records through a separate parent relay. A missing
unfinished record in the early logs does not mean there was no active tool.

[GOMEMLIMIT is a soft per-process budget](https://go.dev/doc/gc-guide#Memory_limit),
not a hard VM limit. Two-second memory samples can miss brief peaks.
`go test -c` retains debug information that ordinary
test execution may omit, so linker costs are not necessarily identical.
Successful preparation does not establish successful live acceptance behaviour.
This investigation identifies the failing operation and a successful memory
setting; it does not identify an introducing commit or establish a Go/SDK
regression. Workflow conclusions alone cannot establish that history.
See the [harness method](README.md) for the commands and acceptance criteria.
