# Conditional Access template API observations

Observed against the lab tenant on 2026-11-17 using curl and the identity_and_access test application.

## Template catalogue

`GET /beta/identity/conditionalAccess/templates` returned 23 templates. The captured response is replayed in `tests/responses/validate_get/get_template_list.json`.

| Template ID | Current name | Relevant fields |
| --- | --- | --- |
| `c7503427-338e-4c5e-902d-abe252abfb43` | Require multifactor authentication for select admins | Renamed; the ID is unchanged. |
| `634b6de7-c38d-4357-a2c7-3842706eedd7` | Require risk remediation for high-risk users | Formerly password change; now uses `riskRemediation`, the MFA authentication strength, and an every-time sign-in frequency. |
| `6acdf4c3-6815-485c-a57d-2c349d517ba0` | Secure account recovery with identity verification (Preview) | Returns `verifiedID` and `urn:user:accountrecovery`. |
| `6bdcbbb7-ebd1-4f8f-a02f-652db0e3665d` | Block high risk agent identities from accessing resources | Returns scenario `aiAgents`, agent risk `high`, and all agent service principals. |

The SDK at v0.168.0 does not include `verifiedID` in its grant-control enum or `aiAgents` in its template-scenario enum. Its generated decoders converted `verifiedID` to the zero-valued `block` control and dropped `aiAgents`. The datasource preserves these API strings before SDK enum conversion, while retaining SDK parsing for other fields. Names and IDs are not used to infer or replace returned values.

## Policy creation is a separate contract

The curl probe posts the template details to `POST /beta/identity/conditionalAccess/policies`, adds a unique display name and report-only state, removes context annotations, and sends only the authentication-strength ID when present.

| Template | Result | Graph inner error |
| --- | --- | --- |
| Risk remediation | HTTP 201; read back successfully | None |
| Account recovery | HTTP 400 | `1249: Account recovery user action currently supports only user risk and sign-in risk conditions.` |
| High-risk agent identities | HTTP 400 | `1267: Cannot create or update Conditional Access policies that target agents. This capability requires a Microsoft Agent 365 or Microsoft 365 E7 license.` |

Removing only `clientAppTypes` from the account-recovery payload still returned HTTP 400, with error `1329` requiring browser and mobile/desktop client types. The original template therefore cannot simply be replayed as a policy in this tenant.

Template discovery succeeds independently of those write restrictions. Acceptance cases 20 and 21 assert the returned preview settings without creating policies or skipping the tests. Supported-template cases retain policy creation and destruction checks. The risk-remediation example includes the authentication strength and session control supplied by the current template.

Policy GETs can return a transient 404 while a newly created policy still exists. Probe cleanup requires an acknowledged HTTP 204 deletion and five consecutive absence checks through both beta and v1.0. A single 404 is insufficient evidence of deletion.

A fresh probe policy returned DELETE statuses `404, 404, 403` before a later HTTP 204 from the same application. Cleanup retries are bounded and restricted to the probe's verified owned objects; error responses are retained in its report.

## Lab probes

Scripts in the local_testing repository, under `environment/windows/msgraph/identity_and_access`:

- `Test-ConditionalAccessTemplateCatalogue.ps1`: catalogue capture and exact-name checks.
- `Test-ConditionalAccessTemplatePolicyRoundTrip.ps1`: report-only/disabled payload probes and owned-policy cleanup.

References: [List templates](https://learn.microsoft.com/en-us/graph/api/conditionalaccessroot-list-templates?view=graph-rest-beta), [grant controls](https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessgrantcontrols?view=graph-rest-beta).
