---
page_title: "microsoft365_graph_beta_device_management_windows_platform_script List Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Lists Windows PowerShell scripts from Microsoft Intune using the /deviceManagement/deviceManagementScripts endpoint. This list resource is used to automatically retrieve all scripts across multiple pages with advanced filtering capabilities for script discovery and import. For full resource details, use Terraform's import functionality with terraform plan -generate-config-out.
---

# microsoft365_graph_beta_device_management_windows_platform_script (List Resource)

Lists Windows PowerShell scripts from Microsoft Intune using the `/deviceManagement/deviceManagementScripts` endpoint. This list resource is used to automatically retrieve all scripts across multiple pages with advanced filtering capabilities for script discovery and import. For full resource details, use Terraform's import functionality with `terraform plan -generate-config-out`.

Lists Windows PowerShell scripts from Microsoft Intune using the `/deviceManagement/deviceManagementScripts` endpoint. Supports filtering by display name, file name, run as account, assignment status, and custom OData queries.

List resources allow you to query and discover existing infrastructure without managing it. This is useful for:
- Finding scripts for import into Terraform
- Discovering scripts by criteria
- Auditing script configuration
- Building dynamic configurations based on existing scripts

## Microsoft Documentation

- [List deviceManagementScripts](https://learn.microsoft.com/en-us/graph/api/intune-shared-devicemanagementscript-list?view=graph-rest-beta)
- [deviceManagementScript resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicemanagementscript?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this list resource:

**Required:**
- `DeviceManagementConfiguration.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.45.0-alpha | Experimental | Initial release |

## Example Usage

### List All Scripts

```terraform
provider "microsoft365" {}

# List all Windows platform scripts
list "microsoft365_graph_beta_device_management_windows_platform_script" "all" {
  provider = microsoft365
  config {}
}
```

### Filter by Display Name

```terraform
provider "microsoft365" {}

# List Windows platform scripts filtered by display name
list "microsoft365_graph_beta_device_management_windows_platform_script" "by_display_name" {
  provider = microsoft365
  config {
    display_name_filter = "Baseline"
  }
}
```

### Filter by Run As Account

```terraform
provider "microsoft365" {}

# List Windows platform scripts running as system
list "microsoft365_graph_beta_device_management_windows_platform_script" "system_scripts" {
  provider = microsoft365
  config {
    run_as_account_filter = "system"
  }
}
```

### Filter by Assignment Status

```terraform
provider "microsoft365" {}

# List only Windows platform scripts that have assignments
list "microsoft365_graph_beta_device_management_windows_platform_script" "assigned_only" {
  provider = microsoft365
  config {
    is_assigned_filter = true
  }
}
```

```terraform
provider "microsoft365" {}

# List only Windows platform scripts without assignments
list "microsoft365_graph_beta_device_management_windows_platform_script" "unassigned_only" {
  provider = microsoft365
  config {
    is_assigned_filter = false
  }
}
```

### Combined Filters

```terraform
provider "microsoft365" {}

# List Windows platform scripts with combined filters
list "microsoft365_graph_beta_device_management_windows_platform_script" "combined" {
  provider = microsoft365
  config {
    display_name_filter   = "Setup"
    run_as_account_filter = "system"
  }
}
```

```terraform
provider "microsoft365" {}

# List assigned Windows platform scripts running as system
list "microsoft365_graph_beta_device_management_windows_platform_script" "assigned_system" {
  provider = microsoft365
  config {
    run_as_account_filter = "system"
    is_assigned_filter    = true
  }
}
```

### Custom OData Filters

```terraform
provider "microsoft365" {}

# Use OData filter with exact match on display name
list "microsoft365_graph_beta_device_management_windows_platform_script" "exact_match" {
  provider = microsoft365
  config {
    odata_filter = "displayName eq 'Windows Baseline Setup'"
  }
}
```

```terraform
provider "microsoft365" {}

# Use OData filter with AND logic
list "microsoft365_graph_beta_device_management_windows_platform_script" "odata_and" {
  provider = microsoft365
  config {
    odata_filter = "runAsAccount eq 'system' and contains(fileName, '.ps1')"
  }
}
```

```terraform
provider "microsoft365" {}

# Use OData filter with OR logic
list "microsoft365_graph_beta_device_management_windows_platform_script" "odata_or" {
  provider = microsoft365
  config {
    odata_filter = "contains(displayName, 'Baseline') or contains(displayName, 'Security')"
  }
}
```

```terraform
provider "microsoft365" {}

# List Windows platform scripts with complex OData filter
list "microsoft365_graph_beta_device_management_windows_platform_script" "odata_complex" {
  provider = microsoft365
  config {
    odata_filter = "runAsAccount eq 'system' and contains(displayName, 'Baseline') and contains(fileName, 'ps1')"
  }
}
```

## Filter Behavior

- **API-level filters**: `display_name_filter`, `file_name_filter`, `run_as_account_filter`, and `odata_filter` are applied at the Microsoft Graph API level.
- **Local filters**: `is_assigned_filter` checks actual script assignments via the `/assignments` endpoint because the API's `isAssigned` field is unreliable.
- **Filter combination**: Multiple filters are combined using AND logic.

## OData Query Patterns

The `odata_filter` parameter supports standard OData query syntax:

### String Functions
- `contains(displayName, 'text')` - Partial match on display name
- `contains(fileName, 'text')` - Partial match on file name

### Comparison Operators
- `eq` - Equals (e.g., `runAsAccount eq 'system'`)
- `ne` - Not equals

### Logical Operators
- `and` - Logical AND
- `or` - Logical OR
- `not` - Logical NOT

### Grouping
- Use parentheses for complex expressions: `(condition1 or condition2) and condition3`

## Run As Account Values

The `run_as_account_filter` parameter accepts the following values:
- `system` - Scripts running in system context
- `user` - Scripts running in user context

## Assignment Filtering

The `is_assigned_filter` parameter:
- Queries the `/assignments` endpoint for each script individually
- Returns accurate assignment status (unlike the API's built-in `isAssigned` field)
- May take 20-30 seconds for tenants with many scripts
- Accepts boolean values: `true` (assigned only) or `false` (unassigned only)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name_filter` (String) Filter scripts by display name using partial matching. Supports the OData `contains` operator. Example: `display_name_filter = "Baseline"` will match "Windows Baseline Script".
- `file_name_filter` (String) Filter scripts by file name using partial matching. Supports the OData `contains` operator. Example: `file_name_filter = "setup.ps1"` will match scripts with "setup.ps1" in the filename.
- `is_assigned_filter` (Boolean) Filter scripts by assignment status. Set to `true` to return only scripts with assignments, `false` for scripts without assignments. This filter queries the assignments endpoint for each script (the API's `isAssigned` field is unreliable) and may take 20-30 seconds for large tenants.
- `odata_filter` (String) Advanced: Custom OData $filter query for complex filtering scenarios. Allows direct control over the API filter expression. Example: `odata_filter = "runAsAccount eq 'system' and contains(displayName, 'Baseline')"`. When specified, this overrides individual filter parameters. See Microsoft Graph API documentation for supported operators and syntax.
- `run_as_account_filter` (String) Filter scripts by execution context. Valid values: `system`, `user`. Example: `run_as_account_filter = "system"` returns only scripts running as system.

