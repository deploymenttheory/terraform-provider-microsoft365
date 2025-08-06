# Microsoft 365 Role Assignment Resource Test Strategy

This document outlines the testing strategy for the `microsoft365_graph_beta_device_management_role_assignment` resource, which manages role assignments in Microsoft Intune through the Microsoft Graph Beta API.

## Testing Architecture

The testing architecture follows a two-tiered approach:

1. **Unit Tests**: Mock-based tests that verify resource functionality without making actual API calls
2. **Acceptance Tests**: Integration tests that verify resource functionality with the real Microsoft Graph API

## Test Configuration Files

We use a minimalist approach with just two Terraform configuration files:

| File | Purpose |
|------|---------|
| `resource_minimal.tf` | Minimal configuration with only required attributes |
| `resource_maximal.tf` | Maximal configuration with all possible attributes |

These files are used by both unit tests and acceptance tests, maintaining consistency across all test scenarios.

## Test Scenarios

The test harness (Terraform's testing framework) supports multi-step tests where different configurations can be applied in sequence. This allows us to test all scenarios using just the two base configuration files:

1. **Create with Minimal Configuration**: Uses `resource_minimal.tf` directly
2. **Create with Maximal Configuration**: Uses `resource_maximal.tf` directly
3. **Update from Minimal to Maximal**: Applies `resource_minimal.tf` first, then applies a modified version of `resource_maximal.tf` in the next step
4. **Update from Maximal to Minimal**: Applies `resource_maximal.tf` first, then applies a modified version of `resource_minimal.tf` in the next step
5. **Delete Minimal Configuration**: Creates a resource with `resource_minimal.tf`, then deletes it
6. **Delete Maximal Configuration**: Creates a resource with `resource_maximal.tf`, then deletes it
7. **Import**: Creates a resource with `resource_minimal.tf`, then imports it
8. **Error Tests**: Modifies `resource_minimal.tf` to create an error condition

For acceptance tests that need dynamic values (like user IDs or group IDs), we use string replacement in the test code to modify these values at runtime.

## Role Assignment Specific Test Cases

### Scope Configuration Tests

1. **AllLicensedUsers Scope**: Tests role assignments scoped to all licensed users
2. **AllDevices Scope**: Tests role assignments scoped to all devices
3. **ResourceScopes**: Tests role assignments with specific resource scope IDs
4. **Scope Validation**: Tests validation of scope configuration (mutual exclusivity, required fields)

### Member Tests

1. **Single Member**: Tests role assignment with one member
2. **Multiple Members**: Tests role assignment with multiple members (users and groups)
3. **Member Updates**: Tests adding/removing members from existing assignments

### Role Definition Tests

1. **Built-in Roles**: Tests with built-in Intune role definitions
2. **Custom Roles**: Tests with custom role definitions
3. **Role Definition References**: Tests referencing role definitions from other resources

## Mock Response Structure

Mock responses are organized in the following structure:

```
tests/responses/
├── validate_create/
│   ├── get_role_assignment_minimal.json
│   ├── get_role_assignment_maximal.json
│   ├── get_role_assignment_all_devices.json
│   └── get_role_assignment_minimal_acceptance.json
├── validate_update/
│   └── [update response files]
└── validate_delete/
    └── [delete response files]
```

## Test Execution

### Unit Tests
```bash
go test -v ./internal/services/resources/device_management/graph_beta/role_assignment/
```

### Acceptance Tests
```bash
TF_ACC=1 go test -v ./internal/services/resources/device_management/graph_beta/role_assignment/
```

## Test Naming Convention

Test functions follow a consistent naming pattern:
- `Test[Unit|Acc]RoleAssignmentResource_[Operation]_[Scenario]`
- Examples: `TestUnitRoleAssignmentResource_Create_Minimal`, `TestAccRoleAssignmentResource_Update_ScopeConfiguration`

## Test Principles

1. **Consistency**: All tests use the same base configuration files
2. **Isolation**: Each test is isolated and doesn't depend on the state from other tests
3. **Verification**: Tests verify both the resource state and the API state (for acceptance tests)
4. **Cleanup**: Acceptance tests clean up all created resources
5. **Error Handling**: Tests include error cases to verify proper error handling
6. **Scope Coverage**: Tests cover all three scope types and their validation rules
7. **Member Management**: Tests verify proper handling of user and group members

## Environment Variables for Acceptance Tests

Required environment variables for acceptance tests:
- `M365_CLIENT_ID`: Azure AD application client ID
- `M365_CLIENT_SECRET`: Azure AD application client secret
- `M365_TENANT_ID`: Azure AD tenant ID
- `M365_AUTH_METHOD`: Authentication method (client_secret)
- `M365_CLOUD`: Microsoft cloud environment (public)

## Test Data Requirements

For acceptance tests, the following test data is required:
- Test user accounts in the tenant
- Access to built-in Intune role definitions
- Proper permissions to create and manage role assignments