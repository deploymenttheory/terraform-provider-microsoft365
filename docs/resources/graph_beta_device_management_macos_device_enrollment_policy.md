---
page_title: "microsoft365_graph_beta_device_management_macos_device_enrollment_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages a macOS Automated Device Enrollment (ADE) profile using the /deviceManagement/configurationPolicies settings catalog endpoint. This is the modern, settings-catalog-backed equivalent of the legacy depMacOSEnrollmentProfile API (see microsoft365_graph_beta_device_management_macos_dep_enrollment_profile), and controls macOS Setup Assistant behavior for devices enrolled via Apple Business Manager / Apple School Manager.
---

# microsoft365_graph_beta_device_management_macos_device_enrollment_policy (Resource)

Manages a macOS Automated Device Enrollment (ADE) profile using the `/deviceManagement/configurationPolicies` settings catalog endpoint. This is the modern, settings-catalog-backed equivalent of the legacy `depMacOSEnrollmentProfile` API (see `microsoft365_graph_beta_device_management_macos_dep_enrollment_profile`), and controls macOS Setup Assistant behavior for devices enrolled via Apple Business Manager / Apple School Manager.

## Microsoft Documentation

- [deviceManagementConfigurationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta)
- [Set up automated device enrollment for macOS](https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-macos)
- [Set up enrollment time grouping in Microsoft Intune](https://learn.microsoft.com/en-us/mem/intune/enrollment/enrollment-time-grouping)
- [enrollmentProfile: setDefaultProfile action](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-enrollmentprofile-setdefaultprofile?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementConfiguration.ReadWrite.All`
- `DeviceManagementServiceConfig.Read.All`
- `Directory.Read.All`
- `Group.Read.All`
- `GroupMember.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.56.0-alpha | Experimental | Initial release |

## Known Issues

### `device_security_group` Requires Delegated Authentication - 07/03/2026

Enrollment time grouping (setting `device_security_group`) is implemented via the dedicated
`setEnrollmentTimeDeviceMembershipTarget` / `clearEnrollmentTimeDeviceMembershipTarget` actions on
`/deviceManagement/configurationPolicies/{id}`, rather than through the settings catalog.

**Issue Summary:**
- These endpoints return an `Internal Server Error - 500` from the Intune backend (`DeviceConfigV2`)
  when called with application permissions (client credentials) - the auth flow this provider
  always uses.
- The identical request succeeds when made with delegated (signed-in user) permissions, e.g. from
  the Intune admin center.
- All other attributes on this resource are unaffected and work correctly with application
  permissions.

**Impact:**
- `device_security_group` cannot be set or changed through this provider when it is configured
  with a service principal / client credentials, which is the only supported configuration.

**Workaround:**
None currently, other than configuring the group's membership target directly through the Intune
admin center, outside of Terraform.

**Reference:**
- Example configuration: [enrollment_time_grouping.tf](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/examples/resources/microsoft365_graph_beta_device_management_macos_device_enrollment_policy/enrollment_time_grouping.tf)

## Example Usage

### Minimal and Maximal

```terraform
# Example 1: Minimal zero-touch macOS ADE enrollment profile (no user affinity, no local
# account creation). await_device_configured must be false when admin_account is omitted.
resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "minimal" {
  name = "macOS ADE - Minimal"

  requires_user_authentication = false
  await_device_configured      = false

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}

# Example 2: Maximal macOS ADE enrollment profile exercising the full settings tree - user
# authentication via Company Portal, a LAPS-style local admin account with a separate primary
# account, locked enrollment, and every Setup Assistant screen toggle.
resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "maximal" {
  name        = "macOS ADE - Maximal"
  description = "macOS ADE enrollment policy exercising the full settings tree"

  # Uncomment to target a specific Apple ABM/ASM token when the tenant has more than one;
  # otherwise this is auto-resolved.
  # dep_onboarding_settings_id = "00000000-0000-0000-0000-000000000000"

  # Makes this the default macOS enrollment profile for the DEP token via the setDefaultProfile
  # action. Only one policy per DEP token can be the default; setting this to true elsewhere
  # supersedes this assignment. There is no "unassign" action - see the resource documentation.
  is_default_policy_assignment = true

  requires_user_authentication                               = true
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = true

  await_device_configured = true

  admin_account = {
    create_local_admin_account = true
    user_name                  = "localadmin"
    full_name                  = "Local Administrator"
    hide_account               = true
    password_rotation_in_days  = 90

    create_local_primary_account = true
    primary_account = {
      prefill_account_info = true
      restrict_editing     = true
      user_name            = "primaryuser"
      full_name            = "Primary User"
    }
  }

  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  # Setup Assistant screen toggles - true hides the pane, false (default) shows it.
  location_services_disabled                = false
  restore_disabled                          = true
  apple_id_disabled                         = true
  terms_and_conditions_disabled             = false
  touch_id_disabled                         = false
  apple_pay_disabled                        = true
  siri_disabled                             = true
  diagnostics_disabled                      = true
  file_vault_disabled                       = false
  icloud_diagnostics_disabled               = true
  icloud_storage_disabled                   = true
  display_tone_setup_disabled               = false
  screen_time_screen_disabled               = true
  privacy_pane_disabled                     = true
  accessibility_screen_disabled             = false
  auto_unlock_with_watch_disabled           = true
  lockdown_mode_disabled                    = true
  software_update_screen_disabled           = false
  software_update_completed_screen_disabled = true
  terms_of_address_screen_disabled          = true
  apple_intelligence_disabled               = false
  os_showcase_screen_disabled               = true
  app_store_disabled                        = false

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
```

### Enrollment Time Grouping (Device Security Group)

```terraform
# Enrollment time grouping (ETG) adds devices to a static Microsoft Entra security group as they
# enroll, via device_security_group. This requires two prerequisites, both shown below:
#
#  1. A static (Assigned) Microsoft Entra security group.
#  2. The "Intune Provisioning Client" service principal (AppId f1346770-5b25-470b-88bd-d5744ab7952c,
#     sometimes shown as "Intune Autopilot ConfidentialClient") set as an OWNER of that group, so
#     Intune can add enrolling devices to it.
#
# ~> Known Microsoft Graph limitation: setting/clearing device_security_group calls the
# setEnrollmentTimeDeviceMembershipTarget / clearEnrollmentTimeDeviceMembershipTarget actions,
# which currently return a 500 error from the Intune backend when called with application
# permissions (client credentials) - the auth flow this provider always uses. These calls succeed
# with delegated (signed-in user) permissions only. See the resource's "Known Issues" section.

data "microsoft365_graph_beta_applications_service_principal" "intune_provisioning_client" {
  app_id = "f1346770-5b25-470b-88bd-d5744ab7952c"
}

resource "microsoft365_graph_beta_groups_group" "macos_ade_enrollment_group" {
  display_name     = "macos-ade-enrollment-time-grouping"
  mail_nickname    = "macosadeenrollmenttimegrouping"
  mail_enabled     = false
  security_enabled = true

  # Membership type: Assigned (static) - omitting "DynamicMembership" from group_types keeps
  # this a static group. Microsoft Entra roles can be assigned to the group: No.
  is_assignable_to_role = false

  description = "Static security group used for macOS ADE enrollment time grouping."

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Give Microsoft Entra ID time to propagate the new group before granting ownership.
resource "time_sleep" "wait_for_group_propagation" {
  depends_on = [
    microsoft365_graph_beta_groups_group.macos_ade_enrollment_group,
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "macos_ade_enrollment_group_owner" {
  group_id          = microsoft365_graph_beta_groups_group.macos_ade_enrollment_group.id
  owner_id          = data.microsoft365_graph_beta_applications_service_principal.intune_provisioning_client.id
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_group_propagation]
}

# Give Microsoft Entra ID time to propagate the new ownership before the policy references the
# group (Intune validates ownership synchronously against Graph, which can lag after the write).
resource "time_sleep" "wait_for_owner_propagation" {
  depends_on = [
    microsoft365_graph_beta_groups_group_owner_assignment.macos_ade_enrollment_group_owner,
  ]

  create_duration = "60s"
}

resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "with_enrollment_time_grouping" {
  name = "macOS ADE - Enrollment Time Grouping"

  requires_user_authentication = false
  await_device_configured      = false

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  # Only valid when the provider is configured with delegated (user) credentials - see the
  # Known Issues note above.
  device_security_group = microsoft365_graph_beta_groups_group.macos_ade_enrollment_group.id

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_owner_propagation]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The display name of the enrollment profile.
- `requires_user_authentication` (Boolean) Whether the enrollment requires user authentication (user affinity). When `false`, the device enrolls without an associated user (shared/kiosk device path). When `true`, exactly one of `enable_authentication_via_company_portal` or `require_company_portal_on_setup_assistant_enrolled_devices` must also be `true` - Microsoft Graph rejects the enrollment profile otherwise.
- `support_department` (String) The department name shown to the user on the Setup Assistant Remote Management pane. Must be between 1 and 125 characters.
- `support_phone_number` (String) The support phone number shown to the user on the Setup Assistant Remote Management pane. Must be between 1 and 125 characters.

### Optional

- `accessibility_screen_disabled` (Boolean) Whether to hide the Accessibility pane in Setup Assistant.
- `admin_account` (Attributes) Local account settings created during Setup Assistant. Required when `await_device_configured` is `true`, and must be omitted when it is `false`. (see [below for nested schema](#nestedatt--admin_account))
- `app_store_disabled` (Boolean) Whether to hide the App Store pane in Setup Assistant.
- `apple_id_disabled` (Boolean) Whether to hide the Apple ID sign-in pane in Setup Assistant.
- `apple_intelligence_disabled` (Boolean) Whether to hide the Apple Intelligence pane in Setup Assistant.
- `apple_pay_disabled` (Boolean) Whether to hide the Apple Pay pane in Setup Assistant.
- `auto_unlock_with_watch_disabled` (Boolean) Whether to hide the Unlock with Apple Watch pane in Setup Assistant.
- `await_device_configured` (Boolean) Whether `admin_account` configures a local account. When `true`, `admin_account` is required; when `false`, it must be omitted and no local account is created. Confirmed against live Intune admin center traffic: the underlying `ade_macos_awaitconfiguration` setting is always sent as active, with the actual create/don't-create choice carried entirely by `admin_account.create_local_admin_account`.
- `dep_onboarding_settings_id` (String) The ID of the Apple ABM/ASM DEP onboarding token (`/deviceManagement/depOnboardingSettings`) that owns this profile. If omitted, it is automatically resolved to the tenant's single Apple ADE/ABM (or ASM) token; if the tenant has more than one Apple token, this must be set explicitly.
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `device_security_group` (String) The ID of the static Microsoft Entra security group to use for enrollment time grouping (the "Device group" tab in the Intune admin center policy wizard). Devices assigned this policy become members of the group as they enroll. This is set via the dedicated `setEnrollmentTimeDeviceMembershipTarget` action on `/deviceManagement/configurationPolicies/{id}` (and cleared via `clearEnrollmentTimeDeviceMembershipTarget` when removed), not via the settings catalog. The group must have the 'Intune Provisioning Client' service principal (AppId: f1346770-5b25-470b-88bd-d5744ab7952c) set as its owner; in some tenants this service principal may appear as 'Intune Autopilot ConfidentialClient'.

~> **Known Microsoft Graph limitation:** as of this writing, `setEnrollmentTimeDeviceMembershipTarget` and `clearEnrollmentTimeDeviceMembershipTarget` return an `Internal Server Error - 500` from the Intune backend (`DeviceConfigV2`) when called with application permissions (client credentials) - the auth flow this provider always uses. The identical request succeeds when made with delegated (signed-in user) permissions, e.g. from the Intune admin center. Until Microsoft resolves this for application permissions, setting `device_security_group` through this provider will fail on `Create` and `Update`; this is a Microsoft Graph service limitation, not a provider defect.
- `diagnostics_disabled` (Boolean) Whether to hide the Diagnostics pane in Setup Assistant.
- `display_tone_setup_disabled` (Boolean) Whether to hide the Appearance (display tone) pane in Setup Assistant.
- `enable_authentication_via_company_portal` (Boolean) Whether Setup Assistant authenticates the user via Company Portal. Only applicable when `requires_user_authentication` is `true`, in which case exactly one of this or `require_company_portal_on_setup_assistant_enrolled_devices` must be `true`. Mutually exclusive with `require_company_portal_on_setup_assistant_enrolled_devices`.
- `file_vault_disabled` (Boolean) Whether to hide the FileVault pane in Setup Assistant.
- `icloud_diagnostics_disabled` (Boolean) Whether to hide the iCloud Analytics pane in Setup Assistant.
- `icloud_storage_disabled` (Boolean) Whether to hide the iCloud Storage pane in Setup Assistant.
- `is_default_policy_assignment` (Boolean) Whether this policy is the default macOS enrollment profile for its `dep_onboarding_settings_id`, set via the dedicated `setDefaultProfile` action. Always reflects the DEP token's actual current default on refresh, regardless of configuration.

~> **No unassign action:** Microsoft Graph does not expose an `unsetDefaultProfile`/`clearDefaultProfile` action - `setDefaultProfile` is the only operation available. Setting this to `false` on a policy that is currently the DEP token's default has no effect on Graph; the next refresh reports `true` again. Only setting a different policy's `is_default_policy_assignment` to `true` changes which profile is the default. A change from `true` to `false` while this policy is still the token's current default can therefore never converge, and the provider rejects the update with a validation error. Promote the replacement policy first - in the same apply, give this policy a `depends_on` for the replacement so the promotion runs first, or apply the promotion separately.
- `location_services_disabled` (Boolean) Whether to hide the Location Services pane in Setup Assistant.
- `lockdown_mode_disabled` (Boolean) Whether to hide the Lockdown Mode pane in Setup Assistant.
- `locked_enrollment_enabled` (Boolean) Whether enrollment is locked to the authorized user/device, preventing the MDM profile from being removed before enrollment completes.
- `os_showcase_screen_disabled` (Boolean) Whether to hide the What's New (OS showcase) pane in Setup Assistant.
- `privacy_pane_disabled` (Boolean) Whether to hide the Privacy pane in Setup Assistant.
- `require_company_portal_on_setup_assistant_enrolled_devices` (Boolean) Whether Company Portal is required on Setup Assistant enrolled devices. Only applicable when `requires_user_authentication` is `true`, in which case exactly one of this or `enable_authentication_via_company_portal` must be `true`. Mutually exclusive with `enable_authentication_via_company_portal`.
- `restore_disabled` (Boolean) Whether to hide the Restore from Backup pane in Setup Assistant.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `screen_time_screen_disabled` (Boolean) Whether to hide the Screen Time pane in Setup Assistant.
- `siri_disabled` (Boolean) Whether to hide the Siri pane in Setup Assistant.
- `software_update_completed_screen_disabled` (Boolean) Whether to hide the post-installation Software Update Completed pane in Setup Assistant.
- `software_update_screen_disabled` (Boolean) Whether to hide the Software Update pane in Setup Assistant.
- `terms_and_conditions_disabled` (Boolean) Whether to hide the Terms and Conditions pane in Setup Assistant.
- `terms_of_address_screen_disabled` (Boolean) Whether to hide the Terms of Address pane in Setup Assistant.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `touch_id_disabled` (Boolean) Whether to hide the Touch ID/Face ID pane in Setup Assistant.

### Read-Only

- `created_date_time` (String) Creation date and time of the policy.
- `id` (String) The unique identifier for this policy.
- `is_assigned` (Boolean) Indicates if the policy is assigned to any scope.
- `last_modified_date_time` (String) Last modification date and time of the policy.
- `platforms` (String) The platforms this policy applies to. Always `macOS`.
- `settings_count` (Number) Number of settings within the policy.
- `technologies` (String) The technology this policy is using. Always `enrollment`.
- `template_family` (String) The template family for this policy (`enrollmentConfiguration`).
- `template_id` (String) The settings catalog template ID used by this policy.

<a id="nestedatt--admin_account"></a>
### Nested Schema for `admin_account`

Required:

- `create_local_admin_account` (Boolean) Whether Setup Assistant creates a local administrator account.
- `create_local_primary_account` (Boolean) Whether Setup Assistant also creates a separate, standard (non-admin) local primary account.

Optional:

- `full_name` (String) The full name for the local administrator account.
- `hide_account` (Boolean) Whether to hide the local administrator account from the login window and Users & Groups.
- `password_rotation_in_days` (Number) Automatic rotation period, in days, for the local administrator account password. `0` disables rotation.
- `primary_account` (Attributes) Standard local account settings. Only applicable when `create_local_primary_account` is `true`. (see [below for nested schema](#nestedatt--admin_account--primary_account))
- `user_name` (String) The account (short) name for the local administrator account.

<a id="nestedatt--admin_account--primary_account"></a>
### Nested Schema for `admin_account.primary_account`

Optional:

- `full_name` (String) The full name to prefill for the primary account. Only applicable when `prefill_account_info` is `true`.
- `prefill_account_info` (Boolean) Whether to prefill the primary account name/full name in Setup Assistant.
- `restrict_editing` (Boolean) Whether to prevent the user from editing the prefilled primary account information. Only applicable when `prefill_account_info` is `true`.
- `user_name` (String) The account (short) name to prefill for the primary account. Only applicable when `prefill_account_info` is `true`.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Settings catalog, not the legacy DEP API**: this resource manages macOS ADE profiles through
  the modern `/deviceManagement/configurationPolicies` settings catalog endpoint. For the legacy
  `depMacOSEnrollmentProfile` API, see
  `microsoft365_graph_beta_device_management_macos_dep_enrollment_profile`.
- **User authentication**: when `requires_user_authentication` is `true`, exactly one of
  `enable_authentication_via_company_portal` or
  `require_company_portal_on_setup_assistant_enrolled_devices` must also be `true`. Microsoft Graph
  rejects the resulting authentication method otherwise.
- **Local account creation**: `admin_account` is required when `await_device_configured` is
  `true`, and must be omitted when it is `false`. `admin_account.primary_account` is only
  applicable when `admin_account.create_local_primary_account` is `true`.
- **Support department/phone**: `support_department` and `support_phone_number` are required and
  must be between 1 and 125 characters - Microsoft Graph rejects empty values.
- **Enrollment time grouping**: see the Known Issues section above and the dedicated
  `device_security_group` prerequisites example.
- **`is_default_policy_assignment` is a singleton per DEP token**: only one macOS enrollment
  policy can be the default for a given `dep_onboarding_settings_id` at a time. Setting this to
  `true` on another policy for the same token supersedes this one's assignment on its next apply.
- **No unassign action**: Microsoft Graph does not expose an `unsetDefaultProfile` or
  `clearDefaultProfile` action - `setDefaultProfile` is the only operation available for this
  relationship. Setting `is_default_policy_assignment` to `false` on a policy that is currently
  the DEP token's default has no effect on Graph; the next refresh reports `true` again. A `true`
  to `false` change while the policy is still the token's current default can therefore never
  converge, so the provider rejects it with a validation error during update. Promote the
  replacement policy first: in the same apply, give the demoted policy a `depends_on` for the
  replacement so the promotion runs first, or apply the promotion separately - once another
  policy is the default, this attribute refreshes to `false` on its own.
- **Drift detection**: `is_default_policy_assignment` is re-derived from the DEP token's current
  default profile on every refresh. If the default is changed outside Terraform (e.g. in the
  Intune admin center), the next plan shows a diff and `apply` restores the configured assignment.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
terraform import microsoft365_graph_beta_device_management_macos_device_enrollment_policy.example 00000000-0000-0000-0000-000000000000
```
