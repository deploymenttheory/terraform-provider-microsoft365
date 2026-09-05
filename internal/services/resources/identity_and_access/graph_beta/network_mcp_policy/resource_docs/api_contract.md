# MCP policy API contract

Observed through the Entra Global Secure Access MCP policy blade and authenticated Graph requests on 2026-09-05/06. All probes used disposable, unassigned `tf-api-probe-` policies and rules. This documents configuration persistence, not traffic evaluation. These types are absent from the provider's generated Graph beta SDK, so requests use the configured Kiota adapter with local models.

Base: `https://graph.microsoft.com/beta/networkaccess/mcpPolicies`

| Method | Relative endpoint | Request | Observed success | Verification |
| --- | --- | --- | --- | --- |
| GET | collection, `/{policyId}` | none | 200, policy objects | Portal capture and item GET |
| POST | collection | `name`, `settings` required | **200**, object with `id` | Minimal payload and immediate GET |
| PATCH | `/{policyId}` | changed name, description, settings.defaultAction | **200**, object | Each changed field checked with GET |
| DELETE | `/{policyId}` | none | 204 | Parent and nested child GET returned 404 |
| GET | `/{policyId}/policyRules`, `/{policyId}/policyRules/{ruleId}` | none | 200, rule objects | Portal capture and item GET |
| POST | `/{policyId}/policyRules` | discriminator, name, priority, settings required | **200**, object with `id` | Minimal payload and immediate GET |
| PATCH | `/{policyId}/policyRules/{ruleId}` | discriminator and changed attributes | **200**, object | Each changed field checked with GET |
| DELETE | `/{policyId}/policyRules/{ruleId}` | none | 204 | Item GET returned 404 |

Policy type: `#microsoft.graph.networkaccess.mcpPolicy`. Rule type: **`#microsoft.graph.networkaccess.mcpPolicyRule`**. Policy creation does not automatically create rules. Rules exist under their parent endpoint; deleting a parent cascades to its children. These resources intentionally manage independent identities and expose no assignment or link operation.

Minimal API payloads (API defaults do not imply Terraform schema defaults):

```json
{"name":"tf-api-probe-policy","settings":{}}
```

```json
{"@odata.type":"#microsoft.graph.networkaccess.mcpPolicyRule","name":"tf-api-probe-rule","priority":1000,"settings":{}}
```

The API defaults policy `defaultAction` and rule `action` to `allow`, and rule `settings.status` to `enabled`. Terraform requires explicit `default_action`, `action`, and `enabled`. Observed policy computed fields are `version` and `lastModifiedDateTime`; the provider sends neither.

## Update behavior

Rule PATCH **without the discriminator can return 200 while applying only the name**. The provider includes the discriminator on every rule POST/PATCH and checks post-write GET against requested fields. Policy and rule descriptions distinguish omitted (preserve), explicit null (clear), and empty string (preserve empty). The provider sends explicit null when configuration removes a description.

Rule name and priority must be unique within a parent. Priorities 100, 65001, and 2147483647 were accepted; 99 and values outside Int32 were rejected. Name, description, priority, action, enabled/status, and supported matching conditions update in place. Changing the parent requires replacement, including a parent whose ID is unknown at plan time.

## MCP matching conditions

API path: `matchingConditions.destinations`. The Terraform `matching_conditions` object exposes these measured fields:

| Terraform | API | Shape / measured values |
| --- | --- | --- |
| server_urls | serverUrls | values list + matchType |
| protocol_versions | protocolVersions | values list + matchType; custom strings accepted |
| insecure_connection | insecureConnection | `excluded` stored by HTTPS selection; inverse `required` shown by the UI |
| missing_prm | missingPrm | `required` stored by missing-PRM selection; inverse `excluded` shown by the UI |
| tool_matching | toolMatching | names match object, methods `call` or `list,call` |
| resource_matching | resourceMatching | names match object, methods `read` |
| prompt_matching | promptMatching | names match object, methods `get` |

Measured matchType values: `exactMatch`, `contains`, `notExactMatch`, `notContains`. Names and URL values preserve input order, case and trailing slashes. Empty arrays remain distinct from null. Multiple primitive groups can be stored together even though the portal limits its editor to one primitive group.

Nested PATCH objects merge. Omitting a primitive condition does not remove it; an explicit null does. The provider sends null for removed supported fields at each nested level. Removing `matching_conditions` sends `matchingConditions: null`, which GET represents as `{ "sources": null, "destinations": null }`.

The API returns additional fields as null, including advanced primitive conditions. Unknown non-null fields or non-null `sources` produce a diagnostic rather than being silently discarded from state. Unknown rule statuses are reported verbatim in diagnostics; they are not coerced to an enabled boolean.

AND/OR semantics, wildcard/regex evaluation, method combinations beyond measured values, case sensitivity during evaluation, and the actual traffic effect remain **unverified**. API acceptance alone is not evidence for these semantics.

## Errors, consistency and permissions

Normal Read removes state only on 404; 400 and 403 keep prior state with diagnostics. Create persists a returned identity before readback, including when readback fails with 404. Update/readback errors keep previous known state. Delete accepts already-absent 404. No asynchronous transitions were observed, so no speculative waiter is used. ETag/concurrency behavior and service quotas remain unverified.

Portal delegated requests succeeded with a Global Administrator in the feature-enabled development tenant. Existing Azure CLI authentication received 403 on this endpoint; a logged-in CLI alone does not establish sufficient Graph scopes. Application authentication and the exact permission/role requirements must be distinguished from portal observations and recorded with Terraform validation results. The candidate application permission is `NetworkAccess.ReadWrite.All`; it is not inferred to work solely from the portal's success.
