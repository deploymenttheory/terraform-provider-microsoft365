# User Resource Tests

This directory contains both unit tests and acceptance tests for the User resource.

## Unit Tests

Unit tests can be run without any special setup:

```bash
go test -v
```

These tests use mocked responses and don't require any real Microsoft 365 resources.

### Unit Test Structure

The unit tests are organized into several focused test functions:

1. `TestUnitUserResource_Minimal` - Tests minimal configuration with only required attributes
2. `TestUnitUserResource_Maximal` - Tests maximal configuration with all possible attributes
3. `TestUnitUserResource_Create` - Tests creating multiple resources simultaneously
4. `TestUnitUserResource_Update` - Tests updating a resource's configuration
5. `TestUnitUserResource_MinimalToMaximal` - Tests updating from minimal to maximal configurations
6. `TestUnitUserResource_Import` - Tests importing an existing resource
7. `TestUnitUserResource_Error` - Tests error handling

To run a specific unit test:

```bash
go test -v -run TestUnitUserResource_Update
```

## Acceptance Tests

Acceptance tests require actual Microsoft 365 resources and credentials. They are designed to verify the resource works correctly against the real Microsoft Graph API.

### Required Environment Variables

To run the acceptance tests, you need to set the following environment variables:

```bash
# Azure AD app registration credentials
export ARM_CLIENT_ID="your-client-id"
export ARM_CLIENT_SECRET="your-client-secret"
export ARM_TENANT_ID="your-tenant-id"

# Test domain for creating users
export TEST_DOMAIN="yourdomain.com"
```

### Running Acceptance Tests

To run the acceptance tests:

```bash
TF_ACC=1 go test -v
```

To run a specific test:

```bash
TF_ACC=1 go test -v -run TestAccUserResource_Minimal
```

### Acceptance Test Structure

The acceptance tests are organized into several focused test functions:

1. `TestAccUserResource_Minimal` - Tests the minimal configuration
2. `TestAccUserResource_Maximal` - Tests the maximal configuration
3. `TestAccUserResource_MinimalToMaximal` - Tests updating from minimal to maximal
4. `TestAccUserResource_MaximalToMinimal` - Tests updating from maximal to minimal

## Important Notes

1. The acceptance tests create real users in your Microsoft 365 tenant. Make sure you have the necessary permissions and are using a test tenant or have cleanup procedures in place.

2. Passwords in the tests must meet Microsoft's password complexity requirements:
   - At least 8 characters
   - At least 3 of the following: lowercase letters, uppercase letters, numbers, and symbols
   - Not commonly used passwords

3. User principal names must be unique within your tenant. The tests use a random prefix with your domain to avoid conflicts.

4. Some attributes may be modified by Microsoft 365 after creation (e.g., mail attributes might be automatically populated). The tests are designed to handle these cases. 