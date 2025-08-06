# Role Assignment Acceptance Test Fixes and Improvements

## Bugs Fixed

### 1. **Incorrect Graph Client Initialization**
- **Problem**: Using `mocks.TestAccProvider.Meta().(*mocks.TestClient)` instead of proper Graph client
- **Fix**: Changed to `acceptance.TestGraphClient()` following the terms_and_conditions pattern
- **Impact**: Proper Graph API client initialization for acceptance tests

### 2. **Missing Dependencies for Group Creation**
- **Problem**: Tests were trying to use user data sources that may not exist
- **Fix**: Created `resource_dependencies.tf` with predefined security groups
- **Impact**: Reliable test dependencies that don't require external user accounts

### 3. **Incomplete Error Handling in Destroy Function**
- **Problem**: Basic error checking without proper Graph error parsing
- **Fix**: Added comprehensive error handling using `errors.GraphError()` pattern
- **Impact**: Better error detection and more reliable destroy verification

### 4. **Missing Dependency Loading in Config Functions**
- **Problem**: Test configurations not loading dependencies properly
- **Fix**: Updated all config functions to load dependencies + config pattern
- **Impact**: Proper resource dependency resolution in tests

## Improvements Made

### 1. **Multiple Maximal Test Types**
Created separate test configurations for each scope type:
- **AllLicensedUsers** (`resource_all_users.tf`) - Uses Endpoint Security Manager role
- **AllDevices** (`resource_all_devices.tf`) - Uses Help Desk Operator role  
- **ResourceScopes** (`resource_resource_scopes.tf`) - Uses Policy and Profile manager role
- **General Maximal** (`resource_maximal.tf`) - Uses Application Manager role

### 2. **Built-in Role ID Usage**
Each test uses different built-in role IDs to keep dependencies simple:
- `fa1d7878-e8cb-41a1-8254-0142355c9f84` - Read Only Operator (minimal)
- `0bd113fe-6be5-400c-a28f-ae5553f9c0be` - Policy and Profile manager (resource scopes)
- `9e0cc482-82df-4ab2-a24c-0c23a3f52e1e` - Help Desk Operator (all devices)
- `c56d53a2-73d0-4502-b6bd-4a9d3dba28d5` - Endpoint Security Manager (all users)
- `c1d9fcbb-cba5-40b0-bf6b-527006585f4b` - Application Manager (maximal)

### 3. **Comprehensive Group Dependencies**
Created 4 test groups with different purposes:
- **Policy Managers** - For policy-related role assignments
- **Device Administrators** - For device management role assignments
- **Application Managers** - For application-related role assignments
- **Security Operators** - For security-related role assignments

### 4. **Enhanced Test Coverage**
Added specific test for each scope type:
- `TestAccRoleAssignmentResource_AllUsersScope` - Tests AllLicensedUsers scope
- `TestAccRoleAssignmentResource_AllDevicesScope` - Tests AllDevices scope
- `TestAccRoleAssignmentResource_ResourceScopes` - Tests ResourceScopes with specific IDs

### 5. **Proper File Organization**
Following the terms_and_conditions pattern:
```
tests/terraform/acceptance/
├── resource_dependencies.tf      # Group dependencies
├── resource_minimal.tf          # Minimal test config
├── resource_maximal.tf          # General maximal config
├── resource_all_users.tf        # AllLicensedUsers scope test
├── resource_all_devices.tf      # AllDevices scope test
└── resource_resource_scopes.tf  # ResourceScopes test
```

### 6. **Additional Response Files**
Created acceptance-specific response files:
- `get_role_assignment_maximal_acceptance.json`
- `get_role_assignment_resource_scopes_acceptance.json`
- `get_role_assignment_all_users_acceptance.json`
- `get_role_assignment_all_devices_acceptance.json`

## Test Scenarios Covered

1. **Lifecycle Test**: Create minimal → Import → Update to maximal
2. **Resource Scopes Test**: Specific resource scope assignments
3. **All Devices Test**: All devices scope assignments
4. **All Users Test**: All licensed users scope assignments

## Configuration Examples

Each configuration uses groups as members and different built-in roles to ensure:
- No external dependencies on user accounts
- Different role permissions for different test scenarios
- Realistic group-based role assignments
- Comprehensive scope type coverage

The fixes ensure the acceptance tests follow the established patterns in the codebase and provide comprehensive coverage of all role assignment scenarios.