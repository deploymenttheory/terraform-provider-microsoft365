---
page_title: "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages an iOS/iPadOS Automated Device Enrollment (ADE) profile using the /deviceManagement/configurationPolicies settings catalog endpoint. This is the modern, settings-catalog-backed equivalent of the legacy depIOSEnrollmentProfile API, and controls iOS/iPadOS Setup Assistant behavior for devices enrolled via Apple Business Manager / Apple School Manager.
---

# microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy (Resource)

Manages an iOS/iPadOS Automated Device Enrollment (ADE) profile using the `/deviceManagement/configurationPolicies` settings catalog endpoint. This is the modern, settings-catalog-backed equivalent of the legacy `depIOSEnrollmentProfile` API, and controls iOS/iPadOS Setup Assistant behavior for devices enrolled via Apple Business Manager / Apple School Manager.

## Microsoft Documentation

- [deviceManagementConfigurationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta)
- [Set up automated device enrollment for iOS/iPadOS](https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-ios)
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
- Example configuration: [enrollment_time_grouping.tf](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/examples/resources/microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy/enrollment_time_grouping.tf)

## Example Usage

### Minimal and Maximal

```terraform
# Example 1: Minimal zero-touch iOS/iPadOS ADE enrollment profile (no user affinity).
resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "minimal" {
  name = "iOS ADE - Minimal"

  requires_user_authentication = false

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}

# Example 2: Maximal iOS/iPadOS ADE enrollment profile exercising the full settings tree - user
# authentication in Setup Assistant with modern authentication and await final configuration,
# locked enrollment, device naming, cellular data plan activation, and every Setup Assistant
# screen toggle.
resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "maximal" {
  name        = "iOS ADE - Maximal"
  description = "iOS/iPadOS ADE enrollment policy exercising the full settings tree"

  # Uncomment to target a specific Apple ABM/ASM token when the tenant has more than one;
  # otherwise this is auto-resolved.
  # dep_onboarding_settings_id = "00000000-0000-0000-0000-000000000000"

  # Makes this the default iOS/iPadOS enrollment profile for the DEP token via the
  # setDefaultProfile action. Only one policy per DEP token can be the default; setting this to
  # true elsewhere supersedes this assignment. There is no "unassign" action - see the resource
  # documentation.
  is_default_policy_assignment = true

  requires_user_authentication                        = true
  require_setup_assistant_with_modern_authentication = true
  await_final_configuration                           = true

  locked_enrollment_enabled = true

  device_name_template         = "{{DEVICETYPE}}-{{SERIAL}}"
  cellular_data_activation_url = "http://activation.carrier.net"

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  # Setup Assistant screen toggles - true hides the pane, false (default) shows it.
  passcode_disabled                         = true
  location_services_disabled                = false
  restore_disabled                          = true
  apple_id_disabled                         = true
  terms_and_conditions_disabled             = false
  touch_id_disabled                         = false
  apple_pay_disabled                        = true
  siri_disabled                             = true
  diagnostics_disabled                      = true
  privacy_pane_disabled                     = true
  restore_from_android_disabled             = true
  imessage_and_facetime_disabled            = true
  screen_time_screen_disabled               = true
  sim_setup_screen_disabled                 = true
  software_update_screen_disabled           = false
  watch_migration_screen_disabled           = true
  appearance_screen_disabled                = false
  device_to_device_migration_disabled       = true
  restore_completed_screen_disabled         = true
  software_update_completed_screen_disabled = true
  get_started_screen_disabled               = false
  action_button_screen_disabled             = true
  safety_screen_disabled                    = true
  terms_of_address_screen_disabled          = true
  apple_intelligence_disabled               = false
  lockdown_mode_disabled                    = true
  app_store_disabled                        = false
  camera_button_screen_disabled             = true
  multitasking_screen_disabled              = true
  os_showcase_screen_disabled               = true
  safety_and_handling_screen_disabled       = true
  web_content_filtering_disabled            = true

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

resource "microsoft365_graph_beta_groups_group" "ios_ade_enrollment_group" {
  display_name     = "ios-ade-enrollment-time-grouping"
  mail_nickname    = "iosadeenrollmenttimegrouping"
  mail_enabled     = false
  security_enabled = true

  # Membership type: Assigned (static) - omitting "DynamicMembership" from group_types keeps
  # this a static group. Microsoft Entra roles can be assigned to the group: No.
  is_assignable_to_role = false

  description = "Static security group used for iOS/iPadOS ADE enrollment time grouping."

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
    microsoft365_graph_beta_groups_group.ios_ade_enrollment_group,
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "ios_ade_enrollment_group_owner" {
  group_id          = microsoft365_graph_beta_groups_group.ios_ade_enrollment_group.id
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
    microsoft365_graph_beta_groups_group_owner_assignment.ios_ade_enrollment_group_owner,
  ]

  create_duration = "60s"
}

resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "with_enrollment_time_grouping" {
  name = "iOS ADE - Enrollment Time Grouping"

  requires_user_authentication = false

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  # Only valid when the provider is configured with delegated (user) credentials - see the
  # Known Issues note above.
  device_security_group = microsoft365_graph_beta_groups_group.ios_ade_enrollment_group.id

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
- `requires_user_authentication` (Boolean) Whether the enrollment requires user authentication (user affinity). When `false`, the device enrolls without an associated user (shared/kiosk device path). When `true`, the authentication flow is selected via `enable_authentication_via_company_portal` or `require_setup_assistant_with_modern_authentication`; when neither is set, the legacy Setup Assistant authentication flow is used.
- `support_department` (String) The department name shown to the user on the Setup Assistant Remote Management pane. Must be between 1 and 125 characters.
- `support_phone_number` (String) The support phone number shown to the user on the Setup Assistant Remote Management pane. Must be between 1 and 50 characters.

### Optional

- `action_button_screen_disabled` (Boolean) Whether to hide the configuration pane for the action button in Setup Assistant. For iOS/iPadOS 17.0 and later.
- `app_store_disabled` (Boolean) Whether to hide the Apple App Store pane in Setup Assistant. For iOS/iPadOS 14.3 and later.
- `appearance_screen_disabled` (Boolean) Whether to hide the appearance setup pane in Setup Assistant. For iOS/iPadOS 13.0 and later.
- `apple_id_disabled` (Boolean) Whether to hide the Apple ID setup pane in Setup Assistant, which gives users the option to sign in with their Apple ID and use iCloud. For iOS/iPadOS 7.0 and later.
- `apple_intelligence_disabled` (Boolean) Whether to hide the Apple Intelligence setup pane in Setup Assistant, where users can configure Apple Intelligence features. For iOS/iPadOS 18.0 and later.
- `apple_pay_disabled` (Boolean) Whether to hide the Apple Pay setup pane in Setup Assistant, which gives users the option to set up Apple Pay on their devices. For iOS/iPadOS 7.0 and later.
- `await_final_configuration` (Boolean) Whether devices are locked in Setup Assistant until all enrollment-time configuration is installed (await final configuration). Only applicable when `require_setup_assistant_with_modern_authentication` is `true`.
- `camera_button_screen_disabled` (Boolean) Whether to hide the camera button pane in Setup Assistant. For iOS/iPadOS 18.0 and later.
- `cellular_data_activation_url` (String) The carrier activation server URL used to activate cellular data plans on eligible devices at enrollment. When omitted, cellular data plan activation is not configured by this profile.
- `dep_onboarding_settings_id` (String) The ID of the Apple ABM/ASM DEP onboarding token (`/deviceManagement/depOnboardingSettings`) that owns this profile. If omitted, it is automatically resolved to the tenant's single Apple ADE/ABM (or ASM) token; if the tenant has more than one Apple token, this must be set explicitly.
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `device_name_template` (String) Device name template applied to supervised devices at enrollment, e.g. `{{DEVICETYPE}}-{{SERIAL}}`. Supports the `{{DEVICETYPE}}` and `{{SERIAL}}` substitution tokens. When omitted, device naming is not managed by this profile.
- `device_security_group` (String) The ID of the static Microsoft Entra security group to use for enrollment time grouping (the "Device group" tab in the Intune admin center policy wizard). Devices assigned this policy become members of the group as they enroll. This is set via the dedicated `setEnrollmentTimeDeviceMembershipTarget` action on `/deviceManagement/configurationPolicies/{id}` (and cleared via `clearEnrollmentTimeDeviceMembershipTarget` when removed), not via the settings catalog. The group must have the 'Intune Provisioning Client' service principal (AppId: f1346770-5b25-470b-88bd-d5744ab7952c) set as its owner; in some tenants this service principal may appear as 'Intune Autopilot ConfidentialClient'.

~> **Known Microsoft Graph limitation:** as of this writing, `setEnrollmentTimeDeviceMembershipTarget` and `clearEnrollmentTimeDeviceMembershipTarget` return an `Internal Server Error - 500` from the Intune backend (`DeviceConfigV2`) when called with application permissions (client credentials) - the auth flow this provider always uses. The identical request succeeds when made with delegated (signed-in user) permissions, e.g. from the Intune admin center. Until Microsoft resolves this for application permissions, setting `device_security_group` through this provider will fail on `Create` and `Update`; this is a Microsoft Graph service limitation, not a provider defect.
- `device_to_device_migration_disabled` (Boolean) Whether to hide the device-to-device migration pane in Setup Assistant. When shown, users can transfer data from an old device to their current device. The option to transfer data directly from a device isn't available for devices running iOS 13 or later.
- `diagnostics_disabled` (Boolean) Whether to hide the diagnostics pane in Setup Assistant, where users can opt in to send diagnostic data to Apple. For iOS/iPadOS 7.0 and later.
- `enable_authentication_via_company_portal` (Boolean) Whether the user authenticates via the Company Portal app instead of Setup Assistant. Only applicable when `requires_user_authentication` is `true`. Mutually exclusive with `require_setup_assistant_with_modern_authentication`.
- `get_started_screen_disabled` (Boolean) Whether to hide the Get Started pane in Setup Assistant.
- `imessage_and_facetime_disabled` (Boolean) Whether to hide the iMessage and FaceTime setup pane in Setup Assistant. For iOS/iPadOS 9.0 and later.
- `is_default_policy_assignment` (Boolean) Whether this policy is the default iOS/iPadOS enrollment profile for its `dep_onboarding_settings_id`, set via the dedicated `setDefaultProfile` action. Always reflects the DEP token's actual current default on refresh, regardless of configuration.

~> **No unassign action:** Microsoft Graph does not expose an `unsetDefaultProfile`/`clearDefaultProfile` action - `setDefaultProfile` is the only operation available. Setting this to `false` on a policy that is currently the DEP token's default has no effect on Graph; the next refresh reports `true` again. Only setting a different policy's `is_default_policy_assignment` to `true` changes which profile is the default. A change from `true` to `false` while this policy is still the token's current default can therefore never converge, and the provider rejects the update with a validation error. Promote the replacement policy first - in the same apply, give this policy a `depends_on` for the replacement so the promotion runs first, or apply the promotion separately.
- `location_services_disabled` (Boolean) Whether to hide the Location Services setup pane in Setup Assistant, where users can enable location services on their device. For iOS/iPadOS 7.0 and later.
- `lockdown_mode_disabled` (Boolean) Whether to hide the Lockdown Mode pane in Setup Assistant.
- `locked_enrollment_enabled` (Boolean) Whether enrollment is locked to the authorized user/device, preventing the MDM profile from being removed before enrollment completes.
- `multitasking_screen_disabled` (Boolean) Whether to hide the multitasking pane in Setup Assistant. For iOS/iPadOS 26.0 and later.
- `os_showcase_screen_disabled` (Boolean) Whether to hide the OS showcase pane in Setup Assistant. For iOS/iPadOS 26.0 and later.
- `passcode_disabled` (Boolean) Whether to hide the passcode and password lock pane in Setup Assistant. When shown, users are prompted for a passcode. Always require a passcode for unsecured devices unless access is controlled in some other way (such as through a kiosk mode configuration that restricts the device to one app). For iOS/iPadOS 7.0 and later.
- `privacy_pane_disabled` (Boolean) Whether to hide the privacy setup pane in Setup Assistant. For iOS/iPadOS 11.3 and later.
- `require_setup_assistant_with_modern_authentication` (Boolean) Whether the user authenticates in Setup Assistant using modern authentication (Microsoft Entra ID). Only applicable when `requires_user_authentication` is `true`. Mutually exclusive with `enable_authentication_via_company_portal`.
- `restore_completed_screen_disabled` (Boolean) Whether to hide the Restore Completed screen shown after a backup and restore is performed during Setup Assistant.
- `restore_disabled` (Boolean) Whether to hide the apps and data (Restore) setup pane in Setup Assistant. When shown, users setting up devices can restore or transfer data from iCloud Backup. For iOS/iPadOS 7.0 and later.
- `restore_from_android_disabled` (Boolean) Whether to hide the Android Migration setup pane in Setup Assistant, meant for previous Android users. When shown, users can migrate data from an Android device. For iOS/iPadOS 9.0 and later.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `safety_and_handling_screen_disabled` (Boolean) Whether to hide the safety and handling pane in Setup Assistant. For iOS/iPadOS 18.4 and later.
- `safety_screen_disabled` (Boolean) Whether to hide the safety (Emergency SOS) setup pane in Setup Assistant. For iOS/iPadOS 16.0 and later.
- `screen_time_screen_disabled` (Boolean) Whether to hide the Screen Time pane in Setup Assistant. For iOS/iPadOS 12.0 and later.
- `sim_setup_screen_disabled` (Boolean) Whether to hide the cellular (SIM Setup) pane in Setup Assistant, where users can add a cellular plan. For iOS/iPadOS 12.0 and later.
- `siri_disabled` (Boolean) Whether to hide the Siri setup pane in Setup Assistant. For iOS/iPadOS 7.0 and later.
- `software_update_completed_screen_disabled` (Boolean) Whether to hide the screen showing all software updates that happen during Setup Assistant.
- `software_update_screen_disabled` (Boolean) Whether to hide the mandatory software update screen in Setup Assistant. For iOS/iPadOS 12.0 and later.
- `terms_and_conditions_disabled` (Boolean) Whether to hide the Apple terms and conditions pane in Setup Assistant. When shown, users are required to accept them. For iOS/iPadOS 7.0 and later.
- `terms_of_address_screen_disabled` (Boolean) Whether to hide the terms of address pane in Setup Assistant, which gives users the option to choose how they want to be addressed throughout the system: feminine, masculine, or neutral. This Apple feature is available for select languages. For iOS/iPadOS 16.0 and later.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `touch_id_disabled` (Boolean) Whether to hide the biometric (Touch ID and Face ID) setup pane in Setup Assistant, which gives users the option to set up fingerprint or facial identification on their devices. For iOS/iPadOS 8.1 and later.
- `watch_migration_screen_disabled` (Boolean) Whether to hide the Apple Watch migration pane in Setup Assistant, where users can migrate data from an Apple Watch. For iOS/iPadOS 11.0 and later.
- `web_content_filtering_disabled` (Boolean) Whether to hide the web content filtering pane in Setup Assistant. For iOS/iPadOS 18.2 and later.

### Read-Only

- `created_date_time` (String) Creation date and time of the policy.
- `id` (String) The unique identifier for this policy.
- `is_assigned` (Boolean) Indicates if the policy is assigned to any scope.
- `last_modified_date_time` (String) Last modification date and time of the policy.
- `platforms` (String) The platforms this policy applies to. Always `iOS`.
- `settings_count` (Number) Number of settings within the policy.
- `technologies` (String) The technology this policy is using. Always `enrollment`.
- `template_family` (String) The template family for this policy (`enrollmentConfiguration`).
- `template_id` (String) The settings catalog template ID used by this policy.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Settings catalog, not the legacy DEP API**: this resource manages iOS/iPadOS ADE profiles
  through the modern `/deviceManagement/configurationPolicies` settings catalog endpoint, not the
  legacy `depIOSEnrollmentProfile` API.
- **User authentication**: when `requires_user_authentication` is `true`, the authentication flow
  is selected with `enable_authentication_via_company_portal` or
  `require_setup_assistant_with_modern_authentication` (mutually exclusive); when neither is set,
  the legacy Setup Assistant authentication flow is used. `await_final_configuration` is only
  applicable with `require_setup_assistant_with_modern_authentication`.
- **Support department/phone**: `support_department` (1-125 characters) and
  `support_phone_number` (1-50 characters) are required - Microsoft Graph rejects empty values.
- **Device naming and cellular activation**: `device_name_template` and
  `cellular_data_activation_url` are optional; when omitted, the corresponding settings catalog
  choices are sent as disabled and the features are not managed by the profile.
- **Enrollment time grouping**: see the Known Issues section above and the dedicated
  `device_security_group` prerequisites example.
- **`is_default_policy_assignment` is a singleton per DEP token**: only one iOS/iPadOS enrollment
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
terraform import microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy.example 00000000-0000-0000-0000-000000000000
```
