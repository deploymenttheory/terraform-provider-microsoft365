---
page_title: "microsoft365_graph_beta_identity_and_access_conditional_access_policy Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft 365 Conditional Access Policies using the /identity/conditionalAccess/policies endpoint. Conditional Access policies define the conditions under which users can access cloud apps.
---

# microsoft365_graph_beta_identity_and_access_conditional_access_policy (Resource)

Manages Microsoft 365 Conditional Access Policies using the `/identity/conditionalAccess/policies` endpoint. Conditional Access policies define the conditions under which users can access cloud apps.

## Microsoft Documentation

- [conditionalAccessPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesspolicy?view=graph-rest-beta)
- [Create conditionalAccessPolicy](https://learn.microsoft.com/en-us/graph/api/conditionalaccessroot-post-policies?view=graph-rest-beta)
- [Update conditionalAccessPolicy](https://learn.microsoft.com/en-us/graph/api/conditionalaccesspolicy-update?view=graph-rest-beta)
- [Delete conditionalAccessPolicy](https://learn.microsoft.com/en-us/graph/api/conditionalaccesspolicy-delete?view=graph-rest-beta)
- [Conditional Access documentation](https://learn.microsoft.com/en-us/azure/active-directory/conditional-access/)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `Policy.ReadWrite.ConditionalAccess`, `Policy.Read.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.19.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "block_legacy_authentication" {
  display_name = "Block Legacy Authentication"
  state        = "enabled"

  conditions = {
    applications = {
      include_applications = ["All"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = ["11111111-1111-1111-1111-111111111111"] # Emergency access group
      exclude_roles = [
        "62e90394-69f5-4237-9190-012177145e10" # Global Administrator
      ]
      exclude_guests_or_external_users = null
    }

    platforms = {
      include_platforms = []
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = [
        "11111111-1111-1111-1111-111111111111" # Trusted office locations
      ]
    }

    client_app_types = ["exchangeActiveSync", "other"]

    devices = {
      device_filter   = null
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = []
    sign_in_risk_levels = []

    authentication_flows = null
  }

  grant_controls = {
    operator          = "OR"
    built_in_controls = ["block"]
  }

  session_controls = null

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "require_mfa_for_admins" {
  display_name = "Require MFA for Admin Roles"
  state        = "enabled"
  conditions = {
    applications = {
      include_applications = ["All"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = ["11111111-1111-1111-1111-111111111111"] # Emergency access group
      include_roles = [
        "62e90394-69f5-4237-9190-012177145e10", # Global Administrator
        "194ae4cb-b126-40b2-bd5b-6091b380977d", # Security Administrator
        "729827e3-9c14-49f7-bb1b-9608f156bbb8"  # Helpdesk Administrator
      ]
      exclude_roles                    = []
      exclude_guests_or_external_users = null
    }

    platforms = {
      include_platforms = []
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = []
    }

    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    devices = {
      device_filter   = null
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = []
    sign_in_risk_levels = []

    authentication_flows = null
  }

  grant_controls = {
    operator          = "AND"
    built_in_controls = ["mfa"]
  }

  session_controls = {
    sign_in_frequency = {
      is_enabled          = true
      type                = "hours"
      value               = 4
      frequency_interval  = "timeBased"
      authentication_type = "primaryAndSecondaryAuthentication"
    }
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "compliant_device_policy" {
  display_name = "Require Compliant or Hybrid Joined Device"
  state        = "enabled"
  conditions = {
    applications = {
      include_applications = ["Office365"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = ["11111111-1111-1111-1111-111111111111"] # Emergency access group
      exclude_roles = [
        "62e90394-69f5-4237-9190-012177145e10" # Global Administrator
      ]
      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["b2bCollaborationGuest", "b2bCollaborationMember"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    platforms = {
      include_platforms = ["windows", "macOS", "iOS", "android"]
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = [
        "11111111-1111-1111-1111-111111111111" # Trusted office locations
      ]
    }

    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True or device.trustType -eq \"Hybrid Azure AD joined\""
      }
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = []
    sign_in_risk_levels = []

    authentication_flows = null
  }

  grant_controls = {
    operator          = "OR"
    built_in_controls = ["compliantDevice", "domainJoinedDevice"]
  }

  session_controls = {
    cloud_app_security = {
      is_enabled              = true
      cloud_app_security_type = "monitorOnly"
    }

    continuous_access_evaluation = {
      mode = "strictLocation"
    }
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "high_risk_sign_in_policy" {
  display_name = "High Risk Sign-in Policy"
  state        = "enabled"
  conditions = {
    applications = {
      include_applications = ["All"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users                    = ["All"]
      exclude_users                    = []
      include_groups                   = []
      exclude_groups                   = ["11111111-1111-1111-1111-111111111111"] # Emergency access group
      exclude_roles                    = []
      exclude_guests_or_external_users = null
    }

    platforms = {
      include_platforms = []
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = []
    }

    client_app_types = ["browser", "mobileAppsAndDesktopClients", "exchangeActiveSync", "other"]

    devices = {
      device_filter   = null
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = []
    sign_in_risk_levels = ["high"]

    authentication_flows = null
  }

  grant_controls = {
    operator          = "AND"
    built_in_controls = ["mfa", "passwordChange"]
  }

  session_controls = {
    sign_in_frequency = {
      is_enabled          = true
      type                = "hours"
      value               = 1
      frequency_interval  = "timeBased"
      authentication_type = "primaryAndSecondaryAuthentication"
    }

    persistent_browser = {
      is_enabled = true
      mode       = "never"
    }
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `conditions` (Attributes) Conditions that must be met for the policy to apply. (see [below for nested schema](#nestedatt--conditions))
- `display_name` (String) The display name for the Conditional Access policy.
- `grant_controls` (Attributes) Controls for granting access. (see [below for nested schema](#nestedatt--grant_controls))
- `state` (String) Specifies the state of the policy. Possible values are: enabled, disabled, enabledForReportingButNotEnforced.

### Optional

- `partial_enablement_strategy` (String) Strategy for partial enablement of the policy.
- `session_controls` (Attributes) Controls for managing user sessions. (see [below for nested schema](#nestedatt--session_controls))
- `template_id` (String) ID of the template this policy is derived from.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The creation date and time of the policy.
- `deleted_date_time` (String) The deletion date and time of the policy, if applicable.
- `id` (String) String (identifier)
- `modified_date_time` (String) The last modified date and time of the policy.

<a id="nestedatt--conditions"></a>
### Nested Schema for `conditions`

Required:

- `applications` (Attributes) Applications and user actions included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--applications))
- `client_app_types` (Set of String) Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, other.
- `sign_in_risk_levels` (Set of String) Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.
- `users` (Attributes) Users, groups, and roles included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--users))

Optional:

- `client_applications` (Attributes) Client applications configuration for the conditional access policy. (see [below for nested schema](#nestedatt--conditions--client_applications))
- `device_states` (Attributes) Device states included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--device_states))
- `devices` (Attributes) Devices included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--devices))
- `locations` (Attributes) Locations included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--locations))
- `platforms` (Attributes) Platforms included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--platforms))
- `service_principal_risk_levels` (Set of String) Service principal risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.
- `times` (Attributes) Times and days when the policy applies. (see [below for nested schema](#nestedatt--conditions--times))
- `user_risk_levels` (Set of String) User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.

<a id="nestedatt--conditions--applications"></a>
### Nested Schema for `conditions.applications`

Required:

- `include_applications` (Set of String) Applications to include in the policy. Can use the special value 'All' to include all applications.

Optional:

- `application_filter` (Attributes) Filter that defines the applications the policy applies to. (see [below for nested schema](#nestedatt--conditions--applications--application_filter))
- `exclude_applications` (Set of String) Applications to exclude from the policy.
- `include_authentication_context_class_references` (Set of String) Authentication context class references to include in the policy.
- `include_user_actions` (Set of String) User actions to include in the policy.

<a id="nestedatt--conditions--applications--application_filter"></a>
### Nested Schema for `conditions.applications.application_filter`

Required:

- `mode` (String) Mode of the filter. Possible values are: include, exclude.
- `rule` (String) Rule syntax for the filter.



<a id="nestedatt--conditions--users"></a>
### Nested Schema for `conditions.users`

Optional:

- `exclude_groups` (Set of String) Groups to exclude from the policy.
- `exclude_guests_or_external_users` (Attributes) Configuration for excluding guests or external users. (see [below for nested schema](#nestedatt--conditions--users--exclude_guests_or_external_users))
- `exclude_roles` (Set of String) Roles to exclude from the policy.
- `exclude_users` (Set of String) Users to exclude from the policy. Can use special values like 'GuestsOrExternalUsers'.
- `include_groups` (Set of String) Groups to include in the policy.
- `include_guests_or_external_users` (Attributes) Configuration for including guests or external users. (see [below for nested schema](#nestedatt--conditions--users--include_guests_or_external_users))
- `include_roles` (Set of String) Roles to include in the policy.
- `include_users` (Set of String) Users to include in the policy. Can use special values like 'All', 'None', or 'GuestsOrExternalUsers'.

<a id="nestedatt--conditions--users--exclude_guests_or_external_users"></a>
### Nested Schema for `conditions.users.exclude_guests_or_external_users`

Required:

- `external_tenants` (Attributes) Configuration for external tenants. (see [below for nested schema](#nestedatt--conditions--users--exclude_guests_or_external_users--external_tenants))
- `guest_or_external_user_types` (String) Types of guests or external users to exclude. Possible values are: internalGuest, b2bCollaborationGuest, b2bCollaborationMember, b2bDirectConnectUser, otherExternalUser, serviceProvider.

<a id="nestedatt--conditions--users--exclude_guests_or_external_users--external_tenants"></a>
### Nested Schema for `conditions.users.exclude_guests_or_external_users.external_tenants`

Required:

- `members` (Set of String) The list of tenant IDs for external tenants.
- `membership_kind` (String) Kind of membership. Possible values are: all, enumerated, unknownFutureValue.



<a id="nestedatt--conditions--users--include_guests_or_external_users"></a>
### Nested Schema for `conditions.users.include_guests_or_external_users`

Required:

- `external_tenants` (Attributes) Configuration for external tenants. (see [below for nested schema](#nestedatt--conditions--users--include_guests_or_external_users--external_tenants))
- `guest_or_external_user_types` (String) Types of guests or external users to include. Possible values are: internalGuest, b2bCollaborationGuest, b2bCollaborationMember, b2bDirectConnectUser, otherExternalUser, serviceProvider.

<a id="nestedatt--conditions--users--include_guests_or_external_users--external_tenants"></a>
### Nested Schema for `conditions.users.include_guests_or_external_users.external_tenants`

Required:

- `members` (Set of String) The list of tenant IDs for external tenants.
- `membership_kind` (String) Kind of membership. Possible values are: all, enumerated, unknownFutureValue.




<a id="nestedatt--conditions--client_applications"></a>
### Nested Schema for `conditions.client_applications`

Required:

- `include_service_principals` (Set of String) Service principals to include in the policy. Can use the special value 'ServicePrincipalsInMyTenant' to include all service principals.

Optional:

- `exclude_service_principals` (Set of String) Service principals to exclude from the policy.


<a id="nestedatt--conditions--device_states"></a>
### Nested Schema for `conditions.device_states`

Optional:

- `exclude_states` (Set of String) Device states to exclude from the policy.
- `include_states` (Set of String) Device states to include in the policy.


<a id="nestedatt--conditions--devices"></a>
### Nested Schema for `conditions.devices`

Optional:

- `device_filter` (Attributes) Filter that defines the devices the policy applies to. (see [below for nested schema](#nestedatt--conditions--devices--device_filter))
- `exclude_device_states` (Set of String) Device states to exclude from the policy.
- `exclude_devices` (Set of String) Devices to exclude from the policy.
- `include_device_states` (Set of String) Device states to include in the policy.
- `include_devices` (Set of String) Devices to include in the policy.

<a id="nestedatt--conditions--devices--device_filter"></a>
### Nested Schema for `conditions.devices.device_filter`

Required:

- `mode` (String) Mode of the filter. Possible values are: include, exclude.
- `rule` (String) Rule syntax for the filter.



<a id="nestedatt--conditions--locations"></a>
### Nested Schema for `conditions.locations`

Required:

- `include_locations` (Set of String) Locations to include in the policy. Can use special values like 'All' or 'AllTrusted'.

Optional:

- `exclude_locations` (Set of String) Locations to exclude from the policy. Can use special values like 'AllTrusted'.


<a id="nestedatt--conditions--platforms"></a>
### Nested Schema for `conditions.platforms`

Required:

- `include_platforms` (Set of String) Platforms to include in the policy.

Optional:

- `exclude_platforms` (Set of String) Platforms to exclude from the policy.


<a id="nestedatt--conditions--times"></a>
### Nested Schema for `conditions.times`

Optional:

- `all_day` (Boolean) Whether the policy applies all day.
- `end_time` (String) End time for the policy.
- `excluded_ranges` (Set of String) Time ranges when the policy does not apply.
- `included_ranges` (Set of String) Time ranges when the policy applies.
- `start_time` (String) Start time for the policy.
- `time_zone` (String) Time zone for the policy times.



<a id="nestedatt--grant_controls"></a>
### Nested Schema for `grant_controls`

Required:

- `operator` (String) Operator to apply to the controls. Possible values are: AND, OR.

Optional:

- `authentication_strength` (Attributes) Authentication strength required for granting access. (see [below for nested schema](#nestedatt--grant_controls--authentication_strength))
- `built_in_controls` (Set of String) List of built-in controls required by the policy. Possible values are: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.
- `custom_authentication_factors` (Set of String) Custom authentication factors for granting access.
- `terms_of_use` (Set of String) Terms of use required for granting access.

<a id="nestedatt--grant_controls--authentication_strength"></a>
### Nested Schema for `grant_controls.authentication_strength`

Required:

- `id` (String) ID of the authentication strength policy.

Optional:

- `allowed_combinations` (Set of String) The allowed authentication method combinations that satisfy the authentication strength policy.
- `description` (String) Description of the authentication strength policy.
- `display_name` (String) Display name of the authentication strength policy.
- `policy_type` (String) Type of the policy. Possible values are: builtIn, custom.
- `requirements_satisfied` (String) Requirements satisfied by the policy.

Read-Only:

- `created_date_time` (String) Creation date and time of the authentication strength policy.
- `modified_date_time` (String) Last modified date and time of the authentication strength policy.



<a id="nestedatt--session_controls"></a>
### Nested Schema for `session_controls`

Optional:

- `application_enforced_restrictions` (Attributes) Application enforced restrictions for the session. (see [below for nested schema](#nestedatt--session_controls--application_enforced_restrictions))
- `cloud_app_security` (Attributes) Cloud app security controls for the session. (see [below for nested schema](#nestedatt--session_controls--cloud_app_security))
- `continuous_access_evaluation` (Attributes) Continuous access evaluation controls for the session. (see [below for nested schema](#nestedatt--session_controls--continuous_access_evaluation))
- `disable_resilience_defaults` (Boolean) Whether to disable resilience defaults.
- `persistent_browser` (Attributes) Persistent browser controls for the session. (see [below for nested schema](#nestedatt--session_controls--persistent_browser))
- `secure_sign_in_session` (Attributes) Secure sign-in session controls. (see [below for nested schema](#nestedatt--session_controls--secure_sign_in_session))
- `sign_in_frequency` (Attributes) Sign-in frequency controls for the session. (see [below for nested schema](#nestedatt--session_controls--sign_in_frequency))

<a id="nestedatt--session_controls--application_enforced_restrictions"></a>
### Nested Schema for `session_controls.application_enforced_restrictions`

Required:

- `is_enabled` (Boolean) Whether application enforced restrictions are enabled.


<a id="nestedatt--session_controls--cloud_app_security"></a>
### Nested Schema for `session_controls.cloud_app_security`

Required:

- `cloud_app_security_type` (String) Type of cloud app security control. Possible values are: blockDownloads, mcasConfigured, monitorOnly, unknownFutureValue.
- `is_enabled` (Boolean) Whether cloud app security controls are enabled.


<a id="nestedatt--session_controls--continuous_access_evaluation"></a>
### Nested Schema for `session_controls.continuous_access_evaluation`

Required:

- `mode` (String) Mode for continuous access evaluation. Possible values are: disabled, basic, strict.


<a id="nestedatt--session_controls--persistent_browser"></a>
### Nested Schema for `session_controls.persistent_browser`

Required:

- `is_enabled` (Boolean) Whether persistent browser controls are enabled.
- `mode` (String) Mode for persistent browser. Possible values are: always, never.


<a id="nestedatt--session_controls--secure_sign_in_session"></a>
### Nested Schema for `session_controls.secure_sign_in_session`

Required:

- `is_enabled` (Boolean) Whether secure sign-in session controls are enabled.


<a id="nestedatt--session_controls--sign_in_frequency"></a>
### Nested Schema for `session_controls.sign_in_frequency`

Required:

- `is_enabled` (Boolean) Whether sign-in frequency controls are enabled.
- `type` (String) Type of sign-in frequency control. Possible values are: days, hours.
- `value` (Number) Value for the sign-in frequency.

Optional:

- `authentication_type` (String) Authentication type for sign-in frequency. Possible values are: primaryAndSecondaryAuthentication, secondaryAuthentication.
- `frequency_interval` (String) Frequency interval for sign-in frequency. Possible values are: timeBased, everyTime.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Policy States
- **enabled**: The policy is active and will be enforced
- **disabled**: The policy exists but is not enforced
- **enabledForReportingButNotEnforced**: The policy will be evaluated and logged but not enforced (report-only mode)

### Applications
- Use `"All"` to target all cloud applications
- Use `"Office365"` to target all Office 365 applications
- Use specific application IDs for targeted policies
- Application filters support complex OData expressions for fine-grained control

### Users and Groups
- Use `"All"` to target all users
- Use `"GuestsOrExternalUsers"` to target external users
- Specify user, group, or role object IDs for targeted policies
- Emergency access accounts should always be excluded from blocking policies

### Locations
- Named locations must be created in Azure AD before referencing
- Use `"All"` for all locations or `"AllTrusted"` for all trusted locations
- IP-based and country-based locations are supported

### Client App Types
- `browser`: Web browsers
- `mobileAppsAndDesktopClients`: Mobile apps and desktop clients
- `exchangeActiveSync`: Exchange ActiveSync clients
- `other`: Other clients including legacy authentication

### Grant Controls
- **Operator**: `AND` requires all controls, `OR` requires any control
- **Built-in Controls**: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`
- **Authentication Strength**: Reference to custom authentication strength policies

### Session Controls
- **Application Restrictions**: Control access to specific applications
- **Cloud App Security**: Integration with Microsoft Defender for Cloud Apps
- **Sign-in Frequency**: Control how often users must re-authenticate
- **Persistent Browser**: Control browser session persistence
- **Continuous Access Evaluation**: Real-time policy evaluation

### Device Filters
- Support complex OData expressions for device-based conditions
- Common filters include device compliance, trust type, and device attributes
- Use `include` mode to target devices matching the filter
- Use `exclude` mode to exclude devices matching the filter

### Risk-based Policies
- **User Risk Levels**: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`
- **Sign-in Risk Levels**: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`
- Requires Azure AD Identity Protection licenses

### Best Practices
- Always exclude emergency access accounts from blocking policies
- Test policies in report-only mode before enabling enforcement
- Use specific targeting rather than broad "All" assignments when possible
- Monitor policy impact through Azure AD sign-in logs
- Implement a phased rollout for new policies
- Document policy purpose and expected behavior

### Common Policy Scenarios
- **Block Legacy Authentication**: Target legacy client app types with block control
- **Require MFA for Admins**: Target administrative roles with MFA requirement
- **Device Compliance**: Require compliant or domain-joined devices for access
- **Location-based Access**: Block or require additional controls based on location
- **Risk-based Access**: Respond to user or sign-in risk with appropriate controls

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_identity_and_access_conditional_access_policy.example conditional-access-policy-id
``` 