---
page_title: "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages a visionOS Automated Device Enrollment (ADE) profile using the /deviceManagement/configurationPolicies settings catalog endpoint. This controls visionOS Setup Assistant behavior for Apple Vision Pro devices enrolled via Apple Business Manager / Apple School Manager.
---

# microsoft365_graph_beta_device_management_visionos_device_enrollment_policy (Resource)

Manages a visionOS Automated Device Enrollment (ADE) profile using the `/deviceManagement/configurationPolicies` settings catalog endpoint. This controls visionOS Setup Assistant behavior for Apple Vision Pro devices enrolled via Apple Business Manager / Apple School Manager.

## Microsoft Documentation

- [deviceManagementConfigurationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta)
- [Set up automated device enrollment for Apple devices](https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-ios)
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

## Example Usage

```terraform
# Example 1: Minimal zero-touch visionOS ADE enrollment profile (no user affinity).
resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "minimal" {
  name = "visionOS ADE - Minimal"

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

# Example 2: Maximal visionOS ADE enrollment profile exercising the full settings tree - user
# affinity, await configuration, locked enrollment, and every Setup Assistant screen toggle.
resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "maximal" {
  name        = "visionOS ADE - Maximal"
  description = "visionOS ADE enrollment policy exercising the full settings tree"

  # Uncomment to target a specific Apple ABM/ASM token when the tenant has more than one;
  # otherwise this is auto-resolved.
  # dep_onboarding_settings_id = "00000000-0000-0000-0000-000000000000"

  # Makes this the default visionOS enrollment profile for the DEP token via the
  # setDefaultProfile action. Only one policy per DEP token can be the default; setting this to
  # true elsewhere supersedes this assignment. There is no "unassign" action - see the resource
  # documentation.
  is_default_policy_assignment = true

  requires_user_authentication = true
  await_device_configured      = true
  locked_enrollment_enabled    = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  # Setup Assistant screen toggles - true hides the pane, false (default) shows it.
  apple_id_disabled               = true
  apple_pay_disabled              = true
  diagnostics_disabled            = true
  get_started_screen_disabled     = false
  apple_intelligence_disabled     = false
  location_services_disabled      = false
  passcode_disabled               = true
  privacy_pane_disabled           = true
  screen_time_screen_disabled     = true
  siri_disabled                   = true
  software_update_screen_disabled = false
  terms_and_conditions_disabled   = false
  tips_screen_disabled            = true
  touch_id_disabled               = false

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The display name of the enrollment profile.
- `requires_user_authentication` (Boolean) Whether the enrollment requires user authentication (user affinity). When `false`, the device enrolls without an associated user (shared/kiosk device path). visionOS uses the basic user affinity setting - unlike iOS/iPadOS, there is no authentication method choice.
- `support_department` (String) The department name shown to the user on the Setup Assistant Remote Management pane. Must be between 1 and 125 characters.
- `support_phone_number` (String) The support phone number shown to the user on the Setup Assistant Remote Management pane. Must be between 1 and 50 characters.

### Optional

- `apple_id_disabled` (Boolean) Whether to hide the Apple ID setup pane in Setup Assistant, which gives users the option to sign in with their Apple ID and use iCloud.
- `apple_intelligence_disabled` (Boolean) Whether to hide the Apple Intelligence setup pane in Setup Assistant, where users can configure Apple Intelligence features.
- `apple_pay_disabled` (Boolean) Whether to hide the Apple Pay setup pane in Setup Assistant, which gives users the option to set up Apple Pay on their devices.
- `await_device_configured` (Boolean) Whether devices are locked in Setup Assistant until all enrollment-time configuration is installed (await configuration).
- `dep_onboarding_settings_id` (String) The ID of the Apple ABM/ASM DEP onboarding token (`/deviceManagement/depOnboardingSettings`) that owns this profile. If omitted, it is automatically resolved to the tenant's single Apple ADE/ABM (or ASM) token; if the tenant has more than one Apple token, this must be set explicitly.
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `device_security_group` (String) The ID of the static Microsoft Entra security group to use for enrollment time grouping (the "Device group" tab in the Intune admin center policy wizard). Devices assigned this policy become members of the group as they enroll. This is set via the dedicated `setEnrollmentTimeDeviceMembershipTarget` action on `/deviceManagement/configurationPolicies/{id}` (and cleared via `clearEnrollmentTimeDeviceMembershipTarget` when removed), not via the settings catalog. The group must have the 'Intune Provisioning Client' service principal (AppId: f1346770-5b25-470b-88bd-d5744ab7952c) set as its owner; in some tenants this service principal may appear as 'Intune Autopilot ConfidentialClient'.

~> **Known Microsoft Graph limitation:** as of this writing, `setEnrollmentTimeDeviceMembershipTarget` and `clearEnrollmentTimeDeviceMembershipTarget` return an `Internal Server Error - 500` from the Intune backend (`DeviceConfigV2`) when called with application permissions (client credentials) - the auth flow this provider always uses. The identical request succeeds when made with delegated (signed-in user) permissions, e.g. from the Intune admin center. Until Microsoft resolves this for application permissions, setting `device_security_group` through this provider will fail on `Create` and `Update`; this is a Microsoft Graph service limitation, not a provider defect.
- `diagnostics_disabled` (Boolean) Whether to hide the diagnostics pane in Setup Assistant, where users can opt in to send diagnostic data to Apple.
- `get_started_screen_disabled` (Boolean) Whether to hide the Get Started pane in Setup Assistant.
- `is_default_policy_assignment` (Boolean) Whether this policy is the default visionOS enrollment profile for its `dep_onboarding_settings_id`, set via the dedicated `setDefaultProfile` action. Always reflects the DEP token's actual current default on refresh, regardless of configuration.

~> **No unassign action:** Microsoft Graph does not expose an `unsetDefaultProfile`/`clearDefaultProfile` action - `setDefaultProfile` is the only operation available. Setting this to `false` on a policy that is currently the DEP token's default has no effect on Graph; the next refresh reports `true` again. Only setting a different policy's `is_default_policy_assignment` to `true` changes which profile is the default. A change from `true` to `false` while this policy is still the token's current default can therefore never converge, and the provider rejects the update with a validation error. Promote the replacement policy first - in the same apply, give this policy a `depends_on` for the replacement so the promotion runs first, or apply the promotion separately.
- `location_services_disabled` (Boolean) Whether to hide the Location Services setup pane in Setup Assistant, where users can enable location services on their device.
- `locked_enrollment_enabled` (Boolean) Whether enrollment is locked to the authorized user/device, preventing the MDM profile from being removed before enrollment completes.
- `passcode_disabled` (Boolean) Whether to hide the passcode and password lock pane in Setup Assistant. When shown, users are prompted for a passcode. Always require a passcode for unsecured devices unless access is controlled in some other way (such as through a kiosk mode configuration that restricts the device to one app).
- `privacy_pane_disabled` (Boolean) Whether to hide the privacy setup pane in Setup Assistant.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `screen_time_screen_disabled` (Boolean) Whether to hide the Screen Time pane in Setup Assistant.
- `siri_disabled` (Boolean) Whether to hide the Siri setup pane in Setup Assistant.
- `software_update_screen_disabled` (Boolean) Whether to hide the mandatory software update screen in Setup Assistant.
- `terms_and_conditions_disabled` (Boolean) Whether to hide the Apple terms and conditions pane in Setup Assistant. When shown, users are required to accept them.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `tips_screen_disabled` (Boolean) Whether to hide the Tips pane in Setup Assistant.
- `touch_id_disabled` (Boolean) Whether to hide the biometric (Optic ID) setup pane in Setup Assistant, which gives users the option to set up biometric identification on their devices.

### Read-Only

- `created_date_time` (String) Creation date and time of the policy.
- `id` (String) The unique identifier for this policy.
- `is_assigned` (Boolean) Indicates if the policy is assigned to any scope.
- `last_modified_date_time` (String) Last modification date and time of the policy.
- `platforms` (String) The platforms this policy applies to. Always `visionOS`.
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

- **Settings catalog**: this resource manages visionOS ADE profiles through the modern
  `/deviceManagement/configurationPolicies` settings catalog endpoint.
- **User authentication**: visionOS uses the basic user affinity setting - unlike iOS/iPadOS,
  there is no authentication method choice or await-final-configuration option; the separate
  `await_device_configured` toggle controls whether devices are locked in Setup Assistant until
  enrollment-time configuration is installed.
- **Support department/phone**: `support_department` (1-125 characters) and
  `support_phone_number` (1-50 characters) are required - Microsoft Graph rejects empty values.
- **Enrollment time grouping**: see the Known Issues section above.
- **`is_default_policy_assignment` is a singleton per DEP token**: only one visionOS enrollment
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
terraform import microsoft365_graph_beta_device_management_visionos_device_enrollment_policy.example 00000000-0000-0000-0000-000000000000
```
