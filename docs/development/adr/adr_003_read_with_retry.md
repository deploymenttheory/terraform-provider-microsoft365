# ADR-003: ReadWithRetry — Centralised Post-Write Read Retry for Terraform Resources

## Status

Accepted

## Date

2026-03-26

## Context

Every Terraform resource that performs a Create or Update must call Read afterwards to
reconcile what the provider wrote with what Terraform stores in state. For a provider
targeting Microsoft Graph this is non-trivial:

* **Transient API errors** — Graph is a globally distributed service. Any single call may
  encounter throttling (429), transient unavailability (503/504), or propagation-related
  not-found (404) immediately after creation. These are not permanent failures and should
  be retried.

* **Eventual consistency** — Microsoft Entra's multi-replica architecture means a 200
  response from a Read immediately after a Create may return pre-write stale data. This
  causes Terraform to report "Provider produced inconsistent result after apply" even though
  the write succeeded. See ADR-002 for full context on this failure mode.

* **Complex reads** — many resources require multiple Graph API calls to assemble a
  complete state object (e.g. assignments, role scope tags, nested policies). This logic
  must not be duplicated across Create and Update — it belongs exclusively in Read.

Without a shared mechanism, each resource's Create and Update would need to independently
implement retry loops, error classification, context deadline handling, and state
propagation. Across a provider with hundreds of resources this produces inconsistent
behaviour, duplicated code, and makes it difficult to reason about or improve retry
behaviour globally.

## Decision Drivers

* Retry logic — including error classification, deadline awareness, and context
  cancellation — must be written once and shared across all resources.
* Read must remain the single source of truth for assembling resource state; Create and
  Update must not duplicate Read logic.
* CRUD functions must remain clean and uniform — a resource author should call one
  function after a write, not implement a retry loop.
* Error classification must distinguish retryable from non-retryable conditions so the
  loop fails fast on permanent errors (400, 401, 403) and retries on transient ones
  (404, 429, 500, 502, 503, 504).
* The mechanism must be context-deadline-aware so it never overruns the operator-configured
  Terraform timeout.
* Eventual consistency (stale 200 reads) must be handleable within the same mechanism
  without requiring a separate code path per resource.

## Considered Options

* **Inline retry loop per resource** — each Create/Update implements its own loop. Maximum
  flexibility but produces duplicated code, inconsistent behaviour, and no central place
  to improve error handling.

* **Single retry in a shared helper, no error classification** — a shared function that
  retries a fixed number of times regardless of error type. Removes duplication but
  treats all errors the same, causing unnecessary retries on permanent failures and
  potential fast-exit on transient ones.

* **Centralised `ReadWithRetry` with error classification and optional consistency
  predicate** — a single function in `internal/services/common/crud` that wraps the
  resource's own Read method, classifies errors as retryable or non-retryable, respects
  the context deadline, propagates state and identity back to the response container,
  and optionally checks a caller-supplied consistency predicate before accepting a
  successful read.

## Decision

Chosen option: "Centralised `ReadWithRetry` with error classification and optional
consistency predicate", because it eliminates duplication, enforces uniform retry
behaviour across the provider, separates concerns cleanly (resource authors write Read
logic; `ReadWithRetry` handles when and how often to call it), and is extensible without
modifying any resource.

## Rationale

### Single responsibility

The Read method on every resource is responsible for assembling complete state from the
Graph API. `ReadWithRetry` is responsible only for deciding when to call Read, how long to
wait between attempts, and when to give up. Neither knows about the other's internals.

### Error classification

`extractErrorFromDiagnostics` parses Terraform diagnostics to recover HTTP status
information, which is then classified by `IsNonRetryableReadError` and
`IsRetryableReadError`:

| Category | Status codes | Behaviour |
|---|---|---|
| Non-retryable | 400, 401, 403, 409, 423 | Fail immediately |
| Retryable | 404, 429, 500, 502, 503, 504 | Wait `RetryInterval`, retry |
| Unknown | (no recognisable code) | Wait `RetryInterval`, retry |

404 is treated as retryable because Graph APIs commonly return 404 immediately after
creation before the object has propagated. This aligns with Microsoft's recommendation
to treat not-found after creation as transient.

### Context deadline awareness

Before starting the loop, `ReadWithRetry` calculates how many retries fit within the
remaining context deadline (with a 1-second safety margin) and caps `MaxRetries` at that
value. Inside the loop every wait uses `select` over `time.After` and `ctx.Done()` so
cancellation is honoured immediately. This ensures `ReadWithRetry` never overruns the
operator's configured `timeout` block.

### State and Identity propagation

`ReadWithRetry` accepts a `StateContainer` interface rather than a concrete response type.
`CreateResponseContainer` and `UpdateResponseContainer` adapt `resource.CreateResponse`
and `resource.UpdateResponse` respectively. This allows a single implementation to serve
both post-create and post-update reads while correctly propagating both `State` and the
optional `ResourceIdentity` back to the caller's response object.

### Operation context propagation

`ReadWithRetry` injects the calling operation name (`create`, `update`) into the context
as `retry_operation`. The resource's Read method can inspect this value to adjust logging
or error handling based on whether it is being called standalone or from within a retry
loop.

### Consistency predicate

When `ConsistencyPredicate` is set, a successful (no-error) read is only accepted if the
predicate confirms the state is consistent with what was written. If the predicate returns
false, `ReadWithRetry` treats the read as if it had not yet succeeded and continues
retrying with the same interval. This handles Microsoft Entra eventual consistency without
any per-resource sleep or separate code path. See ADR-002 for full context.

## Consequences

### Positive

* Retry logic exists in one place — improvements (better error classification, adaptive
  intervals, observability) apply to all resources simultaneously.
* CRUD functions are uniform across the provider: write to the API, set interim state,
  call `ReadWithRetry`. Resource authors do not implement retry loops.
* Error classification is explicit and auditable in one function.
* Context deadline is always respected — no resource can accidentally overrun its timeout.
* Eventual consistency and transient API errors are handled by the same loop via two
  orthogonal mechanisms (error classification and consistency predicate).

### Negative

* `extractErrorFromDiagnostics` uses string matching on diagnostic messages to recover
  HTTP status codes. This is necessary because Terraform diagnostics do not expose
  structured error metadata, but it means classification is fragile to changes in error
  message formatting from the Graph SDK or Kiota.

### Neutral

* Resources that do not need retry (e.g. simple synchronous APIs) still call
  `ReadWithRetry` via the standard pattern — the overhead is one successful read with no
  retries, which is negligible.
* `DefaultReadWithRetryOptions` (30 retries, 2-second interval) is appropriate for most
  Graph APIs. Resources with unusual latency characteristics can override these values.

## Implementation

### Core types

```go
// ReadWithRetryOptions configures retry behaviour
type ReadWithRetryOptions struct {
    MaxRetries           int
    RetryInterval        time.Duration
    Operation            string
    ResourceTypeName     string
    ConsistencyPredicate func(ctx context.Context, state tfsdk.State) bool
}

// StateContainer abstracts CreateResponse and UpdateResponse
type StateContainer interface {
    GetState() tfsdk.State
    SetState(tfsdk.State)
    GetIdentity() any
    SetIdentity(any)
}
```

### Usage pattern in every resource Create and Update

```go
opts := crud.DefaultReadWithRetryOptions()
opts.Operation = constants.TfOperationCreate
opts.ResourceTypeName = ResourceName
// opts.ConsistencyPredicate = myConsistencyPredicate(&object) // set if resource has eventual consistency exposure

err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
if err != nil {
    resp.Diagnostics.AddError("Error reading resource state after create",
        fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
    return
}
```

### Retry loop decision tree

```
attempt N
  ├── ctx cancelled? → return error
  ├── insufficient time for another attempt? → break
  ├── call readFunc
  │   ├── no error
  │   │   ├── ConsistencyPredicate == nil OR predicate returns true → accept state, return nil
  │   │   └── predicate returns false → log, wait RetryInterval, continue
  │   └── error
  │       ├── non-retryable (400/401/403/409/423) → return error immediately
  │       ├── retryable (404/429/5xx) → log, wait RetryInterval, continue
  │       └── unknown → log, wait RetryInterval, continue
└── exhausted → return last error
```

### Files

| File | Purpose |
|---|---|
| `internal/services/common/crud/read_with_retry.go` | `ReadWithRetry`, `ReadWithRetryOptions`, `StateContainer`, `CreateResponseContainer`, `UpdateResponseContainer`, `extractErrorFromDiagnostics` |
| `internal/services/common/errors/kiota/` | `IsRetryableReadError`, `IsNonRetryableReadError`, `GraphErrorInfo` |

### Action Items

* [x] Implement `ReadWithRetry` with error classification and context deadline capping
* [x] Implement `StateContainer` interface and `CreateResponseContainer` / `UpdateResponseContainer` adapters
* [x] Add `ConsistencyPredicate` field and predicate-aware retry path
* [x] Standardise all resource Create and Update methods to use `ReadWithRetry`

## Validation

`ReadWithRetry` is exercised by every resource acceptance test in the provider. Correct
behaviour under transient errors is verified by unit tests using `httpmock` to inject
specific HTTP status codes and confirm retry counts and final outcomes.

## References

* [ADR-002: Eventual Consistency Handling for Microsoft Entra API Reads](./adr_002_eventual_consistency_handling.md)
* [Microsoft: Designing for eventual consistency for Microsoft Entra](https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/)
* [Terraform Plugin Framework: Resource Read](https://developer.hashicorp.com/terraform/plugin/framework/resources/read)
* [Microsoft Graph API: Error responses](https://learn.microsoft.com/en-us/graph/errors)

## Notes

When implementing a new resource, `ReadWithRetry` is always called at the end of Create
and Update — never at the end of Delete, and never inside Read itself.

If a resource's Read requires multiple Graph API calls to assemble state, all of that
logic belongs in the Read method. `ReadWithRetry` will call it as many times as needed;
the resource author does not need to consider retries when writing Read.
