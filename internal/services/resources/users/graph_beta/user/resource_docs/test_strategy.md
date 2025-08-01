# Microsoft 365 User Resource Test Strategy

This document outlines the testing strategy for the `microsoft365_graph_beta_users_user` resource, which manages Microsoft 365 users through the Microsoft Graph Beta API.

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

For acceptance tests that need dynamic values (like user principal names with real domains), we use string replacement in the test code to modify these values at runtime.

## Unit Tests

Unit tests use mock HTTP responses to test the resource's functionality without making actual API calls. The mocks are defined in `mocks/responders.go`.

### Test Cases

| Test | Description |
|------|-------------|
| `TestUnitUserResource_Create_Minimal` | Tests creating a user with minimal configuration |
| `TestUnitUserResource_Create_Maximal` | Tests creating a user with maximal configuration |
| `TestUnitUserResource_Update_MinimalToMaximal` | Tests updating from minimal to maximal configuration |
| `TestUnitUserResource_Update_MaximalToMinimal` | Tests updating from maximal to minimal configuration |
| `TestUnitUserResource_Delete_Minimal` | Tests deleting a user with minimal configuration |
| `TestUnitUserResource_Delete_Maximal` | Tests deleting a user with maximal configuration |
| `TestUnitUserResource_Import` | Tests resource import functionality |
| `TestUnitUserResource_Error` | Tests error handling |

### Running Unit Tests

```bash
cd internal/services/resources/users/graph_beta/user
go test -v
```

## Acceptance Tests

Acceptance tests make actual API calls to verify the resource's functionality with the real Microsoft Graph API.

### Test Cases

| Test | Description |
|------|-------------|
| `TestAccUserResource_Create_Minimal` | Tests creating a user with minimal configuration |
| `TestAccUserResource_Create_Maximal` | Tests creating a user with maximal configuration |
| `TestAccUserResource_Update_MinimalToMaximal` | Tests updating from minimal to maximal configuration |
| `TestAccUserResource_Update_MaximalToMinimal` | Tests updating from maximal to minimal configuration |
| `TestAccUserResource_Delete_Minimal` | Tests deleting a user with minimal configuration |
| `TestAccUserResource_Delete_Maximal` | Tests deleting a user with maximal configuration |
| `TestAccUserResource_Import` | Tests resource import functionality |

### Prerequisites

To run the acceptance tests, you need:

1. A Microsoft 365 tenant with appropriate permissions
2. Valid test domain for creating users

### Environment Variables

| Variable | Description |
|----------|-------------|
| `TF_ACC` | Set to `1` to run acceptance tests |
| `M365_CLIENT_ID` | Azure AD application client ID |
| `M365_CLIENT_SECRET` | Azure AD application client secret |
| `M365_TENANT_ID` | Azure AD tenant ID |
| `TEST_DOMAIN` | Domain for creating test users (e.g., `contoso.com`) |

### Running Acceptance Tests

```bash
# Set required environment variables
export TF_ACC=1
export M365_CLIENT_ID="your-client-id"
export M365_CLIENT_SECRET="your-client-secret"
export M365_CLIENT_ID="your-tenant-id"
export TEST_DOMAIN="your-domain.com"

# Run acceptance tests
cd internal/services/resources/users/graph_beta/user
go test -v -timeout 30m
```

## Implementation Details

### String Replacement for Dynamic Values

For acceptance tests that need dynamic values, we use string replacement to modify the base configuration files:

```go
// Read the template file
content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
if err != nil {
    return ""
}

// Replace the UPN with the test UPN
updated := strings.Replace(string(content), "minimal.user@contoso.com", userPrincipalName, 1)

return updated
```

### Mock HTTP Responses

Unit tests use the `httpmock` package to intercept HTTP requests and provide mock responses:

```go
// Set up mock environment
_, _ = setupMockEnvironment()
defer httpmock.DeactivateAndReset()
```

The mock responses are defined in `mocks/responders.go`.

## Best Practices

1. **DRY Principle**: We follow the Don't Repeat Yourself principle by using just two base configuration files for all test scenarios
2. **Isolation**: Each test is isolated and doesn't depend on the state from other tests
3. **Verification**: Tests verify both the resource state and the API state
4. **Cleanup**: Acceptance tests clean up all created resources
5. **Error Handling**: Tests include error cases to verify proper error handling
6. **Consistent Naming**: Test functions follow a consistent naming pattern:
   - `Test[Unit|Acc]UserResource_[Operation]_[Scenario]` 
   - Examples: `TestUnitUserResource_Create_Minimal`, `TestAccUserResource_Update_MaximalToMinimal`
7. **Helper Functions**: We use helper functions to generate test configurations dynamically 