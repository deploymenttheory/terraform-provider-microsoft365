# Win32 content lifecycle acceptance tests

These tests make real Graph/Azure calls and create/delete applications and assignments.
Use a designated test tenant and an **empty, dedicated test group**. Do not use a
production group. The assignment intent is `available`, not a required deployment.

## Fixtures and credentials

Use the repository's standard `TF_ACC=1`, `M365_CLIENT_ID`, `M365_CLIENT_SECRET`,
`M365_TENANT_ID`, `M365_AUTH_METHOD`, and `M365_CLOUD` configuration. The application
needs the documented application-management permissions. Set
`WIN32_APP_TEST_GROUP_ID` to the empty test group's ID.

On Windows, use Microsoft's Content Prep Tool to package each checked-in fixture:

```powershell
IntuneWinAppUtil.exe -c tests\fixtures\v1 -s setup.cmd -o C:\win32-test\v1 -q
IntuneWinAppUtil.exe -c tests\fixtures\v2 -s setup.cmd -o C:\win32-test\v2 -q
```

Copy the two packages to the test runner (Windows, macOS, or Linux) and set
`WIN32_APP_INTUNEWIN_V1` and `WIN32_APP_INTUNEWIN_V2` to their distinct absolute paths.
The tests create the equivalent plain ZIPs directly from the fixture scripts.
The scripts only create a version marker under `%ProgramData%\TerraformWin32Acceptance`.

From the repository root:

```sh
TF_ACC=1 go test ./internal/services/resources/device_and_app_management/graph_beta/win32_app \
  -run '^TestAccResourceWin32App_' -count=1 -v -timeout=90m
```

## What is checked

Each starting format gets a creation test, a second content version, a switch to
the other format, and another content version. The checks require unchanged app
and assignment IDs, preservation of the actual Graph group assignment, and
exactly one committed content version in Terraform state. An explicit empty
plan after every apply exercises refresh/idempotence. Import and deletion are
also verified. Import verification excludes local sources and the configured
outer filename, which Graph cannot return.

Record the commit SHA, commands, environment/platform, timestamps, and sanitized
pass/fail output in the PR. Do not include credentials, signed upload URLs, or
package encryption keys. Ordinary `go test` skips these tests; skipped tests are
not evidence of tenant validation.

## Device installation remains a separate check

The automated lifecycle tests do not prove that IME can install the content.
Separately deploy each format to an explicitly authorized test Windows device,
check the version marker and Intune installation status, and record sanitized
results. For a v2 installation test, remove the v1 marker or adjust the detection
rule to require v2; the lifecycle fixture's existence rule alone will consider
v1 already installed. Remove the test application/marker afterward.
