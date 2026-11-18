# Delegated permission grants: AzureAD research for issue #4894

Researched 2026-11-18. HashiCorp AzureAD source is pinned to [`e4bf321593a89f1363ef6ebe457f5843b8046882`](https://github.com/hashicorp/terraform-provider-azuread/commit/e4bf321593a89f1363ef6ebe457f5843b8046882). Its Graph dependency is `v0.20260212.1143955`, resolving to SDK commit `7eea1802e3aff42b0c251ae480f5c86fb02cb2a7` ([dependency pin][dependencies]). No applicable `AGENTS.md` was found in the workspace ancestry or this worktree. This is behavioral research, not copied MPL implementation; no tenant operations or resource implementation were performed.

## Resource model and ownership

AzureAD manages **one grant object containing the complete set of scopes**, rather than one Terraform resource per scope. Its schema is:

| Attribute | Cardinality and behavior |
| --- | --- |
| `service_principal_object_id` | Required UUID; replacement on change; Graph `clientId`. |
| `resource_service_principal_object_id` | Required UUID; replacement on change; Graph `resourceId`. |
| `user_object_id` | Optional UUID; replacement on change; Graph `principalId`. Omission selects `AllPrincipals`; presence selects `Principal`. |
| `claim_values` | Required set of strings, minimum one element; mutable; Graph space-separated `scope`. |

There is no separate configurable `consentType`, no nested grant collection, and no `use_existing` switch. All operations default to five minutes ([schema][schema], [documented consent behavior][arguments]). Both service principal fields are **service principal object IDs**, not application/client app IDs or application-registration object IDs. Scope elements are permission **values**, such as `User.Read`, not scope-definition GUIDs. Beta requires a user ID for `Principal`, and null `principalId` for `AllPrincipals` ([beta type][beta-type]). If the new resource exposes consent type explicitly, validate that relationship and replace on changes to consent identity.

AzureAD does not list grants, merge a pre-existing grant, adopt one automatically, or issue an import-required diagnostic during create; after checking the two service principals it POSTs and propagates API failures ([create][create]). Consequently, it provides no useful precedent for safe duplicate adoption. The reviewed beta docs do not establish a precise duplicate-error status or uniqueness guarantee for identical consent tuples; the local runtime probe below observed a duplicate POST returning HTTP 400, `Permission entry already exists.` Treat that as measured behavior, not a universal status guarantee. Treat each remote grant ID as owned by one Terraform resource, and document importing existing consent. Scope PATCH overwrites the entire previous scope string, so splitting the same grant across resources would cause competing updates ([beta update][beta-update]). User consent and tenant-wide consent can coexist for the same client/API; deleting one does not remove the other ([beta delete][beta-delete]).

## Scope conversion and validation

AzureAD joins the set elements with one ASCII space for create/update. Read splits on ASCII spaces, trims surrounding whitespace in each piece, and drops empty pieces; the Terraform set then removes duplicates and ignores ordering ([write conversion][create], [read conversion][read], [split helper][split]). It preserves case. It does not sort the serialized values explicitly.

Its element validator rejects only an empty string; whitespace-only or multi-token elements can pass validation and round-trip into a different set ([validator wrapper][validator-wrapper], [SDK validator][validator]). These are weaknesses to avoid: accept one nonblank scope token per element, preserve case, use set semantics, and sort before joining for deterministic requests. Beta documents a **3,850-character maximum for the complete joined scope string**, not 3,850 per element ([beta type][beta-type]). A response containing repeated spaces or reordered scopes should not cause drift. Null/empty scope is converted to an empty collection on refresh rather than treated as a missing grant ([read][read]).

## IDs and imports

**The remote grant ID is opaque and is not a UUID.** AzureAD v3 uses the state/import form `/oauth2PermissionGrants/{opaque-id}`; the SDK parser declares the final segment user-specified, without UUID validation ([ID parser][id-parser], [ID format][id-format]). Its docs demonstrate a base64url-like final segment ([import example][import-doc]). The importer validates the formatted resource path. Schema migration v0 → v1 wraps the formerly raw remote ID in that path; it does not look up or recreate the grant ([import registration][import-schema], [migration][migration]).

This AzureAD formatting is a provider convention, not Graph's remote ID. For this provider, follow existing import conventions, preserve the returned opaque ID exactly, and use it once in the item URL. Do not reuse a UUID-only importer, lowercase the ID, decode/reconstruct it, or accidentally send the AzureAD path as an item ID. If AzureAD-compatible formatted imports are supported, parse that form deliberately and document it.

## CRUD, responses, and propagation

| Operation | AzureAD behavior | Lesson for this beta resource |
| --- | --- | --- |
| Create | GETs both service principals, POSTs grant, rejects nil model or missing/empty returned ID, stores ID, then reads. | Require a usable returned grant ID. Preflight SP reads are an implementation choice; assess their extra permissions and propagation failures before adopting them. |
| Read | GET by grant ID; 404 clears state, other errors fail; nil successful model fails; refreshes all modeled attributes. | Distinguish an ordinary refresh 404 from a transient 404 immediately after successful create. Do not treat a malformed/nil successful response as deletion. |
| Update | PATCH with only scope populated, then read; no custom polling. | Send only `scope`, representing the complete desired set. A successful 204 has no model to flatten. |
| Delete | DELETE by grant ID; propagates errors; no explicit 404 handling or disappearance wait in the resource. | Use this provider's shared idempotent deletion handling for already-missing objects; do not copy that omission. Delete the grant itself, not the service principal. |

The above comes from the pinned [create][create], [update][update], [read][read], and [delete][delete] implementations. Create alone supplies resource-specific retries: retry 404, and retry 400 only when the OData message indicates a missing referenced object. Other 400s fail. Reads/updates/deletes supply SDK default options; the resource contains no explicit post-write convergence wait. This describes resource-specific behavior, not the SDK's entire generic retry policy. Beta list documentation independently warns about replication delays after create, update, and delete; filtering on `clientId` reduces those delays ([beta list][beta-list]). Carry these cases through existing shared retry/error helpers rather than adding a second retry framework.

AzureAD's client is **v1.0** ([SDK version][sdk-version]); its stable methods POST `/oauth2PermissionGrants` and GET/PATCH/DELETE `/oauth2PermissionGrants/{id}`. The SDK accepts POST 200/201/202/204 but the resource still requires a returned model and ID; PATCH accepts 200/202/204 without decoding a model; DELETE accepts 200/204, not 404 ([POST][sdk-post], [PATCH][sdk-patch], [DELETE][sdk-delete]). Those SDK allowances are not evidence that a bodyless create can be completed reliably.

All required operations are documented on beta: [POST][beta-post] returns a 200-series response **with a grant object**; [GET][beta-get] returns 200 with an object; [PATCH][beta-update] and [DELETE][beta-delete] return 204 without a body. Use `/beta` for all of them. No evidence here requires a stable fallback.

## Historical defects and beta differences

- **#981 / #940:** out-of-band deletion returned 404, but old refresh logic checked for 204 and failed. The [merged fix](https://github.com/hashicorp/terraform-provider-azuread/pull/981) changed the existence check to 404. Preserve the current missing-object behavior and test external deletion; 204 belongs to successful write/delete responses, not a missing GET.
- **#1536:** a scope update failed because serialization emitted immutable `clientId` as an empty string. The [issue](https://github.com/hashicorp/terraform-provider-azuread/issues/1536) records the actual payload/error. [Pandora #4594](https://github.com/hashicorp/pandora/pull/4594) fixed the upstream optionality; the pinned [stable SDK model][sdk-model] now has an optional pointer with omission. Test the serialized PATCH payload, not just model construction. Beta also lists only scope as updateable; the final immutable runtime capture below reproduced the empty-client-ID failure on beta.
- **Beta needs legacy timestamp transport fields:** beta [POST docs][beta-post] and [resource docs][beta-type] say `startTime` and `expiryTime` are required but ignored. AzureAD sends neither, using its stable model; the [beta SDK model][beta-sdk-model] contains both. The [final immutable probe][final-probe] confirms HTTP 400, `Missing property: expiryTime`, for both omitted and explicit-null timestamps, followed by HTTP 201 for explicit RFC3339 values (`startTime: 2000-01-01T00:00:00Z`, `expiryTime: 9999-12-31T23:59:59Z`). **Recommendation: stay on beta and serialize the required legacy fields on CREATE only**, documenting that Graph ignores them. Keep PATCH scope-only. Do not invent an expiration feature or effective grant lifetime. These observations require no stable fallback; independent `startTime` validation was not isolated.
- **Permissions:** AzureAD's [resource docs][permissions] specify `Directory.ReadWrite.All` and Global Administrator. Current beta write docs list `DelegatedPermissionGrant.ReadWrite.All` as least privileged, with `Directory.ReadWrite.All` higher privileged, and several supported delegated directory roles ([beta POST][beta-post]). Use the current beta contract for permissions documentation. Personal Microsoft accounts are unsupported. A grant targeting one user does not itself imply that the provider authenticates as that user.

## Test coverage to carry forward

AzureAD has all-users and nominal single-user create/import acceptance tests, plus GET-based existence checking ([tests][tests]). **The single-user fixture creates a user but never sets `user_object_id`**, so it actually creates all-user consent ([fixture][fixture]). It has no scope-update, replacement, duplicate, whitespace, or disappearance scenario in this file. Do not treat it as sufficient coverage for user consent.

Meaningful cases for this resource: real `Principal` create/import with the correct user; `AllPrincipals` with null principal; scope addition/removal and reordered responses; PATCH omitting every immutable property; opaque non-UUID imports; duplicate create surfaced without takeover; missing dependencies with bounded propagation retries; nil/empty create response; GET 404 versus authentication failure; 204 update/delete; idempotent delete after external removal; both user and admin grants coexisting; and beta timestamp serialization informed by the measured omitted/null/explicit cases. The final probe verified nonempty scope addition/removal with bounded write retries and matched GETs; empty-scope persistence/restoration still needs isolated follow-up. These are recommended checks derived from the documented contract and observed defects, not claims that the AzureAD suite already covers them.

[dependencies]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/go.mod#L3-L12
[schema]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/service_principal_delegated_permission_grant_resource.go#L28-L97
[arguments]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/docs/resources/service_principal_delegated_permission_grant.md#L105-L114
[permissions]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/docs/resources/service_principal_delegated_permission_grant.md#L9-L15
[create]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/service_principal_delegated_permission_grant_resource.go#L101-L162
[update]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/service_principal_delegated_permission_grant_resource.go#L165-L181
[read]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/service_principal_delegated_permission_grant_resource.go#L184-L210
[delete]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/service_principal_delegated_permission_grant_resource.go#L215-L227
[split]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/helpers/tf/slices.go#L34-L42
[validator-wrapper]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/helpers/tf/validation/pluginsdk.go#L260-L267
[validator]: https://github.com/hashicorp/terraform-plugin-sdk/blob/v2.40.1/helper/validation/strings.go#L16-L27
[id-parser]: https://github.com/hashicorp/go-azure-sdk/blob/7eea1802e3aff42b0c251ae480f5c86fb02cb2a7/microsoft-graph/common-types/stable/id_oauth2permissiongrant.go#L27-L40
[id-format]: https://github.com/hashicorp/go-azure-sdk/blob/7eea1802e3aff42b0c251ae480f5c86fb02cb2a7/microsoft-graph/common-types/stable/id_oauth2permissiongrant.go#L85-L96
[import-doc]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/docs/resources/service_principal_delegated_permission_grant.md#L131-L137
[import-schema]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/service_principal_delegated_permission_grant_resource.go#L42-L60
[migration]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/migrations/service_principal_delegated_permission_grant_resource.go#L47-L52
[sdk-version]: https://github.com/hashicorp/go-azure-sdk/blob/7eea1802e3aff42b0c251ae480f5c86fb02cb2a7/microsoft-graph/oauth2permissiongrants/stable/oauth2permissiongrant/version.go#L6
[sdk-post]: https://github.com/hashicorp/go-azure-sdk/blob/7eea1802e3aff42b0c251ae480f5c86fb02cb2a7/microsoft-graph/oauth2permissiongrants/stable/oauth2permissiongrant/method_createoauth2permissiongrant.go#L54-L91
[sdk-patch]: https://github.com/hashicorp/go-azure-sdk/blob/7eea1802e3aff42b0c251ae480f5c86fb02cb2a7/microsoft-graph/oauth2permissiongrants/stable/oauth2permissiongrant/method_updateoauth2permissiongrant.go#L52-L85
[sdk-delete]: https://github.com/hashicorp/go-azure-sdk/blob/7eea1802e3aff42b0c251ae480f5c86fb02cb2a7/microsoft-graph/oauth2permissiongrants/stable/oauth2permissiongrant/method_deleteoauth2permissiongrant.go#L57-L85
[sdk-model]: https://github.com/hashicorp/go-azure-sdk/blob/7eea1802e3aff42b0c251ae480f5c86fb02cb2a7/microsoft-graph/common-types/stable/model_oauth2permissiongrant.go#L15-L40
[beta-sdk-model]: https://github.com/hashicorp/go-azure-sdk/blob/7eea1802e3aff42b0c251ae480f5c86fb02cb2a7/microsoft-graph/common-types/beta/model_oauth2permissiongrant.go#L15-L46
[tests]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/service_principal_delegated_permission_grant_resource_test.go#L23-L68
[fixture]: https://github.com/hashicorp/terraform-provider-azuread/blob/e4bf321593a89f1363ef6ebe457f5843b8046882/internal/services/serviceprincipals/service_principal_delegated_permission_grant_resource_test.go#L118-L137
[beta-type]: https://learn.microsoft.com/en-us/graph/api/resources/oauth2permissiongrant?view=graph-rest-beta
[beta-post]: https://learn.microsoft.com/en-us/graph/api/oauth2permissiongrant-post?view=graph-rest-beta
[beta-get]: https://learn.microsoft.com/en-us/graph/api/oauth2permissiongrant-get?view=graph-rest-beta
[beta-update]: https://learn.microsoft.com/en-us/graph/api/oauth2permissiongrant-update?view=graph-rest-beta
[beta-delete]: https://learn.microsoft.com/en-gb/graph/api/oauth2permissiongrant-delete?view=graph-rest-beta
[beta-list]: https://learn.microsoft.com/en-us/graph/api/oauth2permissiongrant-list?view=graph-rest-beta
[final-probe]: /C:/Users/Dafydd/github/deploymenttheory/local_testing/local_testing/environment/windows/msgraph/applications/oauth2-permission-grant-beta-4894-recorded.json

---

## Local runtime probe: beta permission-scope discovery

Observed 2026-11-18 in a separate, owned-fixture probe. A successful `GET https://graph.microsoft.com/beta/servicePrincipals/{owned-api-sp-id}` returned both `Probe.Read` and `Probe.Write` in **`publishedPermissionScopes`**; top-level `oauth2PermissionScopes` was absent. The application registration retained its definitions in **`api.oauth2PermissionScopes`**. These are distinct beta response paths, so service-principal discovery must inspect `publishedPermissionScopes` rather than the stable service-principal property. Evidence: the `wait-api-scopes` events and API application response in the [local probe report](/C:/Users/Dafydd/github/deploymenttheory/local_testing/local_testing/environment/windows/msgraph/applications/oauth2-permission-grant-beta-4894.json).

The initial discovery loop incorrectly waited on the stable service-principal property. That initial snapshot contained zero `POST oauth2PermissionGrants` events, so it established only the returned beta scope-property name. The discovery correction preceded the separate timestamp probe below.

### Final immutable beta baseline

The [final local probe report][final-probe] records each event with immutable `request_json` and a separately parsed `request` snapshot. It supersedes the preliminary capture for the following facts. All grant endpoints remained on beta.

| Probe | Result |
| --- | --- |
| POST omitting `startTime` and `expiryTime` | HTTP 400, `Request_BadRequest`, `Missing property: expiryTime`; no grant from that request. |
| POST with both timestamps explicitly null | Same HTTP 400 and missing-expiry message. |
| POST with `startTime: 2000-01-01T00:00:00Z` and `expiryTime: 9999-12-31T23:59:59Z` | HTTP 201; returned an opaque, non-UUID grant ID. |
| Duplicate POST | HTTP 400, `Permission entry already exists.` |
| Scope-only PATCH adding `Probe.Write` alongside `Probe.Read` | Initial HTTP 404, then HTTP 204 after bounded retry; GETs included 404 and stale `Probe.Read` before HTTP 200 matched both scopes. |
| Scope-only PATCH removing `Probe.Write`, retaining `Probe.Read` | HTTP 204; a GET initially returned the previous two scopes, then HTTP 200 matched `Probe.Read`. |
| PATCH including `clientId: ""` with `scope` | HTTP 400, `Property cannot be updated: clientId`, reproducing the #1536 failure pattern on beta. |

This verifies ordinary nonempty scope addition and removal on beta, including eventual convergence. A successful transport response alone does not establish convergence; the observed write 404 and stale/missing reads justify bounded operation-specific retries and post-write verification before removing newly created state. Fixture cleanup subsequently completed and `cleanup_verified` is `true`; grant absence was confirmed on three consecutive reads after both owned service principals and application registrations were removed.

Direct grant deletion remains an open lifecycle check: the recorded DELETE returned HTTP 400, `Permission being updated or deleted is not found.`, despite an earlier successful GET. The cleanup report retains that warning. Later grant disappearance following owned-parent cleanup does **not** demonstrate a successful direct grant DELETE. Isolate stale replicas versus a beta operation defect before implementing deletion or considering any API-version fallback.

**Implementation recommendation:** keep all CRUD on beta. Serialize non-null `startTime` and `expiryTime` on CREATE to satisfy its legacy contract; the values above were accepted by the measured probe. Omit them and all consent identity fields from PATCH, sending only `scope`. Document that Graph ignores the legacy timestamps; do not expose or promise an effective expiration feature. Independent `startTime` validation and the necessity of those particular timestamp values were not isolated.

### Superseded preliminary evidence and empty-scope limitation

The [preliminary probe report](/C:/Users/Dafydd/github/deploymenttheory/local_testing/local_testing/environment/windows/msgraph/applications/oauth2-permission-grant-beta-4894.json) reused a mutable request dictionary: all three create-phase request snapshots ended up displaying the final explicit timestamp values. That capture limitation is superseded by the final report's immutable `request_json` for the timestamp, duplicate, nonempty-update, and immutable-client-ID facts above.

The preliminary run also observed empty-scope PATCH 204, later GET 200 with empty scope, and then a restoration PATCH 400, `Permission being updated or deleted is not found.` It also contained intermittent 404/stale reads, so replication effects were not isolated. **Empty-scope stable semantics remain unresolved: neither this sequence nor the final nonempty baseline proves automatic deletion when scopes are cleared.** No expiration semantics are inferred from either probe.
