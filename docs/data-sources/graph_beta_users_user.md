---
page_title: "microsoft365_graph_beta_users_user Data Source - terraform-provider-microsoft365"
subcategory: "Users"

description: |-
  Retrieves Microsoft Entra Users using the /users endpoint. Supports flexible lookup by object ID, display name, employee ID, given name, user principal name, on-premises immutable ID, on-premises distinguished name, or a custom OData query. Can also list all users in the tenant.
---

# microsoft365_graph_beta_users_user (Data Source)

Retrieves Microsoft Entra Users using the `/users` endpoint. Supports flexible lookup by object ID, display name, employee ID, given name, user principal name, on-premises immutable ID, on-premises distinguished name, or a custom OData query. Can also list all users in the tenant.

## Microsoft Documentation

- [user resource type](https://learn.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-beta)
- [List users](https://learn.microsoft.com/en-us/graph/api/user-list?view=graph-rest-beta)
- [Get a user](https://learn.microsoft.com/en-us/graph/api/user-get?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `User.Read.All`
- `Directory.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.55.0-alpha | Experimental | Initial release of users_user data source |

## Example Usage

### Example 1: List all users in the tenant

```terraform
# Example 1: List all users in the tenant
# Results are always returned in the `items` list.

data "microsoft365_graph_beta_users_user" "all" {
  list_all = true
}

# Total number of users returned
output "all_users_count" {
  description = "The total number of users in the tenant"
  value       = length(data.microsoft365_graph_beta_users_user.all.items)
}

# Basic projection of every user
output "all_users_basic_info" {
  description = "Basic information about all users"
  value = [
    for user in data.microsoft365_graph_beta_users_user.all.items : {
      id                  = user.id
      display_name        = user.display_name
      user_principal_name = user.user_principal_name
      mail                = user.mail
    }
  ]
}
```

### Example 2: Look up a user by object_id (with a representative set of outputs)

```terraform
# Example 2: Look up a user by object_id
# This example shows a representative set of the available output attributes.

data "microsoft365_graph_beta_users_user" "by_object_id" {
  object_id = "12345678-1234-1234-1234-123456789012"
}

output "user_id" {
  description = "The unique identifier for the user object"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].id
}

output "display_name" {
  description = "The name displayed in the address book for the user"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].display_name
}

output "user_principal_name" {
  description = "The user principal name (UPN) of the user"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].user_principal_name
}

output "mail" {
  description = "The SMTP address for the user"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].mail
}

output "job_title" {
  description = "The user's job title"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].job_title
}

output "department" {
  description = "The name of the department in which the user works"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].department
}

output "account_enabled" {
  description = "Whether the account is enabled"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].account_enabled
}

output "on_premises_sync_enabled" {
  description = "Whether the user is synced from an on-premises Active Directory"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].on_premises_sync_enabled
}
```

### Example 3: Look up a user by display_name

```terraform
# Example 3: Look up a user by display_name

data "microsoft365_graph_beta_users_user" "by_display_name" {
  display_name = "John Doe"
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_display_name.items[0].id
}

output "user_principal_name" {
  value = data.microsoft365_graph_beta_users_user.by_display_name.items[0].user_principal_name
}
```

### Example 4: Look up a user by employee_id

```terraform
# Example 4: Look up a user by employee_id

data "microsoft365_graph_beta_users_user" "by_employee_id" {
  employee_id = "100200"
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_employee_id.items[0].id
}

output "display_name" {
  value = data.microsoft365_graph_beta_users_user.by_employee_id.items[0].display_name
}
```

### Example 5: Look up users by given_name

```terraform
# Example 5: Look up users by given_name (first name)
# A given name may match multiple users, all returned in `items`.

data "microsoft365_graph_beta_users_user" "by_given_name" {
  given_name = "John"
}

output "matched_users" {
  description = "All users matching the given name"
  value = [
    for user in data.microsoft365_graph_beta_users_user.by_given_name.items : {
      id                  = user.id
      display_name        = user.display_name
      user_principal_name = user.user_principal_name
    }
  ]
}
```

### Example 6: Look up a user by user_principal_name

```terraform
# Example 6: Look up a user by user_principal_name (UPN)

data "microsoft365_graph_beta_users_user" "by_upn" {
  user_principal_name = "user@contoso.com"
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_upn.items[0].id
}

output "display_name" {
  value = data.microsoft365_graph_beta_users_user.by_upn.items[0].display_name
}

output "mail" {
  value = data.microsoft365_graph_beta_users_user.by_upn.items[0].mail
}
```

### Example 7: Look up a user by on_premises_immutable_id

```terraform
# Example 7: Look up a user by on_premises_immutable_id (sourceAnchor)
# Useful for correlating cloud users with on-premises Active Directory accounts.

data "microsoft365_graph_beta_users_user" "by_immutable_id" {
  on_premises_immutable_id = "T0AbQ29udG9zb1VzZXI="
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_immutable_id.items[0].id
}

output "display_name" {
  value = data.microsoft365_graph_beta_users_user.by_immutable_id.items[0].display_name
}

output "on_premises_sam_account_name" {
  value = data.microsoft365_graph_beta_users_user.by_immutable_id.items[0].on_premises_sam_account_name
}
```

### Example 8: Look up a user by on_premises_distinguished_name

```terraform
# Example 8: Look up a user by on_premises_distinguished_name (DN)

data "microsoft365_graph_beta_users_user" "by_dn" {
  on_premises_distinguished_name = "CN=John Doe,OU=Users,DC=contoso,DC=com"
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_dn.items[0].id
}

output "display_name" {
  value = data.microsoft365_graph_beta_users_user.by_dn.items[0].display_name
}

output "on_premises_domain_name" {
  value = data.microsoft365_graph_beta_users_user.by_dn.items[0].on_premises_domain_name
}
```

### Example 9: Look up users using a custom OData query

```terraform
# Example 9: Look up users using a custom OData query
# Use this for advanced filtering when the standard lookup attributes don't fit.

data "microsoft365_graph_beta_users_user" "by_odata_query" {
  odata_query = "accountEnabled eq true and userType eq 'Member'"
}

output "enabled_members" {
  description = "All enabled member users matching the query"
  value = [
    for user in data.microsoft365_graph_beta_users_user.by_odata_query.items : {
      id                  = user.id
      display_name        = user.display_name
      user_principal_name = user.user_principal_name
    }
  ]
}

# Example: More complex OData query using a function
data "microsoft365_graph_beta_users_user" "by_odata_startswith" {
  odata_query = "startswith(displayName, 'A') and accountEnabled eq true"
}

output "users_starting_with_a_count" {
  value = length(data.microsoft365_graph_beta_users_user.by_odata_startswith.items)
}
```

## Notes

- Exactly one of `object_id`, `display_name`, `employee_id`, `given_name`, `user_principal_name`, `on_premises_immutable_id`, `on_premises_distinguished_name`, `odata_query`, or `list_all` must be specified.
- These lookup attributes are mutually exclusive.
- Results are always returned in the `items` list, even when a single user is matched. Reference a single result with `items[0]`.
- For advanced filtering scenarios, use the `odata_query` attribute with custom OData filter expressions.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) The display name of the user. Conflicts with other lookup attributes.
- `employee_id` (String) The employee identifier assigned to the user by the organization. Conflicts with other lookup attributes.
- `given_name` (String) The given name (first name) of the user. Conflicts with other lookup attributes.
- `list_all` (Boolean) Retrieve all users in the tenant. Conflicts with specific lookup attributes.
- `object_id` (String) The unique object identifier of the user in Microsoft Entra ID. Conflicts with other lookup attributes.
- `odata_query` (String) Custom OData filter expression for advanced queries (e.g., `accountEnabled eq true and userType eq 'Member'`). Conflicts with specific lookup attributes.
- `on_premises_distinguished_name` (String) The on-premises Active Directory distinguished name (DN) of the user. Conflicts with other lookup attributes.
- `on_premises_immutable_id` (String) The on-premises immutable ID (sourceAnchor) used to associate an on-premises Active Directory user account. Conflicts with other lookup attributes.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `user_principal_name` (String) The user principal name (UPN) of the user. Conflicts with other lookup attributes.

### Read-Only

- `id` (String) The unique identifier for the data source. This is a placeholder attribute required by Terraform.
- `items` (Attributes List) List of users matching the query criteria. (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `about_me` (String) A freeform text entry field for the user to describe themselves.
- `account_enabled` (Boolean) true if the account is enabled; otherwise, false.
- `age_group` (String) Sets the age group of the user (null, Minor, NotAdult, Adult).
- `business_phones` (List of String) The telephone numbers for the user.
- `city` (String) The city where the user is located.
- `company_name` (String) The company name which the user is associated with.
- `consent_provided_for_minor` (String) Sets whether consent was obtained for minors (null, Granted, Denied, NotRequired).
- `country` (String) The country/region where the user is located.
- `created_date_time` (String) The date and time the user was created.
- `creation_type` (String) Indicates whether the user account was created as a regular school or work account, an external account, etc.
- `deleted_date_time` (String) The date and time the user was deleted.
- `department` (String) The name of the department in which the user works.
- `display_name` (String) The name displayed in the address book for the user.
- `employee_hire_date` (String) The date and time when the user was hired or will start work.
- `employee_id` (String) The employee identifier assigned to the user by the organization.
- `employee_type` (String) Captures enterprise worker type (Employee, Contractor, Consultant, Vendor, etc.).
- `external_user_state` (String) For an external user invited to the tenant, this represents the invitation status (PendingAcceptance, Accepted).
- `external_user_state_change_date_time` (String) Shows the timestamp for the latest change to the external_user_state property.
- `fax_number` (String) The fax number of the user.
- `given_name` (String) The given name (first name) of the user.
- `id` (String) The unique identifier for the user object.
- `job_title` (String) The user's job title.
- `mail` (String) The SMTP address for the user.
- `mail_nickname` (String) The mail alias for the user.
- `mobile_phone` (String) The primary cellular telephone number for the user.
- `office_location` (String) The office location in the user's place of business.
- `on_premises_distinguished_name` (String) Contains the on-premises Active Directory distinguished name (DN).
- `on_premises_domain_name` (String) Contains the on-premises domainFQDN, also called dnsDomainName synchronized from the on-premises directory.
- `on_premises_immutable_id` (String) The on-premises immutable ID (sourceAnchor) used to associate an on-premises Active Directory user account.
- `on_premises_last_sync_date_time` (String) Indicates the last time at which the object was synced with the on-premises directory.
- `on_premises_sam_account_name` (String) Contains the on-premises samAccountName synchronized from the on-premises directory.
- `on_premises_security_identifier` (String) Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud.
- `on_premises_sync_enabled` (Boolean) true if this user object is currently being synced from an on-premises Active Directory.
- `on_premises_user_principal_name` (String) Contains the on-premises userPrincipalName synchronized from the on-premises directory.
- `other_mails` (List of String) A list of additional email addresses for the user.
- `password_policies` (String) Specifies password policies for the user (DisableStrongPassword, DisablePasswordExpiration).
- `postal_code` (String) The postal code for the user's postal address.
- `preferred_data_location` (String) The preferred data location for the user.
- `preferred_language` (String) The preferred language for the user, in ISO 639-1 format.
- `preferred_name` (String) The preferred name for the user.
- `proxy_addresses` (List of String) For example: ["SMTP: bob@contoso.com", "smtp: bob@sales.contoso.com"].
- `security_identifier` (String) Security identifier (SID) of the user, used in Windows scenarios.
- `show_in_address_list` (Boolean) Do not use in Microsoft Graph. Manage this property through the Microsoft 365 admin center instead.
- `sign_in_sessions_valid_from_date_time` (String) Any refresh tokens or session tokens issued before this time are invalid.
- `state` (String) The state or province in the user's address.
- `street_address` (String) The street address of the user's place of business.
- `surname` (String) The user's surname (family name or last name).
- `usage_location` (String) A two letter country code (ISO standard 3166), required for users that are assigned licenses.
- `user_principal_name` (String) The user principal name (UPN) of the user.
- `user_type` (String) A string value that can be used to classify user types in your directory (Member, Guest).
