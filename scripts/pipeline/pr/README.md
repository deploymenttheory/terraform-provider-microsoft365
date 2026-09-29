# PR Go checks

Lint and unit tests are separate workflows. Automatic runs are triggered only by
changed `.go` files. Each workflow compares the checked-out commit with the Git
merge-base of the PR base SHA and selects exact package directories. Root files,
renames and deletions are included; completely deleted packages are skipped.
An empty selection skips Go installation and execution. There is no `./...`
fallback and no dependency-test expansion.

Lint reports findings on changed lines, includes test-file analysis, retains the
configured linters/formatters, and does not rewrite source. It runs one package
per process using the pinned golangci-lint release. Unit tests run all tests in
the selected packages with `TF_ACC=0` for the project's mocks and explicitly skip
`TestAcc*` (HashiCorp treats any nonempty `TF_ACC` as enabled). Failures are
independent of coverage.
Package/test parallelism and `GOMAXPROCS` are one. Resource tests can still compile
much of the provider through the shared mock provider factory.

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
python3 scripts/pipeline/pr/steps/run_lint.py --base-ref origin/main
# Use the exact packages printed by the selector, separated by spaces:
python3 scripts/pipeline/pr/steps/run_tests.py --packages 'internal/client internal/provider'
```

Selection compares committed changes. Commit work intended for comparison first.
Non-Go fixtures, documentation, dependency manifests and workflow-only changes do
not automatically run these Go checks. Scheduled acceptance tests are unchanged.

## Diagnostics

Linux CI wraps commands with `monitor_resources.sh`: output streams directly to
the Actions log and an artifact, `/usr/bin/time -v` records time/peak process RSS,
and memory/disk/process-name samples are emitted every 15 seconds. Reports do not
include environment dumps or process arguments. Lint also writes per-package
SARIF and `results.json` with findings versus execution exit statuses. A runner
shutdown may prevent artifact upload; periodic samples remain in the live log.
Lint has a 45-minute per-package timeout. Lint and test steps have a 50-minute
budget within 60-minute jobs, leaving time for failure reports and cache post-steps.
No speculative memory cap, GC override, or disabled linter is used to conceal a
failed execution. A resource failure must be investigated from the measurements.

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

## Pipeline regression checks

Requires Python 3, Git, Go compatible with `go.mod`, golangci-lint v2.14.0 and
`actionlint`. Tests use isolated repositories and tiny Go modules, not the provider
suite or a tenant. The intentionally failing fixture tests must be caught by the
Python assertions.

```sh
GOMAXPROCS=1 GOFLAGS=-p=1 python3 -m unittest discover -s scripts/pipeline/pr/tests -v
actionlint .github/workflows/go-lint.yml .github/workflows/pr-tests.yml
golangci-lint config verify
bash -n scripts/pipeline/pr/steps/monitor_resources.sh
```
