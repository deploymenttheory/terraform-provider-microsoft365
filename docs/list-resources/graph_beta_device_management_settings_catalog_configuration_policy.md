---
page_title: "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy List Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Lists Settings Catalog configuration policies from Microsoft Intune using the /deviceManagement/configurationPolicies endpoint. This list resource is used to automatically retrieve all policies across multiple pages with advanced filtering capabilities for policy discovery and import. For full resource details including assignments and settings, use Terraform's import functionality with terraform plan -generate-config-out.
---

# microsoft365_graph_beta_device_management_settings_catalog_configuration_policy (List Resource)

Lists Settings Catalog configuration policies from Microsoft Intune using the `/deviceManagement/configurationPolicies` endpoint. This list resource is used to automatically retrieve all policies across multiple pages with advanced filtering capabilities for policy discovery and import. For full resource details including assignments and settings, use Terraform's import functionality with `terraform plan -generate-config-out`.

Lists Settings Catalog configuration policies from Microsoft Intune using the `/deviceManagement/configurationPolicies` endpoint. Supports filtering by name, platform, template family, assignment status, and custom OData queries.

List resources allow you to query and discover existing infrastructure without managing it. This is useful for:
- Finding policies for import into Terraform
- Discovering policies by criteria
- Auditing policy configuration
- Building dynamic configurations based on existing policies

## Microsoft Documentation

- [List configurationPolicies](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfigv2-devicemanagementconfigurationpolicy-list?view=graph-rest-beta)
- [deviceManagementConfigurationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this list resource:

**Required:**
- `DeviceManagementConfiguration.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.40.0-alpha | Experimental | Initial release |

## Example Usage

### List All Policies

```terraform
# List all Settings Catalog configuration policies
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "all" {
  provider = microsoft365
  config {}
}
```

### Filter by Name

```terraform
# List policies with "Kerberos" in the name
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "by_name" {
  provider = microsoft365
  config {
    name_filter = "Kerberos"
  }
}
```

### Filter by Platform

```terraform
# List policies for Windows 10 platform
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "by_platform" {
  provider = microsoft365
  config {
    platform_filter = ["windows10"]
  }
}
```

### Filter by Template Family

```terraform
# List policies from the baseline template family
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "by_template_family" {
  provider = microsoft365
  config {
    template_family_filter = "baseline"
  }
}
```

### Filter by Assignment Status

```terraform
# List only policies that have assignments
# Note: This checks actual assignments via API calls and may take 20-30 seconds for large tenants
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "assigned_only" {
  provider = microsoft365
  config {
    is_assigned_filter = true
  }
}
```

```terraform
# List only policies that do not have assignments
# Note: This checks actual assignments via API calls and may take 20-30 seconds for large tenants
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "unassigned_only" {
  provider = microsoft365
  config {
    is_assigned_filter = false
  }
}
```

### Combined Filters

```terraform
# List policies using multiple filters combined (AND logic)
# This example finds Windows 10 policies with "Defender" in the name
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "combined" {
  provider = microsoft365
  config {
    name_filter     = "Defender"
    platform_filter = ["windows10"]
  }
}
```

```terraform
# List assigned Edge policies (combining name filter with assignment check)
# This is efficient: name filter reduces results first, then assignment check runs on fewer policies
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "assigned_edge" {
  provider = microsoft365
  config {
    name_filter        = "Edge"
    is_assigned_filter = true
  }
}
```

### Custom OData Filters

#### Exact Match

```terraform
# Use custom OData filter for exact name match
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "exact_match" {
  provider = microsoft365
  config {
    odata_filter = "name eq '[Base] Prod | Windows - Settings Catalog | Kerberos ver1.0'"
  }
}
```

#### String Functions

```terraform
# Use OData startsWith function to find policies by name prefix
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "prod_policies" {
  provider = microsoft365
  config {
    odata_filter = "startsWith(name, '[Base] Prod')"
  }
}
```

#### Logical Operators

```terraform
# Use OData AND operator to combine multiple conditions
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "defender_windows10" {
  provider = microsoft365
  config {
    odata_filter = "contains(name, 'Defender') and platforms eq 'windows10'"
  }
}
```

```terraform
# Use OData OR operator to match multiple specific policies
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "kerberos_or_licensing" {
  provider = microsoft365
  config {
    odata_filter = "name eq '[Base] Prod | Windows - Settings Catalog | Kerberos ver1.0' or name eq '[Base] Prod | Windows - Settings Catalog | Licensing ver1.0'"
  }
}
```

#### Nested Properties

```terraform
# Use OData to filter by nested templateReference properties
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "baseline_templates" {
  provider = microsoft365
  config {
    odata_filter = "templateReference/templateFamily eq 'baseline'"
  }
}
```

#### Complex Queries

```terraform
# Use complex OData query with grouping and mixed AND/OR operators
# This finds Windows 10 policies that contain either "Edge" or "Defender" in the name
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "edge_or_defender_windows10" {
  provider = microsoft365
  config {
    odata_filter = "(contains(name, 'Edge') or contains(name, 'Defender')) and platforms eq 'windows10'"
  }
}
```

## Filter Behavior

- **API-level filters**: `name_filter`, `platform_filter`, `template_family_filter`, and `odata_filter` are applied at the Microsoft Graph API level.
- **Local filters**: `is_assigned_filter` checks actual policy assignments via the `/assignments` endpoint because the API's `isAssigned` field is unreliable.
- **Filter combination**: Multiple filters are combined using AND logic.

## OData Query Patterns

The `odata_filter` parameter supports standard OData query syntax:

### String Functions
- `contains(name, 'text')` - Partial match
- `startsWith(name, 'prefix')` - Prefix match
- `endsWith(name, 'suffix')` - Suffix match

### Comparison Operators
- `eq` - Equals
- `ne` - Not equals
- `gt` / `ge` - Greater than / Greater or equal
- `lt` / `le` - Less than / Less or equal

### Logical Operators
- `and` - Logical AND
- `or` - Logical OR
- `not` - Logical NOT

### Grouping
- Use parentheses for complex expressions: `(condition1 or condition2) and condition3`

### Nested Properties
- Use slash notation: `templateReference/templateFamily eq 'baseline'`

## Supported Platforms

The `platform_filter` parameter accepts the following values:
- `windows10`
- `iOS`
- `macOS`
- `android`
- `androidForWork`
- `linux`
- `unknownFutureValue`

## Template Families

Common `template_family_filter` values include:
- `baseline` - Security baselines
- `none` - Custom policies
- `endpointSecurityAntivirus`
- `endpointSecurityDiskEncryption`
- `endpointSecurityFirewall`
- `windowsOsRecoveryPolicies`
- `companyPortal`

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `is_assigned_filter` (Boolean) Filter policies by assignment status. Set to `true` to return only policies with assignments, `false` for policies without assignments. This filter queries the assignments endpoint for each policy (the API's `isAssigned` field is unreliable) and may take 20-30 seconds for large tenants.
- `name_filter` (String) Filter policies by name using partial matching. Supports the OData `contains` operator. Example: `name_filter = "Kerberos"` will match "[Base] Prod | Windows - Settings Catalog | Kerberos ver1.0".
- `odata_filter` (String) Advanced: Custom OData $filter query for complex filtering scenarios. Allows direct control over the API filter expression. Example: `odata_filter = "platforms eq 'windows10' and isAssigned eq true"`. When specified, this overrides individual filter parameters. See Microsoft Graph API documentation for supported operators and syntax.
- `platform_filter` (List of String) Filter policies by platform(s). Valid values: `none`, `android`, `iOS`, `macOS`, `windows10X`, `windows10`, `linux`, `unknownFutureValue`, `androidEnterprise`, `aosp`. Multiple platforms use OR logic. Example: `platform_filter = ["windows10", "macOS"]`.
- `template_family_filter` (String) Filter policies by template family. Valid values: `none`, `endpointSecurityAntivirus`, `endpointSecurityDiskEncryption`, `endpointSecurityFirewall`, `endpointSecurityEndpointDetectionAndResponse`, `endpointSecurityAttackSurfaceReduction`, `endpointSecurityAccountProtection`, `endpointSecurityApplicationControl`, `endpointSecurityEndpointPrivilegeManagement`, `enrollmentConfiguration`, `appQuietTime`, `baseline`, `unknownFutureValue`, `deviceConfigurationScripts`, `deviceConfigurationPolicies`, `windowsOsRecoveryPolicies`, `companyPortal`. Example: `template_family_filter = "baseline"`.

