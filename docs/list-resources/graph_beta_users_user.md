---
page_title: "microsoft365_graph_beta_users_user List Resource - terraform-provider-microsoft365"
subcategory: "Users"

description: |-
  Lists users from Microsoft Entra ID using the /users endpoint. This list resource is used to automatically retrieve all users across multiple pages with advanced filtering capabilities for user discovery and import. For full resource details, use Terraform's import functionality with terraform plan -generate-config-out.
---

# microsoft365_graph_beta_users_user (List Resource)

Lists users from Microsoft Entra ID using the `/users` endpoint. This list resource is used to automatically retrieve all users across multiple pages with advanced filtering capabilities for user discovery and import. For full resource details, use Terraform's import functionality with `terraform plan -generate-config-out`.

Lists users from Microsoft Entra ID using the `/users` endpoint. Supports filtering by display name, user principal name, account status, user type, and custom OData queries.

List resources allow you to query and discover existing infrastructure without managing it. This is useful for:
- Finding users for import into Terraform
- Discovering users by criteria
- Auditing user configuration
- Building dynamic configurations based on existing users

## Microsoft Documentation

- [List users](https://learn.microsoft.com/en-us/graph/api/user-list?view=graph-rest-beta)
- [user resource type](https://learn.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this list resource:

**Required:**
- `User.Read.All`
- `Directory.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.45.0-alpha | Experimental | Initial release |

## Example Usage

### List All Users

```terraform
# List all users
list "microsoft365_graph_beta_users_user" "all" {
  provider = microsoft365
  config {}
}
```

### Filter by Display Name

```terraform
# List users with display name starting with "John"
list "microsoft365_graph_beta_users_user" "by_display_name" {
  provider = microsoft365
  config {
    display_name_filter = "John"
  }
}
```

### Filter by User Principal Name

```terraform
# List users with UPN starting with "admin"
list "microsoft365_graph_beta_users_user" "by_upn" {
  provider = microsoft365
  config {
    user_principal_name_filter = "admin"
  }
}
```

### Filter by Account Status

```terraform
# List only enabled user accounts
list "microsoft365_graph_beta_users_user" "enabled_only" {
  provider = microsoft365
  config {
    account_enabled_filter = true
  }
}
```

```terraform
# List only disabled user accounts
list "microsoft365_graph_beta_users_user" "disabled_only" {
  provider = microsoft365
  config {
    account_enabled_filter = false
  }
}
```

### Filter by User Type

```terraform
# List only member users (excluding guests)
list "microsoft365_graph_beta_users_user" "members_only" {
  provider = microsoft365
  config {
    user_type_filter = "Member"
  }
}
```

```terraform
# List only guest users
list "microsoft365_graph_beta_users_user" "guests_only" {
  provider = microsoft365
  config {
    user_type_filter = "Guest"
  }
}
```

### Combined Filters

```terraform
# List enabled member users
list "microsoft365_graph_beta_users_user" "enabled_members" {
  provider = microsoft365
  config {
    account_enabled_filter = true
    user_type_filter       = "Member"
  }
}
```

### Custom OData Filters

#### Exact Match

```terraform
# Find specific user by exact UPN match
list "microsoft365_graph_beta_users_user" "exact_upn" {
  provider = microsoft365
  config {
    odata_filter = "userPrincipalName eq 'admin@contoso.com'"
  }
}
```

#### String Functions

```terraform
# Find users with display name starting with specific prefix
list "microsoft365_graph_beta_users_user" "name_prefix" {
  provider = microsoft365
  config {
    odata_filter = "startsWith(displayName, 'Adele')"
  }
}
```

#### Logical Operators

```terraform
# Find enabled users with specific job title
list "microsoft365_graph_beta_users_user" "enabled_with_title" {
  provider = microsoft365
  config {
    odata_filter = "accountEnabled eq true and jobTitle eq 'Manager'"
  }
}
```

```terraform
# Find users with either of two job titles
list "microsoft365_graph_beta_users_user" "multiple_titles" {
  provider = microsoft365
  config {
    odata_filter = "jobTitle eq 'Manager' or jobTitle eq 'Director'"
  }
}
```

#### Complex Queries

```terraform
# Complex query combining multiple conditions
list "microsoft365_graph_beta_users_user" "complex" {
  provider = microsoft365
  config {
    odata_filter = "(userType eq 'Member' and accountEnabled eq true) and (startsWith(userPrincipalName, 'admin') or startsWith(userPrincipalName, 'svc'))"
  }
}
```

## Filter Behavior

- **API-level filters**: `display_name_filter`, `user_principal_name_filter`, `account_enabled_filter`, `user_type_filter`, and `odata_filter` are applied at the Microsoft Graph API level.
- **Filter combination**: Multiple filters are combined using AND logic.
- **String matching**: Note that `display_name_filter` and `user_principal_name_filter` use `startsWith` operator as the Microsoft Graph `/users` endpoint does not support `contains` for partial matching.

## OData Query Patterns

The `odata_filter` parameter supports standard OData query syntax:

### String Functions
- `startsWith(property, 'prefix')` - Prefix match (e.g., `startsWith(displayName, 'John')`)
- **Note**: The `/users` endpoint does NOT support `contains()` for substring matching

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

## Supported User Types

The `user_type_filter` parameter accepts the following values:
- `Member` - Regular organizational users
- `Guest` - External/guest users

## Important Notes

### String Filter Limitations
The Microsoft Graph `/users` endpoint has specific limitations on string filtering:
- **Supported**: `startsWith()` for prefix matching
- **NOT Supported**: `contains()` for substring matching
- **NOT Supported**: `endsWith()` for suffix matching

If you need substring matching, use `odata_filter` with multiple `startsWith()` conditions combined with `or`, or retrieve all users and filter client-side.

### Performance Considerations
- Filtering is performed server-side at the Microsoft Graph API level for optimal performance
- Use specific filters to reduce the result set size
- The `odata_filter` parameter provides maximum flexibility for complex queries

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `account_enabled_filter` (Boolean) Filter users by account status. Set to `true` to return only enabled accounts, `false` for disabled accounts. Example: `account_enabled_filter = true`.
- `display_name_filter` (String) Filter users by display name using prefix matching. Supports the OData `startsWith` operator. Example: `display_name_filter = "John"` will match "John Smith" and "Johnny Doe".
- `odata_filter` (String) Advanced: Custom OData $filter query for complex filtering scenarios. Allows direct control over the API filter expression. Example: `odata_filter = "accountEnabled eq true and userType eq 'Member'"`. When specified, this overrides individual filter parameters. See Microsoft Graph API documentation for supported operators and syntax.
- `user_principal_name_filter` (String) Filter users by user principal name using partial matching. Supports the OData `startsWith` operator. Example: `user_principal_name_filter = "admin"` will match "admin@contoso.com" and "admin2@contoso.com".
- `user_type_filter` (String) Filter users by type. Valid values: `Member`, `Guest`. Example: `user_type_filter = "Member"` returns only member users.
