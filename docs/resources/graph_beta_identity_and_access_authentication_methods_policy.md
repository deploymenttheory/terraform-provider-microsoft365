---
page_title: "microsoft365_graph_beta_identity_and_access_authentication_methods_policy Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages the tenant-wide Microsoft Entra authentication methods policy using the Microsoft Graph beta /policies/authenticationMethodsPolicy endpoint and its authenticationMethodConfigurations. The policy controls which authentication methods users can register and use, the registration campaign, system-preferred multifactor authentication and suspicious activity reporting.
  This is a singleton resource: one policy exists per tenant. create adopts the existing policy and applies the configuration; update applies changes. Every setting and every method block is optional. A setting or method block that is not configured keeps its current tenant value and is never written. Each changed method configuration is updated with its own PATCH request, so method configurations that are not in the configuration are left untouched. On destroy, Terraform only removes the resource from state and leaves the tenant policy as it is. Manage only one instance per tenant across all Terraform states.
  This version manages the email, temporary_access_pass, sms, software_oath and hardware_oath method configurations. The fido2, microsoftAuthenticator, voice, x509Certificate, qrCodePin, password, verifiableCredentials, federatedIdentityCredential and external method configurations are not managed by this resource and are never changed by it.
---

# microsoft365_graph_beta_identity_and_access_authentication_methods_policy (Resource)

Manages the tenant-wide Microsoft Entra authentication methods policy using the Microsoft Graph beta `/policies/authenticationMethodsPolicy` endpoint and its `authenticationMethodConfigurations`. The policy controls which authentication methods users can register and use, the registration campaign, system-preferred multifactor authentication and suspicious activity reporting.

This is a **singleton resource**: one policy exists per tenant. `create` adopts the existing policy and applies the configuration; `update` applies changes. Every setting and every method block is optional. A setting or method block that is not configured keeps its current tenant value and is never written. Each changed method configuration is updated with its own PATCH request, so method configurations that are not in the configuration are left untouched. On `destroy`, Terraform only removes the resource from state and leaves the tenant policy as it is. Manage only one instance per tenant across all Terraform states.

This version manages the `email`, `temporary_access_pass`, `sms`, `software_oath` and `hardware_oath` method configurations. The `fido2`, `microsoftAuthenticator`, `voice`, `x509Certificate`, `qrCodePin`, `password`, `verifiableCredentials`, `federatedIdentityCredential` and external method configurations are not managed by this resource and are never changed by it.

## Microsoft Documentation

- [Manage authentication methods for Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity/authentication/concept-authentication-methods-manage)
- [authenticationMethodsPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodspolicy?view=graph-rest-beta)
- [Get authenticationMethodsPolicy](https://learn.microsoft.com/en-us/graph/api/authenticationmethodspolicy-get?view=graph-rest-beta)
- [Update authenticationMethodsPolicy](https://learn.microsoft.com/en-us/graph/api/authenticationmethodspolicy-update?view=graph-rest-beta)
- [Update authenticationMethodConfiguration](https://learn.microsoft.com/en-us/graph/api/authenticationmethodconfiguration-update?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**

- `Policy.Read.AuthenticationMethod`
- `Policy.ReadWrite.AuthenticationMethod`

**Optional:**

- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v1.7.0 | Experimental | Initial release |

## Example Usage

### Adopt the Policy Without Changes

```terraform
# Adopt the tenant authentication methods policy without changing it.
# Every setting and method configuration is read from the tenant.
resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "this" {
}
```

### Target Email One-Time Passcodes at a Group

```terraform
# Enable email one-time passcodes for one group only and leave every other method unchanged.
resource "microsoft365_graph_beta_groups_group" "email_otp" {
  display_name     = "Email OTP users"
  mail_nickname    = "email-otp-users"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "this" {
  email = {
    state                              = "enabled"
    allow_external_id_to_use_email_otp = "default"
    include_targets = [
      {
        id          = microsoft365_graph_beta_groups_group.email_otp.id
        target_type = "group"
      }
    ]
    exclude_targets = []
  }
}
```

### Complete Configuration

```terraform
# Manage the policy-level settings and every supported method configuration.
resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "this" {
  policy_migration_state = "migrationComplete"

  registration_enforcement = {
    authentication_methods_registration_campaign = {
      snooze_duration_in_days                    = 3
      enforce_registration_after_allowed_snoozes = true
      state                                      = "enabled"
      include_targets = [
        {
          id                             = "all_users"
          target_type                    = "group"
          targeted_authentication_method = "microsoftAuthenticator"
        }
      ]
      exclude_targets = []
    }
  }

  report_suspicious_activity_settings = {
    state                = "enabled"
    voice_reporting_code = 0
    include_target = {
      id          = "all_users"
      target_type = "group"
    }
  }

  system_credential_preferences = {
    state = "enabled"
    include_targets = [
      {
        id          = "all_users"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }

  email = {
    state                              = "enabled"
    allow_external_id_to_use_email_otp = "enabled"
    include_targets = [
      {
        id          = "all_users"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }

  temporary_access_pass = {
    state                       = "enabled"
    default_lifetime_in_minutes = 60
    default_length              = 8
    minimum_lifetime_in_minutes = 60
    maximum_lifetime_in_minutes = 480
    is_usable_once              = false
    include_targets = [
      {
        id          = "all_users"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }

  sms = {
    state = "disabled"
    include_targets = [
      {
        id                    = "all_users"
        target_type           = "group"
        is_usable_for_sign_in = false
      }
    ]
    exclude_targets = []
  }

  software_oath = {
    state = "enabled"
    include_targets = [
      {
        id          = "all_users"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }

  hardware_oath = {
    state = "disabled"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `email` (Attributes) Email one-time passcode authentication method configuration (`Email`). When the block is not configured, the method configuration is read from the tenant and is never written. (see [below for nested schema](#nestedatt--email))
- `hardware_oath` (Attributes) Hardware OATH token authentication method configuration (`HardwareOath`). When the block is not configured, the method configuration is read from the tenant and is never written. (see [below for nested schema](#nestedatt--hardware_oath))
- `policy_migration_state` (String) The state of migration of the authentication methods policy from the legacy multifactor authentication and self-service password reset (SSPR) policies. Possible values are: `preMigration`, `migrationInProgress`, `migrationComplete`. When not configured, the current tenant value is kept.
- `registration_enforcement` (Attributes) Enforces registration at sign-in time. When not configured, the current tenant value is kept. (see [below for nested schema](#nestedatt--registration_enforcement))
- `report_suspicious_activity_settings` (Attributes) Settings that let users report suspicious multifactor authentication prompts. When not configured, the current tenant value is kept. (see [below for nested schema](#nestedatt--report_suspicious_activity_settings))
- `sms` (Attributes) Text message authentication method configuration (`Sms`). When the block is not configured, the method configuration is read from the tenant and is never written. (see [below for nested schema](#nestedatt--sms))
- `software_oath` (Attributes) Third-party software OATH token authentication method configuration (`SoftwareOath`). When the block is not configured, the method configuration is read from the tenant and is never written. (see [below for nested schema](#nestedatt--software_oath))
- `system_credential_preferences` (Attributes) System-preferred multifactor authentication settings, which prompt users for the most secure method they registered. When not configured, the current tenant value is kept. (see [below for nested schema](#nestedatt--system_credential_preferences))
- `temporary_access_pass` (Attributes) Temporary Access Pass authentication method configuration (`TemporaryAccessPass`). When the block is not configured, the method configuration is read from the tenant and is never written. (see [below for nested schema](#nestedatt--temporary_access_pass))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `description` (String) The description of the policy. Read-only.
- `display_name` (String) The display name of the policy. Read-only.
- `id` (String) The fixed singleton identifier `authenticationMethodsPolicy`.

<a id="nestedatt--email"></a>
### Nested Schema for `email`

Required:

- `state` (String) Whether the method is enabled in the tenant. Possible values are: `enabled`, `disabled`.

Optional:

- `allow_external_id_to_use_email_otp` (String) Whether external users can use email one-time passcodes for sign-in. Possible values are: `default`, `enabled`, `disabled`. When not configured, the current tenant value is kept.
- `exclude_targets` (Attributes Set) Users and groups excluded from the method. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--email--exclude_targets))
- `include_targets` (Attributes Set) Users and groups that can use email one-time passcodes. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--email--include_targets))

<a id="nestedatt--email--exclude_targets"></a>
### Nested Schema for `email.exclude_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.


<a id="nestedatt--email--include_targets"></a>
### Nested Schema for `email.include_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.

Optional:

- `is_registration_required` (Boolean) Whether users in the target must register the method. Defaults to `false`.



<a id="nestedatt--hardware_oath"></a>
### Nested Schema for `hardware_oath`

Required:

- `state` (String) Whether the method is enabled in the tenant. Possible values are: `enabled`, `disabled`.

Optional:

- `exclude_targets` (Attributes Set) Users and groups excluded from the method. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--hardware_oath--exclude_targets))
- `include_targets` (Attributes Set) Users and groups that can use hardware OATH tokens. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--hardware_oath--include_targets))

<a id="nestedatt--hardware_oath--exclude_targets"></a>
### Nested Schema for `hardware_oath.exclude_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.


<a id="nestedatt--hardware_oath--include_targets"></a>
### Nested Schema for `hardware_oath.include_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.

Optional:

- `is_registration_required` (Boolean) Whether users in the target must register the method. Defaults to `false`.



<a id="nestedatt--registration_enforcement"></a>
### Nested Schema for `registration_enforcement`

Required:

- `authentication_methods_registration_campaign` (Attributes) Run campaigns to remind users to set up targeted authentication methods. (see [below for nested schema](#nestedatt--registration_enforcement--authentication_methods_registration_campaign))

<a id="nestedatt--registration_enforcement--authentication_methods_registration_campaign"></a>
### Nested Schema for `registration_enforcement.authentication_methods_registration_campaign`

Optional:

- `enforce_registration_after_allowed_snoozes` (Boolean) Whether users must register after they use up their allowed snoozes. When not configured, the current tenant value is kept.
- `exclude_targets` (Attributes Set) Users and groups excluded from the registration campaign. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--registration_enforcement--authentication_methods_registration_campaign--exclude_targets))
- `include_targets` (Attributes Set) Users and groups the registration campaign targets. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--registration_enforcement--authentication_methods_registration_campaign--include_targets))
- `snooze_duration_in_days` (Number) The number of days a user can postpone registration. Accepts 0 to 14. When not configured, the current tenant value is kept.
- `state` (String) The state of the registration campaign. Possible values are: `default`, `enabled`, `disabled`. When not configured, the current tenant value is kept.

<a id="nestedatt--registration_enforcement--authentication_methods_registration_campaign--exclude_targets"></a>
### Nested Schema for `registration_enforcement.authentication_methods_registration_campaign.exclude_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.


<a id="nestedatt--registration_enforcement--authentication_methods_registration_campaign--include_targets"></a>
### Nested Schema for `registration_enforcement.authentication_methods_registration_campaign.include_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.
- `targeted_authentication_method` (String) The authentication method users are prompted to register, for example `microsoftAuthenticator`.




<a id="nestedatt--report_suspicious_activity_settings"></a>
### Nested Schema for `report_suspicious_activity_settings`

Optional:

- `include_target` (Attributes) The user or group the setting applies to. When not configured, the current tenant value is kept. (see [below for nested schema](#nestedatt--report_suspicious_activity_settings--include_target))
- `state` (String) The state of suspicious activity reporting. Possible values are: `default`, `enabled`, `disabled`. When not configured, the current tenant value is kept.
- `voice_reporting_code` (Number) The keypad code a user enters during a voice call to report suspicious activity. Accepts 0 to 9; `0` uses the default. When not configured, the current tenant value is kept.

<a id="nestedatt--report_suspicious_activity_settings--include_target"></a>
### Nested Schema for `report_suspicious_activity_settings.include_target`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.



<a id="nestedatt--sms"></a>
### Nested Schema for `sms`

Required:

- `state` (String) Whether the method is enabled in the tenant. Possible values are: `enabled`, `disabled`.

Optional:

- `exclude_targets` (Attributes Set) Users and groups excluded from the method. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--sms--exclude_targets))
- `include_targets` (Attributes Set) Users and groups that can use text messages. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--sms--include_targets))

<a id="nestedatt--sms--exclude_targets"></a>
### Nested Schema for `sms.exclude_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.


<a id="nestedatt--sms--include_targets"></a>
### Nested Schema for `sms.include_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `is_usable_for_sign_in` (Boolean) Whether users in the target can use text messages as a first factor (SMS sign-in).
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.

Optional:

- `is_registration_required` (Boolean) Whether users in the target must register the method. Defaults to `false`.



<a id="nestedatt--software_oath"></a>
### Nested Schema for `software_oath`

Required:

- `state` (String) Whether the method is enabled in the tenant. Possible values are: `enabled`, `disabled`.

Optional:

- `exclude_targets` (Attributes Set) Users and groups excluded from the method. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--software_oath--exclude_targets))
- `include_targets` (Attributes Set) Users and groups that can use software OATH tokens. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--software_oath--include_targets))

<a id="nestedatt--software_oath--exclude_targets"></a>
### Nested Schema for `software_oath.exclude_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.


<a id="nestedatt--software_oath--include_targets"></a>
### Nested Schema for `software_oath.include_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.

Optional:

- `is_registration_required` (Boolean) Whether users in the target must register the method. Defaults to `false`.



<a id="nestedatt--system_credential_preferences"></a>
### Nested Schema for `system_credential_preferences`

Optional:

- `exclude_targets` (Attributes Set) Users and groups excluded from system-preferred multifactor authentication. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--system_credential_preferences--exclude_targets))
- `include_targets` (Attributes Set) Users and groups included in system-preferred multifactor authentication. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--system_credential_preferences--include_targets))
- `state` (String) The state of system-preferred multifactor authentication. Possible values are: `default`, `enabled`, `disabled`. When not configured, the current tenant value is kept.

<a id="nestedatt--system_credential_preferences--exclude_targets"></a>
### Nested Schema for `system_credential_preferences.exclude_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.


<a id="nestedatt--system_credential_preferences--include_targets"></a>
### Nested Schema for `system_credential_preferences.include_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.



<a id="nestedatt--temporary_access_pass"></a>
### Nested Schema for `temporary_access_pass`

Required:

- `state` (String) Whether the method is enabled in the tenant. Possible values are: `enabled`, `disabled`.

Optional:

- `default_length` (Number) The default length in characters of a new pass. Accepts 8 to 48. When not configured, the current tenant value is kept.
- `default_lifetime_in_minutes` (Number) The default lifetime in minutes of a new pass. Accepts 10 to 43200. When not configured, the current tenant value is kept.
- `exclude_targets` (Attributes Set) Users and groups excluded from the method. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--temporary_access_pass--exclude_targets))
- `include_targets` (Attributes Set) Users and groups that can use a Temporary Access Pass. When not configured, the current tenant value is kept; set `[]` to clear it. (see [below for nested schema](#nestedatt--temporary_access_pass--include_targets))
- `is_usable_once` (Boolean) Whether passes are limited to a one-time use. When not configured, the current tenant value is kept.
- `maximum_lifetime_in_minutes` (Number) The maximum lifetime in minutes of any pass. Accepts 10 to 43200. When not configured, the current tenant value is kept.
- `minimum_lifetime_in_minutes` (Number) The minimum lifetime in minutes of any pass. Accepts 10 to 43200. When not configured, the current tenant value is kept.

<a id="nestedatt--temporary_access_pass--exclude_targets"></a>
### Nested Schema for `temporary_access_pass.exclude_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.


<a id="nestedatt--temporary_access_pass--include_targets"></a>
### Nested Schema for `temporary_access_pass.include_targets`

Required:

- `id` (String) The object ID of the user or group, or `all_users` for every user in the tenant.
- `target_type` (String) The type of the target. Possible values are: `user`, `group`.

Optional:

- `is_registration_required` (Boolean) Whether users in the target must register the method. Defaults to `false`.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Singleton Resource**: One authentication methods policy exists per tenant. Manage one instance per tenant across all Terraform states. The ID is always `authenticationMethodsPolicy`.
- **Create and Update Behavior**: Create adopts the existing tenant policy. Policy-level settings are sent in one PATCH to `/policies/authenticationMethodsPolicy`. Each method configuration is sent in its own PATCH to `/policies/authenticationMethodsPolicy/authenticationMethodConfigurations/{id}` with its `@odata.type`. Only settings and method blocks that are configured and changed are written.
- **Destroy Behavior**: Destroy only removes the resource from Terraform state. The tenant policy and its method configurations are left unchanged; they are not reset to defaults.
- **Unconfigured Settings**: Every setting and every method block is optional and computed. A setting or block that is not configured keeps its current tenant value, is shown in state after refresh, and is never written. Removing a block from the configuration stops managing it without changing the tenant.
- **Clearing Targets**: An unset `include_targets` or `exclude_targets` is omitted from the request, because Graph rejects null target lists. Set the attribute to `[]` to clear it.
- **Target IDs**: Targets accept a user or group object ID, or `all_users` for every user. Graph does not validate that a target group exists.
- **Policy-Level Blocks**: Inside a configured policy-level block, nested attributes that are not configured are omitted from the request on create. Configure every attribute of a policy-level block that you want to control.
- **Managed Method Configurations**: This version manages `email`, `temporary_access_pass`, `sms`, `software_oath` and `hardware_oath`. The `Fido2`, `MicrosoftAuthenticator`, `Voice`, `X509Certificate`, `QRCodePin`, `Password`, `VerifiableCredentials`, `FederatedIdentityCredential` and external method configurations are not managed and are never changed by this resource.
- **Voice**: The `Voice` method is not managed because its include targets carry a `verificationMode` that the Graph beta SDK model does not expose, so a PATCH could reset it.

## Import

Import is supported using the following syntax:

```shell
# The authentication methods policy is a tenant singleton. Import it with the literal ID authenticationMethodsPolicy.
terraform import microsoft365_graph_beta_identity_and_access_authentication_methods_policy.this authenticationMethodsPolicy
```
