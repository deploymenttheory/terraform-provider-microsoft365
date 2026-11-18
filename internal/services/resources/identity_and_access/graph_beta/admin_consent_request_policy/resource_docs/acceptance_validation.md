# Tenant verification, issue 4894 ask 5

Validated 2026-11-18 against the M365 lab using the identity-and-access service-domain credentials.

The PowerShell curl probe established beta GET 200, identical-value PUT 204 without a response body, and identical-value PATCH 404. The resource retains beta and overrides generated PATCH request information to PUT; it does not switch API versions or use a separate HTTP client. The captured writable configuration was disabled, all notifications false, duration 0, reviewers empty.

`TestAccResourceAdminConsentRequestPolicy_PreserveTenantConfiguration` passed in 8.88 seconds. It reads the existing policy and configures exactly those writable values, verifies ordinary and generated imports with a no-change plan, then destroys Terraform management. A final Graph read confirms all five writable settings remain unchanged and the policy still exists. The test does not enable the workflow or substitute reviewers. Version is read-only and excluded from preservation checks.

Unit coverage exercises nonempty reviewer updates, version increments, return to false/zero/empty settings, generated import of the zero/empty policy, no-change plans, wrong singleton IDs, state-only destroy without a Graph client, and state preservation on 400/403/404 or an empty response.

The curl script and captured JSON are in the lab `msgraph/identity_and_access` directory as `Test-AdminConsentRequestPolicyBeta.ps1` and `admin-consent-request-policy-beta-4894.json`. Test logs are retained in this branch's ignored `tmp/validation` directory. Nonempty reviewers and enabled workflow changes were validated in unit tests, not changed on the shared lab tenant.
