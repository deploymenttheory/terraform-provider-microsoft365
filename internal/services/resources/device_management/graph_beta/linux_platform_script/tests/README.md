# Linux platform script tests

The harness follows `windows_platform_script`: explicit `resource.UnitTest` and
`resource.Test` cases, numbered Terraform fixtures, HTTP responders under `mocks/`,
shared provider factories, import verification, and acceptance destroy helpers.

| Scenario | Unit | Acceptance |
|---|---|---|
| 001 Minimal configuration and import | Yes | Yes |
| 002 Maximal configuration and import | Yes | Yes |
| 003 Minimal to maximal settings | Yes | Yes |
| 004 Maximal to minimal settings | Yes | Yes |
| 005 One group assignment and import | Yes | Yes |
| 006 Multiple assignments and import | Yes | Yes |
| 007 Grow assignments | Yes | Yes |
| 008 Reduce assignments, then remove all assignments | Yes | Yes |
| 009 API error diagnostics | Yes | — |

The settings scenarios change the script content, execution context, frequency,
and retries. Import verification ignores only `timeouts`. The state mapper tests
also use captured Graph responses to check replacing stale values, Unicode script
content, and preserving existing values when a response cannot be mapped.

Mocks check POST/PUT settings and template references against JSON fixtures.
Settings responses use two-item pages to exercise the existing custom GET helper's
pagination. Assignment fixtures fit within one response page. Unit assignment
coverage includes inclusion/exclusion groups, assignment filters, and all-device
and all-user targets.

Acceptance fixtures use uniquely named policies and newly created empty security
groups. They require no existing group IDs and use the default scope tag. The
fixtures use explicit group resources, consistency waits, and the shared destroy
checks from the reference suite.

## API validation

Curl validation covered creation, policy/settings/assignment reads, PUT updates,
assignment replacement and removal, and deletion. A separate capture of a policy
created by the actual provider confirmed these four settings existed remotely:

- `linux_customconfig_executioncontext`
- `linux_customconfig_executionfrequency`
- `linux_customconfig_executionretries`
- `linux_customconfig_script`

The script decoded to the Terraform fixture exactly. The original acceptance
failures came from `Read` mapping settings into an unused temporary model. The
resource now maps its existing raw settings response into the resource state.

The three-part read remains: policy, settings, assignments. Expanded assignments
were empty on an assigned policy in live testing, so assignment reads use the
dedicated endpoint.

## Running the suite

From the repository root:

```sh
TF_ACC=0 go test ./internal/services/resources/device_management/graph_beta/linux_platform_script \
  -run '^TestUnit' -count=1 -timeout=10m -cover
```

For acceptance tests, export `M365_CLOUD`, `M365_AUTH_METHOD`, `M365_TENANT_ID`,
`M365_CLIENT_ID`, and `M365_CLIENT_SECRET` for the test tenant. The application
needs permission to manage Intune configuration policies and the group
dependencies.

```sh
TF_ACC=1 go test ./internal/services/resources/device_management/graph_beta/linux_platform_script \
  -run '^TestAccResourceLinuxPlatformScript_' -count=1 -timeout=30m -v
```
