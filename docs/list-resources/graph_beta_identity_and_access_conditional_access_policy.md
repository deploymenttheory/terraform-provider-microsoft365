---
page_title: "microsoft365_graph_beta_identity_and_access_conditional_access_policy List Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"

description: |-
  Lists Conditional Access policies from Microsoft Entra ID using the /identity/conditionalAccess/policies endpoint. This list resource is used to automatically retrieve all policies across multiple pages with advanced filtering capabilities for policy discovery and import. For full resource details, use Terraform's import functionality with terraform plan -generate-config-out.
---

# microsoft365_graph_beta_identity_and_access_conditional_access_policy (List Resource)

Lists Conditional Access policies from Microsoft Entra ID using the `/identity/conditionalAccess/policies` endpoint. This list resource is used to automatically retrieve all policies across multiple pages with advanced filtering capabilities for policy discovery and import. For full resource details, use Terraform's import functionality with `terraform plan -generate-config-out`.

Lists Conditional Access policies from Microsoft Entra ID using the `/identity/conditionalAccess/policies` endpoint. Supports filtering by display name, state, and custom OData queries.

List resources allow you to query and discover existing infrastructure without managing it. This is useful for:
- Finding policies for import into Terraform
- Discovering policies by criteria
- Auditing policy configuration
- Building dynamic configurations based on existing policies

## Microsoft Documentation

- [List conditionalAccessPolicies](https://learn.microsoft.com/en-us/graph/api/conditionalaccessroot-list-policies?view=graph-rest-beta)
- [conditionalAccessPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesspolicy?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this list resource:

**Required:**
- `Policy.Read.All`
- `Policy.Read.ConditionalAccess`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.45.0-alpha | Experimental | Initial release |

## Example Usage

### List All Policies

```terraform
# List all Conditional Access policies
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "all" {
  provider = microsoft365
  config {}
}
```

### Filter by Display Name

```terraform
# List policies with "MFA" in the display name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "mfa_policies" {
  provider = microsoft365
  config {
    display_name_filter = "MFA"
  }
}
```

### Filter by State

```terraform
# List only enabled policies
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "enabled_only" {
  provider = microsoft365
  config {
    state_filter = "enabled"
  }
}
```

```terraform
# List only disabled policies
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "disabled_only" {
  provider = microsoft365
  config {
    state_filter = "disabled"
  }
}
```

```terraform
# List policies in report-only mode
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "report_only" {
  provider = microsoft365
  config {
    state_filter = "enabledForReportingButNotEnforced"
  }
}
```

### Combined Filters

```terraform
# List enabled policies with "Admin" in the name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "enabled_admin_policies" {
  provider = microsoft365
  config {
    display_name_filter = "Admin"
    state_filter        = "enabled"
  }
}
```

### Custom OData Filters

#### Exact Match

```terraform
# Find specific policy by exact display name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "exact_name" {
  provider = microsoft365
  config {
    odata_filter = "displayName eq 'Require MFA for Administrators'"
  }
}
```

#### String Functions

```terraform
# Find policies with specific text in display name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "contains_guest" {
  provider = microsoft365
  config {
    odata_filter = "contains(displayName, 'Guest')"
  }
}
```

#### Logical Operators

```terraform
# Find enabled policies with specific display name pattern
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "enabled_with_pattern" {
  provider = microsoft365
  config {
    odata_filter = "state eq 'enabled' and contains(displayName, 'Baseline')"
  }
}
```

```terraform
# Find policies that are either enabled or in report-only mode
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "active_policies" {
  provider = microsoft365
  config {
    odata_filter = "state eq 'enabled' or state eq 'enabledForReportingButNotEnforced'"
  }
}
```

```terraform
# Find policies that are not disabled
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "not_disabled" {
  provider = microsoft365
  config {
    odata_filter = "state ne 'disabled'"
  }
}
```

#### Complex Queries

```terraform
# Complex query combining multiple conditions
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "complex" {
  provider = microsoft365
  config {
    odata_filter = "(state eq 'enabled' or state eq 'enabledForReportingButNotEnforced') and (contains(displayName, 'MFA') or contains(displayName, 'Admin'))"
  }
}
```

## Filter Behavior

- **API-level filters**: `display_name_filter`, `state_filter`, and `odata_filter` are applied at the Microsoft Graph API level.
- **Filter combination**: Multiple filters are combined using AND logic.
- **String matching**: The `display_name_filter` uses the `contains` operator for partial matching.

## OData Query Patterns

The `odata_filter` parameter supports standard OData query syntax:

### String Functions
- `contains(property, 'text')` - Substring match (e.g., `contains(displayName, 'MFA')`)
- `startsWith(property, 'prefix')` - Prefix match
- `endsWith(property, 'suffix')` - Suffix match

### Comparison Operators
- `eq` - Equals
- `ne` - Not equals

### Logical Operators
- `and` - Logical AND
- `or` - Logical OR
- `not` - Logical NOT

### Grouping
- Use parentheses for complex expressions: `(condition1 or condition2) and condition3`

## Supported Policy States

The `state_filter` parameter accepts the following values:
- `enabled` - Policy is actively enforced
- `disabled` - Policy is not enforced
- `enabledForReportingButNotEnforced` - Policy is in report-only mode (logs but doesn't block)

## Important Notes

### Report-Only Mode
Policies in `enabledForReportingButNotEnforced` state are useful for:
- Testing policy impact before enforcement
- Monitoring user sign-in patterns
- Validating policy logic without blocking users

### Policy Discovery
Use the list resource to:
- Discover existing policies before creating new ones
- Identify policies that need updates
- Audit policy configuration across your tenant
- Find policies by naming conventions

### Performance Considerations
- Filtering is performed server-side at the Microsoft Graph API level for optimal performance
- Use specific filters to reduce the result set size
- The `odata_filter` parameter provides maximum flexibility for complex queries

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name_filter` (String) Filter policies by display name using partial matching. Supports the OData `contains` operator. Example: `display_name_filter = "MFA"` will match "Require MFA for Admins".
- `odata_filter` (String) Advanced: Custom OData $filter query for complex filtering scenarios. Allows direct control over the API filter expression. Example: `odata_filter = "state eq 'enabled' and displayName eq 'MFA Policy'"`. When specified, this overrides individual filter parameters. See Microsoft Graph API documentation for supported operators and syntax.
- `state_filter` (String) Filter policies by state. Valid values: `enabled`, `disabled`, `enabledForReportingButNotEnforced`. Example: `state_filter = "enabled"` returns only enabled policies.
