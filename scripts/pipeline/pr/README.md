# PR Go checks

Lint and unit tests are separate workflows. Each workflow compares the checked-out commit with the Git
merge-base of the PR base SHA and selects exact package directories. Root files,
renames and deletions are included; completely deleted packages are skipped.
Planning and factory validation also run when their pipeline files change. An
empty Go selection launches no lint, unit, or race workers. There is no `./...`
fallback. Unit tests include local packages whose test import graph depends on a
changed production package. Test-only changes do not propagate as production
dependencies. Internal fixture and embedded-data changes select their owning Go
packages and dependent tests. Unit workers select packages containing test files;
production-only packages remain in lint analysis and are compiled through their
test consumers. Lint findings remain restricted to changed Go packages and lines.

Lint, unit, and race workers use the acceptance pipeline's bounded Rendezvous
sharder, assigning each selected package exactly once, with at most five packages
per worker by default. Workers use standard `ubuntu-24.04-arm` runners and execute
their packages serially. A failed shard does not cancel other shards. The existing
lint and unit summary checks require every applicable worker to pass.

The production `internal/provider` registry necessarily imports every object kind.
Its lint job runs alone as `✨ Lint Provider Registry`, with `GOMEMLIMIT=20GiB` and
`GOGC=100` to accommodate its measured live analysis data and disk-backed headroom. Go's memory limit is a soft
heap target, not an RSS cap; live analysis data can exceed it. The dedicated job
runs all 23 configured checks with golangci-lint v2.15.0. It also adds 8 GiB of
swap on the ephemeral runner, after requiring 10 GiB of free disk. This provides
disk-backed headroom without changing the runner tier. Ordinary partitions do
not allocate extra swap. Swapping can increase execution time; the existing
timeouts, failure propagation and OOM diagnostics remain enforced. Its success
must be established from the standard runner's diagnostics.
Ordinary package workers keep the default Go garbage-collection settings. The
diagnostics record both settings alongside memory peaks and OOM counters.
When the production registry itself is selected, one standard-runner job compiles
the complete production provider before lint workers start with `go build ./internal/provider`. Its successful post-step
saves the shared Go cache before workers restore it. The compilation and lint
workers use the same `go.mod`, `go.sum` and production-registry cache key; unit
workers use a separate `go.mod`/`go.sum` key so an earlier test worker cannot win
and save a partial production snapshot. Go still validates source hashes and
rebuilds changed packages. Compilation failures fail the required lint gate.
Leaf-package-only lint runs skip this full-provider compilation and use the
ordinary dependency cache, so they cannot populate a partial production snapshot.
Cold compilation and lint analysis each retain their existing 45-minute command
budgets in separate 60-minute jobs; no lint checks are removed.

Manual lint dispatches may set `packages` to a space-separated subset of changed
packages, for example `internal/provider`; PR runs continue checking every
changed package and all configured linters.

Unit and race workers install the latest stable Terraform once, without the CLI
wrapper, and export its path as `TF_ACC_TERRAFORM_PATH` for every test. Tests do
not repeatedly download Terraform while running. In GitHub Actions, test stdout
is streamed in green and failures in red, with a colour reset after every line.
Resource utilisation and system diagnostics remain uncoloured. Set `NO_COLOR`
to disable test colouring.

Lint reports findings on changed lines, includes test-file analysis, retains the
configured linters/formatters, and does not rewrite source. It runs one package
per process using the pinned golangci-lint release. Unit tests run all tests in
the selected packages with `TF_ACC=0` for the project's mocks and explicitly skip
`TestAcc*` (HashiCorp treats any nonempty `TF_ACC` as enabled). Failures are
independent of coverage. Unit execution stops at the first failed or terminated
package and preserves its exit status. Per-package coverage remains available
for diagnostics; merged coverage is produced only after all selected packages pass.
Workers key the shared Go cache on both `go.mod` and `go.sum`. This separates
these compile-heavy jobs from the older `go.sum`-only cache populated before the
SDK build completed, so a successful worker can persist the compiled dependencies.
Cache reuse never skips executing the selected lint checks or tests.

Package/test parallelism and `GOMAXPROCS` are one. Successful unit profiles are
downloaded and merged before the single Codecov upload, coverage gate, and PR
comment. Duplicate coverage blocks are counted once when combining profiles.

## Isolated provider factories

`internal/provider/core` contains the real provider schema, validation and client
configuration. The production provider composes this core with its complete
resource, data source, action, list, and ephemeral registries. Test factories
compose the same core with only the package under test and fixture dependencies.
The shared mocks no longer import the complete production provider.

Each participating test package has a generated `provider_factories_test.go`.
Its local unit and acceptance maps preserve their respective authentication
modes. Test steps, assertions, and lifecycle ordering remain in their existing
files. Generate factories after adding a registration, changing test imports, or
adding a Terraform fixture dependency:

```sh
python3 scripts/pipeline/generate_test_factories.py
python3 scripts/pipeline/generate_test_factories.py --check
```

The generator reads production constructors and type names, package imports, and
Terraform block declarations in Go and Terraform fixtures. Resource and data
source types are resolved separately even when their names match.
List-resource factories also include the matching managed resource required by
the protocol to retrieve its schema and identity. CI fails when
checked-in factories are stale or production type metadata cannot be resolved.
Dependencies assembled by a helper must be represented by an explicit test
import of their registered package or a Terraform declaration in the package's
fixtures. Generated files must not be edited manually.

## Manual fallback

GitHub's Go path trigger can miss changes beyond its 300-file diff limit. Both
workflows support manual dispatch against a selected branch with a `base_ref`
input (default `origin/main`). The input selects the comparison baseline, not a
request to run the whole module. Use an earlier commit to validate a Go change
already present on the selected branch. Manual test runs omit PR-specific Codecov
upload/enforcement and comments; they still fail on test errors.

Local equivalents, from the repository root:

```sh
python3 scripts/pipeline/pr/steps/detect_changes.py --base-ref origin/main
python3 scripts/pipeline/pr/steps/plan_jobs.py --base-ref origin/main
python3 scripts/pipeline/pr/steps/run_lint.py --base-ref origin/main
# Use the exact packages printed by the selector, separated by spaces:
python3 scripts/pipeline/pr/steps/run_tests.py --packages 'internal/client internal/provider'
```

Selection compares committed changes. Commit work intended for comparison first.
Documentation-only changes do not select Go packages. Scheduled acceptance jobs
retain their existing test batches and use the isolated factories as well.

## Diagnostics

Linux CI wraps commands with `monitor_resources.sh`: output streams directly to
the Actions log and an artifact, `/usr/bin/time -v` records time/peak process RSS,
and memory/disk/process-name samples are emitted every 15 seconds. Reports do not
include environment dumps or process arguments. Lint also writes per-package
SARIF and `results.json` with findings versus execution exit statuses. A runner
shutdown may prevent artifact upload; periodic samples remain in the live log.
Lint has a 45-minute per-package timeout. Lint and test steps have a 50-minute
budget within 60-minute jobs, leaving time for failure reports and cache post-steps.
Memory targets and GC overrides are restricted to the dedicated registry worker
and recorded in diagnostics. They never convert a failed execution into success
or disable checks. A resource failure must be investigated from the measurements.

### Initial validation evidence

The [no-Go manual run](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/36590453502)
passed and skipped Go setup, unit tests and race detection. The
[scoped core lint run](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/36590157722)
selected only `internal/client` and `internal/provider`. Client lint passed in
15m33s. During provider lint, sampled linter RSS reached 14.45 GiB on the runner
with 15,947 MiB RAM; its 3,071 MiB swap was almost completely consumed. Disk usage
was 37%. Execution ended with exit 143 before final reports could upload.

These measurements establish memory pressure, not the sender of the termination
signal. They also show that selecting fewer packages alone does not guarantee
that provider analysis fits on this runner. Capacity or GC changes require a
separate measured validation; do not interpret package-selection correctness as
proof that the large provider lint workload succeeds.

The [single-resource unit run](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/36597296286)
passed all 23 B2B management policy unit cases with 79.7% statement coverage and
executed no `TestAcc*` cases. Tests took 76 seconds, while the complete command
including the cold build took 45m33s. Go's build cache was saved successfully.
The [first GC experiment](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/36596573217)
used `GOGC=50` on a temporary branch but hit its earlier 30-minute timeout during
package loading. Analysis had not started, so that result does not establish
whether the setting resolves analysis memory pressure. Its diagnostic artifact
correctly records linter exit 4 as an execution failure.

The [retry with a restored build cache and corrected timeout](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/36603104015)
loaded packages in 1m03s, but analysis again filled RAM and swap and ended with
exit 143. `GOGC=50` therefore did not resolve this workload, and is not included
in the implementation. These provider/resource lint workloads still fail on the
16 GiB runner. Further work should measure a runner with more memory or reduce
the shared analysis dependency graph while retaining the required checks.

## Pipeline regression checks

Requires Python 3, Git, Go compatible with `go.mod`, golangci-lint v2.15.0 and
`actionlint`. Tests use isolated repositories and tiny Go modules, not the provider
suite or a tenant. The intentionally failing fixture tests must be caught by the
Python assertions.

```sh
GOMAXPROCS=1 GOFLAGS=-p=1 python3 -m unittest discover -s scripts/pipeline/pr/tests -v
actionlint .github/workflows/go-lint.yml .github/workflows/pr-tests.yml
golangci-lint config verify
bash -n scripts/pipeline/pr/steps/monitor_resources.sh
```
