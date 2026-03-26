# ADR-002: Eventual Consistency Handling for Microsoft Entra API Reads

## Status

Accepted

## Date

2026-03-26

## Context

Microsoft Entra (and the Microsoft Graph API that exposes it) operates on a multi-region,
multi-replica directory architecture designed for scale and availability. Writes are accepted
immediately and acknowledged with a 2xx response, but replication across replicas is
asynchronous. This means that a read issued immediately after a successful write may be served
from a replica that has not yet received the change — returning pre-write stale data rather
than an error.

For a Terraform provider this is particularly problematic. After a Create or Update, the
provider calls Read to reconcile state. If that read is served from an unconverged replica,
Terraform sees a mismatch between what was planned and what was returned, and reports it as
"Provider produced inconsistent result after apply". This fails the apply and can leave
dangling resources in the tenant if the cleanup destroy also fails.

The problem is not limited to a single resource type. Microsoft's own guidance explicitly
states that application-only access (the access model used by this provider's service
principals) does not guarantee read-after-write consistency by design. Any resource that
writes to Entra and immediately reads back could exhibit this behaviour.

The initial workaround applied to `microsoft365_graph_beta_users_user_license_assignment`
was a fixed `time.Sleep(15 * time.Second)` inserted before `ReadWithRetry`. This violated
Microsoft's own recommendation against fixed sleep delays, paid the cost unconditionally
regardless of whether the replica had already converged, and did not generalise to other
resources or other propagation windows.

## Decision Drivers

* Fixed sleeps are an anti-pattern per Microsoft's eventual consistency guidance — they are
  either too short (edge cases still fail) or too long (unnecessary delay on every apply).
* The existing `ReadWithRetry` mechanism already handles transient API errors via retry
  loops but had no mechanism to re-attempt a read that returned 200 with stale data.
* CRUD functions across the provider are deliberately uniform and clean; consistency
  handling logic should not be scattered through individual Create/Update methods.
* Only resources that are known to have eventual consistency exposure should pay any
  retry cost — resources that are not affected must not be slowed down.
* The approach must be extensible: new resources with eventual consistency characteristics
  should be able to opt in without modifying central infrastructure.

## Considered Options

* **Fixed sleep before ReadWithRetry** — insert `time.Sleep` in each affected resource's
  Create and Update before calling `ReadWithRetry`. Simple but violates Microsoft's guidance,
  unconditionally delays all applies, and does not adapt to actual propagation time.

* **Hard-coded delay inside ReadWithRetry** — add a mandatory initial sleep inside
  `ReadWithRetry` itself, applied to every resource. Centralises the behaviour but imposes
  the cost universally and still does not adapt.

* **Caller-supplied consistency predicate in ReadWithRetry** — extend `ReadWithRetryOptions`
  with an optional `ConsistencyPredicate func(ctx context.Context, state tfsdk.State) bool`.
  When provided, `ReadWithRetry` continues retrying even on a successful (no-error) read
  until the predicate confirms the state is consistent with what was written, or the context
  deadline is reached. Resources with no eventual consistency exposure pass `nil` and retain
  the existing fast-exit behaviour.

## Decision

Chosen option: "Caller-supplied consistency predicate in ReadWithRetry", because it
implements the polling-with-retry pattern recommended by Microsoft, is adaptive to actual
propagation time, is opt-in per resource (zero cost to unaffected resources), and keeps
CRUD functions clean by moving consistency logic into a dedicated `predicate.go` file within
each affected resource package.

## Rationale

Microsoft's guidance on designing for Entra eventual consistency recommends:

1. Trust successful write responses — do not re-read to confirm the write succeeded.
2. Cache identifiers from write responses rather than re-reading objects immediately.
3. When reads are unavoidable, poll with retry, treating stale/not-found responses as
   transient rather than permanent failures.
4. Make retries idempotent.
5. Avoid fixed sleep delays.

The consistency predicate approach satisfies all five points. The write is trusted (the
provider proceeds after a 2xx). The retry loop in `ReadWithRetry` already handles transient
API errors; extending it to also retry on stale-but-successful reads unifies both failure
modes in one place. The predicate is idempotent — it only inspects state, never mutates.
The retry interval reuses the existing `RetryInterval` configured in `ReadWithRetryOptions`
(default 2 seconds), and retries are bounded by the context deadline rather than an
arbitrary sleep.

Placing the predicate in a `predicate.go` file within the resource package keeps all
eventual consistency knowledge co-located with the resource itself. The predicate has full
access to the resource's model struct and can therefore verify the complete read state
against the full expected state captured at write time — not just a single field.

## Consequences

### Positive

* Eliminates fixed sleeps — applies complete as soon as the replica converges, not after
  an arbitrary wait.
* Generalises cleanly — any resource can opt in by adding a `predicate.go` and setting
  `opts.ConsistencyPredicate` in Create and Update.
* Zero cost to unaffected resources — `nil` predicate retains the existing fast-exit
  behaviour with no change to call sites.
* Aligns with Microsoft's published guidance for Entra eventual consistency.
* Consistency logic is self-contained in `predicate.go` per resource, keeping CRUD
  functions uniform across the provider.
* `ReadWithRetry` remains the single location for all post-write read retry logic,
  covering both API errors and eventual consistency in one loop.

### Negative

* None identified that are not better classified as scope limitations or handled by
  existing Terraform reconciliation.

### Neutral

* The predicate applies at write time only. During a standalone refresh or import there
  is no preceding write, so no eventual consistency window exists to guard against.
  This is a design scope boundary, not a gap — by the time a user performs an import
  or refresh, any prior write will have long since propagated.
* Count-based comparison for set fields (e.g. `disabled_plans`) is used as a propagation
  gate, not a value validator. It reliably distinguishes "SKU not yet visible in
  `assignedLicenses`" (0 elements, pre-propagation default) from "SKU visible" (N elements).
  Exact element-wise value correctness is enforced by Terraform's own plan reconciliation
  once the predicate passes — no duplication of that logic is needed in the predicate.
* Each resource that adopts this pattern requires a `predicate.go` file. This is a small
  overhead but makes the contract explicit and discoverable.
* The retry interval for consistency retries is the same as for error retries (default 2s).
  These could be decoupled in future if different resources have very different propagation
  characteristics.

## Implementation

The change consists of three parts:

**1. `internal/services/common/crud/read_with_retry.go`**

`ReadWithRetryOptions` gains a new field:

```go
ConsistencyPredicate func(ctx context.Context, state tfsdk.State) bool
```

The retry loop is updated so that a successful (no-error) read does not immediately
return if the predicate is non-nil and returns false. Instead it logs the inconsistency,
waits `RetryInterval`, and retries — identical behaviour to a retryable API error, but
triggered by the predicate rather than an HTTP status code.

**2. `internal/services/resources/users/graph_beta/license_assignment/predicate.go`**

Contains `licenseAssignmentConsistencyPredicate`, which takes the full
`*UserLicenseAssignmentResourceModel` written at Create/Update time and returns a closure
that verifies the complete read state:

- `user_principal_name` is populated (Computed field, confirms Read executed against API)
- `id` is non-empty (composite `userId_skuId` key resolved)
- `user_id` matches expected
- `sku_id` matches expected
- `disabled_plans` element count matches expected (primary eventual-consistency signal:
  defaults to empty set when the SKU is not yet visible in `assignedLicenses`)

**3. `internal/services/resources/users/graph_beta/license_assignment/crud.go`**

Both Create and Update set `opts.ConsistencyPredicate = licenseAssignmentConsistencyPredicate(&object)`
before calling `ReadWithRetry`. The fixed `time.Sleep` calls are removed.

### Action Items

* [x] Add `ConsistencyPredicate` field to `ReadWithRetryOptions`
* [x] Update `ReadWithRetry` retry loop to honour the predicate on successful reads
* [x] Implement `licenseAssignmentConsistencyPredicate` in `predicate.go`
* [x] Remove `time.Sleep` from `license_assignment` Create and Update
* [x] Validate with acceptance tests for all three license assignment test scenarios

### Timeline

Implemented and validated 2026-03-26.

## Validation

All acceptance tests for `microsoft365_graph_beta_users_user_license_assignment` pass:

- `TestAccResourceUserLicenseAssignment_01_Lifecycle` — minimal assignment, no disabled plans
- `TestAccResourceUserLicenseAssignment_02_Maximal` — assignment with disabled plans
- `TestAccResourceUserLicenseAssignment_03_DisabledPlansLifecycle` — add then remove disabled plans

All six unit tests continue to pass. No regressions in other packages.

## References

* [Microsoft: Designing for eventual consistency for Microsoft Entra](https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/)
* [Microsoft Graph API: user: assignLicense](https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta)
* [Terraform Plugin Framework: Resource Read](https://developer.hashicorp.com/terraform/plugin/framework/resources/read)

## Notes

When adding eventual consistency handling to a new resource:

1. Create `predicate.go` in the resource package.
2. Implement a predicate function that accepts the full expected model and returns a
   `func(ctx context.Context, state tfsdk.State) bool` closure.
3. The closure should verify every field that could be stale — not just the field that
   first exhibited the problem.
4. In Create and Update, set `opts.ConsistencyPredicate` before calling `ReadWithRetry`.
5. Do not add `time.Sleep` — let the predicate and retry loop handle timing adaptively.
