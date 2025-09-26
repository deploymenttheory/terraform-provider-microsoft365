# Provider Comparison in light of terraform-provider-msgraph

In July 2025 microsoft released the [terraform-provider-msgraph](https://github.com/hashicorp/terraform-provider-msgraph) partner provider. This provider is developed by Microsoft and is the official provider for Microsoft Graph API. When choosing between the two providers, it's important to consider the scope and approach of each provider, there are some distinct differences between the two providers and the approaches taken for interacting with Microsoft M365.

## Scope

This projects aim and scope is to cover all aspects of Microsoft 365 workloads including:

- msgraph
- teams.microsoft.com
- exchange.microsoft.com
- sharepointonline.com
- security.microsoft.com
- undocumented api endpoints / related microsoft microservices
- utilities for handing metadata useful for resource lifecycle creation and lifecycle

Out of scope for this project are:

- entra ID resources(it's managed by the azureAD provider)
- operations taken by a secondary service upon primary service (e.g. defining security configuration for managed devices via Defender for Endpoint. When intune handles it.)

As such the scope is broader than the terraform-provider-msgraph provider.

## API Interactions and Developer Experience

The fundamental difference between these providers lies in their approach to API interactions and the level of abstraction provided to users. This section provides a detailed analysis of how each provider handles the complexity of the Microsoft Graph API and the resulting impact on developer productivity.

### Understanding Microsoft Graph API Complexity

Before comparing the providers, it's important to understand the inherent complexity of the Microsoft Graph API that both providers must handle:

**API Characteristics:**
- **Multi-step Operations**: Many business operations require multiple sequential API calls across different endpoints
- **Eventual Consistency**: Microsoft Graph uses eventual consistency, meaning writes may not be immediately visible in reads
- **Complex Request Bodies**: Many operations require nested JSON structures with specific OData type annotations
- **State Dependencies**: Resource creation often requires reading state from multiple related endpoints
- **Error Handling Complexity**: Different endpoints have different retry requirements and error response formats
- **Assignment Management**: Device management resources require separate assignment API calls to be functional

**Real-world Example - Windows Update Ring:**
Creating a functional Windows Update Ring involves:
1. POST to `/deviceManagement/deviceConfigurations` (create policy)
2. POST to `/deviceManagement/deviceConfigurations/{id}/assign` (assign to groups)
3. GET with `$expand=assignments` (verify complete state)
4. Handle eventual consistency between configuration and assignment APIs
5. Manage complex assignment target structures with filter relationships

This complexity exists regardless of which provider you use - the question is how each provider handles it.

### terraform-provider-msgraph Approach: Direct API Exposure

The terraform-provider-msgraph provider is a **thin wrapper** around the Microsoft Graph API that directly exposes the API's complexity to users. This approach prioritizes flexibility and API completeness over developer experience.

#### Resource Architecture
The provider offers four generic resource types that map directly to HTTP operations:

- **`msgraph_resource`** - Generic resource for any Graph API endpoint (POST/GET/PATCH/DELETE)
- **`msgraph_resource_action`** - For performing actions on resources (POST actions like assign, send, etc.)
- **`msgraph_resource_collection`** - For managing reference collections ($ref endpoints)
- **`msgraph_update_resource`** - For updating subset of properties (PATCH operations)

#### Developer Requirements
**Graph API Expertise Required:**
- **Endpoint Knowledge**: Deep understanding of Graph API URL structures (`/deviceManagement/deviceConfigurations`, `/groups/{id}/assignLicense`, etc.)
- **OData Mastery**: Proficiency with OData query parameters (`$expand`, `$select`, `$filter`, `$top`, etc.)
- **Schema Understanding**: Knowledge of complex JSON schema including OData type annotations (`@odata.type`)
- **HTTP Method Mapping**: Understanding when to use POST vs PATCH vs GET for different operations
- **Response Parsing**: Ability to write JMESPath queries for extracting data from API responses

**Manual Orchestration Burden:**
- **Multi-Resource Management**: Complex operations require coordinating multiple Terraform resources
- **Dependency Management**: Manual `depends_on` declarations to ensure proper operation sequencing
- **State Management**: Separate data sources needed to read complete resource state
- **Error Handling**: No built-in retry logic for Graph API-specific issues (throttling, eventual consistency)
- **Lifecycle Management**: Manual handling of create/update/delete operation differences

#### Development Workflow Impact
**Configuration Complexity:**
- Users must translate business requirements into multiple low-level API operations
- Each business operation may require 3-5 separate Terraform resources
- Configuration files become verbose and API-centric rather than business-focused
- Changes require understanding the impact across multiple interdependent resources

**Debugging and Maintenance:**
- Issues require deep Graph API troubleshooting knowledge
- Error messages are raw HTTP/JSON responses without business context
- Updates require understanding how each API endpoint handles partial modifications
- Testing requires knowledge of Graph API test patterns and mock strategies

### This Provider's Approach: Business-Focused Abstraction

This provider uses the Kiota-generated GraphSDKs built from Microsoft's schema to interact with the Graph API, but provides **significant abstraction** that shields users from API complexity while maintaining full functionality.

#### Resource Architecture
The provider offers purpose-built resources that represent complete business operations:

- **Domain-Specific Resources**: Each resource type (e.g., `microsoft365_graph_beta_device_management_windows_update_ring`) represents a complete business workflow
- **Strongly-Typed Schemas**: Terraform configuration uses intuitive field names and validation rather than raw JSON
- **Embedded Relationships**: Related operations (assignments, settings, dependencies) are managed within single resources
- **Business Logic Integration**: Resources understand the business context and handle complex workflows automatically

#### Developer Experience Benefits
**Business Domain Focus:**
- **Declarative Configuration**: Users describe desired end-state rather than API operation sequences
- **Intuitive Field Names**: Configuration uses business-friendly terminology (`allow_windows11_upgrade` vs `allowWindows11Upgrade`)
- **Built-in Validation**: Schema validation catches configuration errors before API calls
- **Contextual Documentation**: Each field includes business context and impact descriptions
- **IDE Support**: Strongly-typed schemas enable autocomplete, validation, and documentation in IDEs

**Automatic Complexity Management:**
- **API Call Chaining**: Single resource operations automatically trigger multiple coordinated API calls
- **State Synchronization**: Built-in retry logic handles eventual consistency across all related endpoints  
- **Error Translation**: Raw Graph API errors are translated into actionable business context
- **Lifecycle Optimization**: Create/update/delete operations use the most efficient API patterns automatically
- **Dependency Resolution**: Resources automatically handle prerequisite operations and timing

#### Technical Implementation
**Under the Hood Automation:**
- **Multi-Endpoint Coordination**: Single Terraform operations may trigger anything from 1-10 Graph API calls across different endpoints
- **Eventual Consistency Handling**: Built-in wait/retry patterns for Microsoft's asynchronous operations
- **Assignment Management**: Automatic construction of complex assignment target objects with proper type annotations
- **State Reconciliation**: Resources automatically detect and correct configuration drift
- **Error Recovery**: Intelligent retry logic with exponential backoff for transient API issues

**Development Workflow Impact:**
**Configuration Simplicity:**
- Business requirements map directly to single resource declarations
- Changes are made at the business logic level rather than API operation level
- Configuration files are concise and focused on business outcomes
- Testing focuses on business functionality rather than API mechanics

**Operational Benefits:**
- **Faster Development**: Developers work at business abstraction level
- **Reduced Errors**: API complexity is encapsulated and tested within the provider
- **Easier Maintenance**: admin focused changes don't require API expertise
- **Better Debugging**: Error messages provide business context and suggested resolutions
- **Consistent Patterns**: All resources follow similar patterns regardless of underlying API complexity

## Detailed Comparison Examples

### Example 1: Group License Assignment

#### This Provider
```hcl
resource "microsoft365_graph_beta_group_license_assignment" "sales_team" {
  group_id = "2243c326-937g-53f0-c9df-2e68f106b901"

  add_licenses = [
    {
      sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
      disabled_plans = ["efb87545-963c-4e0d-99df-69c6916d9eb0"]
    }
  ]

  remove_licenses = ["f30db892-07e9-47e9-9163-06ecf6a79d2c"]
}
```

**Automatically handles:**
- POST /groups/{id}/assignLicense with complex request body construction
- GET /groups/{id}?$select=assignedLicenses for state management  
- Differential license assignment logic for updates
- Proper cleanup of managed licenses on deletion
- Built-in retry for throttling and consistency issues

#### MSGraph Provider Equivalent
```hcl
resource "msgraph_resource" "group_license_assignment" {
  url = "groups/2243c326-937g-53f0-c9df-2e68f106b901/assignLicense"
  
  body = {
    addLicenses = [
      {
        skuId = "6fd2c87f-b296-42f0-b197-1e91e994b900"
        disabledPlans = ["efb87545-963c-4e0d-99df-69c6916d9eb0"]
      }
    ]
    removeLicenses = ["f30db892-07e9-47e9-9163-06ecf6a79d2c"]
  }
}

# Additional resources needed for proper state management
data "msgraph_resource" "group_licenses" {
  url = "groups/2243c326-937g-53f0-c9df-2e68f106b901"
  query_parameters = { "$select" = ["assignedLicenses"] }
}
```

**Users must:**
- Know the exact Graph API endpoint structure
- Understand the assignLicense API's request body format
- Manually handle state reading with separate data sources
- Implement differential update logic themselves
- Handle cleanup logic manually

### Example 2: Windows Update Ring with Assignments

This example demonstrates complex multi-API operations involving device configuration creation, assignment management, and state synchronization across multiple endpoints.

#### This Provider
```hcl
resource "microsoft365_graph_beta_device_management_windows_update_ring" "corporate_updates" {
  display_name                                 = "Corporate Windows Update Ring"
  description                                  = "Managed updates for corporate devices"
  microsoft_update_service_allowed             = true
  drivers_excluded                             = false
  quality_updates_deferral_period_in_days      = 7
  feature_updates_deferral_period_in_days      = 14
  allow_windows11_upgrade                      = false
  skip_checks_before_restart                   = true
  automatic_update_mode                        = "autoInstallAndRebootAtScheduledTime"
  business_ready_updates_only                  = "businessReadyOnly"
  delivery_optimization_mode                   = "httpWithPeeringNat"
  prerelease_features                          = "settingsOnly"
  update_weeks                                 = "firstWeek"
  active_hours_start                           = "09:00:00"
  active_hours_end                             = "17:00:00"
  user_pause_access                            = "disabled"
  user_windows_update_scan_access              = "disabled"
  update_notification_level                    = "defaultNotifications"
  feature_updates_rollback_window_in_days      = 10
  role_scope_tag_ids                           = ["0", "1"]

  deadline_settings = {
    deadline_for_feature_updates_in_days = 7
    deadline_for_quality_updates_in_days = 2
    deadline_grace_period_in_days        = 1
    postpone_reboot_until_after_deadline = true
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    },
    {
      type      = "groupAssignmentTarget"
      group_id  = "55555555-5555-5555-5555-555555555555"
      filter_id = "66666666-6666-6666-6666-666666666666"
      filter_type = "include"
    }
  ]
}
```

**Automatically handles behind the scenes:**

**Create Operation:**
1. **POST** `/deviceManagement/deviceConfigurations` - Creates the Windows Update ring configuration
2. **POST** `/deviceManagement/deviceConfigurations/{id}/assign` - Assigns the policy to specified groups/filters
3. **GET** `/deviceManagement/deviceConfigurations/{id}?$expand=assignments` - Reads back complete state with assignments
4. Built-in retry logic for eventual consistency between configuration and assignment APIs

**Read Operation:**
1. **GET** `/deviceManagement/deviceConfigurations/{id}?$expand=assignments` - Single call retrieves policy + assignments
2. Type validation to ensure resource is `WindowsUpdateForBusinessConfiguration`
3. Automatic state mapping from complex Graph API response to Terraform schema

**Update Operation:**
1. **PATCH** `/deviceManagement/deviceConfigurations/{id}` - Updates the policy configuration
2. **POST** `/deviceManagement/deviceConfigurations/{id}/assign` - Updates assignments (replaces existing)
3. **GET** `/deviceManagement/deviceConfigurations/{id}?$expand=assignments` - Reads back updated state
4. Differential assignment management with proper target construction

**Delete Operation:**
1. **DELETE** `/deviceManagement/deviceConfigurations/{id}` - Single call removes policy and assignments
2. Assignments are automatically cleaned up by the API

#### MSGraph Provider Equivalent
```hcl
# Step 1: Create the Windows Update Ring policy
resource "msgraph_resource" "windows_update_ring" {
  url = "deviceManagement/deviceConfigurations"
  
  body = {
    "@odata.type" = "microsoft.graph.windowsUpdateForBusinessConfiguration"
    displayName = "Corporate Windows Update Ring"
    description = "Managed updates for corporate devices"
    microsoftUpdateServiceAllowed = true
    driversExcluded = false
    qualityUpdatesDeferralPeriodInDays = 7
    featureUpdatesDeferralPeriodInDays = 14
    allowWindows11Upgrade = false
    skipChecksBeforeRestart = true
    automaticUpdateMode = "autoInstallAndRebootAtScheduledTime"
    businessReadyUpdatesOnly = "businessReadyOnly"
    deliveryOptimizationMode = "httpWithPeeringNat"
    prereleaseFeatures = "settingsOnly"
    updateWeeks = "firstWeek"
    activeHoursStart = "09:00:00"
    activeHoursEnd = "17:00:00"
    userPauseAccess = "disabled"
    userWindowsUpdateScanAccess = "disabled"
    updateNotificationLevel = "defaultNotifications"
    featureUpdatesRollbackWindowInDays = 10
    roleScopeTagIds = ["0", "1"]
    deadlineForFeatureUpdatesInDays = 7
    deadlineForQualityUpdatesInDays = 2
    deadlineGracePeriodInDays = 1
    postponeRebootUntilAfterDeadline = true
  }
  
  response_export_values = {
    id = "id"
    all = "@"
  }
}

# Step 2: Manually create assignments (separate API call)
resource "msgraph_resource_action" "assign_update_ring" {
  resource_url = "deviceManagement/deviceConfigurations/${msgraph_resource.windows_update_ring.id}"
  action       = "assign"
  method       = "POST"
  
  body = {
    assignments = [
      {
        "@odata.type" = "microsoft.graph.deviceConfigurationAssignment"
        target = {
          "@odata.type" = "microsoft.graph.groupAssignmentTarget"
          groupId = "44444444-4444-4444-4444-444444444444"
        }
      },
      {
        "@odata.type" = "microsoft.graph.deviceConfigurationAssignment"
        target = {
          "@odata.type" = "microsoft.graph.groupAssignmentTarget"
          groupId = "55555555-5555-5555-5555-555555555555"
          deviceAndAppManagementAssignmentFilterId = "66666666-6666-6666-6666-666666666666"
          deviceAndAppManagementAssignmentFilterType = "include"
        }
      }
    ]
  }
  
  depends_on = [msgraph_resource.windows_update_ring]
}

# Step 3: Separate read for complete state (policy + assignments)  
data "msgraph_resource" "update_ring_with_assignments" {
  url = "deviceManagement/deviceConfigurations/${msgraph_resource.windows_update_ring.id}"
  query_parameters = {
    "$expand" = ["assignments"]
  }
  
  depends_on = [msgraph_resource_action.assign_update_ring]
}

# Step 4: Updates require manual orchestration
resource "msgraph_update_resource" "update_ring_policy" {
  url = "deviceManagement/deviceConfigurations/${msgraph_resource.windows_update_ring.id}"
  body = {
    # Users must manually construct PATCH body with only changed fields
    description = "Updated managed updates for corporate devices"
  }
}

# Step 5: Assignment updates require separate action
resource "msgraph_resource_action" "update_assignments" {
  resource_url = "deviceManagement/deviceConfigurations/${msgraph_resource.windows_update_ring.id}"
  action       = "assign"
  method       = "POST"
  
  body = {
    # Users must reconstruct entire assignments array for updates
    assignments = [/* complete assignment reconstruction required */]
  }
  
  depends_on = [msgraph_update_resource.update_ring_policy]
}
```

**Users must:**
- Understand complex Graph API OData types (`microsoft.graph.windowsUpdateForBusinessConfiguration`, etc.)
- Know the exact API schema including nested deadline settings structure  
- Manually orchestrate multiple API calls across 2-3 different resources
- Handle assignment construction with proper target types and filter relationships
- Understand the dependency chain between policy creation and assignment
- Implement proper error handling and retry logic across multiple resources
- Use separate data sources to read complete state including assignments
- Manually manage the lifecycle of policy updates vs assignment updates
- Understand Graph API expand syntax for reading related data

## Key Differentiators

| Aspect | This Provider | MSGraph Provider |
|--------|--------------|------------------|
| **Abstraction Level** | High-level business operations | Low-level Graph API wrapper |
| **API Knowledge Required** | Minimal - focused on business intent | Extensive - must understand Graph API |
| **Configuration Complexity** | Simple, declarative | Complex, API-centric |
| **Multi-API Operations** | Automatic chaining | Manual orchestration |
| **State Management** | Built-in with retries | Manual implementation |
| **Error Handling** | Comprehensive, contextual | Basic Graph API errors |
| **Resource Lifecycle** | Complete CRUD automation | Manual CRUD construction |
| **Type Safety** | Strongly typed schemas | Dynamic JSON bodies |
| **Learning Curve** | Terraform + business domain | Terraform + Graph API + OData |

## When to Choose Which

### Choose This Provider When:
- You want Infrastructure-as-Code for M365 without deep Graph API knowledge
- You need complex multi-step operations handled automatically
- You prioritize developer productivity and declarative configuration
- You want built-in best practices for error handling and retries
- Your team focuses on business outcomes, not API intricacies

### Choose MSGraph Provider When:
- You need Microsoft's official support
- You're already have strong familiarity with Graph API
- You need maximum flexibility to construct custom API calls  
- You prefer thin abstractions over opinionated frameworks
- You're building simple, single-API-call resources

## Support

One of the primary distinctions between the two providers is that the support for the terraform-provider-msgraph provider is provided by Microsoft. This provider is community supported and is not officially supported by Microsoft. Depending on your use case, and the support you require, this may be a consideration in your choice of provider. However, there's nothing to stop a tf configuration containing both depending on your use case.
