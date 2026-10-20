# Acceptance OOM discovery

This branch investigates the runner failures from the 18 October 2026 acceptance
run. It builds test binaries without executing them. No tenant credentials,
acceptance operations, issue creation, or issue closure are involved.

## Running experiments

The branch-only `Acceptance OOM discovery` workflow reads `experiments.json`.
A push affecting the workflow, matrix, or Python harness starts the configured matrix.
The first two phases used two runners each and overlapped within a four-runner
limit. After both completed, validation uses four independent runners in one
phase. Further phases serialize. Each job has a 75-minute limit.
Change the matrix before pushing the next phase; do not accidentally repeat a
completed phase by pushing harness changes alone. Documentation and saved results
do not trigger experiments.

Each entry specifies an immutable source commit, package, race setting, and one
of these variants:

| Variant | Vet | GOMEMLIMIT |
| --- | --- | --- |
| `baseline` | Go test defaults | Unset |
| `no-vet` | Disabled | Unset |
| `memory-budget` | Go test defaults | `10GiB` |
| `combined` | Disabled | `10GiB` |

All variants use the same instrumented `go test -c`, atomic coverage, `-p=1`,
`-parallel=1`, and `GOMAXPROCS=1`. Each experiment starts with an empty, private
build cache. Go is selected from the examined commit's `go.mod`. The source
checkout and instrumentation checkout are separate.

`tool_exec.py` records tool/package/PID attribution and per-invocation peak RSS
using `wait4`. Vet configs and same-run replay arguments are retained. Tool
version probes are passed through unchanged. From the validation phase onward,
a separate parent relay streams these records, bypassing Go's buffered tool
stderr. In the first two phases a killed invocation's start record can be lost
with the runner, so absent unfinished-tool records do not imply no active tool.
The existing resource monitor
samples approximately every two seconds and streams cgroup OOM counters,
available memory, CPU/I/O pressure, swap, disk and process RSS.

The first phase repeats baseline and vet-disabled builds on the incident commit
twice. Subsequent phases test the memory budget, then validate the successful
configuration against current main and Provider Core with race enabled. The
validation matrix repeats both targets twice with `GOMEMLIMIT=10GiB` and default
vet, following two successful incident builds and the same-run vet replay.

The second memory-budget job also replays its largest completed vet invocation
using retained build inputs, first with `GOMEMLIMIT=10GiB`, then with it unset.
Both replays run in the same 13 GiB hard cgroup limit, with swap disabled, as the
unprivileged runner user. This bounds the control's OOM to the replay service.
Each replay has a five-minute runtime limit. Replay metrics are kept separate
from the preceding build, and the build artifact is uploaded before replay.

## Collecting results

After a run completes:

```sh
python3 scripts/pipeline/oom_discovery/collect.py RUN_ID \
  --output /tmp/oom-discovery-RUN_ID.json
```

The collector uses read-only GitHub API calls, parses stdout, and records check
annotations. It works when the runner terminated before uploading artifacts.
Missing logs remain missing evidence. `candidate_pass` requires a successful job
and command, no OOM-counter increase, and at least 2 GiB sampled available memory.

Artifacts contain raw measurements and summaries; the large compiled test binary
is excluded. The build work directory remains on the runner for same-run vet
replay but is not uploaded. Never execute replay paths from an untrusted artifact.

## Interpretation limits

- `GOMEMLIMIT` is a soft, per-Go-process budget, not a hard limit on the VM.
- An OOM-counter increase proves an OOM kill; exit 143 or lost contact alone does
  not. An unfinished tool record identifies activity, not the kernel's victim.
- Two-second samples can miss transient peaks. Tool RSS measurements complement
  them but do not establish a hard upper bound for the whole process tree.
- Fresh builds and `go test -c` isolate preparation. Compile-only mode retains
  binary debug information that ordinary test execution may omit, so linker
  costs should not be assumed identical to the acceptance workflow.
- The incident source and toolchain can be pinned, but GitHub runner images
  cannot be restored by date. Image versions are recorded for comparison.
- Passing these experiments does not establish that live acceptance tests pass.
  Expected credential skips are outside this investigation.
