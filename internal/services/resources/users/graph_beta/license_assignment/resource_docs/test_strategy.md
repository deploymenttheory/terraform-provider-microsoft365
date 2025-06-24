# User License Assignment Resource Tests

This directory contains both unit tests and acceptance tests for the User License Assignment resource.

## Unit Tests

Unit tests can be run without any special setup:

```bash
go test -v
```

These tests use mocked responses and don't require any real Microsoft 365 resources.

### Unit Test Structure

The unit tests are organized into several focused test functions:

1. `TestUnitUserLicenseAssignmentResource_Minimal` - Tests minimal configuration with only required attributes
2. `TestUnitUserLicenseAssignmentResource_Maximal` - Tests maximal configuration with all possible attributes
3. `TestUnitUserLicenseAssignmentResource_Create` - Tests creating multiple resources simultaneously
4. `TestUnitUserLicenseAssignmentResource_Update` - Tests updating a resource's configuration
5. `TestUnitUserLicenseAssignmentResource_MinimalToMaximal` - Tests updating from minimal to maximal configurations
6. `TestUnitUserLicenseAssignmentResource_Import` - Tests importing an existing resource
7. `TestUnitUserLicenseAssignmentResource_RemoveLicenses` - Tests license removal functionality
8. `TestUnitUserLicenseAssignmentResource_Error` - Tests error handling

To run a specific unit test:

```bash
go test -v -run TestUnitUserLicenseAssignmentResource_Update
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

# Test user IDs (object IDs of existing users)
export TEST_USER_ID_1="user-object-id-1"
export TEST_USER_ID_2="user-object-id-2"

# Test license SKU IDs (must be available in your tenant)
export TEST_LICENSE_SKU_ID_1="license-sku-id-1"
export TEST_LICENSE_SKU_ID_2="license-sku-id-2"

# Optional: Service plan ID to test disabling specific plans
export TEST_SERVICE_PLAN_ID="service-plan-id"
```

### Running Acceptance Tests

To run the acceptance tests:

```bash
TF_ACC=1 go test -v
```

To run a specific test:

```bash
TF_ACC=1 go test -v -run TestAccUserLicenseAssignmentResource_Minimal
```

### Acceptance Test Structure

The acceptance tests are organized into several focused test functions:

1. `TestAccUserLicenseAssignmentResource_Minimal` - Tests the minimal configuration
2. `TestAccUserLicenseAssignmentResource_Maximal` - Tests the maximal configuration
3. `TestAccUserLicenseAssignmentResource_MinimalToMaximal` - Tests updating from minimal to maximal
4. `TestAccUserLicenseAssignmentResource_MaximalToMinimal` - Tests updating from maximal to minimal
5. `TestAccUserLicenseAssignmentResource_RemoveLicenses` - Tests the license removal functionality

### Finding License SKU IDs and Service Plan IDs

To find license SKU IDs and service plan IDs in your tenant, you can use Microsoft Graph Explorer or PowerShell:
