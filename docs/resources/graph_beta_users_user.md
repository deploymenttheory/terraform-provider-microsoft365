---
page_title: "microsoft365_graph_beta_users_user Resource - terraform-provider-microsoft365"
subcategory: "Users"
description: |-
  Manages Microsoft 365 users using the /users endpoint. The user resource lets admins specify user preferences for languages and date/time formats for the user's primary Exchange mailboxes and Microsoft Entra profile. Permissions for this resource are complex and depend on the specific fields you wish tomanage. For more information, see the Microsoft Documentation. https://learn.microsoft.com/en-us/graph/api/user-update?view=graph-rest-beta&tabs=http#permissions-for-specific-scenarios
---

# microsoft365_graph_beta_users_user (Resource)

Manages Microsoft 365 users using the `/users` endpoint. The user resource lets admins specify user preferences for languages and date/time formats for the user's primary Exchange mailboxes and Microsoft Entra profile. Permissions for this resource are complex and depend on the specific fields you wish tomanage. For more information, see the Microsoft Documentation. https://learn.microsoft.com/en-us/graph/api/user-update?view=graph-rest-beta&tabs=http#permissions-for-specific-scenarios

## Microsoft Documentation

- [user resource type](https://learn.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-beta)
- [Create user](https://learn.microsoft.com/en-us/graph/api/user-post-users?view=graph-rest-beta)
- [Update user](https://learn.microsoft.com/en-us/graph/api/user-update?view=graph-rest-beta)
- [Delete user](https://learn.microsoft.com/en-us/graph/api/user-delete?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `User.EnableDisableAccount.All`, `User.ReadWrite.All`, `Directory.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.16.0-alpha | Experimental | Initial release |
| v0.36.0-alpha | Preview | Added support for manager_id and added custom security attributes support|

## Example Usage

### Minimal Example

```terraform
# Minimal example with only required properties
resource "microsoft365_graph_beta_users_user" "minimal_example" {
  display_name        = "John Doe"
  account_enabled     = true
  user_principal_name = "john.doe@contoso.com"
  mail_nickname       = "johndoe"
  hard_delete         = true
  password_profile = {
    password                           = "SecurePassword123!"
    force_change_password_next_sign_in = true
  }
}
```

### Maximal Example

```terraform
resource "microsoft365_graph_beta_users_user" "maximal" {
  account_enabled = true
  hard_delete     = true

  // Identity
  display_name        = "acc-test-user-maximal"
  given_name          = "Maximal"
  surname             = "User"
  user_principal_name = "acc-test-user-maximal@deploymenttheory.com"
  preferred_language  = "en-US"
  password_policies   = "DisablePasswordExpiration"

  // Age and Consent (for minor users)
  age_group                  = "NotAdult"
  consent_provided_for_minor = "Granted"

  // Job Information
  job_title          = "Senior Developer"
  company_name       = "Deployment Theory"
  department         = "Engineering"
  employee_id        = "1234567890"
  employee_type      = "full time"
  employee_hire_date = "2025-11-21T00:00:00Z"
  office_location    = "Building A"
  manager_id         = microsoft365_graph_beta_users_user.dependency_user.id

  // Contact Information
  city            = "Redmond"
  state           = "WA"
  country         = "US"
  street_address  = "123 street"
  postal_code     = "98052"
  usage_location  = "US"
  business_phones = ["+1 425-555-0100"]
  mobile_phone    = "+1 425-555-0101"
  mail            = "acc-test-user-maximal@deploymenttheory.com"
  fax_number      = "+1 425-555-0102"
  mail_nickname   = "acc-test-user-maximal"
  other_mails     = ["acc-test-user-maximal2.other@deploymenttheory.com"]

  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}
```

### Custom Security Attributes Example

```terraform
resource "microsoft365_graph_beta_users_user" "with_custom_security_attributes" {
  display_name        = "Custom Security Attributes User"
  user_principal_name = "custom.sec.user@deploymenttheory.com"
  mail_nickname       = "custom.sec.user"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }

  custom_security_attributes = [
    {
      attribute_set = "Engineering"
      attributes = [
        {
          name          = "Project"
          string_values = ["Baker", "Cascade"]
        },
        {
          name         = "LastTrainingDate"
          string_value = "2024-10-15"
        },
      ]
    },
    {
      attribute_set = "Marketing"
      attributes = [
        {
          name       = "IsContractor"
          bool_value = false
        }
      ]
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_enabled` (Boolean) Set to `true` if the account is enabled; otherwise, `false`. This property is required when a user is created.
- `display_name` (String) The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial, and last name. This property is required when a user is created and it cannot be cleared during updates. Maximum length is 256 characters.
- `mail_nickname` (String) The mail alias for the user. This property must be specified when a user is created. Maximum length is 64 characters.
- `password_profile` (Attributes) Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created. These fields are write-only and used only for initial user provisioning. Password management after user creation should be handled through proper identity management workflows, not Terraform. (see [below for nested schema](#nestedatt--password_profile))
- `user_principal_name` (String) The user principal name (UPN) of the user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. By convention, this should map to the user's email name. The general format is alias@domain, where the domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization. NOTE: This property can't contain accent characters. Only the following characters are allowed: A-Z, a-z, 0-9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see username policies.

### Optional

- `about_me` (String) A freeform text entry field for users to describe themselves.
- `age_group` (String) Sets the age group of the user. Allowed values: `null`, `Minor`, `NotAdult`, `Adult`. Refer to the legal age group property definitions for further information.
- `business_phones` (Set of String) The telephone numbers for the user. NOTE: Although it is a string collection, only one number can be set for this property.
- `city` (String) The city where the user is located. Maximum length is 128 characters.
- `company_name` (String) The name of the company that the user is associated with. This property can be useful for describing the company that an external user comes from. Maximum length is 64 characters.
- `consent_provided_for_minor` (String) Sets whether consent was obtained for minors. Allowed values: `null`, `Granted`, `Denied`, `NotRequired`. Refer to the legal age group property definitions for further information.
- `country` (String) The country or region where the user is located; for example, `US` or `UK`. Maximum length is 128 characters.
- `custom_security_attributes` (Attributes Set) An open complex type that holds the value of a custom security attribute that is assigned to a directory object. Nullable. Returned only on `$select`. Supports `$filter` (eq, ne, not, startsWith). The filter value is case-sensitive. To read this property, the calling app must be assigned the `CustomSecAttributeAssignment.Read.All` permission. To write this property, the calling app must be assigned the `CustomSecAttributeAssignment.ReadWrite.All` permission. (see [below for nested schema](#nestedatt--custom_security_attributes))
- `department` (String) The name of the department in which the user works. Maximum length is 64 characters.
- `employee_hire_date` (String) The date and time when the user was hired or will start work in case of a future hire, in ISO 8601 format and UTC.
- `employee_id` (String) The employee identifier assigned to the user by the organization. Maximum length is 16 characters.
- `employee_type` (String) Captures enterprise worker type. For example, `Employee`, `Contractor`, `Consultant`, or `Vendor`.
- `fax_number` (String) The fax number of the user.
- `given_name` (String) The given name (first name) of the user. Maximum length is 64 characters.
- `hard_delete` (Boolean) When `true`, the user will be permanently deleted (hard delete) during destroy. When `false` (default), the user will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. Note: This field defaults to `false` on import since the API does not return this value.
- `job_title` (String) The user's job title. Maximum length is 128 characters.
- `mail` (String) The SMTP address for the user, for example, `jeff@contoso.com`. Changes to this property also update the user's proxyAddresses collection to include the value as an SMTP address.
- `manager_id` (String) The user ID of the user's manager. Used to set the organizational hierarchy. To update the manager, provide the user ID of the new manager. To remove the manager, set this to an empty string.
- `mobile_phone` (String) The primary cellular telephone number for the user.
- `office_location` (String) The office location in the user's place of business. Maximum length is 128 characters.
- `on_premises_distinguished_name` (String) Contains the on-premises Active Directory `distinguished name` or `DN`. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only.
- `on_premises_domain_name` (String) Contains the on-premises `domainFQDN`, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only.
- `on_premises_immutable_id` (String) This property is used to associate an on-premises Active Directory user account to their Microsoft Entra user object. This property must be specified when creating a new user account in the Graph if you're using a federated domain for the user's userPrincipalName (UPN) property.
- `on_premises_last_sync_date_time` (String) Indicates the last time at which the object was synced with the on-premises directory; for example: `2013-02-16T03:04:54Z`. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.
- `on_premises_sam_account_name` (String) Contains the on-premises `sAMAccountName` synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only.
- `on_premises_security_identifier` (String) Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only.
- `on_premises_sync_enabled` (Boolean) `true` if this user object is currently being synced from an on-premises Active Directory (AD); otherwise, the user isn't being synced and can be managed in Microsoft Entra ID. Read-only. The value is `null` for cloud-only users.
- `on_premises_user_principal_name` (String) Contains the on-premises `userPrincipalName` synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only.
- `other_mails` (Set of String) A list of additional email addresses for the user; for example: `["bob@contoso.com", "Robert@fabrikam.com"]`. NOTE: This property can't contain accent characters. Maximum length per value is 250 characters.
- `password_policies` (String) Specifies password policies for the user. This value is an enumeration with one possible value being `DisableStrongPassword`, which allows weaker passwords than the default policy to be specified. `DisablePasswordExpiration` can also be specified. The two may be specified together; for example: `DisablePasswordExpiration, DisableStrongPassword`.
- `postal_code` (String) The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code. Maximum length is 40 characters.
- `preferred_data_location` (String) The preferred data location for the user. For more information, see OneDrive Online Multi-Geo.
- `preferred_language` (String) The preferred language for the user. The preferred language format is based on RFC 4646. The name combines an ISO 639 two-letter lowercase culture code associated with the language and an ISO 3166 two-letter uppercase subculture code associated with the country or region. Example: `en-US`, or `es-ES`.
- `preferred_name` (String) The preferred name for the user. **Not Supported.** This attribute returns an empty string.
- `proxy_addresses` (Set of String) Email addresses that also represent the user for the same mailbox. For example: `["SMTP: bob@contoso.com", "smtp: bob@sales.contoso.com"]`. Changes to the mail property also update this collection to include the value as an SMTP address. For more information, see mail and proxyAddresses properties. The proxy address prefixed with SMTP (capitalized) is the primary proxy address. This property can't contain accent characters. Read-only in Microsoft Graph; you can only update this property through the Microsoft 365 admin center.
- `security_identifier` (String) Security identifier (SID) of the user, used in Windows scenarios. Read-only.
- `show_in_address_list` (Boolean) `true` if the Outlook global address list should contain this user, otherwise `false`. If not set, this will be treated as `true`. For users invited through the invitation manager, this property will be set to `false`.
- `sign_in_sessions_valid_from_date_time` (String) Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph). If this happens, the application needs to acquire a new refresh token by making a request to the authorize endpoint. Read-only. Use revokeSignInSessions to reset.
- `state` (String) The state or province in the user's address. Maximum length is 128 characters.
- `street_address` (String) The street address of the user's place of business. Maximum length is 1024 characters.
- `surname` (String) The user's surname (family name or last name). Maximum length is 64 characters.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `usage_location` (String) A two-letter country code (ISO standard 3166). Required for users that are assigned licenses due to legal requirements to check for availability of services in countries. Examples include: `US`, `JP`, and `GB`. Not nullable.
- `user_type` (String) A string value that can be used to classify user types in your directory, such as `Member` and `Guest`.

### Read-Only

- `created_date_time` (String) The date and time the user was created, in ISO 8601 format and UTC. Read-only.
- `creation_type` (String) Indicates whether the user account was created through one of the following methods: As a regular school or work account (`null`), as an external account (`Invitation`), as a local account for an Azure Active Directory B2C tenant (`LocalAccount`), through self-service sign-up by an internal user using email verification (`EmailVerified`), or through self-service sign-up by an external user signing up through a link that is part of a user flow (`SelfServiceSignUp`). Read-only.
- `deleted_date_time` (String) The date and time the user was deleted. Read-only.
- `external_user_state` (String) For an external user invited to the tenant using the invitation API, this property represents the invited user's invitation status. For invited users, the state can be `PendingAcceptance` or `Accepted`, or `null` for all other users. Read-only.
- `external_user_state_change_date_time` (String) Shows the timestamp for the latest change to the externalUserState property. Read-only.
- `id` (String) The unique identifier for the user. Read-only.

<a id="nestedatt--password_profile"></a>
### Nested Schema for `password_profile`

Required:

- `force_change_password_next_sign_in` (Boolean, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) true if the user must change their password on the next login; otherwise false. This is a write-only field used only for initial provisioning.
- `password` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) The password for the user. This property is required when a user is created. This is a write-only field used only for initial provisioning - the API never returns password values.

Optional:

- `force_change_password_next_sign_in_with_mfa` (Boolean, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) If true, at next sign-in, the user must perform a multi-factor authentication (MFA) before being forced to change their password. The behavior is identical to forceChangePasswordNextSignIn except that the user is required to first perform a multi-factor authentication before password change. This is a write-only field used only for initial provisioning. Defaults to false if not specified.


<a id="nestedatt--custom_security_attributes"></a>
### Nested Schema for `custom_security_attributes`

Required:

- `attribute_set` (String) The name of the attribute set (e.g., `Engineering`, `Marketing`). This groups related custom security attributes together.
- `attributes` (Attributes Set) The collection of custom security attributes within this attribute set. (see [below for nested schema](#nestedatt--custom_security_attributes--attributes))

<a id="nestedatt--custom_security_attributes--attributes"></a>
### Nested Schema for `custom_security_attributes.attributes`

Required:

- `name` (String) The name of the custom security attribute.

Optional:

- `bool_value` (Boolean) The value if the attribute is a boolean type. Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.
- `int_value` (Number) The value if the attribute is a single-valued integer type. Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.
- `int_values` (Set of Number) The values if the attribute is a multi-valued integer type. Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.
- `string_value` (String) The value if the attribute is a single-valued string type. Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.
- `string_values` (Set of String) The values if the attribute is a multi-valued string type. Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Password Complexity**: Passwords must meet your organization's complexity requirements. Microsoft 365 typically requires a mix of uppercase, lowercase, numbers, and special characters.
- **Required Properties**: The minimum required properties to create a user are `display_name`, `account_enabled`, `user_principal_name`, `mail_nickname`, and `password_profile`.
- **User Principal Name**: The UPN must be unique across your tenant and follow the format `username@domain.com`.
- **Mail Nickname**: This value is used to generate the user's email address if a Microsoft 365 license is assigned.
- **Immutable Properties**: Some properties cannot be changed after creation or are managed by the system.
- **Password Management**: The actual password value is write-only and cannot be read back from the API.
- **Force Password Change**: Use `force_change_password_next_sign_in` to require users to change their password at next login.
- **Identities**: The identities property allows configuring federated authentication methods for the user.
- **External Users**: For guest users, additional properties like `external_user_state` may be relevant.
- **On-Premises Sync**: Properties prefixed with `on_premises_` are typically managed by Azure AD Connect and shouldn't be modified directly.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import an existing user into Terraform
# The import ID format is: {user_id}[:hard_delete=true|false]
#
# Where:
# - {user_id} is the unique identifier for the user
# - hard_delete is optional (defaults to false for soft delete only)

# Basic import (hard_delete defaults to false - soft delete only)
terraform import microsoft365_graph_beta_users_user.example 00000000-0000-0000-0000-000000000000

# Import with hard_delete enabled (permanently deletes on terraform destroy)
terraform import microsoft365_graph_beta_users_user.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
``` 