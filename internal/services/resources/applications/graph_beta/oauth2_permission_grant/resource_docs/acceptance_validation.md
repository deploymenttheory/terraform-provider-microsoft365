# Lab validation, issue 4894 ask 1

Validated 2026-11-18 against the M365 lab using the applications service-domain credentials. No application secrets were created; both owned service principals and the Principal-case user were disabled.

- `Test-OAuth2PermissionGrantBeta.ps1` captures beta POST requirements, duplicate rejection, scope-only PATCH, and direct DELETE/GET absence before deleting parent service principals. The fresh control without an invalid immutable-field PATCH returned DELETE 204 and three subsequent GET 404 responses.
- Principal grant curl reproduction returned `400 Filter is invalid. principalId` twice, followed by GET 200 without changing the grant. A v1.0 deletion control also returned `400 Permission being updated or deleted is not found`; this is not evidence beta lacks support. The provider keeps all grant CRUD on beta and retries these specific replication errors within its timeout.
- Single post-update matching reads were followed by stale scope values during Terraform's immediate refresh. The resource now requires five consecutive matching observations through the shared read-retry helper, resetting the count on a stale observation.
- Final `TestAccResourceOAuth2PermissionGrant_Lifecycle` passed in 290.19 seconds: AllPrincipals 133.38 seconds, Principal 156.81 seconds. Both cover create, import verification, generated import with a no-change plan, add scope, remove scope, no-change plan, and direct grant destruction before parent fixture cleanup. Both `CheckDestroy` checks passed.
- Unit coverage includes exact scope-only PATCH serialization, conditional principal validation, opaque grant IDs, duplicate scope normalization, missing/forbidden/invalid responses, cancellation, deletion verification, and stable-read reset behavior.

Local evidence resides in the lab `msgraph/applications` folder: `oauth2-permission-grant-beta-4894-direct-delete.json`, `oauth2-permission-grant-principal-contract-retry-4894.json`, and `oauth2-permission-grant-principal-cleanup-4894.json`. The final cleanup audit enumerates grant collection pages and confirms no earlier owned grant remains, including cases where deleted-user item lookup produces 400 instead of 404. Initial failures are retained in the probe reports rather than erased.

Acceptance and unit logs are kept in this branch's ignored `tmp/validation` directory. The tenant acceptance result is evidence for the tested lifecycles, not a promise about every permission scope or tenant policy.
