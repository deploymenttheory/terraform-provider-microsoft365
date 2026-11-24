# ADR-004: custom_requests as the Only Non-SDK Graph Transport, and Documented SDK Deviations

## Status

Accepted

## Date

2026-11-24

## Context

The provider calls Microsoft Graph through the generated Go SDKs
(`msgraph-sdk-go` and `msgraph-beta-sdk-go`). Some Graph endpoints cannot be called that
way: the SDK has no model for them, or its model is wrong for what the live API accepts
or returns. Resources that hit these gaps need another way to send the request.

### History of `custom_requests`

* **The Config API.** ShocOne wrote the per-verb request functions (`PostRequest`,
  `PostRequestNoContent`, `GetRequestByResourceId`, `PatchRequestByResourceId`,
  `PutRequestByResourceId`, `DeleteRequestByResourceId`) and their `…RequestConfig`
  structs in late 2024, inside the client package. In June 2025 ShocOne moved them into
  their own package, `internal/services/common/custom_requests`
  ([356f0901](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/356f0901)),
  and the same day reworked the URL template code in `PatchRequestByResourceId` as part
  of a bug-fix and unit-test commit
  ([44bb11a8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/44bb11a8)).
* **kenchan0130's additions.** In November 2025 kenchan0130 added `GetRequest`, a GET
  that sends through `adapter.Send` and returns a Kiota `Parsable` from a caller-supplied
  factory
  ([9e33a72c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9e33a72c)),
  and replaced the `ByIDRequestUrlTemplate` call in `DeleteRequestByResourceId` with an
  inline switch over the resource ID pattern
  ([47ce3efb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/47ce3efb)).
  In September 2026 kenchan0130 added `JSONRequest` (`json_request.go`) inside the
  resource commit "feat: add network explicit forward proxy resource"
  ([f4900a42](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f4900a42)).
  `JSONRequest` sends and decodes `encoding/json` values through the SDK adapter, for
  endpoints that have no generated model.

### The problem

Without one agreed transport, resource packages solved the same problem in different
ways. The AST scan of PR #4903 (head `90a6721b`) found:

* 17 packages that build `RequestInformation` and call `adapter.Send*` themselves.
  Together they carry about 174 hand-written `Serialize`/`GetFieldDeserializers` methods,
  36 discriminator factories and most of 42 `decode…Response` functions.
* Packages using `JSONRequest` (9 on `main` today) alongside packages using the Config
  API, for the same kind of call.
* SDK workarounds explained in free-form comments spread across `requests.go`,
  `construct.go` and `response.go`, with no common format and no stated condition for
  removing them.

The Config API also has defects of its own, which are fixed in #4935 and #4941:

* `GetRequestByResourceId` and `DeleteRequestByResourceId` send through `makeRequest`,
  which uses a bare `http.Client{}` instead of the adapter. The request therefore keeps
  the hard-coded `https://graph.microsoft.com/<version>` base URL, which breaks national
  clouds, and its status code is not checked
  ([#1647](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/1647)).
* The other five functions (`PostRequest`, `PostRequestNoContent`, `GetRequest`,
  `PatchRequestByResourceId` and `PutRequestByResourceId`) send through `adapter.Send*`.
  Kiota replaces the `baseurl` path parameter with the adapter's base URL on every send,
  so they already reach the correct cloud, but they silently ignore `APIVersion`.
* `ByIDRequestUrlTemplate` turns the pattern `/{id}` into `/{<guid>}`, which Kiota
  expands as an empty template variable. No live call reaches this case:
  `PatchRequestByResourceId` and `DeleteRequestByResourceId` handle `/{id}` and `('id')`
  explicitly, and every `GetRequestByResourceId` caller uses `('id')` or no pattern.

### Kiota already supports raw JSON

`JSONRequest` exists because the Config API only accepts Kiota `Parsable` bodies. Kiota
itself does not need that restriction:

* kiota-abstractions-go v1.12.1 provides
  [`RequestInformation.SetStreamContentAndContentType`](https://github.com/microsoft/kiota-abstractions-go/blob/v1.12.1/request_information.go),
  which sets a raw `[]byte` body and its content type.
* kiota-http-go v1.6.1
  [`SendPrimitive(ctx, info, "[]byte", mappings)`](https://github.com/microsoft/kiota-http-go/blob/v1.6.1/nethttp_request_adapter.go)
  returns the raw response body after the adapter's usual status check and error
  mapping, and returns `nil` on 204 (`shouldReturnNil`).

So raw JSON bodies and responses can be added to the Config API without a second
transport.

## Decision Drivers

* Every Graph call must go through the provider's configured adapter, so that the base
  URL (national clouds), authentication, middleware, retry and OData error mapping are the
  same for every resource.
* Request-building code must be written once, in one package, and tested there.
* The SDK is the default. Bypassing it must be the exception, limited to known reasons.
* Every bypass must say why it exists and when it can be removed, so it can be deleted
  when the SDK catches up.
* Resource packages must not carry hand-written Kiota models.

## Considered Options

* **Direct SDK use only.** Use the generated request builders and models everywhere.
  Simplest, but it cannot reach endpoints whose type is missing, cannot send properties
  the model lacks, and fails on responses the model cannot parse.

* **Per-package Kiota models and request building.** Each package builds its own
  `RequestInformation`, defines `Parsable` types for its bodies and responses, and calls
  `adapter.Send*`. This works and keeps the adapter pipeline, but it is the pattern behind
  the 174 hand-written serializer methods and 36 discriminator factories above, and every
  package repeats the URL template, header and error-mapping code.

* **`JSONRequest`.** A single function that marshals `encoding/json` values and sends them
  through the adapter. It removes the hand-written models, but it is a second transport
  next to the Config API, with its own URL template convention, and it was added inside a
  resource commit rather than designed as part of `custom_requests`.

* **The Config API as the only non-SDK transport, with raw JSON support and documented
  deviations.** `custom_requests` gains raw JSON bodies and responses through
  `SetStreamContentAndContentType` and `SendPrimitive(…, "[]byte", …)`, uses the
  adapter's base URL, and becomes the only path for calls the SDK cannot make. A resource
  may use it only for one of four recorded SDK deviations, documented with a fixed comment
  block.

## Decision

Chosen option: "The Config API as the only non-SDK transport, with raw JSON support and
documented deviations", because it keeps one tested request path on the adapter pipeline,
removes per-package Kiota models, and makes every bypass visible and removable.

### 1. The transport rule

* The generated SDK is used wherever it works.
* Every Graph call the SDK cannot make goes through the Config API in
  `internal/services/common/custom_requests`.
* Raw JSON bodies use `RequestInformation.SetStreamContentAndContentType`. Raw JSON
  responses use `adapter.SendPrimitive(ctx, info, "[]byte", mappings)`.
* `JSONRequest` is retired. Its callers move to the Config API (#4946, #4948) and
  `json_request.go` is then deleted.
* Resource packages must not build `RequestInformation` or call `adapter.Send*`
  themselves.

### 2. The four accepted SDK deviations

A resource may use the Config API instead of the SDK only for one of these reasons. Use
the category label exactly as written in the comment block.

| # | Category label | Meaning | Reference package |
|---|---|---|---|
| 1 | `missing type` | The SDK has no model or request builder for the Graph type or endpoint. | `network_content_policy` |
| 2 | `serializer sends read-only fields or nulls` | The SDK model's serializer writes properties the endpoint rejects. | `network_filtering_profile` |
| 3 | `unknown enum value` | The live API returns an enum value the SDK does not know, so deserialisation drops or fails on it. | `on_premises_connector_group` |
| 4 | `missing property` | The SDK model exists but lacks a property the resource must send or read. | `network_forwarding_profile` |

Any other reason is not a deviation: fix the code to use the SDK, or raise it for a new
decision.

### 3. The required comment block

Each package with a deviation documents it in the package doc comment of `crud.go`, the
comment directly above the `package` clause. One block per deviation:

```go
// SDK deviation: <category label>
//
// Evidence:
//   - SDK: <module>@<version>, <package>.<Type or builder> — what is missing or wrong.
//   - Learn: <Microsoft Learn URL>, or "none" with the reason.
//   - Observed: <date>, <HTTP method and path>, the error or behaviour seen.
//
// Removal condition: <the SDK version, type, property or enum value whose arrival makes
// this deviation unnecessary>.
```

* **Category label** is one of the four labels in the table above.
* **Evidence** names the SDK version checked, links the Learn page, and records what was
  observed against the live API. If the evidence comes from reading SDK source rather than
  a live call, say so.
* **Removal condition** names what has to change in the SDK. When a dependency bump
  meets it, the deviation is removed and the package moves back to the SDK.

### 4. Alternatives considered

Per-package Kiota models, `JSONRequest` and direct SDK use are rejected as the general
approach, for the reasons under Considered Options. Direct SDK use remains the default
wherever no deviation applies.

## Worked Examples

The examples below show the comment block as it applies to four real packages. The
evidence is taken from the comments those packages carry today and was checked against
`msgraph-beta-sdk-go` v0.168.0, the version in `go.mod`. Each block replaces the
free-form comment in the package when that package is reworked (RWK tickets).

### Missing type: `network_content_policy`

`internal/services/resources/identity_and_access/graph_beta/network_content_policy`
manages `/networkaccess/filePolicies`. Today it defines its own `contentPolicyRequestBody`
and `contentPolicyResponse` Parsable types and builds `RequestInformation` in
`requests.go`. The beta SDK v0.168.0 has no `filePolicy` model in `models/networkaccess`
and no `filePolicies` request builder in `networkaccess`, and Learn has no
`networkaccess-filepolicy` resource page.

```go
// Package graphBetaNetworkContentPolicy manages Global Secure Access content (file) policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no filePolicy model in
//     models/networkaccess and no filePolicies request builder in networkaccess.
//   - Learn: none. There is no networkaccess-filepolicy resource page; the Global Secure Access
//     API overview does not list file policies.
//   - Observed: the request and response shape for /networkaccess/filePolicies were taken from
//     the Entra admin centre Global Secure Access content policy blade.
//
// Removal condition: the beta SDK ships a networkaccess filePolicy model and a
// /networkAccess/filePolicies request builder.
package graphBetaNetworkContentPolicy
```

### Serializer sends read-only fields or nulls: `network_filtering_profile`

`internal/services/resources/identity_and_access/graph_beta/network_filtering_profile`
sends a hand-written `filteringProfileRequestBody` (see `construct.go`). In the beta SDK
v0.168.0, `networkaccess.Profile.Serialize` writes `description`, `lastModifiedDateTime`
and `version`, and `networkaccess.FilteringProfile.Serialize` writes `createdDateTime` and
`priority`. `createdDateTime`, `lastModifiedDateTime` and `version` are read-only.

```go
// Package graphBetaNetworkFilteringProfile manages Global Secure Access filtering profiles.
//
// SDK deviation: serializer sends read-only fields or nulls
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, networkaccess.FilteringProfile.
//     Its Serialize (with Profile.Serialize) writes the read-only createdDateTime,
//     lastModifiedDateTime and version properties.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringprofile?view=graph-rest-beta
//   - Observed: POST /networkAccess/filteringProfiles with those fields in the body returns
//     "400 Invalid request parameters". The request is accepted with only name, description,
//     priority and state.
//
// Removal condition: the SDK serializer stops writing createdDateTime, lastModifiedDateTime
// and version on create and update.
package graphBetaNetworkFilteringProfile
```

### Unknown enum value: `on_premises_connector_group`

`internal/services/resources/applications/graph_beta/on_premises_connector_group` parses
`region` and `connectorGroupType` as strings (see `response.go`). In the beta SDK
v0.168.0, `models.ParseConnectorGroupRegion` knows only `nam`, `eur`, `aus`, `asia`, `ind`
and `unknownFutureValue`, and returns `nil` for any other value, so the SDK model loses
the region.

```go
// Package graphBetaApplicationsOnPremisesConnectorGroup manages application proxy connector groups.
//
// SDK deviation: unknown enum value
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, models.ConnectorGroupRegion.
//     ParseConnectorGroupRegion knows nam, eur, aus, asia, ind and unknownFutureValue, and
//     returns nil for anything else.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
//     documents connectorGroupType values that are not in the SDK enum either.
//   - Observed: 2026-07-05, GET /beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups
//     returned the default connector group with "region": "japan", which the SDK model drops.
//
// Removal condition: models.ConnectorGroupRegion includes every region Graph returns
// (at least japan), and models.ConnectorGroupType includes the values Learn documents.
package graphBetaApplicationsOnPremisesConnectorGroup
```

### Missing property: `network_forwarding_profile`

`internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile`
sends raw JSON for the forwarding profile v2 properties (see `requests.go` and
`shared_models/graph_beta/identity_and_access/network_forwarding_profile.go`). In the beta
SDK v0.168.0, `networkaccess.ForwardingProfile` (through `ForwardingProfileBase` and
`Profile`) has `associations`, `isCustomProfile`, `priority`, `servicePrincipal`,
`trafficForwardingType`, `description`, `lastModifiedDateTime`, `policies`, `state` and
`version`, but not `clientFallbackAction`, `localNetworkSettings`,
`forwardingSourceMatching` or `assignmentRules`. (`clientFallbackAction` exists only on
`networkaccess.ForwardingRule`.) The Learn resource page does not list them either; the
shared model notes they are not in `$metadata` and that their shapes follow the portal.

```go
// Package graphBetaNetworkForwardingProfile manages custom Global Secure Access forwarding profiles.
//
// SDK deviation: missing property
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, networkaccess.ForwardingProfile
//     has no clientFallbackAction, localNetworkSettings, forwardingSourceMatching or
//     assignmentRules property.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-forwardingprofile?view=graph-rest-beta
//     does not list these properties either; they are not in the beta $metadata.
//   - Observed: the property shapes follow the Entra admin centre's forwarding profile v2
//     requests to /networkAccess/forwardingProfiles; see
//     https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
//
// Removal condition: networkaccess.ForwardingProfile exposes clientFallbackAction,
// localNetworkSettings, forwardingSourceMatching and assignmentRules.
package graphBetaNetworkForwardingProfile
```

## Rationale

### One transport keeps the adapter pipeline

Every request built by the Config API goes through the provider's Kiota adapter, so it
gets the configured base URL, authentication, middleware, retry and ODataError mapping.
Per-package request building and `JSONRequest` can also use the adapter, but each copy is
a place where the base URL, headers or error mapping can drift.

### Raw JSON removes hand-written models

With `SetStreamContentAndContentType` and `SendPrimitive(…, "[]byte", …)`, a package that
needs a type the SDK lacks can send and receive plain Go structs with `encoding/json`
tags. It no longer needs `Serialize`, `GetFieldDeserializers` or discriminator factories.

### The comment block makes deviations removable

Each deviation names the SDK version that was checked and the condition for removing it.
When the SDK is bumped, a search for `SDK deviation:` lists every package to re-check.

## Consequences

### Positive

* One tested request path for every Graph call the SDK cannot make, with national cloud
  support once #4935 lands.
* Hand-written Kiota models and per-package request code are removed from resource
  packages.
* Every SDK bypass is documented in the same place and format, with a removal condition.
* `JSONRequest` and its separate URL template convention are retired.

### Negative

* The Config API has to grow raw JSON bodies, response bodies, headers and paging
  (#4941) before every caller can migrate.
* Packages that currently work through per-package requests or `JSONRequest` must be
  migrated (RWK tickets, #4948).
* Raw JSON structs are not type-checked against Graph metadata. Request-capture and unit
  tests carry that check instead.

### Neutral

* Resources that work with the SDK are unaffected.
* The comment block is documentation. A report-only helper check (#4936) may later look
  for it, but this ADR does not add enforcement.

## Implementation

### Action Items

* [x] Record this decision (#4934).
* [ ] Use the adapter base URL, the Kiota send pipeline and correct by-ID templates in
  every Config function (#4935).
* [ ] Add raw JSON bodies and responses, headers and paging to the Config API (#4941).
* [ ] Migrate the per-package transports and add the comment block to each package that
  keeps a deviation (RWK tickets, #4943 to #4947).
* [ ] Migrate the remaining `JSONRequest` callers and delete `json_request.go` (#4948).

### Files

| File | Purpose |
|---|---|
| `internal/services/common/custom_requests/` | The Config API: the only non-SDK Graph transport |
| `internal/services/common/custom_requests/json_request.go` | `JSONRequest`, retired by this decision |
| `<resource package>/crud.go` | Package doc comment carrying the SDK deviation block |

## Validation

* A code search finds no `adapter.Send*` call and no `RequestInformation` built outside
  `custom_requests` and the generated SDK.
* `json_request.go` no longer exists.
* Every package that uses the Config API instead of the SDK has an `SDK deviation:` block
  in the package doc of its `crud.go`, using one of the four labels.

## References

* [ADR-002: Eventual Consistency Handling for Microsoft Entra API Reads](./adr_002_eventual_consistency_handling.md)
* [ADR-003: ReadWithRetry](./adr_003_read_with_retry.md)
* [Development Guide](../guide.md)
* [kiota-abstractions-go v1.12.1 `request_information.go`](https://github.com/microsoft/kiota-abstractions-go/blob/v1.12.1/request_information.go)
* [kiota-http-go v1.6.1 `nethttp_request_adapter.go`](https://github.com/microsoft/kiota-http-go/blob/v1.6.1/nethttp_request_adapter.go)
* [Microsoft Graph Global Secure Access API overview (beta)](https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-global-secure-access-api-overview?view=graph-rest-beta)
* [Issue #1647: GetRequestByResourceId does not check the status code](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/1647)

## Notes

The four categories are the only accepted reasons for bypassing the SDK. A package that
bypasses the SDK for any other reason, for example because the SDK call was harder to
write, should be moved back to the SDK.

Kiota's JSON serialization writer (kiota-serialization-json-go v1.1.5) skips properties
whose value is `nil`. For the "read-only fields or nulls" category, the evidence should
therefore name the fields that are actually sent, for example read-only values carried
over from a read model, rather than assume every unset property is written as `null`.
